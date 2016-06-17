package graphs

import (
	".."
	"../blocks"
	"fmt"
)

func Sum(id flow.InstanceID) (*flow.Loop, flow.Address) {
	// Create Summation Block
	ins := flow.ParamTypes{"X": flow.NumArray, "Index": flow.Int, "Total": flow.Float}
	outs := flow.ParamTypes{"OUT": flow.Float, "Done": flow.Bool}
	g, _ := flow.NewGraph("array_sum", ins, outs)

	// Create Blocks
	sum, sum_addr := blocks.PlusFloat(0)
	index, index_addr := blocks.Index(0)
	eq, eq_addr := blocks.Equals(0)
	sub, sub_addr := blocks.SubFloat(0)
	ln, ln_addr := blocks.Len(0)
	toflt1, toflt_addr1 := blocks.InttoFloat(0)
	toflt2, toflt_addr2 := blocks.InttoFloat(1)

	// Add Nodes
	err1 := g.AddNode(sum, sum_addr)
	err2 := g.AddNode(index, index_addr)
	err3 := g.AddNode(eq, eq_addr)
	err4 := g.AddNode(sub, sub_addr)
	err5 := g.AddNode(ln, ln_addr)
	err6 := g.AddNode(toflt1, toflt_addr1)
	err7 := g.AddNode(toflt2, toflt_addr2)

	// Create input links
	err8 := g.LinkIn("X", "X", index_addr)         // Connect array to array in index retrieval block
	err9 := g.LinkIn("Index", "Index", index_addr) // Connect index to index of index retrieval block
	err10 := g.LinkIn("Index", "IN", toflt_addr2)
	err11 := g.LinkIn("X", "IN", ln_addr)     // Use the array as the array input of the length block
	err12 := g.LinkIn("Total", "B", sum_addr) // Use the total so far as the B of the addition block

	// Create edges
	err13 := g.AddEdge(ln_addr, "OUT", toflt_addr1, "IN")
	err14 := g.AddEdge(toflt_addr1, "OUT", sub_addr, "A") // Use the output of the length address as the input of the less than block
	err15 := g.AddEdge(toflt_addr2, "OUT", eq_addr, "B")  // Use the output of the length address as the input of the less than block
	err16 := g.AddEdge(index_addr, "OUT", sum_addr, "A")  // Use the output of the index retrieval as the A of the addition block
	err17 := g.AddEdge(sub_addr, "OUT", eq_addr, "A")

	// Create constants
	err18 := g.AddConstant(1.0, sub_addr, "B")

	// Create output links
	err19 := g.LinkOut(sum_addr, "OUT", "OUT") // The sum output is our total
	err20 := g.LinkOut(eq_addr, "OUT", "Done") // The not output is our done

	// Create Loop
	ins = flow.ParamTypes{"X": flow.NumArray}
	outs = flow.ParamTypes{"OUT": flow.Float}
	loop, _ := flow.NewLoop("summation_loop", ins, outs, g)
	loop_addr := flow.Address{"summation_loop", id}

	// Create link inputs
	err21 := loop.LinkIn("X", "X")
	err22 := loop.LinkIn(flow.INDEX_NAME, "Index")

	// Create register
	err23 := loop.AddDefaultRegister("OUT", "Total", 0.0)

	// Create output links
	err24 := loop.LinkOut("OUT", "OUT")
	err25 := loop.LinkOut("Done", flow.DONE_NAME)

	switch {
	case err1 != nil:
		fmt.Println("1: " + err1.Info)
	case err2 != nil:
		fmt.Println("2: " + err2.Info)
	case err3 != nil:
		fmt.Println("3: " + err3.Info)
	case err4 != nil:
		fmt.Println("4: " + err4.Info)
	case err5 != nil:
		fmt.Println("5: " + err5.Info)
	case err6 != nil:
		fmt.Println("6: " + err6.Info)
	case err7 != nil:
		fmt.Println("7: " + err7.Info)
	case err8 != nil:
		fmt.Println("8: " + err8.Info)
	case err9 != nil:
		fmt.Println("9: " + err9.Info)
	case err10 != nil:
		fmt.Println("10: " + err10.Info)
	case err11 != nil:
		fmt.Println("11: " + err11.Info)
	case err12 != nil:
		fmt.Println("12: " + err12.Info)
	case err13 != nil:
		fmt.Println("13: " + err13.Info)
	case err14 != nil:
		fmt.Println("14: " + err14.Info)
	case err15 != nil:
		fmt.Println("15: " + err15.Info)
	case err16 != nil:
		fmt.Println("16: " + err16.Info)
	case err17 != nil:
		fmt.Println("17: " + err17.Info)
	case err18 != nil:
		fmt.Println("18: " + err18.Info)
	case err19 != nil:
		fmt.Println("19: " + err19.Info)
	case err20 != nil:
		fmt.Println("20: " + err20.Info)
	case err21 != nil:
		fmt.Println("21: " + err21.Info)
	case err22 != nil:
		fmt.Println("22: " + err22.Info)
	case err23 != nil:
		fmt.Println("23: " + err23.Info)
	case err24 != nil:
		fmt.Println("24: " + err24.Info)
	case err25 != nil:
		fmt.Println("25: " + err25.Info)
	}
	return loop, loop_addr
}
