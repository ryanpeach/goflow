package blocks

import ".."

func TestUnary(blk flow.FunctionBlock, a, c interface{}, name string) flow.FlowError {
    // Run a Plus block
    f_out := make(chan flow.DataOut)
    f_stop := make(chan bool)
    f_err := make(chan flow.FlowError)

    // Run block and put a timeout on the stop channel
    go blk.Run(flow.ParamValues{"IN": a}, f_out, f_stop, f_err, 0)
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
                if !err.Ok {
                    return err
                }
        }
    }
    
    // Test the output
    if out.Values["OUT"] != c {
        return flow.FlowError{Ok: false, Info: "Returned wrong value.", Addr: addr}
    } else {
        return flow.FlowError{Ok: true, Addr: addr}
    }
}

func TestBinary(blk flow.FunctionBlock, a, b, c interface{}, aN, bN, cN, name string) flow.FlowError {
    
    // Run a Plus block
    f_out := make(chan flow.DataOut)
    f_stop := make(chan bool)
    f_err := make(chan flow.FlowError)

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
                if !err.Ok {
                    return err
                }
        }
    }
    
    // Test the output
    if out.Values[cN] != c {
        return flow.FlowError{Ok: false, Info: "Returned wrong value.", Addr: addr}
    } else {
        return flow.FlowError{Ok: true, Addr: addr}
    }
}