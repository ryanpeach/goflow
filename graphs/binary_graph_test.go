package graphs

import (
    "testing"
    "../blocks"
)

// Logic

func TestNand(t *testing.T) {
    name := "logical_nand"
    //fmt.Println("Testing ", name, "...")
    blk, _ := Nand(0)
    a, b := true, false
    c := !(a && b)
    err := blocks.TestBinary(blk, a, b, c, "A", "B", "OUT", name)
    if err != nil {
        t.Error(err.Info)
    }
}
func BenchmarkNand(b *testing.B) {
    name := "logical_nand"
    //fmt.Println("Testing ", name, "...")
    blk, _ := Nand(0)
    A, B := true, false
    C := !(A && B)
    for i := 0; i < b.N; i++ {
        blocks.TestBinary(blk, A, B, C, "A", "B", "OUT", name)
    }
}
func BenchmarkNand2(b *testing.B) {
    out := func() bool {return !(true && false)}
    for i := 0; i<b.N; i++ {
        go out()
    }
}
/*
func TestNor(t *testing.T) {
    name := "logical_nor"
    fmt.Println("Testing ", name, "...")
    blk := Nor(0)
    a, b := true, false
    c := !(a || b)
    err := testBinary(blk, a, b, c, "A", "B", "OUT", name)
    if !err.Ok {
        t.Error(err.Info)
    }
}


// Comparison
func TestGreaterEquals(t *testing.T) {
    name := "greater_equals"
    fmt.Println("Testing ", name, "...")
    blk := Greater(0)
    a, b := 5, 2
    c := 5 >= 2
    err := testBinary(blk, a, b, c, "A", "B", "OUT", name)
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
    err := testBinary(blk, a, b, c, "A", "B", "OUT", name)
    if !err.Ok {
        t.Error(err.Info)
    }
}

func TestNotEquals(t *testing.T) {
    name := "not_equal_to"
    fmt.Println("Testing ", name, "...")
    blk := NotEquals(0)
    a, b := 5, 2
    c := 5 != 2
    err := testBinary(blk, a, b, c, "A", "B", "OUT", name)
    if !err.Ok {
        t.Error(err.Info)
    }
}
*/