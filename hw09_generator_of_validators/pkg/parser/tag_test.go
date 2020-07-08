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
		expected []tag
	}

	tests := []tc{
		{
			name: "simple",
			in:   `json:"id" validate:"len:36"`,
			expected: []tag{
				{
					tType: tLen,
					value: "36",
				},
			},
		},
		{
			name: "with and",
			in:   `validate:"len:36|min:3"`,
			expected: []tag{
				{
					tType: tLen,
					value: "36",
				},
				{
					tType: tMin,
					value: "3",
				},
			},
		},
		{
			name: "complex",
			in:   `validate:"regexp:\\d+|len:20"`,
			expected: []tag{
				{
					tType: tRegexp,
					value: `\\d+`,
				},
				{
					tType: tLen,
					value: "20",
				},
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			got := parseTags(tst.in, tagToken)

			require.Equal(t, tst.expected, got)
		})

	}
}
