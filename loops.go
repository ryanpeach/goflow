package flow

const (
    INDEX_NAME = "I"
    DONE_NAME = "DONE"
)

type Loop struct {
    name string
    g *Graph
    
    infeed      ParamLstMap    // Connects name of a param in inputs to a ParamAddress of some parameter in some node
    outfeed     ParamMap    // Connects name of a param in outputs to a ParamAddress of some parameter in some node
    inputs    ParamTypes
    outputs   ParamTypes
    
    sources   map[ParamAddress]ParamAddress
    registers NameMap
    initial   ParamValues
}

func NewLoop(name string, inputs, outputs ParamTypes, blk *Graph) (*Loop, *Error) {
    
    // Check that stop_condition has one bool output.
    _, cnd_out := blk.GetParams()
    bool_found := false
    for _, t := range cnd_out {
        if t == Bool {
            bool_found = true
            break
        }
    }
    if !bool_found {
        return nil, &Error{TYPE_ERROR, "Block has no boolean output."}
    }
    
    // Initialize variables
    regs, inits     := make(NameMap), make(ParamValues)
    infeed, outfeed := make(ParamLstMap), make(ParamMap)
    sources         := make(map[ParamAddress]ParamAddress)
    
    // Build Loop
    outLoop := Loop{name, blk, infeed, outfeed, inputs, outputs, sources, regs, inits}
    
    return &outLoop, nil
    
}

// FunctionBlock Fields
func (l Loop) GetName() string {return l.name}
func (l Loop) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    return CopyTypes(l.inputs), CopyTypes(l.outputs)
}

// Loop Fields
func (l Loop) LinkIn(self_param_name string, in_param_name string, in_addr Address) *Error {
    g_ins, _               := l.g.GetParams()
    g_type, g_exists       := g_ins[in_param_name]
    self_type, self_exists := l.inputs[self_param_name]
    in_param               := ParamAddress{l.g.GetName(), Address{l.g.GetName(), 0}, g_type, true}
    _, link_exists         := l.sources[in_param]
    _, other_links         := l.infeed[self_param_name]
    switch {
        case !g_exists:
            return &Error{DNE_ERROR, "out_param_name of inner graph does not exist."}
        case !self_exists && self_param_name != INDEX_NAME && self_param_name != DONE_NAME:
            return &Error{DNE_ERROR, "self_param_name does not exist."}
        case link_exists:
            return &Error{ALREADY_EXISTS_ERROR, "self_param_name is already connected to a parameter."}
        case !CheckSame(g_type, self_type):
            return &Error{TYPE_ERROR, "Types are not compatible."}
        case other_links:
            l.infeed[self_param_name] = append(l.infeed[self_param_name], in_param)
        default:
            l.infeed[self_param_name] = []ParamAddress{in_param}
    }
    l.sources[in_param] = ParamAddress{self_param_name, Address{l.name, 0}, self_type, true}
    return nil
}
func (l Loop) LinkOut(out_addr Address, out_param_name string, self_param_name string) *Error {
    _, g_outs              := l.g.GetParams()
    g_type, g_exists       := g_outs[out_param_name]
    self_type, self_exists := l.outputs[self_param_name]
    _, link_exists         := l.outfeed[self_param_name]
    out_param               := ParamAddress{l.g.GetName(), Address{l.g.GetName(), 0}, g_type, true}
    switch {
        case !g_exists:
            return &Error{DNE_ERROR, "out_param_name of inner graph does not exist."}
        case !self_exists && self_param_name != INDEX_NAME && self_param_name != DONE_NAME:
            return &Error{DNE_ERROR, "self_param_name does not exist."}
        case link_exists:
            return &Error{ALREADY_EXISTS_ERROR, "self_param_name is already connected to a parameter."}
        case !CheckSame(g_type, self_type):
            return &Error{TYPE_ERROR, "Types are not compatible."}
        default:
            l.outfeed[self_param_name] = out_param
    }
    return nil
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
    ins, outs := l.g.GetParams()
    t1, exists1 := ins[in_name]
    t2, exists2 := outs[out_name]
    switch {
        case !exists1:
            return &Error{DNE_ERROR, "in_name is not a parameter of graph."}
        case !exists2:
            return &Error{DNE_ERROR, "out_name is not a parameter of graph."}
        case !CheckType(t1, init) || !CheckType(t2, init):
            return &Error{TYPE_ERROR, "t and init are incompatible types."}
        default:
            _, connected := l.registers[in_name]
            if !connected {
                l.registers[in_name] = out_name
                l.initial[in_name] = init
                in_param := ParamAddress{in_name, Address{l.g.GetName(), 0}, t1, true}
                l.sources[in_param] = ParamAddress{in_name, Address{l.name, 0}, t1, false}
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
    ins, outs := l.g.GetParams()
    t1, exists1 := ins[in_name]
    t2, exists2 := outs[out_name]
    _, feed_exists := l.infeed[in_name]

    // Handle Errors
    switch {
        case !exists1:
            return &Error{DNE_ERROR, "in_name is not a parameter of graph."}
        case !exists2:
            return &Error{DNE_ERROR, "out_name is not a parameter of graph."}
        case !feed_exists:
            return &Error{DNE_ERROR, "in_name must have a feed connected prior to creating a register"}
        case !CheckType(t1, t2):
            return &Error{TYPE_ERROR, "in_name and out_name are incompatible types."}
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
    i_inputs := make(ParamValues)
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
    
    // Handle inputs
    for name, param_lst := range l.infeed {
        val, exists := inputs[name]
        switch {
            case name == INDEX_NAME:
            case !exists:
                err <- NewFlowError(DNE_ERROR, "Not all inputs satisfied.", ADDR)
                return
            default:
                for _, param := range param_lst {
                    i_inputs[param.Name] = val
                }
        }
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
