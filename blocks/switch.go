package blocks

import ".."

func InputSwitch(id flow.InstanceID, t flow.Type) (flow.FunctionBlock, flow.Address) {
	runfunc := func(inputs flow.ParamValues,
		outputs chan flow.ParamValues,
		stop chan bool,
		err chan *flow.Error) {
		out := make(flow.ParamValues)
		if inputs["Condition"].(bool) {
			out["OUT"] = inputs["A"]
		} else {
			out["OUT"] = inputs["B"]
		}
		outputs <- out
	}
	name := "input_switch"
	addr := flow.Address{name, id}
	ins := flow.ParamTypes{"A": t, "B": t, "Condition": flow.Bool}
	outs := flow.ParamTypes{"OUT": t}
	return flow.NewPrimitive(name, runfunc, ins, outs), addr
}

func OutputSwitch(id flow.InstanceID, t flow.Type) (flow.FunctionBlock, flow.Address) {
	runfunc := func(inputs flow.ParamValues,
		outputs chan flow.ParamValues,
		stop chan bool,
		err chan *flow.Error) {
		out := make(flow.ParamValues)
		if inputs["Condition"].(bool) {
			out["A"] = inputs["IN"]
		} else {
			out["B"] = inputs["IN"]
		}
		outputs <- out
	}
	name := "output_switch"
	addr := flow.Address{name, id}
	ins := flow.ParamTypes{"IN": t, "Condition": flow.Bool}
	outs := flow.ParamTypes{"A": t, "B": t}
	return flow.NewPrimitive(name, runfunc, ins, outs), addr
}
