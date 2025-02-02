package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func help() {
	fmt.Printf("Go files copy program\n")
	fmt.Printf("  Usage:\n")
	flag.PrintDefaults()
}

func checkArguments() error {
	stat, err := os.Stat(from)
	if err != nil {
		return err
	}

	if offset > stat.Size() {
		return errors.New("offset is greater that the size of file")
	}
	return nil
}

func main() {
	flag.Parse()
	err := checkArguments()
	if err != nil {
		fmt.Println(err)
		help()
		os.Exit(1)
	}

	err = Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
