package blocks

import ".."

// Creates a variety of blocks for paired operations
func opUnary(id flow.InstanceID, inT flow.Type, outT flow.Type, outname string,
	opfunc func(in flow.ParamValues, out flow.ParamValues)) (flow.FunctionBlock, flow.Address) {
	// Create Plus block
	ins := flow.ParamTypes{"IN": inT}
	outs := flow.ParamTypes{"OUT": outT}
	addr := flow.Address{outname, id}

	// Define the function as a closure
	runfunc := func(inputs flow.ParamValues,
		outputs chan flow.ParamValues,
		stop chan bool,
		err chan *flow.Error) {
		data := make(flow.ParamValues)
		opfunc(inputs, data)
		out := data
		outputs <- out
		return
	}

	// Initialize the block and return
	return flow.NewPrimitive(outname, runfunc, ins, outs), addr
}

// Type conversions
func FloattoInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = int(in["IN"].(float64))
	}
	name := "float_to_int"
	return opUnary(id, flow.Float, flow.Int, name, opfunc)
}
func InttoFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = float64(in["IN"].(int))
	}
	name := "int_to_float"
	return opUnary(id, flow.Int, flow.Float, name, opfunc)
}

// Mathematical
func Inc(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = in["IN"].(int) + 1
	}
	name := "increment"
	return opUnary(id, flow.Int, flow.Int, name, opfunc)
}
func Dec(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = in["IN"].(int) - 1
	}
	name := "decrement"
	return opUnary(id, flow.Int, flow.Int, name, opfunc)
}
func InvFloat(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = -in["IN"].(float64)
	}
	name := "invert_float"
	return opUnary(id, flow.Float, flow.Float, name, opfunc)
}
func InvInt(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = -in["IN"].(int)
	}
	name := "invert_int"
	return opUnary(id, flow.Int, flow.Int, name, opfunc)
}

// Logical
func InvBool(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = !in["IN"].(bool)
	}
	name := "invert_bool"
	return opUnary(id, flow.Bool, flow.Bool, name, opfunc)
}

// Arrays
func Len(id flow.InstanceID) (flow.FunctionBlock, flow.Address) {
	opfunc := func(in flow.ParamValues, out flow.ParamValues) {
		out["OUT"] = len(in["IN"].([]float64))
	}
	name := "array_len"
	return opUnary(id, flow.NumArray, flow.Int, name, opfunc)
}
