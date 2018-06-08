package code

import (
	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

type ImportAlias struct {
	Name    string
	Package string
}
type File struct {
	pkg     string
	jenFile *jen.File

	Code []Code
}

func NewFile(packageName string, code ...Code) *File {
	f := &File{
		pkg:  packageName,
		Code: code,
	}
	f.jenFile = jen.NewFile(packageName)
	return f
}

func (f *File) SetImportAliases(ia []ImportAlias) {
	if f.jenFile == nil {
		return
	}
	for _, i := range ia {
		f.jenFile.ImportAlias(i.Package, i.Name)
	}
}
func (f *File) String() string {
	if f.jenFile == nil {
		return "package " + f.pkg
	}
	for _, c := range f.Code {
		f.jenFile.Add(c.Code())
	}
	return f.jenFile.GoString()
}

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
		[]Code{
			new,
			f.Code[inx],
		}...,
	)
	return nil
}

func (f *File) PrependBefore(c Code, new Code) error {
	inx := -2
	for i, v := range f.Code {
		if v == c {
			inx = i - 1
			break
		}
	}
	if inx == -2 {
		return errors.New("Could not find the code node to append after")
	}
	if inx == 0 {
		f.Code = append([]Code{new}, f.Code...)
		return nil
	}
	f.Code = append(
		f.Code[:inx],
		[]Code{
			new,
			f.Code[inx],
		}...,
	)
	return nil
}
