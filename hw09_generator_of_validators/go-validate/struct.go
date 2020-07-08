package main

import (
	"go/ast"
	"strings"
)

type parsedStruct struct {
	name   string
	short  string
	fields []*field
}

func newStruct(name string) parsedStruct {
	short := strings.ToLower(name)[0:1]
	return parsedStruct{
		name:  name,
		short: short,
	}
}

type field struct {
	name      string
	fieldType *fieldType
	tags      []tag
}

func newField(idents map[string]string, f *ast.Field) (*field, error) {
	if f.Tag == nil {
		return nil, ErrParseFieldTagEmpty
	}
	tags := parseTags(f.Tag.Value)
	if len(tags) == 0 {
		return nil, ErrParseFieldTagEmpty
	}

	fType, err := newType(idents, f.Type)
	if err != nil {
		return nil, err
	}

	return &field{
		name:      f.Names[0].Name,
		fieldType: fType,
		tags:      tags,
	}, nil
}
