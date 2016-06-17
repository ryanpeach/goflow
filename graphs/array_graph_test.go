package graphs

import (
	"../blocks"
	"testing"
)

// Logic

func TestSum(t *testing.T) {
	name := "array_sum"
	//fmt.Println("Testing ", name, "...")
	blk, _ := Sum(0)
	x := []float64{1, 2, 3}
	c := 6.0
	err := blocks.TestUnary(blk, x, c, "X", "OUT", name)
	if err != nil {
		t.Error(err.Info)
	}
}
func BenchmarkSum(b *testing.B) {
	name := "array_sum"
	//fmt.Println("Testing ", name, "...")
	blk, _ := Sum(0)
	x := []float64{1, 2, 3}
	c := 6.0
	for i := 0; i < b.N; i++ {
		err := blocks.TestUnary(blk, x, c, "X", "OUT", name)
		if err != nil {
			b.Error(err.Info)
		}
	}
}
func BenchmarkSum2(b *testing.B) {
	x := []float64{1, 2, 3}
	c := 6.0
	sum := func(out chan float64) {
		val := 0.0
		for _, v := range x {
			val += v
		}
		out <- val
	}
	data := make(chan float64)
	for i := 0; i < b.N; i++ {
		go sum(data)
		if <-data != c {
			b.Error("C != ", c)
		}
	}
}
