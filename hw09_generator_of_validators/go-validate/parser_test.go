package main

import (
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"testing"
)

func TestInspectForStruct(t *testing.T) {
	type tc struct {
		name     string
		in       string
		expected []validatedStruct
	}

	tests := []tc{
		{
			name: "simple",
			in: "package models\n" +
				"type User struct {\n" +
				"ID	int `json:\"user_id\"`\n" +
				"Name string `validate:\"min:10|max:50\"`\n}",
			expected: []validatedStruct{
				{
					name:  "User",
					short: "u",
					fields: []field{
						{
							name: "Name",
							typeStr: fieldType{
								key:   fString,
								value: "string",
							},
							tags: []tag{
								{
									tType: tMin,
									value: "10",
								},
								{
									tType: tMax,
									value: "50",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tst.in, 0)
			require.NoError(t, err)
			got := inspectForStructs(file)
			require.Equal(t, tst.expected, got)
		})
	}
}
