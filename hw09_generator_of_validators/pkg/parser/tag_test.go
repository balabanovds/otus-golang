package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const tagToken = "validate"

func TestParseTags(t *testing.T) {
	type tc struct {
		name     string
		in       string
		expected []Tag
	}

	tests := []tc{
		{
			name: "simple",
			in:   `json:"id" validate:"len:36"`,
			expected: []Tag{
				{
					Type:  TagLen,
					Value: "36",
				},
			},
		},
		{
			name: "with and",
			in:   `validate:"len:36|min:3"`,
			expected: []Tag{
				{
					Type:  TagLen,
					Value: "36",
				},
				{
					Type:  TagMin,
					Value: "3",
				},
			},
		},
		{
			name: "complex",
			in:   `validate:"regexp:\\d+|len:20"`,
			expected: []Tag{
				{
					Type:  TagRegexp,
					Value: `\\d+`,
				},
				{
					Type:  TagLen,
					Value: "20",
				},
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			got := parseTags(tst.in, tagToken, FString)

			require.Equal(t, tst.expected, got)
		})

	}
}
