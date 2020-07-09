package generator

import (
	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
	"os"
	"path/filepath"
	"strings"
)

type templateData struct {
	PackageName string
	Structs     []parser.ParsedStruct
	Funcs       map[string][]interface{}
}

func newTemplateData(packageName string, structs []parser.ParsedStruct) templateData {
	return templateData{
		PackageName: packageName,
		Structs:     structs,
		Funcs:       genFuncsMap(),
	}
}

func Generate(filePath, suffix string, data *parser.ParsedData) error {
	tmpl := newTemplateData(data.PackageName, data.Structs)

	file, err := newFile(filePath, suffix)
	if err != nil {
		return err
	}

	return generateTemplates(tmpl, file)
}

func newFile(filePath, suffix string) (*os.File, error) {
	dir, file := filepath.Split(filePath)

	pos := strings.LastIndex(file, ".go")
	return os.Create(filepath.Join(dir, file[:pos]+"_"+suffix+".go"))
}
