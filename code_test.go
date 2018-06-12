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

func TestFunctionTypeOption(t *testing.T) {
	type args struct {
		tp Type
		m  *FunctionType
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should add the function to the type",
			args: args{
				tp: NewType("Test"),
				m:  NewFunctionType(),
			},
			want: Type{
				Qualifier: "Test",
				Function:  NewFunctionType(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FunctionTypeOption(tt.args.m)
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

func TestNewFunctionType(t *testing.T) {
	type args struct {
		options []FunctionOptions
	}
	tests := []struct {
		name string
		args args
		want *FunctionType
	}{
		{
			name: "Should return a new function type",
			want: &FunctionType{},
		},
		{
			name: "Should return a new function type with options",
			args: args{
				options: []FunctionOptions{
					DocsFunctionOption("Test", "Hello"),
					RecvFunctionOption(NewParameter("hi", NewType("string"))),
				},
			},
			want: &FunctionType{
				docs: []Comment{"Test", "Hello"},
				Recv: NewParameter("hi", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFunctionType(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFunctionType() = %v, want %v", got, tt.want)
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
			},
			want: &Parameter{
				Name: "test",
				Type: NewType("qual"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParameter(tt.args.name, tt.args.tp); !reflect.DeepEqual(got, tt.want) {
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
		want *FieldTags
	}{
		{
			name: "Should return a new field tag with initial key and value pair",
			args: args{
				key:   "json",
				value: "name",
			},
			want: &FieldTags{
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
		want *StructField
	}{
		{
			name: "Should return a new structure field",
			args: args{
				name: "Test",
				tp:   NewType("string"),
			},
			want: &StructField{
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
			want: &StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
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
		tags *FieldTags
		docs []Comment
	}
	tests := []struct {
		name string
		args args
		want *StructField
	}{
		{
			name: "Should return a new structure field with tags",
			args: args{
				name: "Test",
				tp:   NewType("string"),
				tags: NewFieldTags("json", "test"),
			},
			want: &StructField{
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
			want: &StructField{
				Parameter: Parameter{
					Name: "Test",
					Type: NewType("string"),
				},
				Tags: NewFieldTags("json", "test"),
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
					*NewStructField("test", NewType("string")),
				},
			},
			want: &Struct{
				Name: "Test",
				Fields: []StructField{
					*NewStructField("test", NewType("string")),
				},
			},
		},
		{
			name: "Should return a new structure with fields and docs",
			args: args{
				name: "Test",
				fields: []StructField{
					*NewStructField("test", NewType("string")),
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
					*NewStructField("test", NewType("string")),
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

func TestParamsFunctionOption(t *testing.T) {
	type args struct {
		params []Parameter
		mth    *Function
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should add parameters to function",
			args: args{
				params: []Parameter{
					*NewParameter("test", NewType("string")),
				},
				mth: NewFunction("Hi"),
			},
			want: &Function{
				Name: "Hi",
				Params: []Parameter{
					*NewParameter("test", NewType("string")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParamsFunctionOption(tt.args.params...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Function = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestResultsFunctionOption(t *testing.T) {
	type args struct {
		results []Parameter
		mth     *Function
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should add results to function",
			args: args{
				results: []Parameter{
					*NewParameter("test", NewType("string")),
				},
				mth: NewFunction("Hi"),
			},
			want: &Function{
				Name: "Hi",
				Results: []Parameter{
					*NewParameter("test", NewType("string")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResultsFunctionOption(tt.args.results...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Function = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestRecvFunctionOption(t *testing.T) {
	type args struct {
		recv *Parameter
		mth  *Function
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should add receiver to function",
			args: args{
				recv: NewParameter("test", NewType("string")),
				mth:  NewFunction("Hi"),
			},
			want: &Function{
				Name: "Hi",
				Recv: NewParameter("test", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecvFunctionOption(tt.args.recv)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Function = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestBodyFunctionOption(t *testing.T) {
	type args struct {
		body []jen.Code
		mth  *Function
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should add body to function",
			args: args{
				body: []jen.Code{
					jen.Id("print(\"hello\")"),
				},
				mth: NewFunction("Hi"),
			},
			want: &Function{
				Name: "Hi",
				Body: []jen.Code{
					jen.Id("print(\"hello\")"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BodyFunctionOption(tt.args.body...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Function = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestDocsFunctionOption(t *testing.T) {
	type args struct {
		docs []Comment
		mth  *Function
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should add parameters to function",
			args: args{
				docs: []Comment{"Hello", "Hi"},
				mth:  NewFunction("Hi"),
			},
			want: &Function{
				Name: "Hi",
				docs: []Comment{"Hello", "Hi"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DocsFunctionOption(tt.args.docs...)
			got(tt.args.mth)
			if !reflect.DeepEqual(tt.args.mth, tt.want) {
				t.Errorf("Function = %v, want %v", tt.args.mth, tt.want)
			}
		})
	}
}

func TestNewFunction(t *testing.T) {
	type args struct {
		name    string
		options []FunctionOptions
	}
	tests := []struct {
		name string
		args args
		want *Function
	}{
		{
			name: "Should return a new function",
			args: args{
				name: "Test",
			},
			want: &Function{
				Name: "Test",
			},
		},
		{
			name: "Should return a new function with options",
			args: args{
				name: "Test",
				options: []FunctionOptions{
					DocsFunctionOption("Test", "Hello"),
					RecvFunctionOption(NewParameter("hi", NewType("string"))),
				},
			},
			want: &Function{
				Name: "Test",
				docs: []Comment{"Test", "Hello"},
				Recv: NewParameter("hi", NewType("string")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFunction(tt.args.name, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewInterfaceMethod(t *testing.T) {
	type args struct {
		name    string
		options []FunctionOptions
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
				options: []FunctionOptions{
					DocsFunctionOption("Test", "Hello"),
					RecvFunctionOption(NewParameter("hi", NewType("string"))),
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
		RawType   *jen.Statement
		Method    *FunctionType
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
			name: "Should return the correct jen representation of the type if the type is function type",
			fields: fields{
				Method: NewFunctionType(),
			},
			want: jen.Add(jen.Func().Params()),
		},
		{
			name: "Should return the correct jen representation of the type if the type is function type pointer",
			fields: fields{
				Pointer: true,
				Method:  NewFunctionType(),
			},
			want: jen.Id("*").Add(jen.Func().Params()),
		},
		{
			name: "Should return the correct jen representation of the raw type",
			fields: fields{
				RawType: jen.Map(jen.String()).Id("*").Qual("test", "Test"),
			},
			want: jen.Map(jen.String()).Id("*").Qual("test", "Test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Type{
				Import:    tt.fields.Import,
				Function:  tt.fields.Method,
				RawType:   tt.fields.RawType,
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
		RawType   *jen.Statement
		Function  *FunctionType
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
			name: "Should return the correct go source of the type if the type is function type",
			fields: fields{
				Function: NewFunctionType(),
			},
			want: "func()",
		},
		{
			name: "Should return the correct go source of the type if the type is function type pointer",
			fields: fields{
				Pointer:  true,
				Function: NewFunctionType(),
			},
			want: "*func()",
		},
		{
			name: "Should return the correct go source of the type if the type is function type with parameters",
			fields: fields{
				Function: NewFunctionType(
					ParamsFunctionOption(
						*NewParameter("a", NewType("string")),
					),
					ResultsFunctionOption(
						*NewParameter("", NewType("string")),
					),
				),
			},
			want: "func(a string) string",
		},
		{
			name: "Should return the correct go source of the type if the type is function type pointer with parameters",
			fields: fields{
				Pointer: true,
				Function: NewFunctionType(
					ParamsFunctionOption(
						*NewParameter("a", NewType("string")),
					),
					ResultsFunctionOption(
						*NewParameter("", NewType("string")),
					),
				),
			},
			want: "*func(a string) string",
		},
		{
			name: "Should return the correct jen representation of the raw type",
			fields: fields{
				RawType: jen.Map(jen.String()).Id("*").Qual("test", "Test"),
			},
			want: "map[string]*test.Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &Type{
				Import:    tt.fields.Import,
				RawType:   tt.fields.RawType,
				Function:  tt.fields.Function,
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
		Function  *FunctionType
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
				Function:  tt.fields.Function,
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
			want: jen.Add(jen.Comment("Hello").Line()).Var().Id("Test").Add(jen.Id("string")),
		},
		{
			name: "Should return the correct jen representation of a variable with docs and value",
			fields: fields{
				Name:  "Test",
				Type:  NewType("int"),
				docs:  []Comment{"Hello"},
				Value: 4,
			},
			want: jen.Add(jen.Comment("Hello").Line()).Var().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
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
			want: jen.Add(jen.Comment("Hello").Line()).Const().Id("Test").Add(jen.Id("int")).Op("=").Lit(4),
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

func TestConst_AddDocs(t *testing.T) {
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
		},
		{
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
			c := &Const{
				docs:  tt.fields.docs,
				Name:  tt.fields.Name,
				Type:  tt.fields.Type,
				Value: tt.fields.Value,
			}
			c.AddDocs(tt.args.docs...)
			if !reflect.DeepEqual(c.docs, tt.want) {
				t.Errorf("Const.docs = %v, want %v", c.docs, tt.want)
			}
		})
	}
}

func TestParameter_Code(t *testing.T) {
	type fields struct {
		Name string
		Type Type
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name: "Should return the jen representation of the parameter",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
			},
			want: jen.Id("Test").Add(NewType("string").Code()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
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
		Name string
		Type Type
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the jen representation of the parameter",
			fields: fields{
				Name: "Test",
				Type: NewType("string"),
			},
			want: "Test string",
		},
		{
			name: "Should return the jen representation of the parameter with a function type",
			fields: fields{
				Name: "Test",
				// some crazy type...
				Type: NewType(
					"",
					PointerTypeOption(),
					FunctionTypeOption(
						NewFunctionType(
							ParamsFunctionOption(
								*NewParameter("abc", NewType("int", PointerTypeOption())),
								*NewParameter("d", NewType("string")),
							),
							ResultsFunctionOption(
								*NewParameter(
									"r",
									NewType(
										"Test",
										PointerTypeOption(),
										ImportTypeOption(NewImport("abc", "test/abc")),
									),
								),
							),
						),
					),
				),
			},
			want: "Test *func(abc *int, d string) (r *abc.Test)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
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
		{
			name: "Should do nothing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parameter{
				Name: tt.fields.Name,
				Type: tt.fields.Type,
			}
			p.AddDocs(tt.args.docs...)
		})
	}
}

func TestFunction_Code(t *testing.T) {
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
		{
			name: "Should return the correct jen representation of a function",
			fields: fields{
				Name: "Test",
			},
			want: jen.Func().Id("Test").Params().Block(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters",
			fields: fields{
				Name:   "Test",
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: jen.Func().Id("Test").Params(NewParameter("t", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with result",
			fields: fields{
				Name:    "Test",
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Func().Id("Test").Params().Params(NewParameter("r", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters and result",
			fields: fields{
				Name:    "Test",
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Func().Id("Test").Params(
				NewParameter("t", NewType("hi")).Code(),
			).Params(NewParameter("r", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with receiver",
			fields: fields{
				Name: "Test",
				Recv: NewParameter("t", NewType("hi")),
			},
			want: jen.Func().Params(NewParameter("t", NewType("hi")).Code()).Id("Test").Params().Block(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters and receiver",
			fields: fields{
				Name:   "Test",
				Recv:   NewParameter("rcv", NewType("hi")),
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: jen.Func().Params(
				NewParameter("rcv", NewType("hi")).Code(),
			).Id("Test").Params(NewParameter("t", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters and response and receiver",
			fields: fields{
				Name:    "Test",
				Recv:    NewParameter("rcv", NewType("hi")),
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Func().Params(
				NewParameter("rcv", NewType("hi")).Code(),
			).Id("Test").Params(
				NewParameter("t", NewType("hi")).Code(),
			).Params(NewParameter("r", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters and response and receiver with docs",
			fields: fields{
				docs:    []Comment{"Hello", "Hi", "123"},
				Name:    "Test",
				Recv:    NewParameter("rcv", NewType("hi")),
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Add(jen.Comment("Hello").Line(), jen.Comment("Hi").Line(), jen.Comment("123").Line()).Func().Params(
				NewParameter("rcv", NewType("hi")).Code(),
			).Id("Test").Params(
				NewParameter("t", NewType("hi")).Code(),
			).Params(NewParameter("r", NewType("hi")).Code()).Block(),
		},
		{
			name: "Should return the correct jen representation of a function with body",
			fields: fields{
				Name: "Test",
				Body: []jen.Code{jen.Id("fmt.Println(\"Hello\")")},
			},
			want: jen.Func().Id("Test").Params().Block([]jen.Code{jen.Id("fmt.Println(\"Hello\")")}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Function{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Function.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunction_AddDocs(t *testing.T) {
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
		want   []Comment
	}{{
		name: "Should add the docs to the function",
		fields: fields{
			Name: "Test",
		},
		args: args{
			docs: []Comment{"This is some docs", "This is some more docs"},
		},
		want: []Comment{"This is some docs", "This is some more docs"},
	}, {
		name: "Should add the docs to the function with existing docs",
		fields: fields{
			Name: "Test",
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
			m := &Function{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			m.AddDocs(tt.args.docs...)
			if !reflect.DeepEqual(m.docs, tt.want) {
				t.Errorf("Function.docs = %v, want %v", m.docs, tt.want)
			}
		})
	}
}

func TestFunction_String(t *testing.T) {
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
		{
			name: "Should return the correct jen representation of a function",
			fields: fields{
				Name: "Test",
			},
			want: "func Test() {}",
		},
		{
			name: "Should return the correct jen representation of a function with parameters",
			fields: fields{
				Name:   "Test",
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: "func Test(t hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with result",
			fields: fields{
				Name:    "Test",
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "func Test() (r hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with parameters and result",
			fields: fields{
				Name:    "Test",
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "func Test(t hi) (r hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with receiver",
			fields: fields{
				Name: "Test",
				Recv: NewParameter("rcv", NewType("hi")),
			},
			want: "func (rcv hi) Test() {}",
		},
		{
			name: "Should return the correct jen representation of a function with parameters and receiver",
			fields: fields{
				Name:   "Test",
				Recv:   NewParameter("rcv", NewType("hi")),
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: "func (rcv hi) Test(t hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with parameters and response and receiver",
			fields: fields{
				Name:    "Test",
				Recv:    NewParameter("rcv", NewType("hi")),
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "func (rcv hi) Test(t hi) (r hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with parameters and response and receiver with docs",
			fields: fields{
				docs:    []Comment{"Hello", "Hi", "123"},
				Name:    "Test",
				Recv:    NewParameter("rcv", NewType("hi")),
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "// Hello\n// Hi\n// 123\nfunc (rcv hi) Test(t hi) (r hi) {}",
		},
		{
			name: "Should return the correct jen representation of a function with body",
			fields: fields{
				Name: "Test",
				Body: []jen.Code{jen.Id("fmt.Println(\"Hello\")")},
			},
			want: "func Test() {\n\tfmt.Println(\"Hello\")\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Function{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Function.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunction_AddStringBody(t *testing.T) {
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
		want   []jen.Code
	}{
		{
			name: "Should add raw string body to function",
			fields: fields{
				Name: "Test",
			},
			args: args{
				s: "fmt.Println(\"Hello\")",
			},
			want: []jen.Code{jen.Id("fmt.Println(\"Hello\")")},
		},
		{
			name: "Should add raw string body to existing body",
			fields: fields{
				Name: "Test",
				Body: []jen.Code{jen.Id("a := 2\nprint(a)")},
			},
			args: args{
				s: "fmt.Println(\"Hello\")",
			},
			want: []jen.Code{jen.Id("a := 2\nprint(a)"), jen.Id("fmt.Println(\"Hello\")")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Function{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			m.AddStringBody(tt.args.s)
			if !reflect.DeepEqual(m.Body, tt.want) {
				t.Errorf("Function.Body = %v, want %v", m.docs, tt.want)
			}
		})
	}
}

func TestStructField_Code(t *testing.T) {
	type fields struct {
		docs      []Comment
		Parameter Parameter
		Tags      *FieldTags
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name: "Should return the correct jen representative of a structure field",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
			},
			want: jen.Id("test").Add(jen.Id("string")),
		},
		{
			name: "Should return the correct jen representative of a structure field with tags",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
				Tags:      NewFieldTags("json", "test"),
			},
			want: jen.Id("test").Add(jen.Id("string")).Tag(*NewFieldTags("json", "test")),
		},
		{
			name: "Should return the correct jen representative of a structure field with tags and docs",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
				Tags:      NewFieldTags("json", "test"),
				docs:      []Comment{"Hello", "test"},
			},
			want: jen.Add(jen.Comment("Hello").Line(), jen.Comment("test").Line()).Id("test").Add(jen.Id("string")).Tag(*NewFieldTags("json", "test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructField{
				Parameter: tt.fields.Parameter,
				Tags:      tt.fields.Tags,
				docs:      tt.fields.docs,
			}
			if got := s.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructField.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructField_String(t *testing.T) {
	type fields struct {
		docs      []Comment
		Parameter Parameter
		Tags      *FieldTags
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the correct jen representative of a structure field",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
			},
			want: "test string",
		},
		{
			name: "Should return the correct jen representative of a structure field with tags",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
				Tags:      NewFieldTags("json", "test"),
			},
			want: "test string `json:\"test\"`",
		},
		{
			name: "Should return the correct jen representative of a structure field with tags and docs",
			fields: fields{
				Parameter: *NewParameter("test", NewType("string")),
				Tags:      NewFieldTags("json", "test"),
				docs:      []Comment{"Hello", "test"},
			},
			want: "// Hello\n// test\ntest string `json:\"test\"`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructField{
				docs:      tt.fields.docs,
				Parameter: tt.fields.Parameter,
				Tags:      tt.fields.Tags,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("StructField.String() = %v, want %v", got, tt.want)
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
		{
			name: "Should return the correct jen representation of a struct",
			fields: fields{
				Name: "Test",
			},
			want: jen.Type().Id("Test").Struct(),
		},
		{
			name: "Should return the correct jen representation of a struct with fields",
			fields: fields{
				Name:   "Test",
				Fields: []StructField{*NewStructField("test", NewType("string"))},
			},
			want: jen.Type().Id("Test").Struct(NewStructField("test", NewType("string")).Code()),
		},
		{
			name: "Should return the correct jen representation of a struct with fields",
			fields: fields{
				docs:   []Comment{"Hello"},
				Name:   "Test",
				Fields: []StructField{*NewStructField("test", NewType("string"))},
			},
			want: jen.Add(jen.Comment("Hello").Line()).Type().Id("Test").Struct(NewStructField("test", NewType("string")).Code()),
		},
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
		{
			name: "Should return the jen representation of a struct fields list",
			args: args{
				fields: []StructField{
					*NewStructField("test", NewType("string")),
					*NewStructField("abc", NewType("int")),
				},
			},
			wantF: []jen.Code{
				NewStructField("test", NewType("string")).Code(),
				NewStructField("abc", NewType("int")).Code(),
			},
		},
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

		{
			name: "Should return the correct jen representation of a struct",
			fields: fields{
				Name: "Test",
			},
			want: "type Test struct{}",
		},
		{
			name: "Should return the correct jen representation of a struct with fields",
			fields: fields{
				Name:   "Test",
				Fields: []StructField{*NewStructField("test", NewType("string"))},
			},
			want: "type Test struct {\n\ttest string\n}",
		},
		{
			name: "Should return the correct jen representation of a struct with fields",
			fields: fields{
				docs:   []Comment{"Hello"},
				Name:   "Test",
				Fields: []StructField{*NewStructField("test", NewType("string"))},
			},
			want: "// Hello\ntype Test struct {\n\ttest string\n}",
		},
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
		want   []Comment
	}{
		{
			name: "Should add the docs to the structure",
			fields: fields{
				Name: "Test",
			},
			args: args{
				docs: []Comment{"This is some docs", "This is some more docs"},
			},
			want: []Comment{"This is some docs", "This is some more docs"},
		},
		{
			name: "Should add the docs to the structure with existing docs",
			fields: fields{
				Name: "Test",
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
			s := &Struct{
				docs:   tt.fields.docs,
				Name:   tt.fields.Name,
				Fields: tt.fields.Fields,
			}
			s.AddDocs(tt.args.docs...)

			if !reflect.DeepEqual(s.docs, tt.want) {
				t.Errorf("Struct.docs = %v, want %v", s.docs, tt.want)
			}
		})
	}
}

func TestInterfaceMethod_Code(t *testing.T) {
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
		{
			name: "Should return the correct jen representation of an interface method",
			fields: fields{
				Name: "Test",
			},
			want: jen.Id("Test").Params(),
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
			},
			want: jen.Id("Test").Params(NewParameter("p", NewType("string")).Code()),
		},
		{
			name: "Should return the correct jen representation of an interface method with results",
			fields: fields{
				Name: "Test",
				Results: []Parameter{
					*NewParameter("p", NewType("string")),
				},
			},
			want: jen.Id("Test").Params().Params(NewParameter("p", NewType("string")).Code()),
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters and results",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
				Results: []Parameter{
					*NewParameter("r", NewType("string")),
				},
			},
			want: jen.Id("Test").Params(NewParameter("p", NewType("string")).Code()).Params(NewParameter("r", NewType("string")).Code()),
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters and results and docs",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
				Results: []Parameter{
					*NewParameter("r", NewType("string")),
				},
				docs: []Comment{"Test", "Hi"},
			},
			want: jen.Add(jen.Comment("Test").Line(), jen.Comment("Hi").Line()).Id("Test").Params(NewParameter("p", NewType("string")).Code()).Params(NewParameter("r", NewType("string")).Code()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &InterfaceMethod{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceMethod.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceMethod_String(t *testing.T) {
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
		{
			name: "Should return the correct jen representation of an interface method",
			fields: fields{
				Name: "Test",
			},
			want: "Test()",
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
			},
			want: "Test(p string)",
		},
		{
			name: "Should return the correct jen representation of an interface method with results",
			fields: fields{
				Name: "Test",
				Results: []Parameter{
					*NewParameter("p", NewType("string")),
				},
			},
			want: "Test() (p string)",
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters and results",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
				Results: []Parameter{
					*NewParameter("r", NewType("string")),
				},
			},
			want: "Test(p string) (r string)",
		},
		{
			name: "Should return the correct jen representation of an interface method with parameters and results and docs",
			fields: fields{
				Name: "Test",
				Params: []Parameter{
					*NewParameter("p", NewType("string")),
				},
				Results: []Parameter{
					*NewParameter("r", NewType("string")),
				},
				docs: []Comment{"Test", "Hi"},
			},
			want: "// Test\n// Hi\nTest(p string) (r string)",
		},
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
			if got := m.String(); got != tt.want {
				t.Errorf("InterfaceMethod.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceMethod_AddDocs(t *testing.T) {
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
		want   []Comment
	}{
		{
			name: "Should add the docs to the interface method",
			fields: fields{
				Name: "Test",
			},
			args: args{
				docs: []Comment{"This is some docs", "This is some more docs"},
			},
			want: []Comment{"This is some docs", "This is some more docs"},
		}, {
			name: "Should add the docs to the variable with existing docs",
			fields: fields{
				Name: "Test",
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
			m := &InterfaceMethod{
				Name:    tt.fields.Name,
				docs:    tt.fields.docs,
				Recv:    tt.fields.Recv,
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
				Body:    tt.fields.Body,
			}
			m.AddDocs(tt.args.docs...)
			if !reflect.DeepEqual(m.docs, tt.want) {
				t.Errorf("InterfaceMethod.docs = %v, want %v", m.docs, tt.want)
			}
		})
	}
}

func TestFieldTags_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	var tags FieldTags
	tests := []struct {
		name string
		f    *FieldTags
		args args
		want *FieldTags
	}{
		{
			name: "Should set the tag",
			f:    NewFieldTags("test", "123"),
			args: args{
				key:   "other",
				value: "value",
			},
			want: &FieldTags{
				"test":  "123",
				"other": "value",
			},
		},
		{
			name: "Should set the tag if tags nil",
			f:    &tags,
			args: args{
				key:   "other",
				value: "value",
			},
			want: &FieldTags{
				"other": "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.f, tt.want) {
				t.Errorf("FieldTags = %v, want %v", tt.f, tt.want)
			}
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
		{
			name: "Should remove unnecessary tabs",
			args: args{
				c: NewStructWithFields("Hello", []StructField{*NewStructField("s", NewType("string"))}),
			},
			want: "type Hello struct {\n\ts string\n}",
		},
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
		want *jen.Statement
	}{
		{
			name: "Should return add the correct jen representation of a comment list to the code",
			args: args{
				c:    &jen.Statement{},
				docs: []Comment{"Test", "123"},
			},
			want: jen.Add(jen.Comment("Test").Line(), jen.Comment("123").Line()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addDocsCode(tt.args.c, tt.args.docs)
			if !reflect.DeepEqual(tt.args.c, tt.want) {
				t.Errorf("Code = %v, want %v", tt.args.c, tt.want)
			}
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
		{
			name: "Should return the jen representation of a parameter list",
			args: args{
				paramList: []Parameter{
					*NewParameter("a", NewType("string")),
					*NewParameter("b", NewType("int")),
				},
			},
			wantL: []jen.Code{
				NewParameter("a", NewType("string")).Code(),
				NewParameter("b", NewType("int")).Code(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotL := paramsList(tt.args.paramList); !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("paramsList() = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}

func TestFunctionType_Code(t *testing.T) {
	type fields struct {
		Params  []Parameter
		Results []Parameter
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name:   "Should return the correct jen representation of a function",
			fields: fields{},
			want:   jen.Func().Params(),
		},
		{
			name: "Should return the correct jen representation of a function with parameters",
			fields: fields{
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: jen.Func().Params(NewParameter("t", NewType("hi")).Code()),
		},
		{
			name: "Should return the correct jen representation of a function with result",
			fields: fields{
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Func().Params().Params(NewParameter("r", NewType("hi")).Code()),
		},
		{
			name: "Should return the correct jen representation of a function with parameters and result",
			fields: fields{
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: jen.Func().Params(
				NewParameter("t", NewType("hi")).Code(),
			).Params(NewParameter("r", NewType("hi")).Code()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FunctionType{
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
			}
			if got := m.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FunctionType.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionType_String(t *testing.T) {
	type fields struct {
		Params  []Parameter
		Results []Parameter
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Should return the correct jen representation of a function",
			fields: fields{},
			want:   "func()",
		},
		{
			name: "Should return the correct jen representation of a function with parameters",
			fields: fields{
				Params: []Parameter{*NewParameter("t", NewType("hi"))},
			},
			want: "func(t hi)",
		},
		{
			name: "Should return the correct jen representation of a function with result",
			fields: fields{
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "func() (r hi)",
		},
		{
			name: "Should return the correct jen representation of a function with parameters and result",
			fields: fields{
				Params:  []Parameter{*NewParameter("t", NewType("hi"))},
				Results: []Parameter{*NewParameter("r", NewType("hi"))},
			},
			want: "func(t hi) (r hi)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FunctionType{
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("FunctionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionType_AddDocs(t *testing.T) {
	type fields struct {
		Params  []Parameter
		Results []Parameter
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FunctionType{
				Params:  tt.fields.Params,
				Results: tt.fields.Results,
			}
			m.AddDocs(tt.args.docs...)
		})
	}
}

func Test_prepareLines(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should remove unnecessary tabs",
			args: args{
				s: "\ttype Hello struct {\n\t\ts string\n}",
			},
			want: "type Hello struct {\n\ts string\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareLines(tt.args.s); got != tt.want {
				t.Errorf("prepareLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_Code(t *testing.T) {
	type fields struct {
		docs    []Comment
		Name    string
		Methods []InterfaceMethod
	}
	tests := []struct {
		name   string
		fields fields
		want   *jen.Statement
	}{
		{
			name: "Should return the correct jen representation of the interface",
			fields: fields{
				Name: "MyInterface",
			},
			want: jen.Type().Id("MyInterface").Interface(),
		},
		{
			name: "Should return the correct jen representation of the interface with methods",
			fields: fields{
				Name: "MyInterface",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Test"),
					NewInterfaceMethod("Test2"),
				},
			},
			want: jen.Type().Id("MyInterface").Interface(
				func() *jen.Statement {
					im := NewInterfaceMethod("Test")
					return im.Code()
				}(),
				func() *jen.Statement {
					im := NewInterfaceMethod("Test2")
					return im.Code()
				}(),
			),
		},
		{
			name: "Should return the correct jen representation of the interface with methods",
			fields: fields{
				docs: []Comment{"Hello", "Hi"},
				Name: "MyInterface",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Test"),
					NewInterfaceMethod("Test2"),
				},
			},
			want: jen.Add(jen.Comment("Hello").Line(), jen.Comment("Hi").Line()).Type().Id("MyInterface").Interface(
				func() *jen.Statement {
					im := NewInterfaceMethod("Test")
					return im.Code()
				}(),
				func() *jen.Statement {
					im := NewInterfaceMethod("Test2")
					return im.Code()
				}(),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interface{
				docs:    tt.fields.docs,
				Name:    tt.fields.Name,
				Methods: tt.fields.Methods,
			}
			if got := i.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interface.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_String(t *testing.T) {
	type fields struct {
		docs    []Comment
		Name    string
		Methods []InterfaceMethod
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{

		{
			name: "Should return the correct jen representation of the interface",
			fields: fields{
				Name: "MyInterface",
			},
			want: "type MyInterface interface{}",
		},
		{
			name: "Should return the correct jen representation of the interface with methods",
			fields: fields{
				Name: "MyInterface",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Test"),
					NewInterfaceMethod("Test2"),
				},
			},
			want: "type MyInterface interface {\n\tTest()\n\tTest2()\n}",
		},
		{
			name: "Should return the correct jen representation of the interface with methods",
			fields: fields{
				docs: []Comment{"Hello", "Hi"},
				Name: "MyInterface",
				Methods: []InterfaceMethod{
					NewInterfaceMethod("Test"),
					NewInterfaceMethod("Test2"),
				},
			},
			want: "// Hello\n// Hi\ntype MyInterface interface {\n\tTest()\n\tTest2()\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interface{
				docs:    tt.fields.docs,
				Name:    tt.fields.Name,
				Methods: tt.fields.Methods,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("Interface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_AddDocs(t *testing.T) {
	type fields struct {
		docs    []Comment
		Name    string
		Methods []InterfaceMethod
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
			name: "Should add the docs to the interface",
			fields: fields{
				Name: "Test",
			},
			args: args{
				docs: []Comment{"This is some docs", "This is some more docs"},
			},
			want: []Comment{"This is some docs", "This is some more docs"},
		}, {
			name: "Should add the docs to the interface with existing docs",
			fields: fields{
				Name: "Test",
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
			i := &Interface{
				docs:    tt.fields.docs,
				Name:    tt.fields.Name,
				Methods: tt.fields.Methods,
			}
			i.AddDocs(tt.args.docs...)

			if !reflect.DeepEqual(i.docs, tt.want) {
				t.Errorf("Interface.docs = %v, want %v", i.docs, tt.want)
			}
		})
	}
}

func TestNewRawType(t *testing.T) {
	type args struct {
		tp *jen.Statement
	}
	tests := []struct {
		name string
		args args
		want Type
	}{
		{
			name: "Should return a new raw type",
			args: args{
				tp: jen.Map(jen.String()).Id("*").Qual("test", "Test"),
			},
			want: Type{
				RawType: jen.Map(jen.String()).Id("*").Qual("test", "Test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRawType(tt.args.tp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRawType() = %v, want %v", got, tt.want)
			}
		})
	}
}
