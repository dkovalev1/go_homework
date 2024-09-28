package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

// Adding one more layer to allow writing tests
func MyReverse(s string) string {
	return reverse.String(s)
}

func main() {
	in := "Hello, OTUS!"
	out := MyReverse(in)
	fmt.Printf("%s", out)
}
