package code

import (
	"github.com/dave/jennifer/jen"
	"strings"
)

// Code is the interface that all code nodes need to implement.
type Code interface {

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
// 	- MethodTypeOption(m *MethodType) TypeOptions
// 	- PointerTypeOption() TypeOptions
type TypeOptions func(t *Type)

// MethodType is used in Type to specify a method type (e.x func(sting) int).
type MethodType Method

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

	// MethodType is used in Type to specify a method type (e.x func(sting) int)
	// if method is set all other type parameters besides the pointer are ignored.
	Method *MethodType

	// Pointer tells if the type is a pointer.
	Pointer bool

	// Qualifier specifies the qualifier, for simple types like `string` it is the only
	// parameter set on the type.
	Qualifier string
}

// Var represents a variable.
type Var struct {
	// docs stores all the documentation comments of the variable
	docs []Comment

	// Name is the name of the variable (e.x var {name} string)
	Name string

	// Type is the type of the variable.
	Type Type

	// The value of the variable (e.x 2, 2.9, "Some string").
	// currently only supports literal values, the plan is to expend this.
	Value interface{}
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

	// docs are the documentation comments of the struct field.
	docs []Comment

	// Tags represent the field tags.
	Tags *FieldTags
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

type MethodOptions func(m *Method)

type Method struct {
	Name string

	docs []Comment

	Recv *Parameter

	Params []Parameter

	Results []Parameter

	Body []jen.Code
}

type InterfaceMethod Method

type Interface struct {
	docs    []Comment
	Name    string
	Methods []InterfaceMethod
}

func NewImport(alias, path string) *Import {
	return &Import{
		Alias: alias,
		Path:  path,
	}
}

func NewImportWithFilePath(alias, path, filePath string) *Import {
	imp := NewImport(alias, path)
	imp.FilePath = filePath
	return imp
}

func NewComment(s string) Comment {
	return Comment(s)
}

func ImportTypeOption(i *Import) TypeOptions {
	return func(t *Type) {
		t.Import = i
	}
}

func MethodTypeOption(m *MethodType) TypeOptions {
	return func(t *Type) {
		t.Method = m
	}
}

func PointerTypeOption() TypeOptions {
	return func(t *Type) {
		t.Pointer = true
	}
}

func NewType(qualifier string, options ...TypeOptions) Type {
	tp := Type{
		Qualifier: qualifier,
	}
	for _, o := range options {
		o(&tp)
	}
	return tp
}

func NewVar(name string, tp Type, docs ...Comment) *Var {
	return &Var{
		Name: name,
		Type: tp,
		docs: docs,
	}
}

func NewVarWithValue(name string, tp Type, value interface{}, docs ...Comment) *Var {
	v := NewVar(name, tp, docs...)
	v.Value = value
	return v
}

func NewConst(name string, tp Type, value interface{}, docs ...Comment) *Const {
	v := NewVar(name, tp, docs...)
	v.Value = value
	c := Const(*v)
	return &c
}

func NewParameter(name string, tp Type) *Parameter {
	return &Parameter{
		Name: name,
		Type: tp,
	}
}
func NewFieldTags(key, value string) *FieldTags {
	return &FieldTags{
		key: value,
	}
}
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
func NewStructFieldWithTag(name string, tp Type, tags *FieldTags, docs ...Comment) *StructField {
	sf := NewStructField(name, tp, docs...)
	sf.Tags = tags
	return sf
}

func NewStruct(name string, docs ...Comment) *Struct {
	return &Struct{
		Name: name,
		docs: docs,
	}
}

func NewStructWithFields(name string, fields []StructField, docs ...Comment) *Struct {
	st := NewStruct(name, docs...)
	st.Fields = fields
	return st
}
func ParamsMethodOption(params ...Parameter) MethodOptions {
	return func(m *Method) {
		m.Params = params
	}
}
func ResultsMethodOption(results ...Parameter) MethodOptions {
	return func(m *Method) {
		m.Results = results
	}
}
func RecvMethodOption(recv *Parameter) MethodOptions {
	return func(m *Method) {
		m.Recv = recv
	}
}
func BodyMethodOption(body ...jen.Code) MethodOptions {
	return func(m *Method) {
		m.Body = body
	}
}
func DocsMethodOption(docs ...Comment) MethodOptions {
	return func(m *Method) {
		m.docs = docs
	}
}
func NewMethod(name string, options ...MethodOptions) *Method {
	m := &Method{
		Name: name,
	}
	for _, o := range options {
		o(m)
	}
	return m
}

func NewMethodType(options ...MethodOptions) *MethodType {
	m := &Method{}
	for _, o := range options {
		o(m)
	}
	mt := MethodType(*m)
	return &mt
}
func NewInterfaceMethod(name string, options ...MethodOptions) InterfaceMethod {
	m := NewMethod(name, options...)
	return InterfaceMethod(*m)
}

func NewInterface(name string, methods []InterfaceMethod, docs ...Comment) *Interface {
	return &Interface{
		Name:    name,
		Methods: methods,
		docs:    docs,
	}
}
func (c Comment) Code() *jen.Statement {
	return jen.Comment(string(c))
}
func (c Comment) String() string {
	return codeString(&c)
}
func (c Comment) AddDocs(docs ...Comment) {
	return
}

func (t Type) Code() *jen.Statement {
	code := &jen.Statement{}
	if t.Pointer {
		code.Id("*")
	}
	if t.Method != nil {
		code.Add(t.Method.Code())
		return code
	}
	if t.Import != nil {
		code.Qual(t.Import.Path, t.Qualifier)
		return code
	}
	return code.Id(t.Qualifier)
}
func (t Type) String() string {
	if t.Method != nil {
		s := t.Method.String()
		if t.Pointer {
			s = "*" + s
		}
		return s
	}
	return codeString(t)
}
func (t Type) AddDocs(docs ...Comment) {
	return
}

func (v *Var) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, v.docs)
	code.Var().Id(v.Name).Add(v.Type.Code())
	if v.Value != nil {
		code.Op("=").Lit(v.Value)
	}
	return code
}

func (v *Var) String() string {
	return codeString(v)
}

func (v *Var) AddDocs(docs ...Comment) {
	if docs != nil {
		v.docs = append(v.docs, docs...)
	}
}

func (c *Const) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, c.docs)
	code.Const().Id(c.Name).Add(c.Type.Code())
	if c.Value != nil {
		code.Op("=").Lit(c.Value)
	}
	return code
}
func (c *Const) String() string {
	return codeString(c)
}
func (c *Const) AddDocs(docs ...Comment) {
	c.docs = append(c.docs, docs...)
}

func (p *Parameter) Code() *jen.Statement {
	c := &jen.Statement{}
	if p.Name != "" {
		c.Id(p.Name)
	}
	return c.Add(p.Type.Code())
}

func (p *Parameter) String() string {
	// Hack to get the reader to not throw errors in creating string representative of parameters
	code := jen.Func().Id("_").Params(p.Code()).Block()
	s := code.GoString()
	s = strings.TrimPrefix(s, "func _(")
	s = strings.TrimSuffix(s, ") {}")
	// -----
	return s
}

func (p *Parameter) AddDocs(docs ...Comment) {
	return
}

func (m *MethodType) Code() *jen.Statement {
	code := &jen.Statement{}
	code.Func()
	code.Params(paramsList(m.Params)...)
	if m.Results != nil && len(m.Results) > 0 {
		code.Params(paramsList(m.Results)...)
	}
	return code
}
func (m *MethodType) String() string {
	code := m.Code()
	// Hack to get the reader to not throw errors in function types
	fakeStruct := jen.Type().Id("_").Struct(jen.Id("_").Add(code))
	s := fakeStruct.GoString()
	s = strings.TrimPrefix(s, "type _ struct {\n\t_ ")
	s = strings.TrimSuffix(s, "\n}")
	// -----
	return s
}
func (m *MethodType) AddDocs(docs ...Comment) {
	return
}

func (m *Method) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, m.docs)
	code.Func()
	if m.Recv != nil {
		code.Params(m.Recv.Code())
	}
	code.Id(m.Name)
	code.Params(paramsList(m.Params)...)
	if m.Results != nil && len(m.Results) > 0 {
		code.Params(paramsList(m.Results)...)
	}
	return code.Block(m.Body...)
}

func (m *Method) AddDocs(docs ...Comment) {
	m.docs = append(m.docs, docs...)
}

func (m *Method) String() string {
	return codeString(m)
}

func (m *Method) AddStringBody(s string) {
	m.Body = append(m.Body, jen.Id(s))
}

func (s *StructField) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, s.docs)
	code.Id(s.Name).Add(s.Type.Code())
	if s.Tags != nil {
		code.Tag(*s.Tags)
	}
	return code
}
func (s *StructField) String() string {
	code := s.Code()
	// Hack to get the reader to not throw errors in struct fields
	fakeStruct := jen.Type().Id("_").Struct(code)
	str := fakeStruct.GoString()
	str = strings.TrimPrefix(str, "type _ struct {\n\t")
	str = strings.TrimSuffix(str, "\n}")
	// -----
	return prepareLines(str)
}

func (s *Struct) Code() *jen.Statement {
	code := &jen.Statement{}
	addDocsCode(code, s.docs)
	return code.Type().Id(s.Name).Struct(fieldList(s.Fields)...)
}
func fieldList(fields []StructField) (f []jen.Code) {
	for _, p := range fields {
		f = append(f, p.Code())
	}
	return
}

func (s *Struct) String() string {
	return codeString(s)
}

func (s *Struct) AddDocs(docs ...Comment) {
	s.docs = append(s.docs, docs...)
}

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

func (m *InterfaceMethod) String() string {
	code := m.Code()
	// Hack to get the reader to not throw errors in struct fields
	fakeStruct := jen.Type().Id("_").Interface(code)
	str := fakeStruct.GoString()
	str = strings.TrimPrefix(str, "type _ interface {\n\t")
	str = strings.TrimSuffix(str, "\n}")
	// -----
	return prepareLines(str)
}

func (m *InterfaceMethod) AddDocs(docs ...Comment) {
	m.docs = append(m.docs, docs...)
}
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

func (i *Interface) String() string {
	return codeString(i)
}

func (i *Interface) AddDocs(docs ...Comment) {
	i.docs = append(i.docs, docs...)
}

func (f *FieldTags) Set(key, value string) {
	if *f == nil {
		*f = map[string]string{
			key: value,
		}
	}
	(*f)[key] = value
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
	return string(strings.Join(results, "\n"))
}

func paramsList(paramList []Parameter) (l []jen.Code) {
	for _, p := range paramList {
		l = append(l, p.Code())
	}
	return
}
