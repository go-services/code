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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.packageName, tt.args.code...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			if err := f.AppendAfter(tt.args.c, tt.args.new); (err != nil) != tt.wantErr {
				t.Errorf("File.AppendAfter() error = %v, wantErr %v", err, tt.wantErr)
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
	type args struct {
		c   Code
		new Code
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				pkg:     tt.fields.pkg,
				jenFile: tt.fields.jenFile,
				Code:    tt.fields.Code,
			}
			if err := f.PrependBefore(tt.args.c, tt.args.new); (err != nil) != tt.wantErr {
				t.Errorf("File.PrependBefore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
