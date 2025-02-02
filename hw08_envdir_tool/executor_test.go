package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestRunCmd(t *testing.T) {
	ret := RunCmd([]string{"ls"}, nil)
	require.Equal(t, 0, ret)

	ret = RunCmd([]string{"cat", "/wrong/path"}, nil)
	require.NotEqual(t, 0, ret)

	ret = RunCmd([]string{"uname", "-a"}, nil)
	require.Equal(t, 0, ret)

	ret = RunCmd([]string{"bash", "-c", "echo \"testvar=$MYTEST\""}, Environment{"MYTEST": {"test", false}})
	require.Equal(t, 0, ret)
}
