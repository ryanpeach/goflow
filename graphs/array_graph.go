package graphs

import (
    ".."
    "../blocks"
)

func Sum(id flow.InstanceID) (*flow.Loop, flow.Address) {
    // Create Summation Block
    ins  := flow.ParamTypes{"X": flow.NumArray, "Index": flow.Int, "Total": flow.Float}
    outs := flow.ParamTypes{"OUT": flow.Num}
    summation_graph, _ := flow.NewGraph("array_sum", ins, outs)
    summation_addr  := flow.Address{"array_sum", 0}
    sum, sum_addr     := blocks.PlusFloat(0)
    index, index_addr := blocks.Index(0)
    summation_graph.AddNode(sum, sum_addr)
    summation_graph.AddNode(index, index_addr)
    summation_graph.LinkIn("X", "X", index_addr)
    summation_graph.LinkIn("Index", "Index", index_addr)
    summation_graph.AddEdge(index_addr, "OUT", sum_addr, "A")
    summation_graph.LinkIn("Total", "B", sum_addr)
    summation_graph.LinkOut(sum_addr, "OUT", "OUT")
    
    // Create Conditional Block
    ins  = flow.ParamTypes{"X": flow.NumArray, "Index": flow.Int}
    outs = flow.ParamTypes{"OUT": flow.Bool}
    cnd_graph, _ := flow.NewGraph("index_final", ins, outs)
    cnd_addr  := flow.Address{"index_final", 0}
    lt, lt_addr   := blocks.Lesser(0)
    not, not_addr := blocks.InvBool(0)
    ln, ln_addr   := blocks.Len(0)
    cnd_graph.AddNode(lt, lt_addr)
    cnd_graph.AddNode(not, not_addr)
    cnd_graph.AddNode(ln, ln_addr)
    cnd_graph.LinkIn("Index", "A", lt_addr)
    cnd_graph.LinkIn("X", "IN", ln_addr)
    cnd_graph.AddEdge(ln_addr, "OUT", lt_addr, "B")
    cnd_graph.AddEdge(lt_addr, "OUT", not_addr, "IN")
    cnd_graph.LinkOut(not_addr, "OUT", "OUT")
    
    // Create Loop
    ins  = flow.ParamTypes{"X": flow.NumArray}
    outs = flow.ParamTypes{"OUT": flow.Float}
    loop, _ := flow.NewLoop("summation_loop", ins, outs, summation_graph, cnd_graph)
    loop_addr := flow.Address{"summation_loop", id}
    loop.LinkIn("X", "X", cnd_addr)
    loop.LinkIn("X", "X", summation_addr)
    loop.LinkIn(flow.INDEX_NAME, "Index", cnd_addr)
    loop.LinkIn(flow.INDEX_NAME, "Index", blk_addr)
    loop.AddRegister("OUT", "Total", summation_addr, 0)
    loop.LinkOut(cnd_addr, "OUT", DONE_NAME)
    loop.LinkOut(blk_addr, "OUT", "OUT")
    
    return loop, loop_addr
}