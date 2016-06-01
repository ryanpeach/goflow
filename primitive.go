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
func NewAddress(id InstanceID, name string) Address {return Address{name: name, id: id}}
func (a Address) GetName() string {return a.name}
func (a Address) GetID() InstanceID {return a.id}

// Used to declare an error in the flow pipeline
type FlowError struct{
    Ok bool
    Info string
    Addr Address
}

// Used to store the outputs of a FunctionBlock, while keeping it's reference.
type DataOut struct {
    Addr Address
    Values  ParamValues
}

// KeyValues and DataStreams are the types of values and functions
// Used universally inside FunctionBlocks
type Type int
type InstanceID int
type ParamValues map[string]interface{}
type ParamTypes map[string]Type
type DataStream func(inputs ParamValues,
                     outputs chan DataOut,
                     stop chan bool,
                     err chan FlowError)

// The primary interface of the flowchart. Allows running, has a name, and has parameters.
type FunctionBlock interface{
    Run(inputs ParamValues,
        outputs chan DataOut,
        stop chan bool,
        err chan FlowError,
        id InstanceID)
    GetParams() (inputs ParamTypes, outputs ParamTypes)
    GetName() string
}

// A primitive function block that only
// contains a DataStream Function to run
type PrimitiveBlock struct {
    name    string
    fn      DataStream
    inputs  ParamTypes
    outputs ParamTypes
}

// Initializes a FunctionBlock object with given attributes, and an empty parameter list.
// The only way to create Methods's
var nblocks map[string]InstanceID = make(map[string]InstanceID)
func NewPrimitive(name string, function DataStream, inputs ParamTypes, outputs ParamTypes) FunctionBlock {
    nblocks[name] += 1
    return PrimitiveBlock{name: name,
                          fn: function,
                          inputs: inputs,
                          outputs: outputs}
}

// Returns a copy of FunctionBlock's InstanceId
func (m PrimitiveBlock) GetName() string {return m.name}

// Returns copies of all parameters in FunctionBlock
func (m PrimitiveBlock) GetParams() (inputs ParamTypes, outputs ParamTypes) {
    return m.inputs, m.outputs
}

// Run the function
func (m PrimitiveBlock) Run(inputs ParamValues,
                            outputs chan DataOut,
                            stop chan bool,
                            err chan FlowError,
                            id InstanceID) {
    // Check types to ensure inputs are the type defined in input parameters
    ADDR := NewAddress(-1, m.GetName())
    if !CheckTypes(inputs, m.inputs) {
        err <- FlowError{Ok: false, Info: "Inputs are impropper types.", Addr: ADDR}
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
                if CheckTypes(f_return.Values, m.outputs) {           // Check the types with output parameters
                    err <- FlowError{Ok: true}                        // If good, return no error
                    outputs <- DataOut{ADDR, f_return.Values}  // Along with the data
                    return                                            // And stop the function
                } else {
                    fmt.Println(f_return)
                    err <- FlowError{Ok: false, Info: "Wrong output type.", Addr: ADDR}
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

// Types
const (
    String = iota
    Int    = iota
    Float  = iota
    Num    = iota
    Bool   = iota
)

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func CheckTypes(values ParamValues, params ParamTypes) (ok bool) {
    for name, typestr := range params {                             // Iterate through all parameters and get their names and types
        val := values[name]                                      // Get the value of this param from values
        if !CheckType(typestr, val) {  // Check the type based on an empty parameter of type typestr
            fmt.Println(typestr, val)
            return false                                            // If it's not valid, return false
        }
    }
    return true                                                    // If none are valid, return true
}

func CheckType(t Type, val interface{}) bool {
    switch val.(type) {
        case string:
            if t == String {return true}
        case int:
            if t == Int || t == Num {return true}
        case float64:
            if t == Float || t == Num {return true}
        case bool:
            if t == Bool {return true}
    }
    return false
}

// An easy way to initialize a block and get it's channels
func BlockRun(blk FunctionBlock, f_in ParamValues, id InstanceID) (f_out chan DataOut,
                                                                   f_stop chan bool,
                                                                   f_err chan FlowError) {
    // Initialize channels
    f_out  = make(chan DataOut)
    f_stop = make(chan bool)
    f_err  = make(chan FlowError)
        
    // Run in new goroutine
    go blk.Run(f_in, f_out, f_stop, f_err, id)
    return
}

// A Timeout block that can pass to the stop channel
func Timeout(stop chan bool, sleeptime int) {
    time.Sleep(time.Duration(sleeptime))
    stop <- true
}

// Converts a Num type interface to float64 for numeric processing.
func ToNum(n interface{}) float64 {
    switch n.(type) {
        case int:
            return float64(n.(int))
        case float64:
            return n.(float64)
        default:
            panic("Wrong Type in toNum")
    }
}