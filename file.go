package code

import (
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

// ImportAlias is used to specify the import alias in a file
type ImportAlias struct {
	Name string
	Path string
}

// File represents a go source file.
type File struct {
	pkg     string
	jenFile *jen.File

	Code []Code
}

// NewFile creates a new file with the given package name and optional code nodes.
func NewFile(packageName string, code ...Code) *File {
	f := &File{
		pkg:  packageName,
		Code: code,
	}
	f.jenFile = jen.NewFile(packageName)
	return f
}

// NewImportAlias creates a new import alias with the given name and path.
func NewImportAlias(name, path string) ImportAlias {
	return ImportAlias{
		Name: name,
		Path: path,
	}
}

// SetImportAliases sets the files import aliases,
// if the jenFile representation of the file is nil SetImportAliases will do nothing.
func (f *File) SetImportAliases(ia []ImportAlias) {
	if f.jenFile == nil {
		return
	}
	for _, i := range ia {
		f.jenFile.ImportAlias(i.Path, i.Name)
	}
}

// String returns the go source string of the file,
// if the jen representation of the file is nil it will return a basic file with package
func (f *File) String() string {
	if f.jenFile == nil {
		return "package " + f.pkg + "\n"
	}
	for _, c := range f.Code {
		f.jenFile.Add(c.Code())
	}
	return f.jenFile.GoString()
}

// AppendAfter appends a new code node after the given code node.
func (f *File) AppendAfter(c Code, new Code) error {
	inx := -1
	for i, v := range f.Code {
		if v == c {
			inx = i + 1
			break
		}
	}
	if inx == -1 {
		return errors.New("Could not find the code node to append after")
	}
	if inx == len(f.Code) {
		f.Code = append(f.Code, new)
		return nil
	}
	f.Code = append(
		f.Code[:inx],
		append([]Code{new}, f.Code[inx:]...)...,
	)
	return nil
}

// PrependBefore prepends a new code node before the given code node.
func (f *File) PrependBefore(c Code, new Code) error {
	inx := -1
	for i, v := range f.Code {
		if v == c {
			inx = i
			break
		}
	}
	if inx == -1 {
		return errors.New("Could not find the code node to append after")
	}
	if inx == 0 {
		f.Code = append([]Code{new}, f.Code...)
		return nil
	}

	f.Code = append(
		f.Code[:inx],
		append([]Code{new}, f.Code[inx:]...)...,
	)
	return nil
}
