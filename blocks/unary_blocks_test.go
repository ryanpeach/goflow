package blocks

import (
    "fmt"
    "testing"
)

// Testing Type Conversions
func TestFloattoInt(t *testing.T) {
    name := "FloattoInt"
    fmt.Println("Testing ", name, "...")
    blk, _ := FloattoInt(0)
    a := 5.1
    c := int(a)
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInttoFloat(t *testing.T) {
    name := "InttoFloat"
    fmt.Println("Testing ", name, "...")
    blk, _ := InttoFloat(0)
    a := 5
    c := float64(a)
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

// Mathematical
func TestInc(t *testing.T) {
    name := "increment"
    fmt.Println("Testing ", name, "...")
    blk, _ := Inc(0)
    a := 5
    c := a + 1
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestDec(t *testing.T) {
    name := "decrement"
    fmt.Println("Testing ", name, "...")
    blk, _ := Dec(0)
    a := 5
    c := a - 1
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvInt(t *testing.T) {
    name := "invert_int"
    fmt.Println("Testing ", name, "...")
    blk, _ := InvInt(0)
    a := 5
    c := -a
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvFloat(t *testing.T) {
    name := "invert_float"
    fmt.Println("Testing ", name, "...")
    blk, _ := InvFloat(0)
    a := 5.1
    c := -a
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestInvBool(t *testing.T) {
    name := "invert_bool"
    fmt.Println("Testing ", name, "...")
    blk, _ := InvBool(0)
    a := true
    c := !a
    err := TestUnary(blk, a, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}