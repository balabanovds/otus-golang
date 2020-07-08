package parser

import (
	"go/ast"
	"strings"
)

type parsedStruct struct {
	name   string
	short  string
	fields []*field
}

type field struct {
	name      string
	fieldType *fieldType
	tags      []tag
}

func newStruct(name string) parsedStruct {
	short := strings.ToLower(name)[0:1]
	return parsedStruct{
		name:  name,
		short: short,
	}
}

func (s *parsedStruct) addField(data ParsedData, f *ast.Field) error {
	if f.Tag == nil {
		return ErrParseFieldTagEmpty
	}
	tags := parseTags(f.Tag.Value, data.tagToken)
	if len(tags) == 0 {
		return ErrParseFieldTagEmpty
	}

	fType, err := newType(data.Idents, f.Type)
	if err != nil {
		return err
	}

	s.fields = append(s.fields, &field{
		name:      f.Names[0].Name,
		fieldType: fType,
		tags:      tags,
	})

	return nil
}

func (s *parsedStruct) grepTags() map[tagType]struct{} {
	m := make(map[tagType]struct{})

	// add all tags found for field
	for _, f := range s.fields {
		for _, t := range f.tags {
			m[t.tType] = struct{}{}
		}
	}

	return m
}
