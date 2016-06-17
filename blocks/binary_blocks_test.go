package blocks

import (
	"fmt"
	"testing"
)

// Testing Float Numerics
func TestPlusFloat(t *testing.T) {
	name := "PlusFloat"
	fmt.Println("Testing ", name, "...")
	blk, _ := PlusFloat(0)
	a, b := 5.1, 2.2
	c := float64(a + b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestSubFloat(t *testing.T) {
	name := "SubFloat"
	fmt.Println("Testing ", name, "...")
	blk, _ := SubFloat(0)
	a, b := 5.1, 2.2
	c := float64(a - b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestMultFloat(t *testing.T) {
	name := "MultFloat"
	fmt.Println("Testing ", name, "...")
	blk, _ := MultFloat(0)
	a, b := 5.1, 2.2
	c := float64(a * b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestDivFloat(t *testing.T) {
	name := "DivFloat"
	fmt.Println("Testing ", name, "...")
	blk, _ := DivFloat(0)
	a, b := 5.1, 2.2
	c := float64(a / b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}

// Testing Integer Numerics
func TestPlusInt(t *testing.T) {
	name := "PlusInt"
	fmt.Println("Testing ", name, "...")
	blk, _ := PlusInt(0)
	a, b := 5, 2
	c := int(a + b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestSubInt(t *testing.T) {
	name := "SubInt"
	fmt.Println("Testing ", name, "...")
	blk, _ := SubInt(0)
	a, b := 5, 2
	c := int(a - b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestMultInt(t *testing.T) {
	name := "MultInt"
	fmt.Println("Testing ", name, "...")
	blk, _ := MultInt(0)
	a, b := 5, 2
	c := int(a * b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestDivInt(t *testing.T) {
	name := "DivInt"
	fmt.Println("Testing ", name, "...")
	blk, _ := DivInt(0)
	a, b := 5, 2
	c := int(a / b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestMod(t *testing.T) {
	name := "Mod"
	fmt.Println("Testing ", name, "...")
	blk, _ := Mod(0)
	a, b := 5, 2
	c := int(a % b)
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}

// Testing Logical Operators
func TestAnd(t *testing.T) {
	name := "logical_and"
	fmt.Println("Testing ", name, "...")
	blk, _ := And(0)
	a, b := true, false
	c := a && b
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestOr(t *testing.T) {
	name := "logical_or"
	fmt.Println("Testing ", name, "...")
	blk, _ := Or(0)
	a, b := true, false
	c := a || b
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestXor(t *testing.T) {
	name := "logical_xor"
	fmt.Println("Testing ", name, "...")
	blk, _ := Xor(0)
	a, b := true, false
	c := a != b
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}

// Comparison
func TestGreater(t *testing.T) {
	name := "greater_than"
	fmt.Println("Testing ", name, "...")
	blk, _ := Greater(0)
	a, b := 5, 2
	c := 5 > 2
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func TestLesser(t *testing.T) {
	name := "lesser_than"
	fmt.Println("Testing ", name, "...")
	blk, _ := Lesser(0)
	a, b := 5, 2
	c := 5 < 2
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}

func TestEquals(t *testing.T) {
	name := "equal_to"
	fmt.Println("Testing ", name, "...")
	blk, _ := Greater(0)
	a, b := 5, 2
	c := 5 > 2
	err := TestBinary(blk, a, b, c, "A", "B", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}

// Arrays
func TestIndex(t *testing.T) {
	name := "index"
	fmt.Println("Testing ", name, "...")
	blk, _ := Index(0)
	a := []float64{1, 2, 3, 4}
	b := 2
	c := 3.0
	err := TestBinary(blk, a, b, c, "X", "Index", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
