package main

import (
	"bufio"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readLine(fileName string) (string, error) {
	/* this is a learning project and although it was possible to use ReadFile,
	** but there is also a chance to try out the bufio package to read a single line in text mode
	 */
	fp, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	value := ""
	if scanner.Scan() {
		value = scanner.Text()
	}
	value = strings.TrimRight(value, " \t")
	value = strings.ReplaceAll(value, "\x00", "\n")

	return value, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		if !file.IsDir() {
			filePath := dir + "/" + file.Name()
			content, err := readLine(filePath)
			if err != nil {
				return nil, err
			}

			env[file.Name()] = EnvValue{
				Value:      content,
				NeedRemove: len(content) == 0,
			}
		}
	}
	return env, nil
}
