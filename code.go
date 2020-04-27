package code

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

// Code is the interface that all code nodes need to implement.
type Code interface {
	// Docs returns the documentation comments
	Docs() []Comment

	// String returns the string representation of the code.
	String() string

	// Code returns the jen representation of the code.
	Code() *jen.Statement

	// AddDocs adds documentation comments to the code.
	AddDocs(docs ...Comment)
}

// Import represents an import
// this is mostly used for parameters and field types so we know what the package
// of the type is.
type Import struct {
	// FilePath is the path of the package. (e.x /User/Kujtim/go/src/github.com/go-services/code).
	FilePath string

	// Alias is the alias of the import (e.x import my_alias "fmt").
	Alias string

	// Path is the path of the import (e.x import "github.com/go-services/code").
	Path string
}

// Comment represents a code comment, it implements the Code interface.
type Comment string

// TypeOptions is used when you call NewType, it is a handy way to allow multiple configurations
// for a type.
//   tp:= NewType("string", PointerTypeOption())
// This would give you a type of pointer string.
// You can use many build in type option functions like:
//
// 	- ImportTypeOption(i *Import) TypeOptions
// 	- FunctionTypeOption(m *FunctionType) TypeOptions
// 	- PointerTypeOption() TypeOptions
type TypeOptions func(t *Type)

// FunctionType is used in Type to specify a method type (e.x func(sting) int).
type FunctionType Function

// Type represents a type e.x `string`, `context.Context`...
// the type is represented by 2 main parameters.
//
//	1) The Import  e.x `context`
//
//		Import = &Import {
//			Path: "context"
//		}
//
//	2) The Qualifier e.x `Context`
// 		Qualifier = "Context"
// this would give you the representation of `context.Context`.
type Type struct {
	// Import specifies the import of the type, it is used so we know how to call jen.Qual.
	// e.x if you want to specify the type to be `context.Context`.
	// the Import would be
	//   Import{
	//		Path:"context"
	// 	}
	// and the Qualifier = Context
	Import *Import

	// Function is used in Type to specify a method type (e.x func(sting) int)
	// if method is set all other type parameters besides the pointer are ignored.
	Function *FunctionType

	// MapType is used if the type is a map.
	MapType *struct {
		Key   Type
		Value Type
	}

	// RawType is used to specify complex types (e.x map[string]*test.SomeStruct)
	// if the raw type is not nil all the other parameters will be ignored.
	RawType *jen.Statement

	// Pointer tells if the type is a pointer.
	Pointer bool

	// ArrayType tells if the type is an array.
	ArrayType *Type

	// Variadic tells if the type is used for variadic functions
	Variadic bool

	// Qualifier specifies the qualifier, for simple types like `string` it is the only
	// parameter set on the type.
	Qualifier string
}

// Var represents a variable.
type Var struct {
	// Name is the name of the variable (e.x var {name} string)
	Name string

	// Type is the type of the variable.
	Type Type

	// The value of the variable (e.x 2, 2.9, "Some string").
	// currently only supports literal values, the plan is to expend this.
	Value interface{}

	// docs stores all the documentation comments of the variable
	docs []Comment
}

// Const represents a constant, it has the same attributes as the variable only that
// it assumes that the value is always set (there can not be any constant without initial value).
type Const Var

// Parameter represents a parameter, it is used in function parameters, function return parameters.
type Parameter struct {
	Name string
	Type Type
}

// FieldTags represent the structure field tags (e.x `json:"test"`).
type FieldTags map[string]string

// StructField represent a structure field, it uses the parameter representation to represent the name and the type
// the difference between parameter and struct field is that struct fields can have docs and tags.
type StructField struct {
	// Parameter is used for representing the name and type because the go code is the same.
	Parameter

	// Tags represent the field tags.
	Tags *FieldTags

	// docs are the documentation comments of the struct field.
	docs []Comment
}

// Struct represent a structure.
type Struct struct {
	// docs are the documentation comments of the struct field.
	docs []Comment

	// Name represents the name of the structure.
	Name string

	// Fields represents the structure fields.
	Fields []StructField
}

// FunctionOptions is used when you call NewFunction, it is a handy way to allow multiple configurations
// for a function.
//
// Code has some preconfigured option functions you can use to add parameter, results, receiver, body, docs to the function.
type FunctionOptions func(m *Function)

// Function represents a function.
type Function struct {
	// Name is the name of the function.
	Name string

	// Recv is the receiver of the function (e.x func (rcv *MyStruct) name() {}).
	Recv *Parameter

	// Params are the functions parameters.
	Params []Parameter

	// Results are the functions results.
	Results []Parameter

	// Body is the function body.
	Body []jen.Code

	// docs are the function documentation comments.
	docs []Comment
}

// InterfaceMethod is the representation of interface methods (e.x String() string).
type InterfaceMethod Function

// Interface is the representation of an interface.
type Interface struct {
	// Name is the name of the interface.
	Name string

	// Methods are the interface methods, the interface can also have no methods.
	Methods []InterfaceMethod

	// docs are the interface documentation coments.
	docs []Comment
}

// RawCode represents raw lines of code.
type RawCode struct {
	code *jen.Statement
}

// NewImport creates a new Import with the given alias and path.
//
// alias is the alias you want to give a package (e.x import my_alias "fmt")
// path  is the import path (e.x "github.com/go-services/code")
func NewImport(alias, path string) *Import {
	return &Import{
		Alias: alias,
		Path:  path,
	}
}

// NewImportWithFilePath creates a new Import with the given alias and path and filepath.
//
// filePath is the path in the filesystem of the source code.
func NewImportWithFilePath(alias, path, filePath string) *Import {
	imp := NewImport(alias, path)
	imp.FilePath = filePath
	return imp
}

// NewComment creates a Comment with the given string.
func NewComment(s string) Comment {
	return Comment(s)
}

// ImportTypeOption adds the given import to the type.
func ImportTypeOption(i Import) TypeOptions {
	return func(t *Type) {
		t.Import = &i
	}
}

// FunctionTypeOption adds the given function type to the type.
func FunctionTypeOption(m *FunctionType) TypeOptions {
	return func(t *Type) {
		t.Function = m
	}
}

// PointerTypeOption marks the type as a pointer type.
func PointerTypeOption() TypeOptions {
	return func(t *Type) {
		t.Pointer = true
	}
}

// VariadicTypeOption marks the type as variadic.
func VariadicTypeOption() TypeOptions {
	return func(t *Type) {
		t.Variadic = true
	}
}

// ArrayTypeOption marks the type as variadic.
func ArrayTypeOption(tp Type) TypeOptions {
	return func(t *Type) {
		t.ArrayType = &tp
	}
}

// MapTypeOption sets the map type.
func MapTypeOption(key Type, value Type) TypeOptions {
	return func(t *Type) {
		t.MapType = &struct {
			Key   Type
			Value Type
		}{
			Key:   key,
			Value: value,
		}
	}
}

// NewType creates the type with the qualifier and options given.
//
// Options are used so we can create a simple type like `string` and complex types
//
// Code has some build in type options that can be used
//	- ImportTypeOption(i *Import) TypeOptions
//	- FunctionTypeOption(m *FunctionType) TypeOptions
//	- PointerTypeOption() TypeOptions
func NewType(qualifier string, options ...TypeOptions) Type {
	tp := Type{
		Qualifier: qualifier,
	}
	for _, o := range options {
		o(&tp)
	}
	return tp
}

// NewRawType creates a new type with raw jen statement.
func NewRawType(tp *jen.Statement) Type {
	return Type{
		RawType: tp,
	}
}

// NewVar creates a new var with the given name and type,
// there is also an optional list of documentation comments that you can add to the variable
func NewVar(name string, tp Type, docs ...Comment) *Var {
	return &Var{
		Name: name,
		Type: tp,
		docs: docs,
	}
}

// NewVarWithValue creates a new var with the given name type and value.
func NewVarWithValue(name string, tp Type, value interface{}, docs ...Comment) *Var {
	v := NewVar(name, tp, docs...)
	v.Value = value
	return v
}

// NewConst creates a 7new constant with the given name type and value.
func NewConst(name string, tp Type, value interface{}, docs ...Comment) *Const {
	v := NewVar(name, tp, docs...)
	v.Value = value
	c := Const(*v)
	return &c
}

// NewParameter creates a new parameter with the given name and type.
func NewParameter(name string, tp Type) *Parameter {
	return &Parameter{
		Name: name,
		Type: tp,
	}
}

// NewFieldTags creates a new field tags map with initial key and value.
// there is also an optional list of documentation comments that you can add to the variable
func NewFieldTags(key, value string) *FieldTags {
	return &FieldTags{
		key: value,
	}
}

// NewStructField creates a new structure field with the given name and type,
// there is also an optional list of documentation comments that you can add to the variable
func NewStructField(name string, tp Type, docs ...Comment) *StructField {
	pr := &Parameter{
		Name: name,
		Type: tp,
	}
	return &StructField{
		docs:      docs,
		Parameter: *pr,
	}
}

// NewStructFieldWithTag creates a new structure field with the given name type and tags,
// there is also an optional list of documentation comments that you can add to the variable
func NewStructFieldWithTag(name string, tp Type, tags *FieldTags, docs ...Comment) *StructField {
	sf := NewStructField(name, tp, docs...)
	sf.Tags = tags
	return sf
}

// NewStruct creates a new structure with the given name,
// there is also an optional list of documentation comments that you can add to the variable
func NewStruct(name string, docs ...Comment) *Struct {
	return &Struct{
		Name: name,
		docs: docs,
	}
}

// NewStructWithFields creates a new structure with the given name and fields,
// there is also an optional list of documentation comments that you can add to the variable
func NewStructWithFields(name string, fields []StructField, docs ...Comment) *Struct {
	st := NewStruct(name, docs...)
	st.Fields = fields
	return st
}

// ParamsFunctionOption adds given parameters to the function.
func ParamsFunctionOption(params ...Parameter) FunctionOptions {
	return func(f *Function) {
		f.Params = params
	}
}

// ResultsFunctionOption adds given results to the function.
func ResultsFunctionOption(results ...Parameter) FunctionOptions {
	return func(f *Function) {
		f.Results = results
	}
}

// RecvFunctionOption sets the function receiver.
func RecvFunctionOption(recv *Parameter) FunctionOptions {
	return func(f *Function) {
		f.Recv = recv
	}
}

// BodyFunctionOption adds given body code to the function.
func BodyFunctionOption(body ...jen.Code) FunctionOptions {
	return func(f *Function) {
		f.Body = body
	}
}

// DocsFunctionOption adds given docs to the function.
func DocsFunctionOption(docs ...Comment) FunctionOptions {
	return func(f *Function) {
		f.docs = docs
	}
}

// NewFunction creates a new function with the given name and options.
func NewFunction(name string, options ...FunctionOptions) *Function {
	f := &Function{
		Name: name,
	}
	for _, o := range options {
		o(f)
	}
	return f
}

// NewFunctionType creates a new function type with the given options.
func NewFunctionType(options ...FunctionOptions) *FunctionType {
	m := &Function{}
	for _, o := range options {
		o(m)
	}
	mt := FunctionType(*m)
	return &mt
}

// NewInterfaceMethod creates a new interface method with the given name and options.
func NewInterfaceMethod(name string, options ...FunctionOptions) InterfaceMethod {
	m := NewFunction(name, options...)
	return InterfaceMethod(*m)
}

// NewInterface creates a new interface with the given name and methods,
// there is also an optional list of documentation comments that you can add to the variable
func NewInterface(name string, methods []InterfaceMethod, docs ...Comment) *Interface {
	return &Interface{
		Name:    name,
		Methods: methods,
		docs:    docs,
	}
}

func NewRawCode(code *jen.Statement) *RawCode {
	return &RawCode{
		code: code,
	}
}

// Code returns the jen representation of code.
func (c *RawCode) Code() *jen.Statement {
	return c.code
}

// Docs does nothing for the code.
func (c *RawCode) Docs() []Comment {
	return nil
}

// AddDocs does nothing for the code.
// We only implement this so we implement the Code interface.
func (c *RawCode) AddDocs(_ ...Comment) {}

// String returns the go code string of the comment.
func (c *RawCode) String() string {
	return codeString(c)
}

// Code returns the jen representation of the comment.
func (c Comment) Code() *jen.Statement {
	return jen.Comment(string(c))
}

// Docs does nothing for the comment code.
func (c Comment) Docs() []Comment {
	return nil
}

// String returns the go code string of the comment.
func (c Comment) String() string {
	return codeString(&c)
}

// AddDocs does nothing for the comment code.
// We only implement this so we implement the Code interface.
func (c Comment) AddDocs(_ ...Comment) {}

// Code returns the jen representation of the the type.
func (t Type) Code() *jen.Statement {
	code := &jen.Statement{}
	if t.RawType != nil {
		return t.RawType
	}
	if t.Variadic {
		code.Op("...")
	}
	if t.Pointer {
		code.Id("*")
	}
	if t.ArrayType != nil {
		code.Index().Add(t.ArrayType.Code())
		return code
	}
	if t.MapType != nil {
		code.Map(t.MapType.Key.Code()).Add(t.MapType.Value.Code())
		return code
	}
	if t.Function != nil {
		code.Add(t.Function.Code())
		return code
	}
	if t.Import != nil {
		if t.Import.Alias != "" {
			code.Id(t.Import.Alias).Dot(t.Qualifier)
			return code
		}
		code.Qual(t.Import.Path, t.Qualifier)
		return code
	}
	return code.Id(t.Qualifier)
}

// String returns the go code string of the type,
// if the type is a function type the function string tis used.
func (t Type) String() string {
	if t.RawType != nil {
		// Hack to get the reader to not panic for complex types
		v := NewVar("_", t)
		s := v.Code().GoString()
		s = strings.TrimPrefix(s, "var _ ")
		// -----
		return s
	}
	if t.Function != nil {
		s := t.Function.String()
		if t.Variadic {
			s = "..." + s
		}
		if t.Pointer {
			s = "*" + s
		}
		return s
	}
	if t.ArrayType != nil || t.Variadic || t.MapType != nil {
		code := jen.Func().Id("_").Params(t.Code()).Block()
		s := code.GoString()
		s = strings.TrimPrefix(s, "func _(")
		s = strings.TrimSuffix(s, ") {}")
		return s
	}
	return codeString(t)
}

// Docs does nothing for the Type code.
func (t Type) Docs() []Comment {
	return nil
}

// AddDocs does nothing for the type code.
// We only implement this so we implement the Code interface.
func (t Type) AddDocs(_ ...Comment) {}

// Code returns the jen representation of the variable.
func (v *Var) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, v.docs)
	code.Var().Id(v.Name).Add(v.Type.Code())
	if v.Value != nil {
		code.Op("=").Lit(v.Value)
	}
	return code
}

// String returns the go code string of the variable.
func (v *Var) String() string {
	return codeString(v)
}

// Docs returns the docs comments of the variable.
func (v *Var) Docs() []Comment {
	return v.docs
}

// AddDocs adds a list of documentation strings to the variable.
func (v *Var) AddDocs(docs ...Comment) {
	if docs != nil {
		v.docs = append(v.docs, docs...)
	}
}

// Code returns the jen representation of the constant.
func (c *Const) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, c.docs)
	code.Const().Id(c.Name).Add(c.Type.Code())
	if c.Value != nil {
		code.Op("=").Lit(c.Value)
	}
	return code
}

// String returns the go code string of the constant.
func (c *Const) String() string {
	return codeString(c)
}

// Docs returns the docs comments of the constant.
func (c *Const) Docs() []Comment {
	return c.docs
}

// AddDocs adds a list of documentation strings to the constant.
func (c *Const) AddDocs(docs ...Comment) {
	c.docs = append(c.docs, docs...)
}

// Code returns the jen representation of the parameter.
func (p *Parameter) Code() *jen.Statement {
	c := &jen.Statement{}
	if p.Name != "" {
		c.Id(p.Name)
	}
	return c.Add(p.Type.Code())
}

// String returns the go code string of the parameter.
// because the renderer does not render only parameters we create a dummy function to add the parameters to
// than we remove everything besides the parameters.
func (p *Parameter) String() string {
	// Hack to get the reader to not throw errors in creating string representative of parameters
	code := jen.Func().Id("_").Params(p.Code()).Block()
	s := code.GoString()
	s = strings.TrimPrefix(s, "func _(")
	s = strings.TrimSuffix(s, ") {}")
	// -----
	return s
}

// Docs does nothing for the parameter code.
func (p *Parameter) Docs() []Comment {
	return nil
}

// AddDocs does nothing for the parameter code.
// We only implement this so we implement the Code interface.
func (p *Parameter) AddDocs(_ ...Comment) {}

// Code returns the jen representation of the function type.
func (m *FunctionType) Code() *jen.Statement {
	code := &jen.Statement{}
	code.Func()
	code.Params(paramsList(m.Params)...)
	if m.Results != nil && len(m.Results) > 0 {
		code.Params(paramsList(m.Results)...)
	}
	return code
}

// String returns the go code string of the function type.
// because the renderer does not render only function types we create a dummy structure to add the function type field to
// than we remove everything besides the function type.
func (m *FunctionType) String() string {
	code := m.Code()
	// Hack to get the reader to not panic in function types
	fakeStruct := jen.Type().Id("_").Struct(jen.Id("_").Add(code))
	s := fakeStruct.GoString()
	s = strings.TrimPrefix(s, "type _ struct {\n\t_ ")
	s = strings.TrimSuffix(s, "\n}")
	// -----
	return s
}

// Docs does nothing for the function type code.
func (m *FunctionType) Docs() []Comment {
	return nil
}

// AddDocs does nothing for the function type code.
// We only implement this so we implement the Code interface.
func (m *FunctionType) AddDocs(_ ...Comment) {}

// Code returns the jen representation of the function.
func (f *Function) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, f.docs)
	code.Func()
	if f.Recv != nil {
		code.Params(f.Recv.Code())
	}
	code.Id(f.Name)
	code.Params(paramsList(f.Params)...)
	if f.Results != nil && len(f.Results) > 0 {
		code.Params(paramsList(f.Results)...)
	}
	return code.Block(f.Body...)
}

// Docs returns the docs comments of the function.
func (f *Function) Docs() []Comment {
	return f.docs
}

// AddDocs adds a list of documentation strings to the function.
func (f *Function) AddDocs(docs ...Comment) {
	f.docs = append(f.docs, docs...)
}

// String returns the go code string of the function.
func (f *Function) String() string {
	return codeString(f)
}

// AddStringBody adds raw string code to the body of the function.
func (f *Function) AddStringBody(s string) {
	f.Body = append(f.Body, jen.Id(s))
}

// AddParameter adds a new parameter to the function.
func (f *Function) AddParameter(p Parameter) {
	f.Params = append(f.Params, p)
}

// AddResult adds a new result to the function.
func (f *Function) AddResult(p Parameter) {
	f.Results = append(f.Params, p)
}

// Code returns the jen representation of the structure field.
func (s *StructField) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, s.docs)
	code.Id(s.Name).Add(s.Type.Code())
	if s.Tags != nil {
		code.Tag(*s.Tags)
	}
	return code
}

// String returns the go code string of the struct field.
// because the renderer does not render only struct fields we create a dummy structure to add the struct field to
// than we remove everything besides the struct field.
func (s *StructField) String() string {
	code := s.Code()
	// Hack to get the reader to not panic in struct fields
	fakeStruct := jen.Type().Id("_").Struct(code)
	str := fakeStruct.GoString()
	str = strings.TrimPrefix(str, "type _ struct {\n\t")
	str = strings.TrimSuffix(str, "\n}")
	// -----
	return prepareLines(str)
}

// Docs returns the docs comments of the structure field.
func (s *StructField) Docs() []Comment {
	return s.docs
}

// Code returns the jen representation of the structure.
func (s *Struct) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, s.docs)
	return code.Type().Id(s.Name).Struct(fieldList(s.Fields)...)
}

// String returns the go code string of the structure.
func (s *Struct) String() string {
	return codeString(s)
}

// Docs returns the docs comments of the structure.
func (s *Struct) Docs() []Comment {
	return s.docs
}

// AddDocs adds a list of documentation strings to the structure.
func (s *Struct) AddDocs(docs ...Comment) {
	s.docs = append(s.docs, docs...)
}

// Code returns the jen representation of the interface method.
func (m *InterfaceMethod) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, m.docs)
	code.Id(m.Name)
	code.Params(paramsList(m.Params)...)
	if m.Results != nil && len(m.Results) > 0 {
		code.Params(paramsList(m.Results)...)

	}
	return code
}

// String returns the go code string of the interface method.
// because the renderer does not render only interface methods we create a dummy interface to add the interface method to
// than we remove everything besides the interface method.
func (m *InterfaceMethod) String() string {
	code := m.Code()
	// Hack to get the reader to not panic in interface methods.
	fakeInterface := jen.Type().Id("_").Interface(code)
	str := fakeInterface.GoString()
	str = strings.TrimPrefix(str, "type _ interface {\n\t")
	str = strings.TrimSuffix(str, "\n}")
	// -----
	return prepareLines(str)
}

// Docs returns the docs comments of the interface method.
func (m *InterfaceMethod) Docs() []Comment {
	return m.docs
}

// AddDocs adds a list of documentation strings to the interface method.
func (m *InterfaceMethod) AddDocs(docs ...Comment) {
	m.docs = append(m.docs, docs...)
}

// Code returns the jen representation of the interface.
func (i *Interface) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, i.docs)
	code.Type().Id(i.Name).Interface(
		func() []jen.Code {
			var c []jen.Code
			for _, m := range i.Methods {
				c = append(c, m.Code())
			}
			return c
		}()...,
	)
	return code
}

// String returns the go code string of the interface.
func (i *Interface) String() string {
	return codeString(i)
}

// Docs returns the docs comments of the interface.
func (i *Interface) Docs() []Comment {
	return i.docs
}

// AddDocs adds a list of documentation strings to the interface.
func (i *Interface) AddDocs(docs ...Comment) {
	i.docs = append(i.docs, docs...)
}

// AddMethod a method to the method list
func (i *Interface) AddMethod(m InterfaceMethod) {
	i.Methods = append(i.Methods, m)
}

// Set is used to set an existing or new field tag.
func (f *FieldTags) Set(key, value string) {
	if *f == nil {
		*f = map[string]string{
			key: value,
		}
	}
	(*f)[key] = value
}

func fieldList(fields []StructField) (f []jen.Code) {
	for _, p := range fields {
		f = append(f, p.Code())
	}
	return
}

func codeString(c Code) string {
	return c.Code().GoString()
}

func addDocsCode(c *jen.Statement, docs []Comment) {
	for _, d := range docs {
		c.Add(d.Code().Line())
	}
}

func prepareLines(s string) string {
	lines := strings.Split(s, "\n")
	var results []string
	// fixes the unnecessary tab in the beginning
	for _, l := range lines {
		results = append(results, strings.TrimPrefix(l, "\t"))
	}
	return strings.Join(results, "\n")
}

func paramsList(paramList []Parameter) (l []jen.Code) {
	for _, p := range paramList {
		l = append(l, p.Code())
	}
	return
}
