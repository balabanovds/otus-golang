package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	tests := []test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
		{
			input:    "ф2",
			expected: "фф",
		},
		{
			input:    "😝2",
			expected: "😝😝",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			result, err := Unpack(tt.input)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnpackWithEscape(t *testing.T) {
	//t.Skip() // Remove if task with asterisk completed

	tests := []test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			result, err := Unpack(tt.input)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnpackWithSomeOtherChars(t *testing.T) {
	tests := []test{
		{
			input:    `<3*2`,
			expected: `<<<**`,
		},
		{
			input:    `аб3в4гг`,
			expected: `абббввввгг`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			result, err := Unpack(tt.input)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, result)
		})
	}
}
