package graphs

import (
    ".."
    "../blocks"
)

const (
    CND = true
    BLK = false
)

func Sum(id flow.InstanceID) (flow.Graph, flow.Address) {
    // Create Summation Block
    ins  := flow.ParamTypes{"X": flow.NumArray, "Index": flow.Int, "Total": flow.Num}
    outs := flow.ParamTypes{"OUT": flow.Num}
    summation_graph := flow.NewGraph("array_sum", ins, outs)
    summation_addr  := flow.Address{"array_sum", 0}
    sum, sum_addr     := blocks.Sum(0)
    index, index_addr := blocks.Index(0)
    summation_graph.AddNode(sum, sum_addr)
    summation_graph.AddNode(index, index_addr)
    summation_graph.LinkIn("X", "IN", index_addr)
    summation_graph.AddEdge("OUT", index_addr, "A", sum_addr)
    summation_graph.LinkIn("Total", "B", sum_addr)
    summation_graph.LinkOut(sum_addr, "OUT", "OUT")
    
    // Create Conditional Block
    ins  := flow.ParamTypes{"X": flow.NumArray, "Index": flow.Int}
    outs := flow.ParamTypes{"OUT": flow.Bool}
    cnd_graph := flow.NewGraph("index_final", ins, outs)
    cnd_addr  := flow.Address{"index_final", 0}
    lt, lt_addr   := blocks.Lesser(0)
    not, not_addr := blocks.InvBool(0)
    ln, ln_addr   := blocks.Len(0)
    cnd_graph.AddNode(lt, lt_addr)
    cnd_graph.AddNode(not, not_addr)
    cnd_graph.AddNode(ln, ln_addr)
    cnd_graph.LinkIn("Index", "A", lt_addr)
    cnd_graph.LinkIn("X", "IN", ln_addr)
    cnd_graph.AddEdge("OUT", ln_addr, "B", lt_addr)
    cnd_graph.AddEdge("OUT", lt_addr, "IN", not_addr)
    cnd_graph.LinkOut(not_addr, "OUT", "OUT")
    
    // Create Loop
    ins  := flow.ParamTypes{"X": flow.NumArray}
    outs := flow.ParamTypes{"OUT": flow.Bool}
    loop := flow.NewLoop("summation_loop", ins, outs, summation_graph, cnd_graph)
    loop_addr := flow.Address{"summation_loop", id}
    loop.LinkIn("X", "X", CND)
    loop.LinkIn("X", "X", BLK)
    loop.LinkIn(INDEX_NAME, "Index", CND)
    loop.AddRegister("OUT", "Total", BLK, 0)
    
    return loop, addr
}