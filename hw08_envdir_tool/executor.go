package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func removeEnvVar(s []string, deleteKey string) []string {
	if s == nil {
		return nil
	}

	for k, v := range s {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if key == deleteKey {
			return append(s[:k], s[k+1:]...)
		}
	}
	return s
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	execCmd.Env = os.Environ()
	for key, value := range env {
		if key == "" {
			fmt.Println("Error: environment variable key is empty")
			return 1
		}
		if value.NeedRemove {
			execCmd.Env = removeEnvVar(execCmd.Env, key)
		} else {
			execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}

	if err := execCmd.Run(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}
