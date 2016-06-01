package flow

type Parameter struct {
    name string
    t    Type
    addr Address
}
func NewParameter(name string, t Type, addr Address) Parameter {
    return Parameter{name: name, t: t, addr: addr}
}

type InstanceMap map[Address]FunctionBlock
type BlockMap map[string]FunctionBlock
type EdgeMap map[Parameter][]Parameter
type ParamMap map[string]Parameter

type Graph struct {
    name        string
    nodes       InstanceMap
    edges       EdgeMap
    inputs      ParamMap  // FIXME: Can these be ParamValues instead?
    outputs     ParamMap  // FIXME: Can these be ParamValues instead?
}

func NewGraph(name string, inputs, outputs ParamTypes) Graph {
    self_addr := NewAddress(-1, name)
    ins, outs := make(ParamMap), make(ParamMap)
    for name, t := range inputs {
        ins[name] = NewParameter(name, t, self_addr)
    }
    for name, t := range outputs {
        outs[name] = NewParameter(name, t, self_addr)
    }
    return Graph{name: name, nodes: make(InstanceMap), edges: make(EdgeMap), inputs: ins, outputs: outs}
}

func (g Graph) AddNode(blk FunctionBlock, addr Address) (ok bool) {
    _, exists := g.nodes[addr]
    if !exists {
        g.nodes[addr] = blk
        ok = true
    } else {
        ok = false
    }
    return
}

// out_addr[out_param_name] -> in_addr[in_param_name]
func (g Graph) AddEdge(out_addr Address, out_param_name string,
                       in_addr Address, in_param_name string) (ok bool) {
    ok = false                                // Default assume something went wrong
    out_blk, out_exists := g.nodes[out_addr]  // Check if out_addr exists in nodes, and get it's parameters
    in_blk, in_exists   := g.nodes[in_addr]   // Check if in_addr exists in nodes, and get it's parameters
    
    if in_exists && out_exists {              // If both exist
        _, out_params := out_blk.GetParams()  // Get the output parameters of out_blk
        in_params, _  := in_blk.GetParams()   // Get the input parameters of in_blk
        t_out, out_p_exists := out_params[out_param_name]   // Get the types of the parameters
        t_in, in_p_exists   := in_params[in_param_name]     // ...
        in_param  := Parameter{name: in_param_name, addr: in_addr, t: t_in}    // Create an parameter struct for indexing
        out_param := Parameter{name: out_param_name, addr: out_addr, t: t_out} // ...
        if (in_p_exists && out_p_exists && t_out == t_in) {                    // Check if both exists and types are the same
            g.edges[out_param] = append(g.edges[out_param], in_param)          // Append the new link to the edges under the out_param
            ok = true                                                          // Return true
        }
    }
    return
}

// self[self_param_name] -> in_addr[in_param_name]
func (g Graph) LinkIn(self_param_name string, in_param_name string, in_addr Address) (ok bool) {
    self_addr := NewAddress(-1, g.GetName())
    return g.AddEdge(self_addr, self_param_name,
                       in_addr, in_param_name)
}

// out_addr[out_param_name] -> self[self_param_name]
func (g Graph) LinkOut(out_addr Address, out_param_name string, self_param_name string) (ok bool) {
    self_addr := NewAddress(-1, g.GetName())
    return g.AddEdge(out_addr, out_param_name,
                       self_addr, self_param_name)
}

// Returns copies of all parameters in FunctionBlock
func (g Graph) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    ins, outs := make(ParamTypes), make(ParamTypes)
    for name, param := range g.inputs {
        ins[name] = param.t
    }
    for name, param := range g.outputs {
        outs[name] = param.t
    }
    return ins, outs
}

// Returns a copy of FunctionBlock's InstanceId
func (g Graph) GetName() string {return g.name}

func (g Graph) Run(inputs ParamValues,
                   outputs chan DataOut,
                   stop chan bool,
                   err chan FlowError, id InstanceID) {
    ADDR := NewAddress(-1, g.GetName())
    
    // Check types to ensure inputs are the type defined in input parameters
    if !g.checkTypes(inputs, g.inputs) {
      err <- FlowError{Ok: false, Info: "Inputs are impropper types.", Addr: ADDR}
      return
    }

    // Declare variables
    running       := true                           // When this turns to false, the process stops
    all_waiting   := make(InstanceMap)              // A map of all blocks waiting for inputs
    all_running   := make(InstanceMap)              // A map of all blocks we are waiting to return data
    all_suspended := make(InstanceMap)              // A map of all blocks with blocked outputs waiting to be shifted down the graph
    all_data_in   := make(map[Address]ParamValues)  // A map of all data waiting at the inputs of each block
    all_data_out  := make(map[Address]ParamValues)  // A map of all data waiting at the outputs of each block
    all_stops     := make(map[Address](chan bool))  // A map of all stop channels passed to each running block
    flow_errs     := make(chan FlowError)           // A channel passed to each running block to send back errors
    data_flow     := make(chan DataOut)             // A channel passed to each running block to send back return data
    graph_out     := make(ParamValues)              // The output for the entire graph

    // Create some functions for simplified code structure
    
    // Stops all children blocks
    stopAll := func() {
        // Push stop down to all subfunctions
        for _, val := range all_stops {
            val <- true
        }
        running = false  // Stop this block as well
    }
    
    // Pushes an error up the pipeline and stops all blocks
    pushError := func(info string) {
        flow_errs <- FlowError{Ok: false, Info: info, Addr: ADDR}
        stopAll()
    }
    
    // Adds data to all_data_in, creates ParamValues struct if necessary.
    handleInput := func(param Parameter, val interface{}) (ok bool) {
        if CheckType(param.t, val) {                                   // Check the type of param relative to val
            ok = true                                                  // If it is ok, then return true
            addr := param.addr                                         // Get the address of the parameter
            val, exists := all_data_in[addr]                           // Check if the addr already exists
            if !exists {                                               // If the parameter exists
                all_data_in[param.addr] = ParamValues{param.name: val} // Create new map and add addr,val
            } else {                                                   // Otherwise
                all_data_in[param.addr][param.name] = val              // Add addr,val to the preexisting map
            }
        } else {                                                       // If type is not ok.
            ok = false                                                 // Return false
            pushError("Input is not the right type.")
        }
        return
    }
    
    // Adds data to graph_out, pushes error if type is wrong or if out already set
    handleOutput := func(param Parameter, val interface{}) (ok bool) {
        if CheckType(param.t, val) {                            // Check the type of param relative to val
            ok = true                                           // If it is ok, then return true
            val, exists := graph_out[param.name]                // Check if the addr already exists
            if !exists {                                        // If the parameter exists
                graph_out[param.name] = val                     // Add it to the map
            } else {                                            // Otherwise
                ok = false                                      // Return false
                pushError("Output already existed")             // Push an error
            }
        } else {                                                // If type is not ok.
            ok = false                                          // Return false
            pushError("Output is not the right type.")          // Push an error
        }
        return
    }
    
    // Iterates through all given inputs and adds them to method's all_data_ins.
    loadvars := func() {
        // Main loop
        for name, val := range inputs {         // Iterate through the names/values given in function parameters
            param, exists := g.inputs[name]     // Lookup this parameter in the graph inputs
            if exists {                         // If the parameter does exist
                handleInput(param, val)         // Add the value to all_data_in
            } else {                            // Otherwise, error
                pushError("Input parameter does not exist.")
                return
            }
        }
    }

    // Iterate through all blocks that are waiting
    // to see if all of their inputs have been set.
    // If so, it runs them...
    // Deleting them from waiting, and placing them in running.
    checkWaiting := func() (ran bool) {
        
        // Runs a block and moves it from waiting to running, catalogues all channels
        blockRun := func(addr Address, blk FunctionBlock, f_in ParamValues) {
            f_stop := make(chan bool)                                // Make a stop channel
            go blk.Run(f_in, data_flow, f_stop, flow_errs, addr.id)           // Run the block as a goroutine
            delete(all_waiting, addr)                                // Delete the block from waiting
            all_running[addr] = blk                                  // Add the block to running
            all_stops[addr] = f_stop                                 // Add stop channel to map
        }
        
        // Main loop
        ran = false
        for addr, blk := range all_waiting {                // Iterate through all waiting
            in_params, _ := blk.GetParams()                 // Get the inputs from the block
            in_vals, val_exists := all_data_in[addr]        // Get their stored values
            if val_exists {                                 // If any values are set
                ready := true                               // Declare a variable to log whether or not the block is ready, default true
                f_in  := make(ParamValues)                  // Declare a variable to use as the block input after verifying all it's inputs
                for name, t := range in_params {            // Check if all parameters are ready by iterating through them
                    val, param_set := in_vals[name]         // Get the value of the parameter from in_vals
                    if !param_set {                         // If parameter not set, we set ready to false and break
                        ready = false
                        break
                    } else {                                // Otherwise
                        if CheckType(t, val) {              // Check the type of the found value against the parameter type
                            f_in[name] = val                // If it is valid, we begin setting our input parameter
                        } else {                            // If the type is wrong, we throw an error, and everything stops as the graph is broken
                            pushError("Input parameter is not the right type.")
                            return
                        }
                    }
                }
                
                if ready {                     // Check to see that ready made it through all parameters without changing values
                    blockRun(addr, blk, f_in)  // If so, then run the block.
                    ran = true                 // Indicate that we have indeed ran at least one block.
                }
            }
        }
        return
    }
    
    // Monitor the data_flow channel for the next incoming data.
    // Blocks until some packet is received, either data, stop, or error 
    checkRunning := func() (found bool) {
        
        // Adds data to all_data_out
        // Deletes from all_running and adds to all_waiting
        handleOutput := func(vals DataOut) {
            V    := vals.Values                     // Get values for easy access
            addr := vals.Addr                       // Get address for easy access
            blk  := all_running[addr]               // Retrieve the block from running
            _, out_params := blk.GetParams()        // Get blk output param types
            if CheckTypes(V, out_params) {          // If the types are all compatable
                all_data_out[addr] = V              // Set the output data
                delete(all_running, addr)           // Delete block from running
                all_suspended[addr] = blk           // Add block to suspended
                delete(all_stops, addr)             // Delete channels
            } else {                                // If types are incompatible, push an error, graph is broken
                pushError("Output parameter not the right type.")
            }
        }
        
        // Wait for some data input
        found = false
        if running {                            // Do not begin waiting if this block is not running
            select {                            // Wait for input
                case data := <- data_flow:      // If there is data input
                    handleOutput(data)          // Handle it
                    found = true                // Declare data was found
                case e := <- flow_errs:         // If it is an error
                    if !e.Ok {                  // Check to see if it's dangerous
                        pushError(e.Info)       // If it is dangerous, push the error
                    }
                case <-stop:                    // If a stop is received
                    stopAll()                   // Stop all processes
            }
        }
        return
    }
    
    // Shift outputs to inputs based on graph, and also to graph_out
    shiftData := func() (success bool) {
        success = false                                         // Test that at least one item was restored
        var restore_lst []Address                               // Create a list of Addresses to move from suspended to waiting after move
        for addr, blk := range all_suspended {                  // Iterate through all suspended blocks
            flag := false                                       // Create a flag to indicate whether or not this block was left incomplete
            _, out_p_lst := blk.GetParams()                     // Get the parameters of the block
            blk_out, blk_exists := all_data_out[addr]           // Check that block has any outputs at all
            if blk_exists {                                     // If so...
                for name, t := range out_p_lst {                     // Iterate through all output parameters
                    out_p := Parameter{addr: addr, name: name, t: t} // Create the full parameter structure for indexing
                    val, v_exists := blk_out[name]                   // Check that the value exists
                    if v_exists {                               // If so...
                        p_flag := false                         // Create a flag to indicate whether or not this parameter was left incomplete
                        for _, in_p := range g.edges[out_p] {   // Iterate through all linked inputs
                            if in_p.addr.name != g.name {  // Don't handle outputs that relate to the graph outputs
                                ok := handleInput(in_p, val)    // Check the type and add it to all_data_in
                                if !ok {return}                 // Handle errors
                            } else {                            // If the input parameter is this graph's output
                                ok := handleOutput(in_p, val)   // Check the type and add it to graph_out
                                if !ok {return}                 // Handle errors
                                p_flag = true                   // Don't restore if the address is self
                            }    
                        }
                        if !p_flag {                            // If no flags were raised
                            delete(all_data_out[addr], name)    // Delete the parameter from all_data_out
                        } else { flag = true }                  // If any were then pass the flag up to the block flag
                    }
                }
            }
            if !flag {
                restore_lst = append(restore_lst, addr)         // Add this to list to restore if all items were moved successfully.
            }                
        }
        
        // Restore the successful blocks to waiting, check again just to make sure all outputs were removed.
        // FIXME: Is the second check necessary?
        for _, addr := range restore_lst {      // Iterate through all the successful blocks
            val, exists := all_data_out[addr]   // Check if any of the outputs still exist
            if (len(val) == 0) || !exists {     // If it's empty or non-existant. FIXME: Is this valid?
                blk := all_suspended[addr]      // Get the block from suspended
                delete(all_suspended, addr)     // Delete it from suspended
                all_waiting[addr] = blk         // And add it to waiting
                success = true                  // Indicate at least one was successfully moved
            }
        }
        return
    }
    
    checkDone := func() (done bool) {
        for name, _ := range g.outputs {                // Iterate through all output parameters
            _, exists := graph_out[name]                // Check if each parameters exist in graph_out
            if !exists {return false}                   // If any does not exist, immediately return false
        }
        return true                                     // If you pass the loop, all exist, return true
    }
    
    // Main Logic
    for running {
        loadvars()                   // Begin by loading all inputs into the data maps and adding all blocks to waiting
        checkWaiting()               // Then run waiting blocks that have all inputs availible
        checkRunning()               // Then, wait for some return on a running block
        done  := checkDone()         // See if the output map contains enough data to be done
        if done {                    // If we are done
            stopAll()                // Stop all processes
            outADDR := NewAddress(id, g.GetName())
            outputs <- DataOut{Addr: outADDR, Values: graph_out} // And return our outputs
        } else {                     // If we are not done
            shiftData()              // Try to shift outputs to linked inputs
        }
    }
}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func (g Graph) checkTypes(values ParamValues, params ParamMap) (ok bool) {
    for name, param := range params {
        val, exists := values[name]
        if !exists || !CheckType(param.t, val) {
            return false
        }
    }
    return true
}