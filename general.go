package flow

import (
    "os"
    "log"
    "io/ioutil"
    "fmt"
    "reflect"
)

// KeyValues and DataStreams are the types of values and functions
// Used universally inside FunctionBlocks
type Type string
type InstanceID int
type NameMap map[string]string
type ParamValues map[string]interface{}
type ParamTypes  map[string]Type
type DataStream func(inputs ParamValues,
                     outputs chan ParamValues,
                     stop chan bool,
                     err chan *FlowError)
type InstanceMap map[Address]FunctionBlock
type EdgeMap map[ParamAddress][]ParamAddress
type BlockMap map[string]FunctionBlock
type ParamMap map[string]ParamAddress
type ParamLstMap map[string][]ParamAddress

type Address struct {
    Name string
    ID   InstanceID
}

type ParamAddress struct {
    Name     string
    Addr     Address
    T        Type
    is_input bool
}

const (
// Errors
    NIL                  = iota
    STOPPING             = iota // Used to declare a stopping error
    NOT_INPUT_ERROR      = iota
    TYPE_ERROR           = iota
    DNE_ERROR            = iota
    ALREADY_EXISTS_ERROR = iota
    NOT_READY_ERROR      = iota
    VALUE_ERROR          = iota
)

// Used to declare an error in the flow pipeline
type Error struct{
    Class int
    Info  string
}
type FlowError struct{
    *Error
    Addr Address
}
func (e Error) Error() string {
    return e.Info
}
func NewFlowError(Class int, Info string, Addr Address) *FlowError {
    return &FlowError{&Error{Class, Info}, Addr}
}

// Types
const (
    Float    Type = "Float"
    String   Type = "String"
    Int      Type = "Int"
    Num      Type = "Num"
    Bool     Type = "Bool"
    NumArray Type = "NumArray"
)
var Types = map[Type][]reflect.Type {
    String: {reflect.TypeOf("")},
    Int:    {reflect.TypeOf(5)},
    Float:  {reflect.TypeOf(5.1)},
    Num:    {reflect.TypeOf(5), reflect.TypeOf(5.1)},
    Bool:   {reflect.TypeOf(true)},
    NumArray: {reflect.TypeOf([]float64{})}}


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
    T, exists := Types[t]
    if exists {
        for _, t := range T {
            if t == reflect.TypeOf(val) {
                return true
            }
        }
    }
    return false
}

// Adds a new name or appends new type compatibilities to a preexisting type.
func AddType(newName Type, compatible []reflect.Type) {
    _, exists := Types[newName]
    if !exists {
        Types[newName] = compatible
    } else {
        Types[newName] = append(Types[newName], compatible...)
    }
}

func CheckSame(t1, t2 Type) bool {
    switch {
        case t1 == t2: return true
        case t1 == Num && t2 == Int:   return true
        case t1 == Num && t2 == Float: return true
        case t2 == Num && t1 == Int:   return true
        case t2 == Num && t1 == Float: return true
    }
    return false
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

// Supporting function for copying ParamTypes map
func CopyTypes(X ParamTypes) (out ParamTypes) {
    out = make(ParamTypes)
    for name, t := range X {
        out[name] = t
    }
    return
}
func CopyValues(X ParamValues) (out ParamValues) {
    out = make(ParamValues)
    for name, t := range X {
        out[name] = t
    }
    return
}

func checkInputs(inputs ParamValues, req_inputs ParamTypes) (ok bool) {
    for name, _ := range req_inputs {
        _, exists := inputs[name]
        if !exists {
            return false
        }
    }
    return true
}