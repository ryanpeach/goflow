package flow

import (
    "os"
    "log"
    "io/ioutil"
)

// Source: http://changelog.ca/log/2015/03/09/golang
func CreateLogger(logMode, tag string) *log.Logger {
    switch logMode {
        case "none":
            out := ioutil.Discard
            return log.New(out, tag, log.Lshortfile)
        case "screen":
            out := os.Stdout
            return log.New(out, tag, log.Lshortfile)
        default:
            out, err := os.Create(logMode)
            if nil != err {
                panic(err.Error())
            } else {
                return log.New(out, tag, log.Lshortfile)
            }
    }
}

type InstanceMap map[Address]FunctionBlock
type EdgeMap map[ParamAddress][]ParamAddress
type BlockMap map[string]FunctionBlock
type ParamMap map[string]ParamAddress
type ParamLstMap map[string][]ParamAddress

type Graph struct {
    // Block Data
    name        string
    nodes       InstanceMap
    edges       EdgeMap

    // Block Inputs
    infeed      ParamLstMap    // Connects name of a param in inputs to a ParamAddress of some parameter in some node
    outfeed     ParamLstMap    // Connects name of a param in outputs to a ParamAddress of some parameter in some node
    inputs      ParamTypes
    outputs     ParamTypes
}

func NewGraph(name string, inputs, outputs ParamTypes) Graph {
    nodes, edges    := make(InstanceMap), make(EdgeMap)
    infeed, outfeed := make(ParamLstMap), make(ParamLstMap)
    return Graph{name, nodes, edges, infeed, outfeed, inputs, outputs}
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
    //logger.Println("Adding Edge: ", out_addr, out_param_name, " -> ", in_addr, in_param_name)
    ok = false
    out_param, out_exists := g.FindParam(out_param_name, out_addr)  // Get the output parameters of out_blk
    in_param, in_exists := g.FindParam(in_param_name, in_addr)   // Get the input parameters of in_blk
    if in_exists && out_exists {              // If both exist
        if CheckCompatibility(in_param.T, out_param.T) && in_param.is_input && !out_param.is_input {
            g.edges[out_param] = append(g.edges[out_param], in_param)          // Append the new link to the edges under the out_param
            //logger.Println("Edge Added: ", g.edges[out_param])
            ok = true
        }
    }
    return
}

func (g Graph) FindParam(name string, addr Address) (param ParamAddress, exists bool) {
    in_params, out_params := g.nodes[addr].GetParams()
    in_t, in_exists := in_params[name]
    out_t, out_exists := out_params[name]
    switch {
        case in_exists == out_exists:
            exists = false
        case in_exists:
            param = ParamAddress{name, addr, in_t, true}
            exists = true
        case out_exists:
            param = ParamAddress{name, addr, out_t, false}
            exists = true
    }
    return
}

// self[self_param_name] -> in_addr[in_param_name]
func (g Graph) LinkIn(self_param_name string, in_param_name string, in_addr Address) (ok bool) {
    param, exists := g.FindParam(in_param_name, in_addr)
    if exists && CheckCompatibility(g.inputs[self_param_name], param.T) {
        g.infeed[self_param_name] = append(g.infeed[self_param_name], ParamAddress{in_param_name, in_addr, param.T, true})
        return true
    } else {
        return false
    }
}

// out_addr[out_param_name] -> self[self_param_name]
func (g Graph) LinkOut(out_addr Address, out_param_name string, self_param_name string) (ok bool) {
    param, exists := g.FindParam(out_param_name, out_addr)
    if exists && CheckCompatibility(g.outputs[self_param_name], param.T) {
        g.outfeed[self_param_name] = append(g.outfeed[self_param_name], ParamAddress{out_param_name, out_addr, param.T, false})
        return true
    } else {
        return false
    }
}

// Returns copies of all parameters in FunctionBlock
func (g Graph) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    return g.inputs, g.outputs
}

// Returns a copy of FunctionBlock's InstanceId
func (g Graph) GetName() string {return g.name}

func (g Graph) Run(inputs ParamValues,
                   outputs chan DataOut,
                   stop chan bool,
                   err chan FlowError, id InstanceID) {
    logger := CreateLogger("none", "[INFO]")
    
    ADDR := Address{g.GetName(), id}
    logger.Println("Running Graph: ", ADDR)
    
    // Check types to ensure inputs are the type defined in input parameters
    if !CheckTypes(inputs, g.inputs) {
      err <- FlowError{Ok: false, Info: "Inputs are impropper types.", Addr: ADDR}
      return
    }

    // Declare variables
    running       := true                           // When this turns to false, the process stops
    all_waiting   := make(InstanceMap)              // A map of all blocks waiting for inputs
    all_running   := make(InstanceMap)              // A map of all blocks we are waiting to return data
    all_suspended := make(InstanceMap)              // A map of all blocks with blocked outputs waiting to be shifted down the graph
    all_data_in   := make(map[ParamAddress]interface{})  // A map of all data waiting at the inputs of each block
    all_data_out  := make(map[ParamAddress]interface{})  // A map of all data waiting at the outputs of each block
    all_stops     := make(map[Address](chan bool))  // A map of all stop channels passed to each running block
    flow_errs     := make(chan FlowError)           // A channel passed to each running block to send back errors
    data_flow     := make(chan DataOut)             // A channel passed to each running block to send back return data
    graph_out     := make(ParamValues)              // The output for the entire graph

    // Put all nodes in waiting
    logger.Println("Putting all nodes in waiting.")
    for addr, blk := range g.nodes {
        all_waiting[addr] = blk
    }

    // Create some functions for simplified code structure
    logger.Println("Defining Functions")
    // Stops all children blocks
    stopAll := func() {
        logger.Println("Stopping all.")
        // Push stop down to all subfunctions
        for _, val := range all_stops {
            val <- true
        }
        running = false  // Stop this block as well
    }
    
    // Pushes an error up the pipeline and stops all blocks
    pushError := func(info string) {
        logger.Println("Pushing Error: ", info)
        flow_errs <- FlowError{Ok: false, Info: info, Addr: ADDR}
        stopAll()
    }
    
    // Adds data to all_data_in, creates ParamValues struct if necessary.
    handleInput := func(param ParamAddress, val interface{}) (ok bool) {
        logger.Println("Handling Inputs.")
        logger.Println("Param: ", param, "Val: ", val)
        _ , param_exists := all_data_in[param]                        // Get the input value and check if it exists
        if CheckType(param.T, val) || !param_exists {                  // Check the type of param relative to val and check if it exists
            ok = true                                                  // If it is ok, then return true
            all_data_in[param] = val                                   // Add addr,val to the preexisting map
        } else {                                                       // If type is not ok or param is not in all_data_in
            ok = false                                                 // Return false
            pushError("Input is not the right type.")
        }
        logger.Println("Check: ", all_data_in[param], all_data_in[param] == val)
        return
    }
    
    // Adds data to graph_out, pushes error if type is wrong or if out already set
    // Adds data to all_data_out
    // Deletes from all_running and adds to all_waiting
    handleOutput := func(vals DataOut) {
        logger.Println("Handling Output: ", vals)
        V    := vals.Values                     // Get values for easy access
        addr := vals.Addr                       // Get address for easy access
        blk  := all_running[addr]               // Retrieve the block from running
        _, out_params := blk.GetParams()        // Get blk output param types
        if CheckTypes(V, out_params) {          // If the types are all compatable
            for param_name, t := range out_params {
                param := ParamAddress{param_name, addr, t, false}
                val, val_exists := V[param_name] // Set the output data
                if val_exists {                  // Only set val if it exists
                    all_data_out[param] = val
                } else {
                    pushError("All Data Output not present.")
                    return
                }
            }
            delete(all_running, addr)           // Delete block from running
            all_suspended[addr] = blk           // Add block to suspended
            delete(all_stops, addr)             // Delete channels
        } else {                                // If types are incompatible, push an error, graph is broken
            pushError("Output parameter not the right type.")
        }
    }
    
    // Iterates through all given inputs and adds them to method's all_data_ins.
    loadvars := func() {
        // Main loop
        logger.Println("Loading Variables.")
        for name , param_lst := range g.infeed {         // Iterate through the names/values given in function parameters
            for _, node_param := range param_lst {
                val, exists := inputs[name]          // Lookup this parameter in the graph inputs
                if exists {                          // If the parameter does exist
                    handleInput(node_param, val)     // Add the value to all_data_in
                    logger.Println(node_param, val, all_data_in[node_param])
                } else {                             // Otherwise, error
                    pushError("Input parameter does not exist.")
                    return
                }
            }
        }
    }

    // Iterate through all blocks that are waiting
    // to see if all of their inputs have been set.
    // If so, it runs them...
    // Deleting them from waiting, and placing them in running.
    checkWaiting := func() (ran bool) {
        logger.Println("Checking Waiting.")
        
        //logger.Println("Data In: ", all_data_in)
        // Runs a block and moves it from waiting to running, catalogues all channels
        blockRun := func(addr Address, blk FunctionBlock, f_in ParamValues) {
            logger.Println("Running ", blk.GetName())
            f_stop := make(chan bool)                                // Make a stop channel
            go blk.Run(f_in, data_flow, f_stop, flow_errs, addr.ID)  // Run the block as a goroutine
            delete(all_waiting, addr)                                // Delete the block from waiting
            all_running[addr] = blk                                  // Add the block to running
            all_stops[addr] = f_stop                            // Add stop channel to map
            
            // Delete all inputs from all_data_in
            for param_name, _ := range f_in {
                param, _ := g.FindParam(param_name, addr)
                delete(all_data_in, param)
            }
            
        }
        
        // Main loop
        ran = false
        for addr, blk := range all_waiting {                // Iterate through all waiting
            in_params, _ := blk.GetParams()                 // Get the inputs from the block
            f_in := make(ParamValues)
            ready := true
            for param_name, t := range in_params {
                param := ParamAddress{param_name, addr, t, true}
                in_val, val_exists := all_data_in[param]        // Get their stored values
                if val_exists {
                    if !CheckType(t, in_val) {
                        logger.Println(param_name, t, in_val, val_exists)
                        pushError("Input parameter is not the right type.")
                        return false
                    } else {
                        f_in[param_name] = in_val
                    }
                } else {
                    ready = false
                }
            }
            if ready {
                blockRun(addr, blk, f_in)  // If so, then run the block.
                ran = true                 // Indicate that we have indeed ran at least one block.
            }
        }
        return ran
    }
    
    // Monitor the data_flow channel for the next incoming data.
    // Blocks until some packet is received, either data, stop, or error 
    checkRunning := func() (found bool) {
        logger.Println("Checking running")
        // Wait for some data input
        found = false
        done := false
        for running && !done {                  // Do not begin waiting if this block is not running
            logger.Println("Waiting for data input.")
            select {                            // Wait for input
                case data := <- data_flow:      // If there is data input
                    handleOutput(data)          // Handle it
                    found = true                // Declare data was found
                    done = true
                case e := <- flow_errs:         // If it is an error
                    logger.Println("Error Returned: ", e)
                    if !e.Ok {                  // Check to see if it's dangerous
                        pushError(e.Info)       // If it is dangerous, push the error
                        done = true
                    }
                case <-stop:                    // If a stop is received
                    stopAll()                   // Stop all processes
                    done = true
            }
        }
        return
    }
    
    // Iterate through all blocks that are suspended
    // to see if all of their outputs have been set.
    // If so, it runs them...
    // Deleting them from waiting, and placing them in running.
    checkSuspended := func() {
        for addr, blk := range all_suspended {                         // Iterate through all suspended blocks
            _, out_p_map := blk.GetParams()                            // Get the parameters of the block
            ready := true
            for name, t := range out_p_map {
                param := ParamAddress{name, addr, t, false}
                _, exists := all_data_out[param]
                if exists {
                    ready = false
                }
            }
            if ready {
                delete(all_suspended, addr)
                all_waiting[addr] = blk
            }
        }
    }
    
    // Shift outputs to inputs based on graph, and also to graph_out
    shiftData := func() (success bool) {
        
        // Shift data from all_data_out to all_data_in
        logger.Println("Shifting Data")
        for out_p, val := range all_data_out {
            logger.Println("HERE", out_p, g.edges[out_p])
            for _, in_p := range g.edges[out_p] {              // Iterate through all linked inputs
                ok := handleInput(in_p, val)                   // Check the type and add it to all_data_in
                if ok {
                    delete(all_data_out, out_p)                // Delete the parameter from all_data_out
                } else {
                    return false                               // Handle errors
                }    
            }
        }
        
        // Shift data to the graph_out by iterating through the outfeed
        logger.Println("Shifting Data to Output")
        for self_param_name, param_lst := range g.outfeed {
            for _, node_param := range param_lst {
                val, exists := all_data_out[node_param]
                if exists {
                    graph_out[self_param_name] = val
                    logger.Println("Graph Out: ", val, graph_out[self_param_name])
                }
            }
        }
        return true
    }
    
    // Returns true if all parameters in g.outputs referenced in graph_out
    checkDone := func() (done bool) {
        logger.Println("Checking Done")
        for name, _ := range g.outfeed {                // Iterate through all output parameters
            _, exists := graph_out[name]                // Check if each parameters exist in graph_out
            if !exists {return false}                   // If any does not exist, immediately return false
        }
        logger.Println("DONE!!!")
        return true                                     // If you pass the loop, all exist, return true
    }
    
    logger.Println("Done Defining Functions")
    
    // Main Logic
    loadvars()                   // Begin by loading all inputs into the data maps and adding all blocks to waiting
    logger.Println("Done Loading")
    for running {
        checkWaiting()               // Then run waiting blocks that have all inputs availible
        checkRunning()               // Then, wait for some return on a running block
        shiftData()                  // Try to shift outputs to linked inputs
        if checkDone() {             // See if the output map contains enough data to be done
            stopAll()                // Stop all processes
            outputs <- DataOut{Addr: ADDR, Values: graph_out} // And return our outputs
        } else {                     // If we are not done
            checkSuspended()             // Then, move from the outputs to the inputs and graph_out, and make methods waiting again
        }    
    }
}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func (g Graph) checkTypes(values ParamValues, params ParamMap) (ok bool) {
    for name, param := range params {
        val, exists := values[name]
        if !exists || !CheckType(param.T, val) {
            return false
        }
    }
    return true
}