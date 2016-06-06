package flow

import (
    "os"
    "log"
    "io/ioutil"
    "fmt"
)

// KeyValues and DataStreams are the types of values and functions
// Used universally inside FunctionBlocks
type Type int
type InstanceID int
type ParamValues map[string]interface{}
type ParamTypes  map[string]Type
type DataStream func(inputs ParamValues,
                     outputs chan DataOut,
                     stop chan bool,
                     err chan FlowError)
type InstanceMap map[Address]FunctionBlock
type EdgeMap map[ParamAddress][]ParamAddress
type BlockMap map[string]FunctionBlock
type ParamMap map[string]ParamAddress
type ParamLstMap map[string][]ParamAddress

// Errors
const (
    STOPPING = "STOP" // Used to declare a stopping error
    NOT_INPUT_ERROR = "Parameter is not an input."
    TYPE_ERROR = "Type Check did not confirm compatablitiy."
    LINK_EXISTS_ERROR = "Link already exists for that input."
    DNE_ERROR = "Parameter does not exist."
)

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
    fmt.Println("Wrong Type: ", val, t)
    return false
}

func CheckCompatibility(t1, t2 Type) bool {
    return t1 == t2
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

// Source: http://changelog.ca/log/2015/03/09/golang
func CreateLogger(logMode, tag string) *log.Logger {
    switch logMode {
        case "none":
            out := ioutil.Discard
            return log.New(out, tag, log.Lshortfile)
        case "screen":
            out := os.Stdout
            return log.New(out, tag, log.Lshortfile)
        default:
            out, err := os.Create(logMode)
            if nil != err {
                panic(err.Error())
            } else {
                return log.New(out, tag, log.Lshortfile)
            }
    }
}