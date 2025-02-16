package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *TelnetClientImpl) Close() error {
	return c.conn.Close()
}

func (c *TelnetClientImpl) Send() error {
	buf := make([]byte, 1024)
	n, err := c.in.Read(buf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "EOF\n")
		}
		return err
	}
	_, err = c.conn.Write(buf[:n])
	return err
}

func (c *TelnetClientImpl) Receive() error {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "connection closed\n")
		}
		return err
	}
	_, err = c.out.Write(buf[:n])
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
