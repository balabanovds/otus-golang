package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("expected as argument a go file")
		os.Exit(1)
	}

	filepath := os.Args[1]

	if fi, err := os.Stat(filepath); err != nil || !fi.Mode().IsRegular() {
		log.Fatalln("wrong file")
	}

	data, err := parse(filepath)
	if err != nil {
		log.Fatalln("parser:", err)
	}

	if err := generate(filepath, data); err != nil {
		log.Fatalln("generator:", err)
	}
}
