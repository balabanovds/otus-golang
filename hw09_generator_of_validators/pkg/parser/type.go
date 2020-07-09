package parser

import (
	"go/ast"
)

type FType int

const (
	FUnknown FType = iota
	FString
	FInt
	FSliceString
	FSliceInt
)

type FieldType struct {
	Key      FType
	Declared string
	Root     string
}

func newType(idents map[string]string, e ast.Expr) (*FieldType, error) {
	var declared, root string

	switch t := e.(type) {
	case *ast.Ident:
		declared = t.Name
	case *ast.ArrayType:
		declared = "[]" + t.Elt.(*ast.Ident).Name
	}

	key := getType(declared)
	if key == FUnknown {
		var ok bool
		root, ok = idents[declared]
		if !ok {
			return nil, ErrParseFailed
		}
		key = getType(root)
	} else {
		root = declared
	}

	return &FieldType{key, declared, root}, nil
}

func getType(key string) FType {
	switch key {
	case "string":
		return FString
	case "int":
		return FInt
	case "[]string":
		return FSliceString
	case "[]int":
		return FSliceInt
	default:
		return FUnknown
	}

}
