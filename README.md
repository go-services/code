# Code
Code is a small wrapper around [jen](https://github.com/dave/jennifer) that allows `go-services` to generate gocode easier.

It has a very friendly api 
```go
package main

import (
    "github.com/go-services/code"
    "fmt"
)

func main() {
    structure := code.NewStruct("MyStruct")
    
    fmt.Println(structure) // 	type MyStruct struct{}

    structure.Fields = []code.StructField{
        code.NewStructField("Name", code.NewType("string")),
    }
    fmt.Println(structure)  // 	type MyStruct struct {
                            //  	Name string
                            //  }

}
```