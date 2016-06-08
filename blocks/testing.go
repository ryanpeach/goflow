package blocks

import ".."

func TestUnary(blk flow.FunctionBlock, a, c interface{}, nA, nC, name string) *flow.FlowError {
    // Run a Plus block
    f_out := make(chan flow.DataOut)
    f_stop := make(chan bool)
    f_err := make(chan *flow.FlowError)

    // Run block and put a timeout on the stop channel
    go blk.Run(flow.ParamValues{nA: a}, f_out, f_stop, f_err, 0)
    //go flow.Timeout(f_stop, 100000)
    addr := flow.Address{blk.GetName(), 0}
    
    // Wait for output or error
    var out flow.DataOut
    var cont bool = true
    for cont {
        select {
            case out = <-f_out:
                cont = false
            case err := <-f_err:
                return err
        }
    }
    
    // Test the output
    if out.Values[nC] != c {
        return flow.NewFlowError(flow.VALUE_ERROR, "Returned wrong value.", addr)
    } else {
        return nil
    }
}

func TestBinary(blk flow.FunctionBlock, a, b, c interface{}, aN, bN, cN, name string) *flow.FlowError {
    
    // Run a Plus block
    f_out := make(chan flow.DataOut)
    f_stop := make(chan bool)
    f_err := make(chan *flow.FlowError)

    // Run block and put a timeout on the stop channel
    go blk.Run(flow.ParamValues{aN: a, bN: b}, f_out, f_stop, f_err, 0)
    //go flow.Timeout(f_stop, 100000)
    addr := flow.Address{blk.GetName(), 0}
    
    // Wait for output or error
    var out flow.DataOut
    var cont bool = true
    for cont {
        select {
            case out = <-f_out:
                cont = false
            case err := <-f_err:
                return err
        }
    }
    
    // Test the output
    if out.Values[cN] != c {
        return flow.NewFlowError(flow.VALUE_ERROR, "Returned wrong value.", addr)
    } else {
        return nil
    }
}