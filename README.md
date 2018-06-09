# Code [![Go Report Card](https://goreportcard.com/badge/github.com/go-services/code)](https://goreportcard.com/report/github.com/go-services/code) [![Coverage Status](https://coveralls.io/repos/github/go-services/code/badge.svg?branch=master)](https://coveralls.io/github/go-services/code?branch=master) [![Build Status](https://travis-ci.org/go-services/code.svg?branch=master)](https://travis-ci.org/go-services/code)
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