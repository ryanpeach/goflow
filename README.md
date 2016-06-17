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

params_in and params_out are maps linking names to types, as defined by enum constant in flow.

Now run it like this.

    inputs  := map[string]interface{}{"IN": float64(2)}
    outputs := make(chan map[string]interface{})
    stop    := make(chan bool)
    err     := make(chan flow.FlowError)

    go sqrtblk.Run(inputs, outputs, stop, err, 0)

The output and err channel you read from, the stop channel you may write to.

    out := (<-outputs)["OUT"]

Run checks the types of the inputs and the outputs to makes sure everything is implemented correctly.

And now you have the sqrt.

## Graphs

Graphs are function blocks which contain input parameters, output parameters, other function blocks (nodes), and edges connecting parameters (either it's own, or to function blocks).

### Rules
 * A block only implements once all input data has been set.
 * A block returns all output data at once, and then terminates.
 * Outputs connect to multiple inputs. Inputs never connect to outputs.
 * The Graph moves all data from outputs to linked inputs, once outputs are clear a block may run again.
 * The Graph removes all data from inputs upon running a block, the block instance may not be run again until it returns an output
 * The Graph returns immediately upon all output parameters are determined.
 * Graph output parameters reflect the most recent value of their source.

### Example:

Let's say you want to create a graph representing the Nand function.

First, create two ParamType structures, one for graph inputs and the other for outputs, containing the name of parameters, and their type selected from flow constants.

    ins  := flow.ParamTypes{"A": flow.Bool, "B": flow.Bool}
    outs := flow.ParamTypes{"OUT": flow.Bool}

Now create your graph:

    graph := flow.NewGraph("logical_nand", ins, outs)

Instantiate some function blocks to use as nodes:
    
    and, and_addr := blocks.And(0)
    not, not_addr := blocks.InvBool(0)

And add them as nodes in the graph using the Graph.AddNode(FunctionBlock, Address) method:

    graph.AddNode(and, and_addr)
    graph.AddNode(not, not_addr)

Link the graph inputs to some set of node inputs:

    graph.LinkIn("A", "A", and_addr)
    graph.LinkIn("B", "B", and_addr)
    
Link the graph outputs to some node output:
    
    graph.LinkOut(not_addr, "OUT", "OUT")
    
Connect the nodes together following the Nand pattern, the output of the And to the input of the Not:

    graph.AddEdge(and_addr, "OUT", not_addr, "IN")
    
We are done! From here, you can run the graph as a FunctionBlock!

### Notes
I admit, this is verbose, but here's the deal. Because it's made like this, with the graph structure I will soon implement and describe created, which is ran and read from the exact same way (they both use the same interface), you can call long strings of processes. And, an AI program can create graphs intelligently by calling functions like AddNode, AddEdge, RemoveEdge, RemoveNode. I will let you know more once I have implemented that, but that is how it works.

### Improvements

Graph v2.0 removed the memory maps in favor of a goroutine graph structure with edges being channels. This makes golang's internal channel waiting system handle dataflow. This improved speeds from the Nand function from 252640 ns/op to 110463 ns/op, a 250% increase!
### Roadmap
 - [x] Primitive Blocks
 - [x] Graphs
 - [x] Loops
 - [x] Switches (Implemented but Untested in Graphs)
 - [x] Custom Types
 - [ ] Webapp GUI (In progress)
 - [ ] Easy Accessors and Python Port
 - [ ] Remove String Maps in favor of Indexing (Maybe)
