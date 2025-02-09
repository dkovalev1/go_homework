package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestReadDir(t *testing.T) {
	// Place your code here
	env, err := ReadDir("./testdata/env")
	require.NoError(t, err)
	require.Equal(t, "", env["EMPTY"].Value)
	require.True(t, env["EMPTY"].NeedRemove)

	require.Equal(t, "bar", env["BAR"].Value)
	require.False(t, env["BAR"].NeedRemove)
}

func TestReadLine(t *testing.T) {
	// Place your code here
	value, err := readLine("./testdata/env/FOO")
	require.NoError(t, err)
	require.Equal(t, "   foo\nwith new line", value)
}

func TestReadLineNoFile(t *testing.T) {
	_, err := readLine("./testdata/env/noSuchFile")
	require.Error(t, err)
}

func TestReadLineEmptyFile(t *testing.T) {
	value, err := readLine("./testdata/env/EMPTY")
	require.NoError(t, err)
	require.Equal(t, "", value)
}
