package blocks

import (
	".."
	"fmt"
	"testing"
)

// Testing Input Switch
func TestInputSwitch(t *testing.T) {
	name := "InputSwitch"
	fmt.Println("Testing ", name, "...")
	blk, _ := InputSwitch(0, flow.Int)
	a, b, cnd := 5.1, 2.2, true
	in := flow.ParamValues{"A": a, "B": b, "Condition": cnd}
	out, err := RunBlock(blk, in)
	switch {
	case err != nil:
		t.Error(err.Info)
	case out["OUT"] != a:
		t.Error("Not the right value")
	}
}

// Testing Output Switch
func TestOutputSwitch(t *testing.T) {
	name := "OutputSwitch"
	fmt.Println("Testing ", name, "...")
	blk, _ := OutputSwitch(0, flow.Int)
	a, cnd := 5.1, true
	in := flow.ParamValues{"IN": a, "Condition": cnd}
	out, err := RunBlock(blk, in)
	b_exists := out["B"]
	switch {
	case err != nil:
		t.Error(err.Info)
	case b_exists:
		t.Error("B should not return a value.")
	case out["A"] != a:
		t.Error("Not the right value.")
	}
}
