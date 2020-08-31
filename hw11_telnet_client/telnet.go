package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type Client struct {
	conn       net.Conn
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	info       io.Writer
	cancelFunc context.CancelFunc
}

func NewTelnetClient(
	address string,
	timeout time.Duration,
	in io.ReadCloser,
	out io.Writer,
	cancelFunc context.CancelFunc) TelnetClient {
	return &Client{
		address:    address,
		timeout:    timeout,
		in:         in,
		out:        out,
		info:       os.Stderr,
		cancelFunc: cancelFunc,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err == nil {
		c.conn = conn
		c.printInfo("connected to " + c.address)
	}

	return err
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := c.conn.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			c.printInfo("connection closed by remote peer")

			return err
		}
	}
	c.printInfo("EOF")
	c.cancelFunc()

	return nil
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.conn)

	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) printInfo(msg string) {
	fmt.Fprintln(c.info, ">> "+msg)
}
