package code

import (
	"github.com/dave/jennifer/jen"
	"strings"
)

// Code is the interface that all code nodes need to implement
type Code interface {
	String() string
	Code() *jen.Statement
	AddDocs(docs ...Comment)
}

type Import struct {
	FilePath string
	Alias    string
	Path     string
}
type Comment string

type TypeOptions func(t *Type)

type Type struct {
	Import *Import

	// this is for method type types. e.x  func(string) string
	Method    *Method
	Pointer   bool
	Qualifier string
}

type Var struct {
	docs  []Comment
	Name  string
	Type  Type
	Value interface{}
}

type Const struct {
	Var
}

type Parameter struct {
	docs []Comment
	Name string
	Type Type
}

type FieldTags map[string]string

type StructField struct {
	Parameter
	Tags FieldTags
}

type Struct struct {
	docs   []Comment
	Name   string
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

func MethodTypeOption(m *Method) TypeOptions {
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
	return &Const{Var: *v}
}

func NewParameter(name string, tp Type, docs ...Comment) *Parameter {
	return &Parameter{
		Name: name,
		Type: tp,
		docs: docs,
	}
}
func NewFieldTags(key, value string) FieldTags {
	return FieldTags(map[string]string{key: value})
}
func NewStructField(name string, tp Type, docs ...Comment) StructField {
	pr := &Parameter{
		Name: name,
		Type: tp,
		docs: docs,
	}
	return StructField{
		Parameter: *pr,
	}
}
func NewStructFieldWithTag(name string, tp Type, tags FieldTags, docs ...Comment) StructField {
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
func ParamsMethodOption(params []Parameter) MethodOptions {
	return func(m *Method) {
		m.Params = params
	}
}
func ResultsMethodOption(results []Parameter) MethodOptions {
	return func(m *Method) {
		m.Results = results
	}
}
func RecvMethodOption(recv *Parameter) MethodOptions {
	return func(m *Method) {
		m.Recv = recv
	}
}
func BodyMethodOption(body []jen.Code) MethodOptions {
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
func (c *Comment) Code() *jen.Statement {
	return jen.Comment(string(*c))
}
func (c *Comment) String() string {
	return codeString(c)
}
func (c *Comment) AddDocs(docs ...Comment) {
	return
}

func (t *Type) Code() *jen.Statement {
	code := jen.Empty()
	if t.Pointer {
		code.Id("*")
	}
	if t.Import != nil {
		code.Qual(t.Import.Path, t.Qualifier)
		return code
	}
	return code.Id(t.Qualifier)
}
func (t *Type) String() string {
	return codeString(t)
}
func (t *Type) AddDocs(docs ...Comment) {
	return
}

func (v *Var) Code() *jen.Statement {
	code := jen.Empty()
	addDocsCode(code, v.docs)
	code.Line().Var().Id(v.Name).Add(v.Type.Code())
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
	code := jen.Empty()
	addDocsCode(code, c.docs)
	code.Line().Const().Id(c.Name).Add(c.Type.Code())
	if c.Value != nil {
		code.Op("=").Lit(c.Value)
	}
	return code
}

func (p *Parameter) Code() *jen.Statement {
	return jen.Id(p.Name).Add(p.Type.Code())
}

func (p *Parameter) String() string {
	return codeString(p)
}

func (p *Parameter) AddDocs(docs ...Comment) {
	if docs != nil {
		p.docs = append(p.docs, docs...)
	}
}

func (m *Method) Code() *jen.Statement {
	code := jen.Empty()
	addDocsCode(code, m.docs)
	code.Func()
	if m.Recv != nil {
		code.Params(m.Recv.Code())
	}
	code.Id(m.Name)
	code.Params(paramsList(m.Params)...)
	code.Params(paramsList(m.Results)...)
	code.Block(m.Body...)
	return code
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
	code := jen.Empty()
	addDocsCode(code, s.docs)
	return code.Id(s.Name).Add(s.Type.Code()).Tag(s.Tags)
}

func (s *Struct) Code() *jen.Statement {
	code := jen.Empty()
	addDocsCode(code, s.docs)
	code.Type().Id(s.Name)
	code.Struct(fieldList(s.Fields)...)
	return code
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
	code := jen.Empty()
	addDocsCode(code, m.docs)
	code.Id(m.Name)
	code.Params(paramsList(m.Params)...)
	code.Params(paramsList(m.Results)...)
	return code
}
func (f *FieldTags) Set(key, value string) {
	if f == nil {
		*f = map[string]string{
			key: value,
		}
	}
	(*f)[key] = value
}
func codeString(c Code) string {
	lines := strings.Split(c.Code().GoString(), "\n")
	var results []string
	// fixes the unnecessary tab in the beginning
	for _, l := range lines {
		results = append(results, strings.TrimPrefix(l, "\t"))
	}
	return string(strings.Join(results, "\n"))
}

func addDocsCode(c *jen.Statement, docs []Comment) {
	for _, d := range docs {
		c.Add(d.Code())
	}
}

func paramsList(paramList []Parameter) (l []jen.Code) {
	for _, p := range paramList {
		l = append(l, p.Code())
	}
	return
}
