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

	packageName, structs, err := parse(os.Args[1])
	if err != nil {
		log.Fatalln("parser:", err)
	}

	if err := generate(packageName, structs); err != nil {
		log.Fatalln("generator:", err)
	}
}
