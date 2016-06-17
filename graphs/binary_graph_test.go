package graphs

import (
	"../blocks"
	"testing"
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
		err := blocks.TestBinary(blk, A, B, C, "A", "B", "OUT", name)
		if err != nil {
			b.Error(err.Info)
		}
	}
}
func BenchmarkNand2(b *testing.B) {
	x, y := true, false
	z := true
	nand := func(out chan bool) {
		out <- !(x && y)
	}
	data := make(chan bool)
	for i := 0; i < b.N; i++ {
		go nand(data)
		if <-data != z {
			b.Error("Z != ", z)
		}
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
