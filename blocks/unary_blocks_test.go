package flow

import (
    "fmt"
    "testing"
)

func testUnary(blk FunctionBlock, a, c interface{}, name string) FlowError {
    // Run a Plus block
    f_out := make(chan DataOut)
    f_stop := make(chan bool)
    f_err := make(chan FlowError)

    // Run block and put a timeout on the stop channel
    go blk.Run(ParamValues{"IN": a}, f_out, f_stop, f_err)
    go Timeout(f_stop, 100000)
    
    // Wait for output or error
    var out DataOut
    var cont bool = true
    for cont {
        select {
            case out = <-f_out:
                cont = false
            case err := <-f_err:
                if !err.Ok {
                    return err
                }
            case <-f_stop:
                return FlowError{Ok: false, Info: "Timeout", Addr: blk.GetAddr()}
        }
    }
    
    // Test the output
    if out.Values["OUT"] != c {
        return FlowError{Ok: false, Info: "Returned wrong value.", Addr: blk.GetAddr()}
    } else {
        return FlowError{Ok: true, Addr: blk.GetAddr()}
    }
}

// Testing Type Conversions
func TestFloattoInt(t *testing.T) {
    name := "FloattoInt"
    fmt.Println("Testing ", name, "...")
    blk := FloattoInt(0)
    a := 5.1
    c := int(a)
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInttoFloat(t *testing.T) {
    name := "InttoFloat"
    fmt.Println("Testing ", name, "...")
    blk := InttoFloat(0)
    a := 5
    c := float64(a)
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

// Mathematical
func TestInc(t *testing.T) {
    name := "increment"
    fmt.Println("Testing ", name, "...")
    blk := Inc(0)
    a := 5
    c := a + 1
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestDec(t *testing.T) {
    name := "decrement"
    fmt.Println("Testing ", name, "...")
    blk := Dec(0)
    a := 5
    c := a - 1
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvInt(t *testing.T) {
    name := "invert_int"
    fmt.Println("Testing ", name, "...")
    blk := InvInt(0)
    a := 5
    c := -a
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvFloat(t *testing.T) {
    name := "invert_float"
    fmt.Println("Testing ", name, "...")
    blk := InvFloat(0)
    a := 5.1
    c := -a
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvBool(t *testing.T) {
    name := "invert_bool"
    fmt.Println("Testing ", name, "...")
    blk := InvBool(0)
    a := true
    c := !a
    err := testUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}