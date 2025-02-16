package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10, "Connection timeout in seconds")
}

func usage() {
}

func doTelnet(host string) error {
	client := NewTelnetClient(host, timeout, os.Stdin, os.Stdout)

	// Handle SIGINT (Ctrl-C) and SIGTERM properly
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := client.Connect()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	fmt.Fprintf(os.Stderr, "connected to: %s\n", host)

	exit := make(chan struct{})

	// Send data
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:

				err := client.Send()
				if err != nil {
					exit <- struct{}{}
					return
				}
			}
		}
	}()

	// Receive data
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Receive()
				if err != nil {
					exit <- struct{}{}
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "Received SIGINT, closing connection...")
	case <-exit:
	}
	client.Close()
	return nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 2 {
		usage()
		os.Exit(1)
	}

	host := net.JoinHostPort(flag.Args()[0], flag.Args()[1])

	err := doTelnet(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
