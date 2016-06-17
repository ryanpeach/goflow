package blocks

import ".."

// Creates a variety of blocks for paired operations
func opBinary(addr flow.Address, aT,bT,cT flow.Type, aN,bN,cN string,
              outname string, opfunc func(in flow.ParamValues, out flow.ParamValues)) flow.FunctionBlock {
    // Create Plus block
    ins := flow.ParamTypes{aN: aT, bN: bT}
    outs := flow.ParamTypes{cN: cT}
    
    // Define the function as a closure
    runfunc := func(inputs flow.ParamValues,
                     outputs chan flow.ParamValues,
                     stop chan bool,
                     err chan *flow.FlowError) {
        data := make(flow.ParamValues)
        opfunc(inputs, data)
        outputs <- data
        return
    }
    
    // Initialize the block and return
    return flow.NewPrimitive(addr.Name, runfunc, ins, outs)
}

// Numeric Float Functions
func PlusFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) + in["B"].(float64)
    }
    name := "numeric_plus_float"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Float,flow.Float,flow.Float,"A","B","OUT",name,opfunc), addr
}
func SubFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) - in["B"].(float64)
    }
    name := "numeric_subtract_float"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Float,flow.Float,flow.Float,"A","B","OUT",name,opfunc), addr
}
func MultFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) * in["B"].(float64)
    }
    name := "numeric_multiply_float"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Float,flow.Float,flow.Float,"A","B","OUT",name,opfunc), addr
}
func DivFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) / in["B"].(float64)
    }
    name := "numeric_divide_float"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Float,flow.Float,flow.Float,"A","B","OUT",name,opfunc), addr
}

// Numeric Int Functions
func PlusInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) + in["B"].(int)
    }
    name := "numeric_plus_int"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Int,flow.Int,flow.Int,"A","B","OUT",name,opfunc), addr
}
func SubInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) - in["B"].(int)
    }
    name := "numeric_subtract_int"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Int,flow.Int,flow.Int,"A","B","OUT",name,opfunc), addr
}
func MultInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) * in["B"].(int)
    }
    name := "numeric_multiply_int"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Int,flow.Int,flow.Int,"A","B","OUT",name,opfunc), addr
}
func DivInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) / in["B"].(int)
    }
    name := "numeric_divide_int"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Int,flow.Int,flow.Int,"A","B","OUT",name,opfunc), addr
}
func Mod(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) % in["B"].(int)
    }
    name := "numeric_mod_int"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Int,flow.Int,flow.Int,"A","B","OUT",name,opfunc), addr
}

// Boolean Logic Functions
func And(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) && in["B"].(bool)
    }
    name := "logical_and"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,"A","B","OUT",name,opfunc), addr
}
func Or(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) || in["B"].(bool)
    }
    name := "logical_or"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,"A","B","OUT",name,opfunc), addr
}
func Xor(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) != in["B"].(bool)
    }
    name := "logical_xor"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,"A","B","OUT",name,opfunc), addr
}


// Comparison
func Greater(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) > flow.ToNum(in["B"])
    }
    name := "greater_than"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,"A","B","OUT",name,opfunc), addr
}
func Lesser(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) < flow.ToNum(in["B"])
    }
    name := "lesser_than"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,"A","B","OUT",name,opfunc), addr
}
func Equals(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) == flow.ToNum(in["B"])
    }
    name := "equals"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,"A","B","OUT",name,opfunc), addr
}


// Array's
func Index(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = (in["X"].([]float64))[in["Index"].(int)]
    }
    name := "index"
    addr := flow.Address{name, id}
    return opBinary(addr,flow.NumArray,flow.Int,flow.Float,"X","Index","OUT",name,opfunc), addr
}