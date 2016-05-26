package flow

// Creates a variety of blocks for paired operations
func opBinary(addr Address, aT,bT,cT TypeStr, outname string, opfunc func(in ParamValues, out *ParamValues)) FunctionBlock {
    // Create Plus block
    ins := ParamTypes{"A": aT, "B": bT}
    outs := ParamTypes{"OUT": cT}
    
    // Define the function as a closure
    runfunc := func(inputs ParamValues,
                     outputs chan DataOut,
                     stop chan bool,
                     err chan FlowError) {
        data := make(ParamValues)
        opfunc(inputs, &data)
        out := DataOut{Addr: addr, Values: data}
        outputs <- out
        return
    }
    
    // Initialize the block and return
    return NewPrimitive(addr.name, runfunc, ins, outs)
}

// Numeric Float Functions
func PlusFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(float64) + in["B"].(float64)
    }
    name := "numeric_plus_float"
    addr := Address{id: id, name: name}
    return opBinary(addr,"float","float","float",name,opfunc)
}
func SubFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(float64) - in["B"].(float64)
    }
    name := "numeric_subtract_float"
    addr := Address{id: id, name: name}
    return opBinary(addr,"float","float","float",name,opfunc)
}
func MultFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(float64) * in["B"].(float64)
    }
    name := "numeric_multiply_float"
    addr := Address{id: id, name: name}
    return opBinary(addr,"float","float","float",name,opfunc)
}
func DivFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(float64) / in["B"].(float64)
    }
    name := "numeric_divide_float"
    addr := Address{id: id, name: name}
    return opBinary(addr,"float","float","float",name,opfunc)
}

// Numeric Int Functions
func PlusInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(int) + in["B"].(int)
    }
    name := "numeric_plus_int"
    addr := Address{id: id, name: name}
    return opBinary(addr,"int","int","int",name,opfunc)
}
func SubInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(int) - in["B"].(int)
    }
    name := "numeric_subtract_int"
    addr := Address{id: id, name: name}
    return opBinary(addr,"int","int","int",name,opfunc)
}
func MultInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(int) * in["B"].(int)
    }
    name := "numeric_multiply_int"
    addr := Address{id: id, name: name}
    return opBinary(addr,"int","int","int",name,opfunc)
}
func DivInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(int) / in["B"].(int)
    }
    name := "numeric_divide_int"
    addr := Address{id: id, name: name}
    return opBinary(addr,"int","int","int",name,opfunc)
}
func Mod(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(int) % in["B"].(int)
    }
    name := "numeric_mod_int"
    addr := Address{id: id, name: name}
    return opBinary(addr,"int","int","int",name,opfunc)
}

// Boolean Logic Functions
func And(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(bool) && in["B"].(bool)
    }
    name := "logical_and"
    addr := Address{id: id, name: name}
    return opBinary(addr,"bool","bool","bool",name,opfunc)
}
func Or(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(bool) || in["B"].(bool)
    }
    name := "logical_or"
    addr := Address{id: id, name: name}
    return opBinary(addr,"bool","bool","bool",name,opfunc)
}
func Xor(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["A"].(bool) != in["B"].(bool)
    }
    name := "logical_xor"
    addr := Address{id: id, name: name}
    return opBinary(addr,"bool","bool","bool",name,opfunc)
}


// Comparison
func Greater(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) > toNum(in["B"])
    }
    name := "greater_than"
    addr := Address{id: id, name: name}
    return opBinary(addr,"num","num","bool",name,opfunc)
}
func Lesser(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) < toNum(in["B"])
    }
    name := "lesser_than"
    addr := Address{id: id, name: name}
    return opBinary(addr,"num","num","bool",name,opfunc)
}
func Equals(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = toNum(in["A"]) == toNum(in["B"])
    }
    name := "equals"
    addr := Address{id: id, name: name}
    return opBinary(addr,"num","num","bool",name,opfunc)
}
