package logger

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestLogger(t *testing.T) {
	t.Run("Debug", func(t *testing.T) {
		l := New("debug")
		l.Debug("test")

		require.Equal(t, l.logLevel, DEBUG)
	})

	t.Run("Info", func(t *testing.T) {
		l := New("info")
		l.Info("test")

		require.Equal(t, l.logLevel, INFO)
	})

	t.Run("Warn", func(t *testing.T) {
		l := New("warn")
		l.Warn("test")

		require.Equal(t, l.logLevel, WARN)
	})

	t.Run("Error", func(t *testing.T) {
		l := New("error")
		l.Error("test")

		require.Equal(t, l.logLevel, ERROR)
	})

	t.Run("Default", func(t *testing.T) {
		l := New("default")
		l.Info("test")

		require.Equal(t, l.logLevel, INFO)
	})
}
