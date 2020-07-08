package parser

import (
	"go/ast"
)

type fType int

const (
	fUnknown fType = iota
	fString
	fInt
	fSliceString
	fSliceInt
)

type fieldType struct {
	key      fType
	declared string
	root     string
}

func newType(idents map[string]string, e ast.Expr) (*fieldType, error) {
	var declared, root string

	switch t := e.(type) {
	case *ast.Ident:
		declared = t.Name
	case *ast.ArrayType:
		declared = "[]" + t.Elt.(*ast.Ident).Name
	}

	key := getType(declared)
	if key == fUnknown {
		var ok bool
		root, ok = idents[declared]
		if !ok {
			return nil, ErrParseFailed
		}
		key = getType(root)
	} else {
		root = declared
	}

	return &fieldType{key, declared, root}, nil
}

func getType(key string) fType {
	switch key {
	case "string":
		return fString
	case "int":
		return fInt
	case "[]string":
		return fSliceString
	case "[]int":
		return fSliceInt
	default:
		return fUnknown
	}

}
