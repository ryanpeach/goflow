package flow

import (
    "fmt"
    "testing"
)

func testBinary(blk FunctionBlock, a, b, c interface{}, name string) FlowError {
    // Run a Plus block
    f_out := make(chan DataOut)
    f_stop := make(chan bool)
    f_err := make(chan FlowError)

    // Run block and put a timeout on the stop channel
    go blk.Run(ParamValues{"A": a, "B": b}, f_out, f_stop, f_err)
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

// Testing Float Numerics
func TestPlusFloat(t *testing.T) {
    name := "PlusFloat"
    fmt.Println("Testing ", name, "...")
    blk := PlusFloat(0)
    a, b := 5.1, 2.2
    c := float64(a + b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestSubFloat(t *testing.T) {
    name := "SubFloat"
    fmt.Println("Testing ", name, "...")
    blk := SubFloat(0)
    a, b := 5.1, 2.2
    c := float64(a - b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestMultFloat(t *testing.T) {
    name := "MultFloat"
    fmt.Println("Testing ", name, "...")
    blk := MultFloat(0)
    a, b := 5.1, 2.2
    c := float64(a * b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestDivFloat(t *testing.T) {
    name := "DivFloat"
    fmt.Println("Testing ", name, "...")
    blk := DivFloat(0)
    a, b := 5.1, 2.2
    c := float64(a / b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}


// Testing Integer Numerics
func TestPlusInt(t *testing.T) {
    name := "PlusInt"
    fmt.Println("Testing ", name, "...")
    blk := PlusInt(0)
    a, b := 5, 2
    c := int(a + b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestSubInt(t *testing.T) {
    name := "SubInt"
    fmt.Println("Testing ", name, "...")
    blk := SubInt(0)
    a, b := 5, 2
    c := int(a - b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestMultInt(t *testing.T) {
    name := "MultInt"
    fmt.Println("Testing ", name, "...")
    blk := MultInt(0)
    a, b := 5, 2
    c := int(a * b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestDivInt(t *testing.T) {
    name := "DivInt"
    fmt.Println("Testing ", name, "...")
    blk := DivInt(0)
    a, b := 5, 2
    c := int(a / b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestMod(t *testing.T) {
    name := "Mod"
    fmt.Println("Testing ", name, "...")
    blk := Mod(0)
    a, b := 5, 2
    c := int(a % b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

// Testing Logical Operators
func TestAnd(t *testing.T) {
    name := "logical_and"
    fmt.Println("Testing ", name, "...")
    blk := And(0)
    a, b := true, false
    c := a && b
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestOr(t *testing.T) {
    name := "logical_or"
    fmt.Println("Testing ", name, "...")
    blk := Or(0)
    a, b := true, false
    c := a || b
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestXor(t *testing.T) {
    name := "logical_xor"
    fmt.Println("Testing ", name, "...")
    blk := Xor(0)
    a, b := true, false
    c := a != b
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

// Comparison
func TestGreater(t *testing.T) {
    name := "greater_than"
    fmt.Println("Testing ", name, "...")
    blk := Greater(0)
    a, b := 5, 2
    c := 5 > 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestLesser(t *testing.T) {
    name := "lesser_than"
    fmt.Println("Testing ", name, "...")
    blk := Lesser(0)
    a, b := 5, 2
    c := 5 < 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

func TestEquals(t *testing.T) {
    name := "equal_to"
    fmt.Println("Testing ", name, "...")
    blk := Greater(0)
    a, b := 5, 2
    c := 5 > 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}