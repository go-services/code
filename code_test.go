package code

import (
	"reflect"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestNewImport(t *testing.T) {
	type args struct {
		alias string
		path  string
	}
	tests := []struct {
		name string
		args args
		want *Import
	}{
		{
			name: "Should create a new import",
			args: args{
				alias: "code",
				path:  "github.com/go-services/code",
			},
			want: &Import{
				Alias: "code",
				Path:  "github.com/go-services/code",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImport(tt.args.alias, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewImportWithFilePath(t *testing.T) {
	type args struct {
		alias    string
		path     string
		filePath string
	}
	tests := []struct {
		name string
		args args
		want *Import
	}{
		{
			name: "Should create a new import",
			args: args{
				alias:    "code",
				path:     "github.com/go-services/code",
				filePath: "path/to/go/root/src/github.com/go-services/code",
			},
			want: &Import{
				Alias:    "code",
				Path:     "github.com/go-services/code",
				FilePath: "path/to/go/root/src/github.com/go-services/code",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImportWithFilePath(tt.args.alias, tt.args.path, tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImportWithFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewComment(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Comment
	}{
		{
			name: "Should return a new comment",
			args: args{
				s: "My comment here",
			},
			want: Comment("My comment here"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewComment(tt.args.s); got != tt.want {
				t.Errorf("NewComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImportTypeOption(t *testing.T) {
	type args struct {
		i  *Import
		tp Type
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should add the import to the type",
			args: args{
				i:  NewImport("test", "test"),
				tp: NewType("Test"),
			},
			want: Type{
				Qualifier: "Test",
				Import:    NewImport("test", "test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ImportTypeOption(tt.args.i)
			got(&tt.args.tp)
			if !reflect.DeepEqual(tt.args.tp, tt.want) {
				t.Errorf("Type = %v, want %v", tt.args.tp, tt.want)
			}
		})
	}
}

func TestMethodTypeOption(t *testing.T) {
	type args struct {
		tp Type
		m  *MethodType
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should add the method to the type",
			args: args{
				tp: NewType("Test"),
				m:  NewMethodType(),
			},
			want: Type{
				Qualifier: "Test",
				Method:    NewMethodType(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MethodTypeOption(tt.args.m)
			got(&tt.args.tp)
			if !reflect.DeepEqual(tt.args.tp, tt.want) {
				t.Errorf("Type = %v, want %v", tt.args.tp, tt.want)
			}
		})
	}
}

func TestPointerTypeOption(t *testing.T) {
	type args struct {
		tp Type
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should set pointer to true",
			args: args{
				tp: NewType("test"),
			},
			want: Type{
				Qualifier: "test",
				Pointer:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PointerTypeOption()
			got(&tt.args.tp)
			if !reflect.DeepEqual(tt.args.tp, tt.want) {
				t.Errorf("Type = %v, want %v", tt.args.tp, tt.want)
			}
		})
	}
}

func TestNewType(t *testing.T) {
	type args struct {
		qualifier string
		options   []TypeOptions
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should return a new type",
			args: args{
				qualifier: "test",
			},
			want: Type{
				Qualifier: "test",
			},
		},
		{
			name: "Should return a new type with options set",
			args: args{
				qualifier: "test",
				options: []TypeOptions{
					PointerTypeOption(),
					ImportTypeOption(NewImport("test", "test")),
				},
			},
			want: Type{
				Qualifier: "test",
				Pointer:   true,
				Import:    NewImport("test", "test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewType(tt.args.qualifier, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVar(t *testing.T) {
	type args struct {
		name string
		tp   Type
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want *Var
	}{
		{
			name: "Should return a new variable",
			args: args{
				name: "test",
				tp:   NewType("Qual"),
			},
			want: &Var{
				Name: "test",
				Type: NewType("Qual"),
			},
		},
		{
			name: "Should return a new variable with docs",
			args: args{
				name: "test",
				tp:   NewType("Qual"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Var{
				Name: "test",
				Type: NewType("Qual"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVar(tt.args.name, tt.args.tp, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVarWithValue(t *testing.T) {
	type args struct {
		name  string
		tp    Type
		value interface{}
		docs  []Comment
	}
	tests := []struct {
		name string
		args args
		want *Var
	}{
		{
			name: "Should return a new variable with value",
			args: args{
				name:  "test",
				tp:    NewType("Qual"),
				value: 2,
			},
			want: &Var{
				Name:  "test",
				Type:  NewType("Qual"),
				Value: 2,
			},
		},
		{
			name: "Should return a new variable with value and docs",
			args: args{
				name:  "test",
				tp:    NewType("Qual"),
				value: 2,
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Var{
				Name:  "test",
				Type:  NewType("Qual"),
				Value: 2,
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVarWithValue(tt.args.name, tt.args.tp, tt.args.value, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVarWithValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConst(t *testing.T) {
	type args struct {
		name  string
		tp    Type
		value interface{}
		docs  []Comment
	}
	tests := []struct {
		name string
		args args
		want *Const
	}{
		{
			name: "Should return a new constant",
			args: args{
				name:  "test",
				tp:    NewType("Qual"),
				value: 2,
			},
			want: &Const{
				Name:  "test",
				Type:  NewType("Qual"),
				Value: 2,
			},
		},
		{
			name: "Should return a new constant with docs",
			args: args{
				name:  "test",
				tp:    NewType("Qual"),
				value: 2,
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Const{
				Name:  "test",
				Type:  NewType("Qual"),
				Value: 2,
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConst(tt.args.name, tt.args.tp, tt.args.value, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewParameter(t *testing.T) {
	type args struct {
		name string
		tp   Type
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want *Parameter
	}{
		{
			name: "Should return a new parameter",
			args: args{
				name: "test",
				tp:   NewType("qual"),
			},
			want: &Parameter{
				Name: "test",
				Type: NewType("qual"),
			},
		},
		{
			name: "Should return a new parameter with docs",
			args: args{
				name: "test",
				tp:   NewType("qual"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Parameter{
				Name: "test",
				Type: NewType("qual"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParameter(tt.args.name, tt.args.tp, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFieldTags(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want FieldTags
	}{
		{
			name: "Should return a new field tag with initial key and value pair",
			args: args{
				key:   "json",
				value: "name",
			},
			want: FieldTags{
				"json": "name",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFieldTags(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFieldTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStructField(t *testing.T) {
	type args struct {
		name string
		tp   Type
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want StructField
	}{
		{
			name: "Should return a new structure field",
			args: args{
				name: "Test",
				tp:   NewType("string"),
			},
			want: StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
				},
			},
		},
		{
			name: "Should return a new structure field with docs",
			args: args{
				name: "Test",
				tp:   NewType("string"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
					docs: []Comment{
						"Some",
						"Comments",
						"Go Here",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStructField(tt.args.name, tt.args.tp, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStructField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStructFieldWithTag(t *testing.T) {
	type args struct {
		name string
		tp   Type
		tags FieldTags
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want StructField
	}{
		{
			name: "Should return a new structure field with tags",
			args: args{
				name: "Test",
				tp:   NewType("string"),
				tags: NewFieldTags("json", "test"),
			},
			want: StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
				},
				Tags: NewFieldTags("json", "test"),
			},
		},
		{
			name: "Should return a new structure field with tags and docs",
			args: args{
				name: "Test",
				tp:   NewType("string"),
				tags: NewFieldTags("json", "test"),
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
					docs: []Comment{
						"Some",
						"Comments",
						"Go Here",
					},
				},
				Tags: NewFieldTags("json", "test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStructFieldWithTag(tt.args.name, tt.args.tp, tt.args.tags, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStructFieldWithTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStruct(t *testing.T) {
	type args struct {
		name string
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want *Struct
	}{
		{
			name: "Should return a new structure",
			args: args{
				name: "Test",
			},
			want: &Struct{Name: "Test"},
		},
		{
			name: "Should return a new structure with docs",
			args: args{
				name: "Test",
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Struct{
				Name: "Test",
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStruct(tt.args.name, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStructWithFields(t *testing.T) {
	type args struct {
		name   string
		fields []StructField
		docs   []Comment
	}
	tests := []struct {
		name string
		args args
		want *Struct
	}{
		{
			name: "Should return a new structure with fields",
			args: args{
				name: "Test",
				fields: []StructField{
					NewStructField("test", NewType("string")),
				},
			},
			want: &Struct{
				Name: "Test",
				Fields: []StructField{
					NewStructField("test", NewType("string")),
				},
			},
		},
		{
			name: "Should return a new structure with fields and docs",
			args: args{
				name: "Test",
				fields: []StructField{
					NewStructField("test", NewType("string")),
				},
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
			want: &Struct{
				Name: "Test",
				Fields: []StructField{
					NewStructField("test", NewType("string")),
				},
				docs: []Comment{
					"Some",
					"Comments",
					"Go Here",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStructWithFields(tt.args.name, tt.args.fields, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStructWithFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamsMethodOption(t *testing.T) {
	type args struct {
		params []Parameter
		mth    *Method
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should add parameters to method",
			args: args{
				params: []Parameter{
					*NewParameter("test", NewType("string")),
				},
				mth: NewMethod("Hi"),
			},
			want: &Method{
				Name: "Hi",
				Params: []Parameter{
					*NewParameter("test", NewType("string")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParamsMethodOption(tt.args.params...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Method = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestResultsMethodOption(t *testing.T) {
	type args struct {
		results []Parameter
		mth     *Method
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should add results to method",
			args: args{
				results: []Parameter{
					*NewParameter("test", NewType("string")),
				},
				mth: NewMethod("Hi"),
			},
			want: &Method{
				Name: "Hi",
				Results: []Parameter{
					*NewParameter("test", NewType("string")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResultsMethodOption(tt.args.results...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Method = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestRecvMethodOption(t *testing.T) {
	type args struct {
		recv *Parameter
		mth  *Method
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should add receiver to method",
			args: args{
				recv: NewParameter("test", NewType("string")),
				mth:  NewMethod("Hi"),
			},
			want: &Method{
				Name: "Hi",
				Recv: NewParameter("test", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecvMethodOption(tt.args.recv)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Method = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestBodyMethodOption(t *testing.T) {
	type args struct {
		body []jen.Code
		mth  *Method
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should add body to method",
			args: args{
				body: []jen.Code{
					jen.Id("print(\"hello\")"),
				},
				mth: NewMethod("Hi"),
			},
			want: &Method{
				Name: "Hi",
				Body: []jen.Code{
					jen.Id("print(\"hello\")"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BodyMethodOption(tt.args.body)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Method = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestDocsMethodOption(t *testing.T) {
	type args struct {
		docs []Comment
		mth  *Method
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should add parameters to method",
			args: args{
				docs: []Comment{"Hello", "Hi"},
				mth:  NewMethod("Hi"),
			},
			want: &Method{
				Name: "Hi",
				docs: []Comment{"Hello", "Hi"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DocsMethodOption(tt.args.docs...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Method = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestNewMethod(t *testing.T) {
	type args struct {
		name    string
		options []MethodOptions
	}
	tests := []struct {
		name string
		args args
		want *Method
	}{
		{
			name: "Should return a new method",
			args: args{
				name: "Test",
			},
			want: &Method{
				Name: "Test",
			},
		},
		{
			name: "Should return a new method with options",
			args: args{
				name: "Test",
				options: []MethodOptions{
					DocsMethodOption("Test", "Hello"),
					RecvMethodOption(NewParameter("hi", NewType("string"))),
				},
			},
			want: &Method{
				Name: "Test",
				docs: []Comment{"Test", "Hello"},
				Recv: NewParameter("hi", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMethod(tt.args.name, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInterfaceMethod(t *testing.T) {
	type args struct {
		name    string
		options []MethodOptions
	}
	tests := []struct {
		name string
		args args
		want InterfaceMethod
	}{
		{
			name: "Should return a new interface method",
			args: args{
				name: "Hi",
			},
			want: InterfaceMethod{
				Name: "Hi",
			},
		},
		{
			name: "Should return a new interface method with options",
			args: args{
				name: "Hi",
				options: []MethodOptions{
					DocsMethodOption("Test", "Hello"),
					RecvMethodOption(NewParameter("hi", NewType("string"))),
				},
			},
			want: InterfaceMethod{
				Name: "Hi",
				docs: []Comment{"Test", "Hello"},
				Recv: NewParameter("hi", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInterfaceMethod(tt.args.name, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInterfaceMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInterface(t *testing.T) {
	type args struct {
		name    string
		methods []InterfaceMethod
		docs    []Comment
	}
	tests := []struct {
		name string
		args args
		want *Interface
	}{
		{
			name: "Should return a new interface",
			args: args{
				name: "Test",
			},
			want: &Interface{
				Name: "Test",
			},
		},
		{
			name: "Should return a new interface with methods",
			args: args{
				name: "Test",
				methods: []InterfaceMethod{
					NewInterfaceMethod("Hi"),
				},
			},
			want: &Interface{
				Name: "Test",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Hi"),
				},
			},
		},
		{
			name: "Should return a new interface with methods and docs",
			args: args{
				name: "Test",
				methods: []InterfaceMethod{
					NewInterfaceMethod("Hi"),
				},
				docs: []Comment{"Hi", "Hello"},
			},
			want: &Interface{
				Name: "Test",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Hi"),
				},
				docs: []Comment{"Hi", "Hello"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInterface(tt.args.name, tt.args.methods, tt.args.docs...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_Code(t *testing.T) {
	tests := []struct {
		name string
		c    Comment
		want *jen.Statement
	}{
		{
			name: "Should return the jen representation of a comment",
			c:    "Test",
			want: jen.Comment("Test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Comment.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_String(t *testing.T) {
	tests := []struct {
		name string
		c    Comment
		want string
	}{
		{
			name: "Should return the jen representation of a comment",
			c:    "Test",
			want: "// Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Comment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComment_AddDocs(t *testing.T) {
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name string
		c    Comment
		args args
	}{
		{
			name: "Should do nothing",
			c:    "Test",
			args: args{
				docs: []Comment{
					"Hello",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AddDocs(tt.args.docs...)
		})
	}
}

func TestType_Code(t *testing.T) {
	type fields struct {
		Import    *Import
		Method    *MethodType
		Pointer   bool
		Qualifier string
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name: "Should return the correct jen representation of the type if the type is only a qualifier",
			fields: fields{
				Qualifier: "string",
			},
			want: jen.Id("string"),
		},
		{
			name: "Should return the correct jen representation of the type if the type is pointer qualifier",
			fields: fields{
				Pointer:   true,
				Qualifier: "string",
			},
			want: jen.Id("*").Id("string"),
		},
		{
			name: "Should return the correct jen representation of the type if the type is import qualifier",
			fields: fields{
				Import: &Import{
					Alias: "test",
					Path:  "test/test/test",
				},
				Qualifier: "Test",
			},
			want: jen.Qual("test/test/test", "Test"),
		},
		{
			name: "Should return the correct jen representation of the type if the type is pointer import qualifier",
			fields: fields{
				Pointer: true,
				Import: &Import{
					Path: "test/test/test",
				},
				Qualifier: "Test",
			},
			want: jen.Id("*").Qual("test/test/test", "Test"),
		},
		{
			name: "Should return the correct jen representation of the type if the type is method type",
			fields: fields{
				Method: NewMethodType(),
			},
			want: jen.Add(jen.Func().Params().Params()),
		},
		{
			name: "Should return the correct jen representation of the type if the type is method type pointer",
			fields: fields{
				Pointer: true,
				Method:  NewMethodType(),
			},
			want: jen.Id("*").Add(jen.Func().Params().Params()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Type{
				Import:    tt.fields.Import,
				Method:    tt.fields.Method,
				Pointer:   tt.fields.Pointer,
				Qualifier: tt.fields.Qualifier,
			}
			if got := tp.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Type.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_String(t *testing.T) {
	type fields struct {
		Import    *Import
		Method    *MethodType
		Pointer   bool
		Qualifier string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the correct go source of the type if the type is only a qualifier",
			fields: fields{
				Qualifier: "string",
			},
			want: "string",
		},
		{
			name: "Should return the correct go source of the type if the type is pointer qualifier",
			fields: fields{
				Pointer:   true,
				Qualifier: "string",
			},
			want: "*string",
		},
		{
			name: "Should return the correct go source of the type if the type is import qualifier",
			fields: fields{
				Import: &Import{
					Alias: "test",
					Path:  "test/test/test",
				},
				Qualifier: "Test",
			},
			want: "test.Test",
		},
		{
			name: "Should return the correct go source of the type if the type is pointer import qualifier",
			fields: fields{
				Pointer: true,
				Import: &Import{
					Path: "test/test/test",
				},
				Qualifier: "Test",
			},
			want: "*test.Test",
		},
		{
			name: "Should return the correct go source of the type if the type is method type",
			fields: fields{
				Method: NewMethodType(),
			},
			want: "func()",
		},
		{
			name: "Should return the correct go source of the type if the type is method type pointer",
			fields: fields{
				Pointer: true,
				Method:  NewMethodType(),
			},
			want: "*func()",
		},
		{
			name: "Should return the correct go source of the type if the type is method type with parameters",
			fields: fields{
				Method: NewMethodType(
					ParamsMethodOption(
						*NewParameter("a", NewType("string")),
					),
					ResultsMethodOption(
						*NewParameter("", NewType("string")),
					),
				),
			},
			want: "func(a string) string",
		},
		{
			name: "Should return the correct go source of the type if the type is method type pointer with parameters",
			fields: fields{
				Pointer: true,
				Method: NewMethodType(
					ParamsMethodOption(
						*NewParameter("a", NewType("string")),
					),
					ResultsMethodOption(
						*NewParameter("", NewType("string")),
					),
				),
			},
			want: "*func(a string) string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Type{
				Import:    tt.fields.Import,
				Method:    tt.fields.Method,
				Pointer:   tt.fields.Pointer,
				Qualifier: tt.fields.Qualifier,
			}
			if got := tp.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_AddDocs(t *testing.T) {
	type fields struct {
		Import    *Import
		Method    *MethodType
		Pointer   bool
		Qualifier string
	}
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should do nothing",
			args: args{
				docs: []Comment{"test", "hello 123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Type{
				Import:    tt.fields.Import,
				Method:    tt.fields.Method,
				Pointer:   tt.fields.Pointer,
				Qualifier: tt.fields.Qualifier,
			}
			tp.AddDocs(tt.args.docs...)
		})
	}
}

func TestVar_Code(t *testing.T) {
	type fields struct {
		docs  []Comment
		Name  string
		Type  Type
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name: "Should return the correct jen representation of a variable",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
			},
			want: jen.Var().Id("Test").Add(jen.Id("string")),
		},
		{
			name: "Should return the correct jen representation of a variable with value",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				Value: 4,
			},
			want: jen.Var().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
		},
		{
			name: "Should return the correct jen representation of a variable with docs",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
				docs: []Comment{"Hello"},
			},
			want: jen.Add(jen.Comment("Hello")).Line().Var().Id("Test").Add(jen.Id("string")),
		},
		{
			name: "Should return the correct jen representation of a variable with docs and value",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				docs:  []Comment{"Hello"},
				Value: 4,
			},
			want: jen.Add(jen.Comment("Hello")).Line().Var().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Var{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := v.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Var.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_String(t *testing.T) {
	type fields struct {
		docs  []Comment
		Name  string
		Type  Type
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the correct go string of a variable",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
			},
			want: "var Test string",
		},
		{
			name: "Should return the correct go string of a variable with value",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				Value: 4,
			},
			want: "var Test int = 4",
		},
		{
			name: "Should return the correct go string of a variable with docs",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
				docs: []Comment{"Hello"},
			},
			want: "// Hello\nvar Test string",
		},
		{
			name: "Should return the correct go string of a variable with docs and value",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				docs:  []Comment{"Hello"},
				Value: 4,
			},
			want: "// Hello\nvar Test int = 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Var{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Var.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVar_AddDocs(t *testing.T) {
	type fields struct {
		docs  []Comment
		Name  string
		Type  Type
		Value interface{}
	}
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Comment
	}{
		{
			name: "Should add the docs to the variable",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
			},
			args: args{
				docs: []Comment{"This is some docs", "This is some more docs"},
			},
			want: []Comment{"This is some docs", "This is some more docs"},
		}, {
			name: "Should add the docs to the variable with existing docs",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
				docs: []Comment{"Some initial docs"},
			},
			args: args{
				docs: []Comment{"This is some docs", "This is some more docs"},
			},
			want: []Comment{"Some initial docs", "This is some docs", "This is some more docs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Var{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			v.AddDocs(tt.args.docs...)
			if !reflect.DeepEqual(v.docs, tt.want) {
				t.Errorf("Var.docs = %v, want %v", v.docs, tt.want)
			}
		})
	}
}

func TestConst_Code(t *testing.T) {
	type fields struct {
		docs  []Comment
		Name  string
		Type  Type
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{

		{
			name: "Should return the correct jen representation of a constant",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				Value: 4,
			},
			want: jen.Const().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
		},

		{
			name: "Should return the correct jen representation of a constant with docs",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				docs:  []Comment{"Hello"},
				Value: 4,
			},
			want: jen.Add(jen.Comment("Hello")).Line().Const().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Const{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := c.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Const.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestConst_String(t *testing.T) {
	type fields struct {
		docs  []Comment
		Name  string
		Type  Type
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the correct go string of a constant",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				Value: 4,
			},
			want: "const Test int = 4",
		},
		{
			name: "Should return the correct go string of a constant with docs",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				docs:  []Comment{"Hello"},
				Value: 4,
			},
			want: "// Hello\nconst Test int = 4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Const{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Const.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParameter_Code(t *testing.T) {
	type fields struct {
		docs []Comment
		Name string
		Type Type
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
				docs: tt.fields.docs,
				Name: tt.fields.Name,
				Type: tt.fields.Type,
			}
			if got := p.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parameter.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParameter_String(t *testing.T) {
	type fields struct {
		docs []Comment
		Name string
		Type Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
				docs: tt.fields.docs,
				Name: tt.fields.Name,
				Type: tt.fields.Type,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("Parameter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParameter_AddDocs(t *testing.T) {
	type fields struct {
		docs []Comment
		Name string
		Type Type
	}
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
				docs: tt.fields.docs,
				Name: tt.fields.Name,
				Type: tt.fields.Type,
			}
			p.AddDocs(tt.args.docs...)
		})
	}
}

func TestMethod_Code(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Recv    *Parameter
		Params  []Parameter
		Results []Parameter
		Body    []jen.Code
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Method{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Method.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethod_AddDocs(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Recv    *Parameter
		Params  []Parameter
		Results []Parameter
		Body    []jen.Code
	}
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Method{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			m.AddDocs(tt.args.docs...)
		})
	}
}

func TestMethod_String(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Recv    *Parameter
		Params  []Parameter
		Results []Parameter
		Body    []jen.Code
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Method{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Method.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethod_AddStringBody(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Recv    *Parameter
		Params  []Parameter
		Results []Parameter
		Body    []jen.Code
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Method{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			m.AddStringBody(tt.args.s)
		})
	}
}

func TestStructField_Code(t *testing.T) {
	type fields struct {
		Parameter Parameter
		Tags      FieldTags
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructField{
				Parameter: tt.fields.Parameter,
				Tags:      tt.fields.Tags,
			}
			if got := s.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructField.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct_Code(t *testing.T) {
	type fields struct {
		docs   []Comment
		Name   string
		Fields []StructField
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Struct{
				docs:   tt.fields.docs,
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			if got := s.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Struct.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fieldList(t *testing.T) {
	type args struct {
		fields []StructField
	}
	tests := []struct {
		name  string
		args  args
		wantF []jen.Code
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotF := fieldList(tt.args.fields); !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("fieldList() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestStruct_String(t *testing.T) {
	type fields struct {
		docs   []Comment
		Name   string
		Fields []StructField
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Struct{
				docs:   tt.fields.docs,
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("Struct.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct_AddDocs(t *testing.T) {
	type fields struct {
		docs   []Comment
		Name   string
		Fields []StructField
	}
	type args struct {
		docs []Comment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Struct{
				docs:   tt.fields.docs,
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			s.AddDocs(tt.args.docs...)
		})
	}
}

func TestInterfaceMethod_Code(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Recv    *Parameter
		Params  []Parameter
		Results []Parameter
		Body    []jen.Code
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InterfaceMethod{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceMethod.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldTags_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		f    *FieldTags
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Set(tt.args.key, tt.args.value)
		})
	}
}

func Test_codeString(t *testing.T) {
	type args struct {
		c Code
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codeString(tt.args.c); got != tt.want {
				t.Errorf("codeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addDocsCode(t *testing.T) {
	type args struct {
		c    *jen.Statement
		docs []Comment
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addDocsCode(tt.args.c, tt.args.docs)
		})
	}
}

func Test_paramsList(t *testing.T) {
	type args struct {
		paramList []Parameter
	}
	tests := []struct {
		name  string
		args  args
		wantL []jen.Code
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotL := paramsList(tt.args.paramList); !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("paramsList() = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}

func TestNewTypeMethod(t *testing.T) {
	type args struct {
		options []MethodOptions
	}
	tests := []struct {
		name string
		args args
		want *MethodType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMethodType(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTypeMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeMethod_Code(t *testing.T) {
	type fields struct {
		Name    string
		docs    []Comment
		Params  []Parameter
		Results []Parameter
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MethodType{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodType.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
