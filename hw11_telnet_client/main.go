package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

const (
	defaultTimeout = 10
)

var timeout time.Duration

func init() {
	pflag.DurationVar(&timeout, "timeout", defaultTimeout*time.Second, "connection timeout in seconds")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [options] host port\n", os.Args[0])
		pflag.PrintDefaults()
	}
}

func main() {
	pflag.Parse()

	if len(pflag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "error: expected two arguments")
		pflag.Usage()
		os.Exit(1)
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	client := NewTelnetClient(
		net.JoinHostPort(pflag.Arg(0), pflag.Arg(1)),
		timeout,
		os.Stdin,
		os.Stdout,
		cancelFn,
	)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := client.Receive(); err != nil {
			cancelFn()
		}
	}()

	go func() {
		if err := client.Send(); err != nil {
			cancelFn()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	select {
	case <-sigCh:
		cancelFn()
	case <-ctx.Done():
		close(sigCh)
	}
}
