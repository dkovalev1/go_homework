package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint
)

type BufferCloser struct {
	*bytes.Buffer
	closed bool
}

func (b *BufferCloser) Close() error {
	b.closed = true
	return nil
}

func (b *BufferCloser) Read(p []byte) (n int, err error) {
	if b.closed {
		return 0, io.EOF
	}
	return b.Buffer.Read(p)
}

func NewBufferCloser(data string) *BufferCloser {
	return &BufferCloser{Buffer: bytes.NewBufferString(data)}
}

func TestTelnetClient(t *testing.T) {
	t.Run("Connect", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		c := NewTelnetClient(l.Addr().String(), 0, nil, nil)
		require.NoError(t, c.Connect())
	})

	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("EOF", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := NewBufferCloser("hello world\n")
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, in, out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			require.NoError(t, client.Send())
			require.NoError(t, client.Receive())

			require.NoError(t, in.Close())

			err = client.Send()
			require.ErrorIs(t, err, io.EOF)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			nr, err := conn.Read(request)
			require.NoError(t, err)
			require.NotEqual(t, 0, nr)

			nw, err := conn.Write(request)
			require.NoError(t, err)
			require.NotEqual(t, 0, nw)
		}()

		wg.Wait()
	})
}
