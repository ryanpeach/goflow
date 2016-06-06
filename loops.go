package flow

const (
    New = "New"
    Done = "Done"
)

func NewLoop(name string, blk FunctionBlock, stop_condition FunctionBlock) Graph {
    NewGraph(, inputs, outputs ParamTypes)
}

type Loop struct {
    name string
    blk  FunctionBlock
    cnd  FunctionBlock  // The stop condition
    inputs ParamTypes
    outputs ParamTypes
    registers map[ParamAddress]ParamAddress
    in_feed map[string][]ParamAddress
    out_feed map[string]ParamAddress
}
func (l Loop) GetParams() (inputs ParamTypes, outputs ParamTypes) {return l.inputs, l.outputs}
func (l Loop) GetName() string {return l.blk.GetName() + "Loop"}

const (
    I_NAME = "I"
    DONE_NAME = "DONE"
)

func (l Loop) Run(inputs ParamValues, outputs chan DataOut, stop chan bool, err chan FlowError, id InstanceID) {
    Nodes    := map[string]FunctionBlock{l.blk.GetName(): l.blk, l.cnd.GetName(): l.cnd}
    ADDR     := Address{Name: l.name, ID: id}
    I        := ParamAddress{Name: I_NAME, Addr: ADDR, T: Int, is_input = false}
    Done     := ParamAddress{Name: DONE_NAME, Addr: ADDR, T: Bool, is_input = true}
    data_out := map[ParamAddress]interface{}
    data_in  := map[ParamAddress]interface{}
    done, i  := false, 0
    running := true
    
    // Reads the inputs into the data_in map
    loadvars := func() {
        for name, param_lst := range in_feed {
            for param := range param_lst {
                val, exists := inputs[name]
                switch {
                    case !exists:
                        pushError(DNE_ERROR)
                    case !CheckType(param, val):
                        pushError(TYPE_ERROR)
                    default:
                        data_in[param] = val
                }
            }
        }
    }
    
    // Puts the iteration number into a paramter I
    updateI := func(i int) {
        data_out[I] = i
    }
    
    // Puts the Done value from data_in in parameter done
    updateDone := func() {
        done = data_in[DONE]
    }
    
    passError := func(e FlowError) {
        err <- e
        if !e.Ok {
            stopAll()
        }
    }
    
    stopAll := func() {
        blk_stop <- true
        cnd_stop <- true
        running = false
    }
    
    // Used to listen for outputs and respond either by passing errors or storing them in data_out
    handleOutput := func(blk_out, cnd_out chan DataOut, blk_stop, cnd_stop chan bool, f_err chan FlowError) {
        storeValues := func(d DataOut) {
            for name, val := range d.Values {
                blk := Nodes[d.Addr.Name]
                _, out_params := blk.GetParams() 
                t, t_exists := out_params[name]
                if t_exists {
                    param = ParamAddress{Name: name, Addr: d.Addr, T: t, is_input: false}
                    if CheckType(param, val) {
                        data_out[param] = val
                    } else {
                        passError(FlowError{Ok: false, Info: TYPE_ERROR, Addr: ADDR})
                    }
                }
            }
        }
        
        blk_found, cnd_found := false, false
        for !blk_found && !cnd_found && running {
            select {
                case temp_out := <-blk_out:
                    storeValues(temp_out)
                    blk_found = true
                case temp_out := <-cnd_out:
                    storeValues(temp_out)
                    cnd_found = true
                case temp_err := <-f_err:
                    passError(temp_err)
                case <-stop:
                    passError(FlowError{Ok: false, Info: STOPPING, Addr: addr}
            }
        }
    }
    
    // Used to pass values from outputs to inputs based on registers
    shiftLoopValues := func() {
        for param, val := range(data_out) {
            new_param, exists := l.registers[param]
            if exists && CheckType(new_param, val) {
                data_in[new_param] = val
                delete(data_out, param)
            }
        }
    }
    
    getIns := func() (cnd_val ParamValues, blk_val ParamValues) {
        cnd_in, _ := l.cnd.GetParams()
        blk_in, _ := l.blk.GetParams()
        
        get := func(param_lst) {
            out_val := make(ParamValues)
            for name, t := range cnd_in {
                param := ParamAddress{Name: name, Addr: Address{Name: cnd.GetName(), ID: 0}, is_input: true}
                val, exists := data_in[param]
                if exists {
                    out_val[name] = val
                }
            }
        }
        
        return get(cnd_in), get(blk_in)
    }
    getOut := func() DataOut {
        out := make(ParamValues)
        for name, param := range l.out_feed {
            val, exists := data_out[param]
            _, is_output := l.outputs[name]
            chk := CheckType(val, param)
            switch {
                case !exists || !is_output || !chk:
                    return DataOut{Addr: ADDR, Values: make(ParamValues)}
                default:
                    out[name] = data_out[param]
            }
        }
        return DataOut{Addr: ADDR, Values: out}
    }
    
    // Initialize variables
    blk_out, blk_stop, f_err := make(chan DataOut), make(chan bool), make(chan FlowError)
    cnd_out, cnd_stop        := make(chan DataOut), make(chan bool)
    
    // Loop until done or error
    loadvars()
    blk_in, cnd_in := getIns()
    for !done and running {
        go l.blk.Run(blk_in, blk_out, blk_stop, f_err, 0)
        go l.cnd.Run(cnd_in, cnd_out, cnd_stop, f_err, 0)
        handleOutput(blk_out, cnd_out, blk_stop, cnd_stop, f_err)
        shiftLoopValues()
        
        i += 1
        updateI(i)
        updateDone()
    }
    
    // Return values
    out := getOut()
    outputs <- out
    return
}

func NewLoop(blk FunctionBlock, stop_condition FunctionBlock) (out FunctionBlock, err FlowError) {
    b_in, b_out := blk.GetParams()
    
    // Check that b_out has a Done output
    done_type, done_exists := b_out[Done]
    if (!done_exists || done_type != Bool) {
        err = FlowError{Ok: false, Info: "Must have a Done output of type Bool."}
        out = nil
        return
    }
    
    // Check that outputs has a NewX for every X in inputs
    out_loops := make(map[string]bool)
    for name, t_in := range b_in {
        t_out, exists := b_out[New + name]
        if (!exists || t_out != t_in) {
            err = FlowError{Ok: false, Info: name + " must have an output named " + New + name + " of same type."}
            out = nil
            return
        }
        out_loops[name] = true
    }
    
    // Create Loop inputs and outputs without looping outputs
    inputs, outputs := copyTypes(b_in), make(ParamTypes)
    for name, t := range b_out {
        if !out_loops[name] && name != Done {
            outputs[name] = t
        }
    }
    
    // Create loop block and return
    out = Loop{blk: blk, inputs: inputs, outputs: outputs}
    err = FlowError{Ok: true}
    return
}

func copyTypes(p ParamTypes) ParamTypes {
    out := make(ParamTypes)
    for name, t := range p {out[name] = t}
    return out
}
func copyValues(p ParamValues) ParamValues {
    out := make(ParamValues)
    for name, t := range p {out[name] = t}
    return out
}