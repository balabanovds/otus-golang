package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var (
	ErrStructsNotFound = errors.New("parser: structs not found")
)

type validatedStruct struct {
	name   string
	short  string
	fields []field
}

func newStruct(name string) validatedStruct {
	short := strings.ToLower(name)[0:1]
	return validatedStruct{
		name:  name,
		short: short,
	}
}

type field struct {
	name    string
	typeStr fieldType
	tags    []tag
}

func newField(f *ast.Field) (field, bool) {
	if f.Tag == nil {
		return field{}, false
	}
	typeStr, ok := newType(f.Type)
	if !ok {
		return field{}, false
	}

	tags := parseTags(f.Tag.Value)
	if len(tags) == 0 {
		return field{}, false
	}

	return field{
		name:    f.Names[0].Name,
		typeStr: typeStr,
		tags:    tags,
	}, true
}

func parse(fileName string) (packageName string, structs []validatedStruct, err error) {
	fset := token.NewFileSet()
	parseFile, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return
	}

	packageName = parseFile.Name.Name

	structs = inspectForStructs(parseFile)

	if len(structs) == 0 {
		err = ErrStructsNotFound
	}

	return
}

func inspectForStructs(file *ast.File) []validatedStruct {
	var structName string
	var structs []validatedStruct

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			structName = x.Name
		case *ast.StructType:
			vs := newStruct(structName)

			for _, f := range x.Fields.List {

				vField, ok := newField(f)
				if !ok {
					continue
				}

				vs.fields = append(vs.fields, vField)

			}

			if len(vs.fields) != 0 {
				structs = append(structs, vs)
			}
		}
		return true
	})

	return structs
}
