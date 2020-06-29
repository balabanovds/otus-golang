package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedMap = Environment{
	"BAR": "bar",
	"FOO": `   foo
with new line`,
	"HELLO": `"hello"`,
	"UNSET": "",
}

func TestReadDir(t *testing.T) {
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
			expected: expectedMap["BAR"],
		},
		{
			name:     "FOO: replace NUL char with 'new line'",
			in:       "   foo\x00with new line",
			expected: expectedMap["FOO"],
		},
		{
			name:     "HELLO: get value with quotes",
			in:       `"hello"`,
			expected: expectedMap["HELLO"],
		},
		{
			name:     "UNSET: empty line",
			in:       "",
			expected: expectedMap["UNSET"],
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			b := bytes.NewBufferString(tst.in)
			got, err := readValue(b)
			assert.NoError(t, err)
			assert.Equal(t, tst.expected, got)
		})

	}
}
