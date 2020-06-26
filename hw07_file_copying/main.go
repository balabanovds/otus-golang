package main

import (
	"flag"
	"log"
	"os"
)

var (
	from, to                 string
	limit, offset, chunkSize int64
	progress                 bool
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from (mandatory)")
	flag.StringVar(&to, "to", "", "file to write to (mandatory)")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
	flag.Int64Var(&chunkSize, "bs", 1, "chuck of bytes to copy within one iteration")
	flag.BoolVar(&progress, "progress", false, "show progress in stdout")
}

func main() {
	flag.Parse()

	if len(from) == 0 || len(to) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := Copy(from, to, offset, limit); err != nil {
		log.Fatalln(err)
	}
}
