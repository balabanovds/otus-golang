package main

import (
	"fmt"
	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
	"log"
	"os"
)

const tagToken = "validate"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected as argument a go file")
		os.Exit(1)
	}

	filepath := os.Args[1]

	if fi, err := os.Stat(filepath); err != nil || !fi.Mode().IsRegular() {
		log.Fatalln("wrong file")
	}

	data, err := parser.Parse(filepath, tagToken)
	if err != nil {
		log.Fatalln("parser:", err)
	}

	if err := generate(filepath, data); err != nil {
		log.Fatalln("generator:", err)
	}
}
