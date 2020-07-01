// +build !windows

package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	BAR string = "bar"
	FOO string = `   foo
with new line`
	HELLO string = `"hello"`
	UNSET string = ""
)

func TestReadDir(t *testing.T) {
	var expectedMap = Environment{
		"BAR":   BAR,
		"FOO":   FOO,
		"HELLO": HELLO,
		"UNSET": UNSET,
	}
	const dir = "./testdata/env"
	env, err := ReadDir(dir)
	assert.NoError(t, err)

	assert.Len(t, env, 4)
	assert.Equal(t, expectedMap, env)
}

func TestReadValue(t *testing.T) {
	type testCase struct {
		name     string
		in       string
		expected string
		err      error
	}

	tests := []testCase{
		{
			name: "BAR: get only first string",
			in: `bar
PLEASE IGNORE SECOND LINE
`,
			expected: BAR,
		},
		{
			name:     "FOO: replace NUL char with 'new line'",
			in:       "   foo\x00with new line",
			expected: FOO,
		},
		{
			name:     "HELLO: get value with quotes",
			in:       `"hello"`,
			expected: HELLO,
		},
		{
			name:     "UNSET: empty line",
			in:       "",
			expected: UNSET,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			b := bytes.NewBufferString(tst.in)
			got, err := readValue(b)
			require.NoError(t, err)
			assert.Equal(t, tst.expected, got)
		})

	}
}
