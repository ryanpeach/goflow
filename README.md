# goflow
A LabVIEW and TensorFlow Inspired Graph-Based Programming Environment for AI handled within the Go Programming Language.

Originally written for github.com/ryanpeach/SurvivalAI

This package is still in development, so this is just a basic tutorial and roadmap.

## FunctionBlocks
All elements in this library implement an interface type FunctionBlock:

    type FunctionBlock interface{
        Run(inputs ParamValues,
            outputs chan DataOut,
            stop chan bool,
            err chan FlowError,
            id InstanceID)
        GetParams() (inputs ParamTypes, outputs ParamTypes)
        GetName() string
    }
    
They can be run by passing the appropriate data and channels to the Run function.

1. inputs: All inputs are required to begin running a block.
2. outputs: Outputs are passed upon successful completion
3. stop: Stop is an input command to call for the immediate termination of this and all subblocks.
3. err: Errors are outputs passed upon either a critical or non-critical error, or just upon stopping to give extra info.
   Errors can be either Ok, or not Ok, to differentiate between stopping errors and warnings.
4. id: ID identifies the instance of this block call within a greater structure, like a graph. This is useful if there is more than one instance of a block in a graph, because there might still only be one copy of this block in memory.

GetParams will return maps linking parameter names to their type for both inputs and outputs, these are important for the system to know what types to share with the block automatically.
GetName returns the unique name of this block.

###Warning

* Every block name must be unique, as well as every parameter name within a given block.
* Always send command to stop all subblocks before returning an output or stopping a block itself.
* Do not wait to receive feedback from subblocks upon receiving a stop command before termination.
* Always return a FlowError with Info: StopInfo upon receiving a stop command and terminating block execution.

## Primitives

Primitives are blocks which have been written by a human user in code. Some default blocks and examples are provided in flow/blocks.

### Example:
Lets say you want a function block for Square root.

Construct it from the blocks library.

    sqrtblk = blocks.Sqrt()

Lets say this is the only function you need, no graph. You can read from a database that Sqrt accepts a float64 value named "IN" and returns a float64 value named "OUT", or you can access that data by calling:

    params_in, params_out := sqrtblk.GetParams()

params_in and params_out are Parameter structs with methods GetName() and GetType().

So run it like this.

    inputs  := map[string]interface{}{"IN": float64(2)}
    outputs := make(chan map[string]interface{})
    stop    := make(chan bool)
    err     := make(chan flow.FlowError)

    go sqrtblk.Run(inputs, outputs, stop, err)

The output and err channel you read from, the stop channel you may write to.

    out := (<-outputs)["OUT"]

Run checks the types of the inputs and the outputs to makes sure everything is implemented correctly.

And now you have the sqrt.

### Notes
I admit, this is verbose, but here's the deal. Because it's made like this, with the graph structure I will soon implement and describe created, which is ran and read from the exact same way (they both use the same interface), you can call long strings of processes. And, an AI program can create graphs intelligently by calling functions like AddNode, AddEdge, RemoveEdge, RemoveNode. I will let you know more once I have implemented that, but that is how it works.
