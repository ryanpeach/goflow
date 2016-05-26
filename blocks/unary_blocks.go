package flow

// Creates a variety of blocks for paired operations
func opUnary(id InstanceID, inT TypeStr, outT TypeStr, outname string, opfunc func(in ParamValues, out *ParamValues)) FunctionBlock {
    // Create Plus block
    ins := ParamTypes{"IN": inT}
    outs := ParamTypes{"OUT": outT}
    
    // Define the function as a closure
    runfunc := func(inputs ParamValues,
                     outputs chan DataOut,
                     stop chan bool,
                     err chan FlowError) {
        data := make(ParamValues)
        opfunc(inputs, &data)
        addr := Address{id: id, name: outname}
        out := DataOut{Addr: addr, Values: data}
        outputs <- out
        return
    }
    
    // Initialize the block and return
    return PrimitiveBlock{name: outname, fn: runfunc, id: id, inputs: ins, outputs: outs}
}

// Type conversions
func FloattoInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = int(in["IN"].(float64))
    }
    name := "float_to_int"
    return opUnary(id,"float","int",name,opfunc)
}
func InttoFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = float64(in["IN"].(int))
    }
    name := "int_to_float"
    return opUnary(id,"int","float",name,opfunc)
}

// Mathematical
func Inc(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["IN"].(int) + 1
    }
    name := "increment"
    return opUnary(id,"int","int",name,opfunc)
}
func Dec(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = in["IN"].(int) - 1
    }
    name := "decrement"
    return opUnary(id,"int","int",name,opfunc)
}
func InvFloat(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = -in["IN"].(float64)
    }
    name := "invert_float"
    return opUnary(id,"float","float",name,opfunc)
}
func InvInt(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = -in["IN"].(int)
    }
    name := "invert_int"
    return opUnary(id,"int","int",name,opfunc)
}

// Logical
func InvBool(id InstanceID) FunctionBlock {
    opfunc := func(in ParamValues, out *ParamValues) {
        (*out)["OUT"] = !in["IN"].(bool)
    }
    name := "invert_bool"
    return opUnary(id,"bool","bool",name,opfunc)
}