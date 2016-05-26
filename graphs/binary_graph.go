package flow

func Nand(id InstanceID) FunctionBlock {
    // Create Nodes
    and, not := And(0), InvBool(0)
    A := InstanceMap{0: and}
    N := InstanceMap{0: not}
    nodes := InstanceMap{Address{"And",0}: A, Address{"Not",0}: N}
    
    // Create Input Outputs
    and_ins, and_outs := and.GetParams()
    not_ins, not_outs := not.GetParams()
    in_A := NewParameter("A", "bool", and.GetAddr())
    in_B := and_ins["B"]
    out_C := not_outs["OUT"]
    inputs := ParamMap{"A": in_A, "B": in_B}
    outputs := ParamMap{"OUT": out_C}
    
    // Create Edges
    edges := EdgeMap{and_outs["OUT"]: []Parameter{not_ins["IN"]}}
    
    return NewGraph("logical_nand", id, nodes, edges, inputs, outputs)
}
    

/*
func Nor(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = !(in["A"].(bool) || in["B"].(bool))
    }
    name := "logical_nor"
    return opBinary(id,"bool","bool","bool",name,opfunc)
}

// Comparison
func GreaterEquals(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) >= toNum(in["B"])
    }
    name := "greater_equals"
    return opBinary(id,"num","num","bool",name,opfunc)
}
func LesserEquals(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) <= toNum(in["B"])
    }
    name := "lesser_equals"
    return opBinary(id,"num","num","bool",name,opfunc)
}
func NotEquals(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) != toNum(in["B"])
    }
    name := "not_equals"
    return opBinary(id,"num","num","bool",name,opfunc)
}*/