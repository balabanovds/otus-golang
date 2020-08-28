package main

import (
	"bufio"
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
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err == nil {
		c.conn = conn
		fmt.Fprintf(os.Stderr, ">> connected to %s\n", c.address)
	}

	return err
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := c.conn.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			fmt.Fprintln(os.Stderr, ">> connection closed by remote peer")

			return err
		}
	}
	fmt.Fprintln(os.Stderr, ">> EOF")

	return c.Close()
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.out, c.conn)

	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
