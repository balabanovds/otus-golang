package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var (
	ErrParseFailed        = errors.New("parser: failed to parse")
	ErrParseStructs404    = errors.New("parser: structs not found")
	ErrParseFieldTagEmpty = errors.New("parser: filed tag empty")
)

type ParsedData struct {
	PackageName string
	tagToken    string
	fSet        *token.FileSet
	Structs     []ParsedStruct
	Idents      map[string]string
	Tags        map[TagType]struct{} // this needed to generate imports
}

func newData(packageName, tagToken string) *ParsedData {
	return &ParsedData{
		PackageName: packageName,
		tagToken:    tagToken,
		Idents:      make(map[string]string),
		Tags:        make(map[TagType]struct{}),
	}
}

func (d *ParsedData) addIdent(name string, ident *ast.Ident) {
	d.Idents[name] = ident.Name
}

func (d *ParsedData) addStruct(name string, s *ast.StructType) error {
	vs := newStruct(name)

	for _, f := range s.Fields.List {
		err := vs.addField(*d, f)
		if err != nil {
			// just skip field if tags empty
			if errors.Is(err, ErrParseFieldTagEmpty) {
				continue
			}
			if errors.Is(err, ErrParseFieldUnknownType) {
				// we just print here a warning that we faced with unknown field type
				log.Printf("at %s: %v\n", d.fSet.Position(f.Pos()), err)
				continue
			}
			// if smth weird
			return err
		}

	}

	if len(vs.Fields) != 0 {
		d.Structs = append(d.Structs, vs)
		for tag := range vs.grepTags() {
			d.Tags[tag] = struct{}{}
		}
	}

	return nil
}

// Parse go file for Idents and struct fields containing tagToken
func Parse(fileName, tagToken string) (*ParsedData, error) {
	fset := token.NewFileSet()
	parseFile, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return nil, err
	}

	data := newData(parseFile.Name.Name, tagToken)
	data.fSet = fset

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

	if len(data.Structs) == 0 {
		return nil, ErrParseStructs404
	}

	return data, nil
}
