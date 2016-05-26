package flow

import (
    "fmt"
    "time"
)

const (
  STOPPING = "STOP" // Used to declare a stopping error
)

type Address struct {
    name string
    id InstanceID
}
func (a Address) GetName() string {return a.name}
func (a Address) GetID() string {return a.id}

// Used to declare an error in the flow pipeline
type FlowError struct{
    Ok bool
    Info string
    Addr Address
}

// Used to represent a parameter to a FunctionBlock
// Everything is private, as this struct is immutable
type Parameter struct {
    addr      Address
    paramName string
    paramType TypeStr
}
func (p Parameter) GetName() string {return p.paramName}
func (p Parameter) GetType() TypeStr {return p.paramType}
func (p Parameter) GetBlock() Address {return p.addr}

func NewParameter(name string, t TypeStr, addr Address) Parameter {
    return Parameter{addr: addr, paramType: t, paramName: name}
}

// Used to store the outputs of a FunctionBlock, while keeping it's reference.
type DataOut struct {
    Addr Address
    Values  ParamValues
}

// KeyValues and DataStreams are the types of values and functions
// Used universally inside FunctionBlocks
type TypeStr string
type InstanceID int
type ParamValues map[string]interface{}
type ParamTypes map[string]TypeStr
type ParamMap map[string]Parameter
type DataStream func(inputs ParamValues,
                     outputs chan DataOut,
                     stop chan bool,
                     err chan FlowError)

// The primary interface of the flowchart. Allows running, has a name, and has parameters.
type FunctionBlock interface{
    Run(inputs ParamValues,
        outputs chan DataOut,
        stop chan bool,
        err chan FlowError)
    GetParams() (inputs ParamMap, outputs ParamMap)
    GetAddr() Address
}

// A primitive function block that only
// contains a DataStream Function to run
type PrimitiveBlock struct {
    addr    Address
    fn      DataStream
    inputs  ParamTypes
    outputs ParamTypes
}

// Initializes a FunctionBlock object with given attributes, and an empty parameter list.
// The only way to create Methods's
func NewPrimitive(name string, function DataStream, inputs ParamTypes, outputs ParamTypes) FunctionBlock {
    return PrimitiveBlock{name: name,
                          fn: function,
                          inputs: inputs,
                          outputs: outputs}
}

// Returns a copy of FunctionBlock's InstanceId
func (m PrimitiveBlock) GetAddr() Address {return m.addr}

// Returns copies of all parameters in FunctionBlock
func (m PrimitiveBlock) GetParams() (inputs ParamMap, outputs ParamMap) {
    inputs = make(ParamMap, 0, len(m.inputs))
    for name, t := range m.inputs {
        inputs[name] = NewParameter(name, t, m.GetAddr())
    }

    outputs = make(ParamMap, 0, len(m.outputs))
    for name, t := range m.outputs {
        outputs[name] = NewParameter(name, t, m.GetAddr())
    }
    return
}

// Run the function
func (m PrimitiveBlock) Run(inputs ParamValues,
                            outputs chan DataOut,
                            stop chan bool,
                            err chan FlowError) {
    // Check types to ensure inputs are the type defined in input parameters
    if CheckTypes(inputs, m.inputs) {
        err <- FlowError{Ok: false, Info: "Inputs are impropper types.", Addr: m.GetAddr()}
        return
    }

    // Duplicate the given channel to pass to the enclosed function
    // Run the function
    f_err  := make(chan FlowError)
    f_out  := make(chan DataOut)
    f_stop := make(chan bool)
    go m.fn(inputs, f_out, f_stop, f_err)

    // Wait for a stop or an output
    for {
        select {
            case f_return := <-f_out:                                 // If an output is returned
                if CheckTypes(f_return.Values, m.outputs) {         // Check the types with output parameters
                    err <- FlowError {Ok: true}                       // If good, return no error
                    outputs <- DataOut{m.GetAddr(), f_return.Values}  // Along with the data
                    return                                            // And stop the function
                } else {
                    fmt.Println(f_return)
                    err <- FlowError{Ok: false, Info: "Wrong output type.", Addr: m.GetAddr()}
                    return
                }
            case <-stop:                              // If commanded to stop externally
                f_stop <- true                        // Pass it on to subfunction
                return                                // And stop immediately
            case temp_err := <-f_err:                 // If there is an error, save it
                if !temp_err.Ok {                     // See if it is bad
                    err <- temp_err                   // If it is bad, pass it up the chain
                    return                            // And stop the function
                }
        }
    }
}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func CheckTypes(values ParamValues, params ParamTypes) (ok bool) {
    for name, typestr := range params {                             // Iterate through all parameters and get their names and types
        val, _ := values[name]                                      // Get the value of this param from values
        if !CheckType(NewParameter(name,typestr,Address{}), val) {  // Check the type based on an empty parameter of type typestr
            return false                                            // If it's not valid, return false
        }
    }
    return true                                                     // If none are valid, return true
}

func CheckType(param Parameter, val interface{}) bool {
    t := param.GetType()
    switch val.(type) {
        case string:
            if t != "string" {return false}
        case int:
            if t != "int" && t != "num" {return false}
        case float64:
            if t != "float" && t != "num" {return false}
        case bool:
            if t != "bool" {return false}
        default:
            return true
    }
}

func BlockRun(blk FunctionBlock, f_in ParamValues, f_stop chan bool) (f_out chan DataOut,
                                                                      f_stop chan bool,
                                                                      f_err chan FlowError) {
    // Initialize channels
    f_out  := chan DataOut
    f_stop := chan bool
    f_err  := chan FlowError
        
    // Run in new goroutine
    go blk.Run(f_in, f_out, f_stop, f_err)
        
    return f_out, f_stop, f_err
}

func Timeout(stop chan bool, sleeptime int) {
    time.Sleep(time.Duration(sleeptime))
    stop <- true
}

func toNum(n interface{}) float64 {
    switch n.(type) {
        case int:
            return float64(n.(int))
        case float64:
            return n.(float64)
        default:
            panic("Wrong Type in toNum")
    }
}