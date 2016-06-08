package graphs

import (
    "testing"
    "../blocks"
)

// Logic

func TestSum(t *testing.T) {
    name := "array_sum"
    //fmt.Println("Testing ", name, "...")
    blk, _ := Sum(0)
    x := []float64{1,2,3}
    c := 6.0
    err := blocks.TestUnary(blk, x, c, "X", "OUT", name)
    if err != nil {
        t.Error(err.Info)
    }
}