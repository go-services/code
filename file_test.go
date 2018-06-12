package code

import (
	"reflect"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestNewFile(t *testing.T) {
	type args struct {
		packageName string
		code        []Code
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		{
			name: "Should return a new file",
			args: args{
				packageName: "test",
			},
			want: &File{
				pkg:     "test",
				jenFile: jen.NewFile("test"),
			},
		},
		{
			name: "Should return a new file with code",
			args: args{
				packageName: "test",
				code:        []Code{NewStruct("Hello")},
			},
			want: &File{
				pkg:     "test",
				Code:    []Code{NewStruct("Hello")},
				jenFile: jen.NewFile("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.packageName, tt.args.code...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewImportAlias(t *testing.T) {
	type args struct {
		name string
		path string
	}
	tests := []struct {
		name string
		args args
		want ImportAlias
	}{
		{
			name: "Should return a new import alias",
			args: args{
				name: "fmt_alias",
				path: "fmt",
			},
			want: ImportAlias{
				Name: "fmt_alias",
				Path: "fmt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImportAlias(tt.args.name, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImportAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_SetImportAliases(t *testing.T) {
	type fields struct {
		pkg     string
		jenFile *jen.File
		Code    []Code
	}
	type args struct {
		ia []ImportAlias
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should set the aliases",
			fields: fields{
				pkg:     "test",
				jenFile: jen.NewFile("test"),
			},
			args: args{
				ia: []ImportAlias{
					NewImportAlias("fmt_alias", "fmt"),
				},
			},
		},
		{
			name: "Should do nothing",
			fields: fields{
				pkg: "test",
			},
			args: args{
				ia: []ImportAlias{
					NewImportAlias("fmt_alias", "fmt"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			f.SetImportAliases(tt.args.ia)
		})
	}
}

func TestFile_String(t *testing.T) {
	type fields struct {
		pkg     string
		jenFile *jen.File
		Code    []Code
	}
	tests := []struct {
		name    string
		fields  fields
		aliases []ImportAlias
		want    string
	}{
		{
			name: "Should return the correct string of the file",
			fields: fields{
				pkg:     "awesome_package",
				jenFile: jen.NewFile("awesome_package"),
			},
			want: "package awesome_package\n",
		},
		{
			name: "Should return only empty file with package if jen file is nil",
			fields: fields{
				pkg: "awesome_package",
			},
			want: "package awesome_package\n",
		},
		{
			name: "Should return the correct string of the file",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewInterface("SomeInterface", nil),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			want: "package awesome_package\n\ntype SomeInterface interface{}\n",
		},
		{
			name: "Should return the correct string of the file",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewInterface("SomeInterface", nil),
					NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			want: "package awesome_package\n\nimport \"fmt\"\n\ntype SomeInterface interface{}\n\nfunc MyMethod() {\n\tfmt.Println(\"Hello World\")\n}\n",
		},
		{
			name: "Should return the correct string of the file",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewInterface("SomeInterface", nil),
					NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			aliases: []ImportAlias{
				NewImportAlias("fmt_alias", "fmt"),
			},
			want: "package awesome_package\n\nimport fmt_alias \"fmt\"\n\ntype SomeInterface interface{}\n\nfunc MyMethod() {\n\tfmt_alias.Println(\"Hello World\")\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			if len(tt.aliases) > 0 {
				f.SetImportAliases(tt.aliases)
			}
			if got := f.String(); got != tt.want {
				t.Errorf("File.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_AppendAfter(t *testing.T) {
	type fields struct {
		pkg     string
		jenFile *jen.File
		Code    []Code
	}
	type args struct {
		c   Code
		new Code
	}
	inf := NewInterface("SomeInterface", nil)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Code
		wantErr bool
	}{
		{
			name: "Should add the code after the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					inf,
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				inf,
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should add the code after the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					inf,
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should add the code after the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					inf,
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				inf,
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should add the code after the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					inf,
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should return error if the given code is not found",
			fields: fields{
				pkg:     "awesome_package",
				Code:    []Code{},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want:    []Code{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			err := f.AppendAfter(tt.args.c, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("File.AppendAfter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if (err == nil) && !reflect.DeepEqual(f.Code, tt.want) {
				t.Errorf("File.Code = %v, want %v", f.Code, tt.want)
			}
		})
	}
}

func TestFile_PrependBefore(t *testing.T) {
	type fields struct {
		pkg     string
		jenFile *jen.File
		Code    []Code
	}
	inf := NewInterface("SomeInterface", nil)
	type args struct {
		c   Code
		new Code
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Code
		wantErr bool
	}{

		{
			name: "Should add the code before the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					inf,
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
			},
		},
		{
			name: "Should add the code before the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					inf,
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
			},
		},
		{
			name: "Should add the code before the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					inf,
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should add the code before the given code element",
			fields: fields{
				pkg: "awesome_package",
				Code: []Code{
					NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
					inf,
					NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want: []Code{
				NewFunction("MyMethod1", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
				inf,
				NewFunction("MyMethod2", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
		},
		{
			name: "Should return error if the given code is not found",
			fields: fields{
				pkg:     "awesome_package",
				Code:    []Code{},
				jenFile: jen.NewFile("awesome_package"),
			},
			args: args{
				c:   inf,
				new: NewFunction("MyMethod", BodyFunctionOption(jen.Qual("fmt", "Println").Call(jen.Lit("Hello World")))),
			},
			want:    []Code{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			err := f.PrependBefore(tt.args.c, tt.args.new)
			if (err != nil) != tt.wantErr {
				t.Errorf("File.PrependBefore() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (err == nil) && !reflect.DeepEqual(f.Code, tt.want) {
				t.Errorf("File.Code = %v, want %v", f.Code, tt.want)
			}
		})
	}
}
