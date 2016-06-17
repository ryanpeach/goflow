package flow

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// KeyValues and DataStreams are the types of values and functions
// Used universally inside FunctionBlocks
type Type string                        // Used to represent parameter types
type InstanceID int                     // Used to separate blocks of the same name in graphs from one another
type NameMap map[string]string          // Used primarily in loops since they only contain one block
type ParamValues map[string]interface{} // Used to map parameter names to their values for quick lookup
type ParamTypes map[string]Type         // Used to map parameter names to their types for quick lookup
type DataStream func(inputs ParamValues,
	outputs chan ParamValues,
	stop chan bool,
						err chan *Error) // The primary run function for primitive blocks
type InstanceMap map[Address]FunctionBlock   // Maps addresses to blocks in a graph
type EdgeMap map[ParamAddress][]ParamAddress // Maps paramters to multiple parameters
type BlockMap map[string]FunctionBlock       // Maps a name to a function block
type ParamMap map[string]ParamAddress        // Maps a name to a parameter
type ParamLstMap map[string][]ParamAddress   // Maps a name to a list of parameters

// Supporting function for copying ParamTypes map
func (X ParamTypes) Copy() (out ParamTypes) {
	out = make(ParamTypes)
	for name, t := range X {
		out[name] = t
	}
	return
}

// Supporting function for copying ParamValues map
func (X ParamValues) Copy() (out ParamValues) {
	out = make(ParamValues)
	for name, t := range X {
		out[name] = t
	}
	return
}

// Used in maps to reference FunctionBlocks
type Address struct {
	Name string
	ID   InstanceID
}

// Used in maps to reference FunctionBlock parameters
// FIXME: Is this still needed with the new Graph structure?
type ParamAddress struct {
	Name     string
	Addr     Address
	T        Type
	is_input bool
}

// Error constants:
const (
	NIL                  = iota // No error
	STOPPING             = iota // Used to signal that a block is stopping
	NOT_INPUT_ERROR      = iota // Used to signal that a parameter is not an input
	TYPE_ERROR           = iota // Used to signal that two parameters or a param and a value are incompatible types
	DNE_ERROR            = iota // Something Does not Exist
	ALREADY_EXISTS_ERROR = iota // Something already exists
	NOT_READY_ERROR      = iota // Not ready to do what you wanted
	VALUE_ERROR          = iota // Value is not acceptable
)

// Used to declare a general error.
// Class: Contains a simple class as declared in the const above for easy code handling.
// Info:  Also contains a detailed string for user understanding.
type Error struct {
	Class int
	Info  string
}

// Used to declare an error while keeping
// the address of the block which returned it.
// Used as the primary error type in Graphs.
// Inherits Error.
type FlowError struct {
	*Error
	Addr Address
}

// The error interface's required function.
func (e Error) Error() string {
	return e.Info
}

// Easily create a flow error without first creating an Error struct.
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

// A map of Type objects linked to the reflect types that are valid for them.
var Types = map[Type][]reflect.Type{
	String:   {reflect.TypeOf("")},
	Int:      {reflect.TypeOf(5)},
	Float:    {reflect.TypeOf(5.1)},
	Num:      {reflect.TypeOf(5), reflect.TypeOf(5.1)},
	Bool:     {reflect.TypeOf(true)},
	NumArray: {reflect.TypeOf([]float64{})}}

// Checks if all keys in params are present in values
// And that all values are of their appropriate types as labeled in in params
func CheckTypes(values ParamValues, params ParamTypes) *Error {
	for name, typestr := range params { // Iterate through all parameters and get their names and types
		val, exists := values[name] // Get the value of this param from values
		switch {
		case !exists:
			return &Error{DNE_ERROR, "Some param does not exist in values."}
		case !CheckType(typestr, val):
			return &Error{TYPE_ERROR, "Not all types are compatible."}
		}
	}
	return nil // If none are valid, return true
}

// Checks if type t is compatible with val.
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
func AddType(newName Type, compatible ...reflect.Type) {
	_, exists := Types[newName]
	if !exists {
		Types[newName] = compatible
	} else {
		Types[newName] = append(Types[newName], compatible...)
	}
}

func CheckSame(t1, t2 Type) bool {
	switch {
	case t1 == t2:
		return true
	case t1 == Num && t2 == Int:
		return true
	case t1 == Num && t2 == Float:
		return true
	case t2 == Num && t1 == Int:
		return true
	case t2 == Num && t1 == Float:
		return true
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

// Creates a logger for logging data.
// If logmode is "none" it does not print to screen or to file.
// If logmode is "screen" it prints to the screen.
// All other logmodes define the file location.
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
