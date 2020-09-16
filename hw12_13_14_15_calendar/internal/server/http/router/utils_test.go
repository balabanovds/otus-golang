package router

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestSplitPath(t *testing.T) {
	tests := []struct {
		path string
		tail string
		head string
	}{
		{"", "/", ""},
		{"/", "/", ""},
		{"foo", "/", "foo"},
		{"/foo/bar", "foo", "/bar"},
		{"/foo/bar/", "foo", "/bar"},
	}

	for _, tst := range tests {
		t.Run(tst.path, func(t *testing.T) {
			tail, head := splitPath(tst.path)
			require.Equal(t, tst.tail, tail)
			require.Equal(t, tst.head, head)
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expInt       int
		expTail      string
		expErrSubstr string
	}{
		{
			name:         "simple ID",
			path:         "/123",
			expInt:       123,
			expTail:      "",
			expErrSubstr: "",
		},
		{
			name:         "error",
			path:         "/123a",
			expInt:       0,
			expErrSubstr: "invalid syntax",
		},
		{
			name:         "middle path ID",
			path:         "/123/foo",
			expInt:       123,
			expTail:      "/foo",
			expErrSubstr: "",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			id, tail, err := parseInt(tst.path)
			if tst.expErrSubstr != "" && err != nil {
				require.True(t, strings.Contains(err.Error(), tst.expErrSubstr))
				return
			}
			require.NoError(t, err)
			require.Equal(t, tst.expInt, id)
			require.Equal(t, tst.expTail, tail)
		})
	}
}
