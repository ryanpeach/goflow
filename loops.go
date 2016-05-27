package flow

const (
    New = "New"
    Done = "Done"
)

type Loop struct {
    blk FunctionBlock
    inputs ParamTypes
    outputs ParamTypes
}
func (l Loop) GetParams() (inputs ParamTypes, outputs ParamTypes) {return l.inputs, l.outputs}
func (l Loop) GetName() string {return l.blk.GetName() + "Loop"}

func (l Loop) Run(inputs ParamValues, outputs chan DataOut, stop chan bool, err chan FlowError, id InstanceID) {
    // Used to pass values from outputs to inputs
    passLoopValues := func(d DataOut) (rIn ParamValues, out ParamValues, done bool) {
        done = d.Values[Done].(bool)
        out = make(ParamValues)
        rIn = make(ParamValues)
        for name, _ := range l.inputs {
            val, exists := d.Values[New+name]
            if exists {
                rIn[name] = val
            } else {
                out[name] = val
            }
        }
        return
    }
    
    // Initialize variables
    var out_vals ParamValues
    addr     := NewAddress(id, l.GetName())
    thisIn   := copyValues(inputs)
    done     := false
    f_out, f_stop, f_err := make(chan DataOut), make(chan bool), make(chan FlowError)
    
    // Loop until done
    for !done {
        go l.blk.Run(thisIn, f_out, f_stop, f_err, id)
        select {
            case temp_out := <-f_out:
                thisIn, out_vals, done = passLoopValues(temp_out)
            case temp_err := <-f_err:
                f_stop <- true
                err <- temp_err
                return
            case <-stop:
                f_stop <- true
                err <- FlowError{Ok: false, Info: StopInfo, Addr: addr}
                return
        }
    }
    
    // Return values
    out := DataOut{Addr: addr, Values: out_vals}
    outputs <- out
    return
}

func NewLoop(blk FunctionBlock) (out FunctionBlock, err FlowError) {
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