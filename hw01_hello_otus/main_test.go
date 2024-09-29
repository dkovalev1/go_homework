package main

import (
	"testing"
)

func TestMyReverse(t *testing.T) {
	src := "Hello"
	expected := "olleH"
	actual := MyReverse(src)

	if expected != actual {
		t.Fatalf("Expected: %s, got: %s", expected, actual)
	}
}
