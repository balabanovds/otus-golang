package main

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
	key   fType
	value string
}

func newType(e ast.Expr) (fieldType, bool) {
	var value string

	switch t := e.(type) {
	case *ast.Ident:
		value = t.Name
	case *ast.ArrayType:
		value = "[]" + t.Elt.(*ast.Ident).Name
	}

	key := getType(value)
	if key == fUnknown {
		return fieldType{}, false
	}

	return fieldType{key, value}, true
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
