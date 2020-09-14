package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestTelnetClient(
	addr string,
	in io.ReadCloser,
	out io.Writer,
	info io.Writer,
) (TelnetClient, error) {

	timeout, err := time.ParseDuration("10s")
	if err != nil {
		return nil, err
	}

	if info == nil {
		info = os.Stderr
	}

	return &Client{
		address:    addr,
		timeout:    timeout,
		in:         in,
		out:        out,
		info:       info,
		cancelFunc: func() {},
	}, nil
}

func TestTelnetClient(t *testing.T) {
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

			client, err := newTestTelnetClient(
				l.Addr().String(),
				ioutil.NopCloser(in),
				out,
				ioutil.Discard)
			require.NoError(t, err)
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

	t.Run("test EOF", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			pr, pw := io.Pipe()
			info := &bytes.Buffer{}

			client, err := newTestTelnetClient(
				l.Addr().String(),
				pr,
				ioutil.Discard,
				info)
			require.NoError(t, err)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			err = pw.Close()
			require.NoError(t, err)
			err = client.Send()
			require.NoError(t, err)
			expected := fmt.Sprintf(">> connected to %s\n>> EOF\n", l.Addr().String())

			require.Equal(t, expected, info.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
		}()

		wg.Wait()
	})
}
