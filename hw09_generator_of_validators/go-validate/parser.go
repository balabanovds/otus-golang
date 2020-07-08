package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

var (
	ErrParseFailed        = errors.New("parser: failed to parse")
	ErrParseStructs404    = errors.New("parser: structs not found")
	ErrParseFieldTagEmpty = errors.New("parser: filed tag empty")
)

type parsedData struct {
	packageName string
	structs     map[string]parsedStruct
	idents      map[string]string
	tags        map[tagType]struct{} // this needed to generate imports
}

func newData(packageName string) *parsedData {
	return &parsedData{
		packageName: packageName,
		structs:     make(map[string]parsedStruct),
		idents:      make(map[string]string),
		tags:        make(map[tagType]struct{}),
	}
}

func (d *parsedData) addIdent(name string, ident *ast.Ident) {
	d.idents[name] = ident.Name
}

func (d *parsedData) addStruct(name string, s *ast.StructType) error {
	vs := newStruct(name)

	for _, f := range s.Fields.List {
		vField, err := newField(d.idents, f)
		if err != nil {
			// just skip field if tags empty
			if errors.Is(err, ErrParseFieldTagEmpty) {
				continue
			}
			// if smth weird
			return err
		}
		vs.fields = append(vs.fields, vField)

		// add all tags found for field
		for _, tag := range vField.tags {
			d.tags[tag.tType] = struct{}{}
		}
	}

	if len(vs.fields) != 0 {
		d.structs[name] = vs
	}

	return nil
}

func parse(fileName string) (*parsedData, error) {
	fset := token.NewFileSet()
	parseFile, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return nil, err
	}

	data := newData(parseFile.Name.Name)

	// first - we collect all idents
	for name, obj := range parseFile.Scope.Objects {
		switch x := obj.Decl.(*ast.TypeSpec).Type.(type) {
		case *ast.Ident:
			data.addIdent(name, x)
		}
	}

	// last - we run ones again to collect all structs
	for name, obj := range parseFile.Scope.Objects {
		switch x := obj.Decl.(*ast.TypeSpec).Type.(type) {
		case *ast.StructType:
			err = data.addStruct(name, x)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(data.structs) == 0 {
		return nil, ErrParseStructs404
	}

	return data, nil
}

//
//func inspectForStructs(file *ast.File) []parsedStruct {
//	var structName string
//	var structs []parsedStruct
//
//	ast.Inspect(file, func(n ast.Node) bool {
//		switch x := n.(type) {
//		case *ast.Ident:
//			structName = x.Name
//		case *ast.StructType:
//			vs := newStruct(structName)
//
//			for _, f := range x.Fields.List {
//
//				vField, ok := newField(f)
//				if !ok {
//					continue
//				}
//
//				vs.fields = append(vs.fields, vField)
//
//			}
//
//			if len(vs.fields) != 0 {
//				structs = append(structs, vs)
//			}
//		}
//		return true
//	})
//
//	return structs
//}
