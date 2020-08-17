package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/internal/generator"
	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
)

const (
	tagToken       = "validate"
	filenameSuffix = "validation_generated"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected as argument a go file")
		os.Exit(1)
	}

	filePath := os.Args[1]

	if fi, err := os.Stat(filePath); err != nil || !fi.Mode().IsRegular() {
		log.Fatalln("wrong file")
	}

	if ext := filepath.Ext(filePath); ext != ".go" {
		log.Fatalln("wrong file extension")
	}

	data, err := parser.Parse(filePath, tagToken)
	if err != nil {
		log.Fatalln("parser:", err)
	}

	if err := generator.Generate(filePath, filenameSuffix, data); err != nil {
		log.Fatalln("generator:", err)
	}
}
