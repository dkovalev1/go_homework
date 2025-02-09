package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Usage: %s <env-path> <command> [args]\n", os.Args[0])
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}
	envPath := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envPath)
	if err != nil {
		fmt.Printf("failed to prepare env: %v\n", err)
		os.Exit(1)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}
