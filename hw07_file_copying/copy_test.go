package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

const inputFile = "testdata/input.txt"

func compareFiles(file1, file2 string) bool {
	fp1, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer fp1.Close()

	fp2, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer fp2.Close()

	bufsize := 256
	for {
		b1 := make([]byte, bufsize)
		_, err1 := fp1.Read(b1)
		b2 := make([]byte, bufsize)
		_, err2 := fp2.Read(b2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			}

			// File error, will treate is as a diff for the sake of test, but could be proper error reporting here
			return false
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func TestCopy_0_0(t *testing.T) {
	testCopy(t, 0, 0)
}

func TestCopy_0_10(t *testing.T) {
	testCopy(t, 0, 10)
}

func TestCopy_0_1000(t *testing.T) {
	testCopy(t, 0, 1000)
}

func TestCopy_0_10000(t *testing.T) {
	testCopy(t, 0, 10000)
}

func TestCopy_100_1000(t *testing.T) {
	testCopy(t, 100, 1000)
}

func TestCopy_6000_1000(t *testing.T) {
	testCopy(t, 6000, 1000)
}

func testCopy(t *testing.T, offset, limit int64) {
	t.Helper()

	test := fmt.Sprintf("testdata\\test_offset%d_limit%d.txt", offset, limit)
	etalon := fmt.Sprintf("testdata\\out_offset%d_limit%d.txt", offset, limit)

	err := Copy(inputFile, test, offset, limit)
	require.NoError(t, err)

	defer func() {
		os.Remove(test)
	}()

	require.True(t, compareFiles(test, etalon))
}
