package flow

import (
    "fmt"
    "testing"
)

// Logic
/*
func TestNand(t *testing.T) {
    name := "logical_nand"
    fmt.Println("Testing ", name, "...")
    blk := Nand(0)
    a, b := true, false
    c := !(a && b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestNor(t *testing.T) {
    name := "logical_nor"
    fmt.Println("Testing ", name, "...")
    blk := Nor(0)
    a, b := true, false
    c := !(a || b)
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
*/

// Comparison
func TestGreaterEquals(t *testing.T) {
    name := "greater_equals"
    fmt.Println("Testing ", name, "...")
    blk := Greater(0)
    a, b := 5, 2
    c := 5 >= 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
func TestLesserEquals(t *testing.T) {
    name := "lesser_equals"
    fmt.Println("Testing ", name, "...")
    blk := Lesser(0)
    a, b := 5, 2
    c := 5 <= 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
/*
func TestNotEquals(t *testing.T) {
    name := "not_equal_to"
    fmt.Println("Testing ", name, "...")
    blk := NotEquals(0)
    a, b := 5, 2
    c := 5 != 2
    err := testBinary(blk, a, b, c, name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
*/