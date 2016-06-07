package flow

const (
    INDEX_NAME = "I"
    DONE_NAME = "DONE"
)

type Loop struct {
    g *Graph
    blk, cnd  Address
    registers map[ParamAddress]ParamAddress
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
    inputs[INDEX_NAME] = Int
    outputs[DONE_NAME] = Bool
    
    // Initialize variables
    regs := make(map[ParamAddress]ParamAddress)
    blk_addr, cnd_addr := Address{blk.GetName(), 0}, Address{stop_condition.GetName(), 0}
    
    // Build Graph
    graph, err1 := NewGraph(name, inputs, outputs)
    outLoop := Loop{graph, blk_addr, cnd_addr, regs}
    err2 := outLoop.AddNode(blk, blk_addr)
    err3 := outLoop.AddNode(stop_condition, cnd_addr)
    
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

// Inhereted Fields
func (l Loop) AddNode(blk FunctionBlock, addr Address) (ok *Error) {return l.g.AddNode(blk, addr)}
func (l Loop) AddEdge(out_addr Address, out_param_name string,
                       in_addr Address, in_param_name string) (ok bool) {return l.g.AddEdge(out_addr, out_param_name, in_addr, in_param_name)}
func (l Loop) GetName() string {return l.g.GetName()}
func (l Loop) LinkIn(self_param_name string, in_param_name string, in_addr Address) (ok bool) {
    return l.g.LinkIn(self_param_name, in_param_name, in_addr)
}
func (l Loop) LinkOut(out_addr Address, out_param_name string, self_param_name string) (ok bool) {
    return l.g.LinkOut(out_addr, out_param_name, self_param_name)
}

// Gets parameters, but ignores Index and Done
func (l Loop) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    for name, t := range l.g.inputs {
        if name != INDEX_NAME {
            inputs[name] = t
        }
    }
    for name, t := range l.g.outputs {
        if name != DONE_NAME {
            outputs[name] = t
        }
    }
    return inputs, outputs
}

// Novel Methods
func (l Loop) AddRegister(out_name, in_name string, addr Address, init interface{}) *Error {
    return nil
}

func (l Loop) Run(inputs ParamValues, outputs chan DataOut, stop chan bool, err chan *FlowError, id InstanceID) {
    return
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
