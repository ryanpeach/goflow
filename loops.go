package flow

const (
    INDEX_NAME = "I"
    DONE_NAME = "DONE"
)

type Loop struct {
    g *Graph
    blk, cnd  Address
    inputs    ParamTypes
    outputs   ParamTypes
    registers NameMap
    initial   ParamValues
}

func NewLoop(name string, inputs, outputs ParamTypes, blk, stop_condition FunctionBlock) (*Loop, *Error) {
    
    // Check that stop_condition has one bool output.
    _, cnd_out := stop_condition.GetParams()
    bool_found := false
    for _, t := range cnd_out {
        if t == Bool {
            bool_found = true
            break
        }
    }
    if !bool_found {
        return nil, &Error{TYPE_ERROR, "Stop Condition has no boolean output."}
    }
    
    // Add Done and Index as outputs and inputs
    g_inputs, g_outputs := CopyTypes(inputs), CopyTypes(outputs)
    g_inputs[INDEX_NAME] = Int
    g_outputs[DONE_NAME] = Bool
    
    // Initialize variables
    regs, inits := make(NameMap), make(ParamValues)
    blk_addr, cnd_addr := Address{blk.GetName(), 0}, Address{stop_condition.GetName(), 0}
    
    // Build Graph
    graph, err1 := NewGraph(name, g_inputs, g_outputs)
    outLoop := Loop{graph, blk_addr, cnd_addr, inputs, outputs, regs, inits}
    err2 := outLoop.g.AddNode(blk, blk_addr)
    err3 := outLoop.g.AddNode(stop_condition, cnd_addr)
    
    // Output handling errors
    switch {
        case err1 != nil:
            return nil, err1
        case err2 != nil:
            return nil, &Error{err2.Class, "Blk could not be added to graph."}
        case err3 != nil:
            return nil, &Error{err3.Class, "Stop Condition could not be added to graph."}
    }
    return &outLoop, nil
    
}

// FunctionBlock Fields
func (l Loop) GetName() string {return l.g.GetName()}
func (l Loop) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    return l.inputs, l.outputs
}

// Inhereted Fields
func (l Loop) AddEdge(out_addr Address, out_param_name string,
                       in_addr Address, in_param_name string) (ok bool) {return l.g.AddEdge(out_addr, out_param_name, in_addr, in_param_name)}
func (l Loop) LinkIn(self_param_name string, in_param_name string, in_addr Address) (ok bool) {
    return l.g.LinkIn(self_param_name, in_param_name, in_addr)
}
func (l Loop) LinkOut(out_addr Address, out_param_name string, self_param_name string) (ok bool) {
    return l.g.LinkOut(out_addr, out_param_name, self_param_name)
}

// Adds parameter "name" of Type "t" to self as an input (if is_input) or output (if !is_input).
// Also adds it to the graph as a feed.
func (l Loop) AddFeed(name string, t Type, is_input bool) (err *Error) {
    // This is the function to be called twice
    wrapper := func(X ParamTypes) *Error {
        t2, exists := X[name]
        if exists {
            if !CheckSame(t, t2) {
                return &Error{TYPE_ERROR, "This parameter already exists in a different type."}
            } else {
                return &Error{ALREADY_EXISTS_ERROR, "This parameter already exists."}
            }
        } else {
            err := l.g.AddFeed(name, t, is_input)                   // Add the feed to the graph
            if err == nil || err.Class != ALREADY_EXISTS_ERROR {   // If there is no error, or the error is just that the param already exists
                X[name] = t                                         // Then add the feed to this loop
            } else {
                return err                                          // Otherwise return an error
            }
        }
        return nil
    }
    
    if is_input {
        err = wrapper(l.inputs)
    } else {
        err = wrapper(l.outputs)
    }
    return
}
// --------------- Novel Methods --------------

// Connects out_name parameter of the inner graph to in_name parameter of the inner graph.
// Creates a default parameter value for the input
// Assumes a feed_input does not exist for the loop and instead uses a default value
// Can create an out_feed, will return an error if out_feed name is taken but type is different
// Will return an error if the in_name parameter is already connected to a register
func (l Loop) AddDefaultRegister(out_name, in_name string, t Type, init interface{}) *Error {
    err1 := l.g.AddFeed(in_name, t, true) // Create an input if one does not already exist
    err2 := l.g.AddFeed(out_name, t, false)
    if err1 != nil && err1.Class == ALREADY_EXISTS_ERROR {
        delete(l.inputs, in_name) // If it already existed, remove it from loop inputs so a default value may be used instead
    }
    switch {
        case err2 != nil && err2.Class != ALREADY_EXISTS_ERROR:
            return err2
        case !CheckType(t, init):
            return &Error{TYPE_ERROR, "t and init are incompatible types."}
        default:
            _, connected := l.registers[in_name]
            if !connected {
                l.registers[in_name] = out_name
                l.initial[in_name] = init
            } else {
                return &Error{ALREADY_EXISTS_ERROR, "Connection to input already exists."}
            }
    }
    return nil
}

// Will return an error if input in_name does not exist in either loop or graph
// Can create an out_feed, will return an error if out_feed name is taken but type is different
// Will return an error if the in_name parameter is already connected to a register
func (l Loop) AddRegister(out_name, in_name string, t Type) *Error {
    // Check input feed
    _, in_exists1 := l.inputs[in_name]
    _, in_exists2 := l.g.inputs[in_name]
    err := l.g.AddFeed(out_name, t, false)

    // Handle Errors
    switch {
        case !in_exists1 || !in_exists2:
            return &Error{DNE_ERROR, "Input does not exist."}
        case err != nil && err.Class != ALREADY_EXISTS_ERROR:
            return err
        default:
            _, connected := l.registers[in_name]
            if !connected {
                l.registers[in_name] = out_name
            } else {
                return &Error{ALREADY_EXISTS_ERROR, "Connection to input already exists."}
            }
    }
    return nil
}

func (l Loop) Run(inputs ParamValues, outputs chan DataOut, stop chan bool, err chan *FlowError, id InstanceID) {
    // Declare variables
    ADDR     := Address{l.GetName(), id}
    data_out := make(ParamValues)
    all_done := false
    loop_i   := 0
    i_inputs := CopyValues(inputs)
    i_out    := make(chan DataOut)
    i_stop   := make(chan bool)
    i_err    := make(chan *FlowError)
    
    // Copy output values to data_out and i_inputs
    handleOutput := func(out_vals ParamValues) {
        // Copy output values to data_out
        for name, val := range out_vals {
            _, exists := l.outputs[name]
            if exists {data_out[name] = val}
            if name == DONE_NAME {all_done = val.(bool)}
        }
        
        // Copy output values to i_inputs
        for in_name, out_name := range l.registers {
            val, exists := out_vals[out_name]
            if exists {
                i_inputs[in_name] = val
            }
        }
    }
    

    // Check that all inputs are satisfied
    chk_exists := checkInputs(inputs, l.inputs)
    chk_types  := CheckTypes(inputs, l.inputs)
    switch {
        case !chk_exists:
            err <- NewFlowError(DNE_ERROR, "Not all inputs satisfied.", ADDR)
            return
        case !chk_types:
            err <- NewFlowError(TYPE_ERROR, "Inputs are impropper types.", ADDR)
            return
    }
    
    // Copy initial values
    for name, val := range l.initial {
        i_inputs[name] = val
    }
    
    // Run main loop until done is set
    for !all_done {
        i_inputs[INDEX_NAME] = loop_i                  // Update index input
        go l.g.Run(i_inputs, i_out, i_stop, i_err, 0)  // Run once
        select {
            case data_out := <- i_out:                 // Listen for data
                handleOutput(data_out.Values)
            case <-stop:                               // Listen for external stop command
                i_stop <- true
                return
            case temp_err := <- i_err:                 // Listen for internal error
                err <- temp_err
                i_stop <- true
                return
        }
        loop_i += 1                                    // Iterate index value
    }
}

// func (l Loop) Run(inputs ParamValues, outputs chan DataOut, stop chan bool, err chan *FlowError, id InstanceID) {
//     Nodes    := map[string]FunctionBlock{l.blk.GetName(): l.blk, l.cnd.GetName(): l.cnd}
//     ADDR     := Address{Name: l.name, ID: id}
//     I        := ParamAddress{Name: INDEX_NAME, Addr: ADDR, T: Int, is_input: false}
//     DONE     := ParamAddress{Name: DONE_NAME, Addr: ADDR, T: Bool, is_input: true}
//     data_out := make(map[ParamAddress]interface{})
//     data_in  := make(map[ParamAddress]interface{})
//     done, i  := false, 0
//     running  := true
//     blk_out, blk_stop, f_err := make(chan DataOut), make(chan bool), make(chan FlowError)
//     cnd_out, cnd_stop        := make(chan DataOut), make(chan bool)
    
//     stopAll := func() {
//         blk_stop <- true
//         cnd_stop <- true
//         running = false
//     }
    
//     passError := func(e FlowError) {
//         err <- e
//         if !e.Ok {
//             stopAll()
//         }
//     }
    
//     // Reads the inputs into the data_in map
//     loadvars := func() {
//         for name, param_lst := range l.in_feed {
//             for _, param := range param_lst {
//                 val, exists := inputs[name]
//                 switch {
//                     case !exists:
//                         passError(FlowError{false, DNE_ERROR, ADDR})
//                     case !CheckType(param.T, val):
//                         passError(FlowError{false, TYPE_ERROR, ADDR})
//                     default:
//                         data_in[param] = val
//                 }
//             }
//         }
//     }
    
//     // Puts the iteration number into a paramter I
//     updateI := func(i int) {
//         data_out[I] = i
//     }
    
//     // Puts the Done value from data_in in parameter done
//     updateDone := func() {
//         done = data_in[DONE].(bool)
//     }
    
//     // Used to listen for outputs and respond either by passing errors or storing them in data_out
//     handleOutput := func(blk_out, cnd_out chan DataOut, blk_stop, cnd_stop chan bool, f_err chan FlowError) {
//         storeValues := func(d DataOut) {
//             for name, val := range d.Values {
//                 blk := Nodes[d.Addr.Name]
//                 _, out_params := blk.GetParams() 
//                 t, t_exists := out_params[name]
//                 if t_exists {
//                     param := ParamAddress{Name: name, Addr: d.Addr, T: t, is_input: false}
//                     if CheckType(param.T, val) {
//                         data_out[param] = val
//                     } else {
//                         passError(FlowError{Ok: false, Info: TYPE_ERROR, Addr: ADDR})
//                     }
//                 }
//             }
//         }
        
//         blk_found, cnd_found := false, false
//         for !blk_found && !cnd_found && running {
//             select {
//                 case temp_out := <-blk_out:
//                     storeValues(temp_out)
//                     blk_found = true
//                 case temp_out := <-cnd_out:
//                     storeValues(temp_out)
//                     cnd_found = true
//                 case temp_err := <-f_err:
//                     passError(temp_err)
//                 case <-stop:
//                     passError(FlowError{Ok: false, Info: STOPPING, Addr: ADDR})
//             }
//         }
//     }
    
//     // Used to pass values from outputs to inputs based on registers
//     shiftLoopValues := func() {
//         for param, val := range(data_out) {
//             new_param, exists := l.registers[param]
//             if exists && CheckType(new_param.T, val) {
//                 data_in[new_param] = val
//                 delete(data_out, param)
//             }
//         }
//     }
    
//     getIns := func() (cnd_val ParamValues, blk_val ParamValues) {
//         cnd_in, _ := l.cnd.GetParams()
//         blk_in, _ := l.blk.GetParams()
        
//         get := func(params ParamTypes) ParamValues {
//             out_val := make(ParamValues)
//             for name, _ := range params {
//                 param := ParamAddress{Name: name, Addr: Address{Name: l.cnd.GetName(), ID: 0}, is_input: true}
//                 val, exists := data_in[param]
//                 if exists {
//                     out_val[name] = val
//                 }
//             }
//             return out_val
//         }
        
//         return get(cnd_in), get(blk_in)
//     }
    
//     getOut := func() DataOut {
//         out := make(ParamValues)
//         for name, param := range l.out_feed {
//             val, exists := data_out[param]
//             _, is_output := l.outputs[name]
//             chk := CheckType(param.T, val)
//             switch {
//                 case !exists || !is_output || !chk:
//                     return DataOut{Addr: ADDR, Values: make(ParamValues)}
//                 default:
//                     out[name] = data_out[param]
//             }
//         }
//         return DataOut{Addr: ADDR, Values: out}
//     }
    
//     // Loop until done or error
//     loadvars()
//     blk_in, cnd_in := getIns()
//     for !done && running {
//         go l.blk.Run(blk_in, blk_out, blk_stop, f_err, 0)
//         go l.cnd.Run(cnd_in, cnd_out, cnd_stop, f_err, 0)
//         handleOutput(blk_out, cnd_out, blk_stop, cnd_stop, f_err)
//         shiftLoopValues()
        
//         i += 1
//         updateI(i)
//         updateDone()
//     }
    
//     // Return values
//     out := getOut()
//     outputs <- out
//     return
// }
