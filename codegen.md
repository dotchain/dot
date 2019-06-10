# Code generation in DOT

Golang does not natively support generics yet. So code generation
is required to generate the boilerplate implementation for values and
collections.

All the code generation for DOT is implemented by the
[dotc](https://godoc.org/github.com/dotchain/dotc) package.

For example, the following code generates the
[Value](https://godoc.org/github.com/dotchain/dot/changes#Value) and
[Collection](https://godoc.org/github.com/dotchain/dot/changes#Collection)
interfaces.  In addition, the corresponding **Stream** implementations
are produced as well.

For the TODO MVC app, the following code generates the required Value
types and Stream implementations.

```go global

// import fmt
// import github.com/dotchain/dot/x/dotc
// import io/ioutil

func main() {
	code, err := info.Generate()
        if err != nil {
        	panic(err)
        }
        err = ioutil.WriteFile("example/generated.go", []byte(code), 0644)
        if err != nil {
        	panic(err)
        }
}

var info = dotc.Info{
	Package: "example",
        Structs: []dotc.Struct{{
        	Recv: "t",
                Type: "Todo",
                Fields: []dotc.Field{{
                	Name: "Complete",
                        Key: "complete",
                        Type: "bool",
                }, {
                	Name: "Description",
                        Key: "desc",
                        Type: "string",
                }},
        }},
        StructStreams: []dotc.Struct{{
        	Recv: "t",
                Type: "Todo",
                Fields: []dotc.Field{{
                	Name: "Complete",
                        Key: "complete",
                        Type: "bool",
                }, {
                	Name: "Description",
                        Key: "desc",
                        Type: "string",
                }},
        }},
        Slices: []dotc.Slice{{
        	Recv: "t",
               	Type: "TodoList",
               	ElemType: "Todo",
        }},
        SliceStreams: []dotc.Slice{{
        	Recv: "t",
               	Type: "TodoList",
               	ElemType: "Todo",
        }},
}
```

As the example above shows, code generation requires explicitly
specifying the types needed.  It is possible to guess these directly
from the sources but that is not yet implemented

## Markdown to code

This markdown can be **executed** to spit out the generated code like
so:

```sh
$ go get github.com/tvastar/test/cmd/testmd
$ testmd -pkg main codegen.md
```
