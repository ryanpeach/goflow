package blocks

import ".."

// Creates a variety of blocks for paired operations
func opBinary(addr flow.Address, aT,bT,cT flow.TypeStr,
              outname string, opfunc func(in flow.ParamValues, out flow.ParamValues)) flow.FunctionBlock {
    // Create Plus block
    ins := flow.ParamTypes{"A": aT, "B": bT}
    outs := flow.ParamTypes{"OUT": cT}
    
    // Define the function as a closure
    runfunc := func(inputs flow.ParamValues,
                     outputs chan flow.DataOut,
                     stop chan bool,
                     err chan flow.FlowError) {
        data := make(flow.ParamValues)
        opfunc(inputs, data)
        out := flow.DataOut{Addr: addr, Values: data}
        outputs <- out
        return
    }
    
    // Initialize the block and return
    return flow.NewPrimitive(addr.GetName(), runfunc, ins, outs)
}

// Numeric Float Functions
func PlusFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) + in["B"].(float64)
    }
    name := "numeric_plus_float"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Float,flow.Float,flow.Float,name,opfunc)
}
func SubFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) - in["B"].(float64)
    }
    name := "numeric_subtract_float"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Float,flow.Float,flow.Float,name,opfunc)
}
func MultFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) * in["B"].(float64)
    }
    name := "numeric_multiply_float"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Float,flow.Float,flow.Float,name,opfunc)
}
func DivFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(float64) / in["B"].(float64)
    }
    name := "numeric_divide_float"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Float,flow.Float,flow.Float,name,opfunc)
}

// Numeric Int Functions
func PlusInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) + in["B"].(int)
    }
    name := "numeric_plus_int"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Int,flow.Int,flow.Int,name,opfunc)
}
func SubInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) - in["B"].(int)
    }
    name := "numeric_subtract_int"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Int,flow.Int,flow.Int,name,opfunc)
}
func MultInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) * in["B"].(int)
    }
    name := "numeric_multiply_int"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Int,flow.Int,flow.Int,name,opfunc)
}
func DivInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) / in["B"].(int)
    }
    name := "numeric_divide_int"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Int,flow.Int,flow.Int,name,opfunc)
}
func Mod(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(int) % in["B"].(int)
    }
    name := "numeric_mod_int"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Int,flow.Int,flow.Int,name,opfunc)
}

// Boolean Logic Functions
func And(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) && in["B"].(bool)
    }
    name := "logical_and"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,name,opfunc)
}
func Or(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) || in["B"].(bool)
    }
    name := "logical_or"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,name,opfunc)
}
func Xor(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["A"].(bool) != in["B"].(bool)
    }
    name := "logical_xor"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Bool,flow.Bool,flow.Bool,name,opfunc)
}


// Comparison
func Greater(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) > flow.ToNum(in["B"])
    }
    name := "greater_than"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,name,opfunc)
}
func Lesser(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) < flow.ToNum(in["B"])
    }
    name := "lesser_than"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,name,opfunc)
}
func Equals(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = flow.ToNum(in["A"]) == flow.ToNum(in["B"])
    }
    name := "equals"
    addr := flow.NewAddress(id, name)
    return opBinary(addr,flow.Num,flow.Num,flow.Bool,name,opfunc)
}
