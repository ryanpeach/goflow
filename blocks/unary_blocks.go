package blocks

import ".."

// Creates a variety of blocks for paired operations
func opUnary(id flow.InstanceID, inT flow.TypeStr, outT flow.TypeStr, outname string,
             opfunc func(in flow.ParamValues, out flow.ParamValues)) flow.FunctionBlock {
    // Create Plus block
    ins := flow.ParamTypes{"IN": inT}
    outs := flow.ParamTypes{"OUT": outT}
    
    // Define the function as a closure
    runfunc := func(inputs flow.ParamValues,
                     outputs chan flow.DataOut,
                     stop chan bool,
                     err chan flow.FlowError) {
        data := make(flow.ParamValues)
        opfunc(inputs, data)
        addr := flow.NewAddress(id, outname)
        out := flow.DataOut{Addr: addr, Values: data}
        outputs <- out
        return
    }
    
    // Initialize the block and return
    return flow.NewPrimitive(outname, runfunc, ins, outs)
}

// Type conversions
func FloattoInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = int(in["IN"].(float64))
    }
    name := "float_to_int"
    return opUnary(id,"float","int",name,opfunc)
}
func InttoFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = float64(in["IN"].(int))
    }
    name := "int_to_float"
    return opUnary(id,"int","float",name,opfunc)
}

// Mathematical
func Inc(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["IN"].(int) + 1
    }
    name := "increment"
    return opUnary(id,"int","int",name,opfunc)
}
func Dec(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = in["IN"].(int) - 1
    }
    name := "decrement"
    return opUnary(id,"int","int",name,opfunc)
}
func InvFloat(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = -in["IN"].(float64)
    }
    name := "invert_float"
    return opUnary(id,"float","float",name,opfunc)
}
func InvInt(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = -in["IN"].(int)
    }
    name := "invert_int"
    return opUnary(id,"int","int",name,opfunc)
}

// Logical
func InvBool(id flow.InstanceID) flow.FunctionBlock {
    opfunc := func(in flow.ParamValues, out flow.ParamValues) {
        out["OUT"] = !in["IN"].(bool)
    }
    name := "invert_bool"
    return opUnary(id,"bool","bool",name,opfunc)
}