package flow

type InstanceMap map[Address]FunctionBlock
type BlockMap map[string]FunctionBlock
type EdgeMap map[Parameter][]Parameter
type ReadyMap map[string]bool

type Graph struct {
    addr        Address
    nodes       InstanceMap
    edges       EdgeMap
    inputs      ParamMap
    outputs     ParamMap
}

func NewGraph(name string, id InstanceID, nodes InstanceMap, edges EdgeMap,
              inputs ParamMap, outputs ParamMap) FunctionBlock {
    return Graph{name: name, id: id, nodes: nodes, edges: edges, inputs: inputs, outputs: outputs}
}

// Returns copies of all parameters in FunctionBlock
func (g Graph) GetParams() (inputs ParamMap, outputs ParamMap) {
    return g.inputs, g.outputs
}

// Returns a copy of FunctionBlock's InstanceId
func (g Graph) GetAddr() Address {return g.addr}

func (g Graph) Run(inputs ParamValues,
                   outputs chan ParamValues,
                   stop chan bool,
                   err chan FlowError) {
    // Check types to ensure inputs are the type defined in input parameters
    if !g.checkTypes(inputs, g.input_types) {
      return FlowError{Ok: false, Info: "Inputs are impropper types.", Addr: g}
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
        for name, val := range all_stops {
            val <- true
        }
        running = false  // Stop this block as well
    }
    
    // Pushes an error up the pipeline and stops all blocks
    pushError := func(info string) {
        all_errs <- FlowError{Ok: false, Info: info, Addr: g.GetAddr()}
        stopAll()
    }
    
    // Adds data to all_data_in, creates ParamValues struct if necessary.
    handleInput := func(param Parameter, val interface{}) (ok bool) {
        if CheckType(param, val) {                              // Check the type of param relative to val
            ok = true                                           // If it is ok, then return true
            addr := param.GetAddr()                             // Get the address of the parameter
            val, exists := all_data_in[addr]                    // Check if the addr already exists
            if !exists {                                        // If the parameter exists
                all_data_in[addr] = ParamValues{addr.name: val} // Create new map and add addr,val
            } else {                                            // Otherwise
                all_data_in[addr][addr.name] = val              // Add addr,val to the preexisting map
            }
        } else {                                                // If type is not ok.
            ok = false                                          // Return false
            pushError("Input is not the right type.")
        }
        return
    }
    
    // Adds data to graph_out, pushes error if type is wrong or if out already set
    handleOutput := func(param Parameter, val interface{}) (ok bool) {
        if CheckType(param, val) {                              // Check the type of param relative to val
            ok = true                                           // If it is ok, then return true
            name := param.GetName()                             // Get the address of the parameter
            val, exists := graph_out[name]                      // Check if the addr already exists
            if !exists {                                        // If the parameter exists
                graph_out[name] = val                           // Add it to the map
            } else {                                            // Otherwise
                ok = false                                      // Return false
                pushError("Output already existed")             // Push an error
            }
        } else {                                                // If type is not ok.
            ok = false                                          // Return false
            pushError("Output is not the right type.")          // Push an error
        }
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
        blockRun := func(blk FunctionBlock, f_in ParamValues) {
            addr := blk.GetAddr()                           // Get address for indexing
            f_stop := make(chan bool)                       // Make a stop channel
            go blk.Run(f_in, data_flow, f_stop, flow_errs)  // Run the block as a goroutine
            delete(all_waiting[addr])                       // Delete the block from waiting
            all_running[addr] = blk                         // Add the block to running
            all_stops[addr] = f_stop                        // Add stop channel to map
        }
        
        // Main loop
        ran = false
        for addr, blk := range all_waiting {                // Iterate through all waiting
            in_params, _ := blk.GetParams()                 // Get the inputs from the block
            in_vals, val_exists := all_data_in[addr]        // Get their stored values
            if val_exists {                                 // If any values are set
                ready := true                               // Declare a variable to log whether or not the block is ready, default true
                f_in  := make(ParamValues)                  // Declare a variable to use as the block input after verifying all it's inputs
                for _, param := range in_params {           // Check if all parameters are ready by iterating through them
                    name := param.GetName()                 // Get the name of the parameter to use as a key
                    val, param_set := in_vals[name]         // Get the value of the parameter from in_vals
                    if !param_set {ready = false; break;}   // If parameter not set, we set ready to false and break
                    else {                                  // Otherwise
                        if CheckType(param, val) {          // Check the type of the found value against the parameter type
                            f_in[name] = val                // If it is valid, we begin setting our input parameter
                        } else {                            // If the type is wrong, we throw an error, and everything stops as the graph is broken
                            pushError("Input parameter is not the right type.")
                            return
                        }
                    }
                }
                
                if ready {               // Check to see that ready made it through all parameters without changing values
                    blockRun(blk, f_in)  // If so, then run the block.
                    ran = true           // Indicate that we have indeed ran at least one block.
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
        handleOutput := func(vals DataOut} {
            V    := vals.Values                    // Get values for easy access
            addr := vals.Addr                      // Get address for easy access
            blk  := all_running[addr]              // Retrieve the block from running
            if CheckTypes(V, all_outputs[addr]) {  // If the types are all compatable
                all_data_out[addr] = V             // Set the output data
                delete(all_running[addr])          // Delete block from running
                all_suspended[addr] = blk          // Add block to suspended
                delete(all_stops[addr])            // Delete channels
            } else {                               // If types are incompatible, push an error, graph is broken
                pushError("Output parameter not the right type.")
            }
        }
        
        // Wait for some data input
        found = false
        select running {                    // Do not begin waiting if this block is not running
            case data := <- data_flow:      // If there is data input
                handleOutput(data)          // Handle it
                found = true                // Declare data was found
            case e := <- flow_err:          // If it is an error
                if !e.Ok {                  // Check to see if it's dangerous
                    pushError(e.Info)       // If it is dangerous, push the error
                }
            case <-stop:                    // If a stop is received
                allStop()                   // Stop all processes
        }
        return
    }
    
    // Shift outputs to inputs based on graph, and also to graph_out
    shiftData := func() (success bool) {
        success = false                                         // Test that at least one item was restored
        restore_lst = make([]Address)                           // Create a list of Addresses to move from suspended to waiting after move
        for addr, blk := range all_suspended {                  // Iterate through all suspended blocks
            flag = false                                        // Create a flag to indicate whether or not this block was left incomplete
            in_p_lst, out_p_lst := blk.GetParams()              // Get the parameters of the block
            blk_out, blk_exists := all_data_out[addr]           // Check that block has any outputs at all
            if blk_exists {                                     // If so...
                for _, out_p := range out_p_lst {               // Iterate through all output parameters
                    name := out_p.GetName()                     // Get the parameters name
                    val, v_exists := blk_out[name]              // Check that the value exists
                    if v_exists {                               // If so...
                        p_flag := false                         // Create a flag to indicate whether or not this parameter was left incomplete
                        for _, in_p := range g.edges[out_p] {   // Iterate through all linked inputs
                            if in_p.GetAddr() != g.GetAddr() {  // Don't handle outputs that relate to the graph outputs
                                ok := handleInput(in_p, val)    // Check the type and add it to all_data_in
                                if !ok {return}                 // Handle errors
                            } else {                            // If the input parameter is this graph's output
                                ok := handleOutput(in_p, val)   // Check the type and add it to graph_out
                                if !ok {return}                 // Handle errors
                                p_flag = true                   // Don't restore if the address is self
                            }    
                        }
                        if !p_flag := {                         // If no flags were raised
                            delete(all_data_out[addr][name])    // Delete the parameter from all_data_out
                        } else { flag = true }                  // If any were then pass the flag up to the block flag
                    }
                }
            }
            if !flag {
                append(restore_lst, addr)                       // Add this to list to restore if all items were moved successfully.
            }                
        }
        
        // Restore the successful blocks to waiting, check again just to make sure all outputs were removed.
        // FIXME: Is the second check necessary?
        for _, addr := range restore_lst {      // Iterate through all the successful blocks
            val, exists := all_data_out[addr]   // Check if any of the outputs still exist
            if (len(val) == 0) || !exists {     // If it's empty or non-existant. FIXME: Is this valid?
                delete(all_suspended[addr])     // Delete it from suspended
                all_waiting[addr] = blk         // And add it to waiting
                success = true                  // Indicate at least one was successfully moved
            }
        }
        return
    }
    
    checkDone := func() (done bool) {
        for name, param := range g.outputs {            // Iterate through all output parameters
            _, exists := graph_out[name]                // Check if each parameters exist in graph_out
            if !exists {return false}                   // If any does not exist, immediately return false
        }
        return true                                     // If you pass the loop, all exist, return true
    }
    
    // Main Logic
    for running {
        loadvars()                   // Begin by loading all inputs into the data maps and adding all blocks to waiting
        checkWaiting()               // Then run waiting blocks that have all inputs availible
        found := checkRunning()      // Then, wait for some return on a running block
        done  := checkDone()         // See if the output map contains enough data to be done
        if done {                    // If we are done
            allStop()                // Stop all processes
            outputs <- graph_out     // And return our outputs
        } else {                     // If we are not done
            shiftData()              // Try to shift outputs to linked inputs
        }
    }
}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func (g Graph) checkTypes(values ParamValues, params ParamTypes) (ok bool) {
    var val interface{}
    for name, kind := range params {
        val, exists = values[name]
        switch x := val.(type) {
            case !exists:
                return false
            case x != kind:
                return false
        }
    }
    return true
}