package flow

type Node struct {
    f       FunctionBlock
    inputs  map[string]*InParameter
    outputs map[string]*OutParameter
}

func (n Node) Run(stop chan bool, err chan *FlowError, id InstanceID) {
    logger  := CreateLogger("none", "[INFO]")
    blk_ins := make(ParamValues)
    logger.Println(n.f.GetName(), "\tReading Params... ")
    for name, in_param := range n.inputs {
        blk_ins[name] = <-in_param.val
        logger.Println(n.f.GetName(), "Found: ", name)
    }
    
    logger.Println(n.f.GetName(), "\tRunning... ")
    blk_outs := make(chan ParamValues, 1)
    blk_stop := make(chan bool, 1)
    blk_err  := make(chan *FlowError, 1)
    go n.f.Run(blk_ins, blk_outs, blk_stop, blk_err, id)
    
    logger.Println(n.f.GetName(), "\tWaiting... ")
    select {
        case out := <- blk_outs:
            for name, out_param := range n.outputs {
                val, exists := out[name]
                if exists {
                    out_param.PassValue(val)
                }
            }
        case <-stop:
            blk_stop <- true
        case temp := <-blk_err:
            err <- temp
    }
    logger.Println(n.f.GetName(), "\tDone!")
    return
}

type InParameter struct {
    t    Type
    val  chan interface{}
    source Edge
}

type OutParameter struct {
    t     Type
    edges []*InParameter
}

func (o OutParameter) PassValue(val interface{}) {
    for _, in_param := range o.edges {
        in_param.val <- val
    }
}

type Constant struct {
    t    Type
    val  interface{}
    edge *InParameter
}

func (c Constant) PassValue(val interface{}) {
    c.edge.val <- val
}

type Edge interface {
    PassValue(val interface{})
}

type Graph struct {
    name    string
    nodes   map[Address]*Node
    consts  []*Constant
    inputs  map[string]*OutParameter
    outputs map[string]*InParameter
}

func createInParams(inputs ParamTypes) (map[string]*InParameter) {
    ins := make(map[string]*InParameter, len(inputs))
    for name, t := range inputs {
        ins[name] = &InParameter{t, make(chan interface{}, 1), nil}
    }
    return ins
}

func createOutParams(outputs ParamTypes) (map[string]*OutParameter) {
    outs  := make(map[string]*OutParameter, len(outputs))
    for name, t := range outputs {
        outs[name] = &OutParameter{t, make([]*InParameter, 0)}
    }
    return outs
}

func NewGraph(name string, inputs, outputs ParamTypes) (*Graph, *Error) {
    // Handle Errors
    nilGraph := &Graph{}
    switch {
        case len(inputs) == 0:
            return nilGraph, &Error{DNE_ERROR, "Inputs has a length of Zero."}
        case len(outputs) == 0:
            return nilGraph, &Error{DNE_ERROR, "Outputs has a length of Zero."}
    }
    
    // Create input and output parameter structures from map
    ins  := createOutParams(inputs)
    outs := createInParams(outputs)
    
    // Create placeholders for nodes and constants
    nodes  := make(map[Address]*Node)
    consts := make([]*Constant, 0)
    return &Graph{name, nodes, consts, ins, outs}, nil
}

func (g Graph) FindInParam(param_name string, param_addr Address) (*InParameter, *Error) {
    nd, nd_exists := g.nodes[param_addr]
    if nd_exists {
        param, p_exists   := nd.inputs[param_name]
        if p_exists {
            return param, nil
        } else {
            return nil, &Error{DNE_ERROR, "Parameter does not exist."}
        }
    } else {
        return nil, &Error{DNE_ERROR, "Node does not exist."}
    }
}

func (g Graph) FindOutParam(param_name string, param_addr Address) (*OutParameter, *Error) {
    nd, nd_exists := g.nodes[param_addr]
    if nd_exists {
        param, p_exists   := nd.outputs[param_name]
        if p_exists {
            return param, nil
        } else {
            return nil, &Error{DNE_ERROR, "Parameter does not exist."}
        }
    } else {
        return nil, &Error{DNE_ERROR, "Node does not exist."}
    }
}

func (g *Graph) AddConstant(val interface{}, param_addr Address, param_name string) *Error {
    in_param, err := g.FindInParam(param_name, param_addr)
    switch {
        case err != nil:
            return err
        case in_param.source != nil:
            return &Error{ALREADY_EXISTS_ERROR, "Parameter is already linked."}
        case !CheckType(in_param.t, val):
            return &Error{TYPE_ERROR, "Parameter is not the same type as val."}
        default:
            new_const := &Constant{in_param.t, val, in_param}
            g.consts = append(g.consts, new_const)
            in_param.source = new_const
    }
    return nil
}

func (g *Graph) AddNode(blk FunctionBlock, addr Address) *Error {
    _, exists := g.nodes[addr]
    if !exists {
        in_map, out_map := blk.GetParams()
        inputs  := createInParams(in_map)
        outputs := createOutParams(out_map)
        g.nodes[addr] = &Node{blk, inputs, outputs}
        return nil
    } else {
        return &Error{ALREADY_EXISTS_ERROR, "blk is already a node in Graph."}
    }
}

// out_addr[out_param_name] -> in_addr[in_param_name]
func (g *Graph) AddEdge(out_addr Address, out_param_name string,
                       in_addr Address, in_param_name string) *Error {
    out_param, out_err := g.FindOutParam(out_param_name, out_addr)
    in_param,  in_err  := g.FindInParam(in_param_name, in_addr)
    switch {
        case out_err != nil:
            return out_err
        case in_err != nil:
            return in_err
        case in_param.source != nil:
            return &Error{ALREADY_EXISTS_ERROR, "in_param already has a source."}
        case !CheckSame(out_param.t, in_param.t):
            return &Error{TYPE_ERROR, "in_param and out_param incompatible types."}
        default:
            out_param.edges = append(out_param.edges, in_param)       // Add the input as an edge for the output
            in_param.source = out_param             // Set the input source
            return nil                              // No error
    }
}

// self[self_param_name] -> in_addr[in_param_name]
func (g *Graph) LinkIn(self_param_name string, in_param_name string, in_addr Address) *Error {
    in_param, err           := g.FindInParam(in_param_name, in_addr)
    self_param, self_exists := g.inputs[self_param_name]
    switch {
        case err != nil:
            return err
        case !self_exists:
            return &Error{DNE_ERROR, "Self param does not exist."}
        default:
            self_param.edges = append(self_param.edges, in_param)
            return nil
    }
}

// out_addr[out_param_name] -> self[self_param_name]
func (g *Graph) LinkOut(out_addr Address, out_param_name string, self_param_name string) *Error {
    out_param, err          := g.FindOutParam(out_param_name, out_addr)
    self_param, self_exists := g.outputs[self_param_name]
    switch {
        case err != nil:
            return err
        case !self_exists:
            return &Error{DNE_ERROR, "Self param does not exist."}
        case self_param.source != nil:
            return &Error{ALREADY_EXISTS_ERROR, "Self param already has source."}
        default:
            out_param.edges = append(out_param.edges, self_param)
            self_param.source = out_param  // Connect out param to self output
            return nil                     // Return no error
    }
}

// Returns copies of all parameters in FunctionBlock
func (g Graph) GetParams() (inputs, outputs ParamTypes) {
    inputs = make(ParamTypes)
    for name, param := range g.inputs {
        inputs[name] = param.t
    }
    outputs = make(ParamTypes)
    for name, param := range g.outputs {
        outputs[name] = param.t
    }
    return
}

// Returns a copy of FunctionBlock's InstanceId
func (g Graph) GetName() string {return g.name}

func (g Graph) Run(inputs ParamValues,
                   outputs chan ParamValues,
                   stop chan bool,
                   err chan *FlowError, id InstanceID) {
    
    ADDR := Address{g.GetName(), id}
    logger   := CreateLogger("none", "[INFO]")
    
    // Pass all inputs to input parameters
    logger.Println("Passing Inputs... ", inputs)
    for name, val := range inputs {
        param_in, exists := g.inputs[name]
        if exists {
            param_in.PassValue(val)
        } else {
            err <- &FlowError{&Error{DNE_ERROR, "Not all inputs fulfilled."}, ADDR}
            logger.Println("Not all inputs fulfilled.")
            return
        }
    }
    
    // Load Constants
    logger.Println("Passing Constants...")
    logger.Println("Constants: ", g.consts)
    for _, c := range g.consts {
        logger.Println("Passing... ", c)
        c.PassValue(c.val)
    }
    
    // Run all nodes
    logger.Println("Starting Nodes...")
    all_stop := make([](chan bool), 0, len(g.nodes))
    blk_err  := make(chan *FlowError, 1)
    for addr, nd := range g.nodes {
        blk_stop := make(chan bool, 1)
        go nd.Run(blk_stop, blk_err, addr.ID)
        all_stop = append(all_stop, blk_stop)
    }
    
    allStop := func() {
        logger.Println("Stopping...")
        for _, blk_stop := range all_stop {
            blk_stop <- true
        }
    }

    // Wait for all output parameters to be set
    logger.Println("Waiting...")
    data_out := make(ParamValues)
    for name, out_param := range g.outputs {
        logger.Println(name)
        select {
            case <-stop:
                allStop()
                return
            case temp_err := <-blk_err:
                err <- temp_err
                logger.Println(temp_err)
                allStop()
                return
            case temp := <-out_param.val:
                logger.Println(temp)
                data_out[name] = temp
        }
        logger.Println("-------------------------------")
        logger.Println(data_out)
    }
    
    // If you made it this far, return the output
    outputs <- data_out
    allStop()
    return
}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func (g Graph) checkTypes(values ParamValues, params ParamMap) (ok bool) {
    for name, param := range params {
        val, exists := values[name]
        if !exists || !CheckType(param.T, val) {
            return false
        }
    }
    return true
}