// Package code is a small wrapper for https://github.com/dave/jennifer this package is used to easily generate go code
// it also has a very straight forward api.
//
// e.x to create a structure:
//    structure := code.NewStruct("MyStruct")
//
//    fmt.Println(structure) // 	type MyStruct struct{}
//
// to add fields to the structure
//    structure.Fields = []code.StructField{
//        code.NewStructField("Name", code.NewType("string")),
//    }
//    fmt.Println(structure)  // 	type MyStruct struct {
//                            //  		Name string
//                            //  	}
package code
