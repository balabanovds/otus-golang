package parser

import (
	"errors"
	"go/ast"
	"strings"
)

var (
	ErrParseTagsNotValid     = errors.New("parse: tags not valid for field")
	ErrParseFieldUnknownType = errors.New("parse: field of unknown type")
)

type ParsedStruct struct {
	Name   string
	Short  string
	Fields []*Field
}

func newStruct(name string) ParsedStruct {
	short := strings.ToLower(name)[0:1]
	return ParsedStruct{
		Name:  name,
		Short: short,
	}
}

func (s *ParsedStruct) addField(data ParsedData, f *ast.Field) error {
	if f.Tag == nil {
		return ErrParseFieldTagEmpty
	}

	fType, err := newType(data.Idents, f.Type)
	if err != nil {
		return err
	}

	tags := parseTags(f.Tag.Value, data.tagToken, fType.Key)
	if len(tags) == 0 {
		return ErrParseFieldTagEmpty
	}

	field := newField(f.Names[0].Name, fType, tags)

	err = field.validate()
	if err != nil {
		return nil
	}

	s.Fields = append(s.Fields, field)

	return nil
}

func (s *ParsedStruct) grepTags() map[TagType]struct{} {
	m := make(map[TagType]struct{})

	// add all tags found for field
	for _, f := range s.Fields {
		for _, t := range f.Tags {
			m[t.Type] = struct{}{}
		}
	}

	return m
}
