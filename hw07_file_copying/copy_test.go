package main

import (
	"crypto/md5"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inputFile := path.Join("testdata", "input.txt")

	type testCase struct {
		name          string
		inputFile     string
		referenceFile string
		offset        int64
		limit         int64
		expectedError error
	}

	tests := []testCase{
		{
			name:          "offset 0 limit 0",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset0_limit0.txt"),
			offset:        0,
			limit:         0,
		},
		{
			name:          "offset 0 limit 10",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset0_limit10.txt"),
			offset:        0,
			limit:         10,
		},
		{
			name:          "offset 0 limit 1000",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset0_limit1000.txt"),
			offset:        0,
			limit:         1000,
		},
		{
			name:          "offset 0 limit 10000",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset0_limit10000.txt"),
			offset:        0,
			limit:         10000,
		},
		{
			name:          "offset 100 limit 1000",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset100_limit1000.txt"),
			offset:        100,
			limit:         1000,
		},
		{
			name:          "offset 6000 limit 1000",
			inputFile:     inputFile,
			referenceFile: refFile(t, "out_offset6000_limit1000.txt"),
			offset:        6000,
			limit:         1000,
		},
		{
			name:          "copy from irregular file",
			inputFile:     "/dev/urandom",
			expectedError: ErrUnsupportedFile,
		},
		{
			name:          "offset is greater than file size",
			inputFile:     inputFile,
			offset:        100000,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "not found source file",
			inputFile:     inputFile + "404",
			expectedError: ErrFileNotExists,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			dst, err := ioutil.TempFile("/tmp", "godd")
			assert.NoError(t, err)
			defer os.Remove(dst.Name())

			err = Copy(tst.inputFile, dst.Name(), tst.offset, tst.limit)
			if err == nil {
				refMD5 := md5sum(t, tst.referenceFile)
				dstMD5 := md5sum(t, dst.Name())

				assert.Equal(t, refMD5, dstMD5)
			} else {
				require.EqualError(t, err, tst.expectedError.Error())
			}

			err = dst.Close()
			assert.NoError(t, err)
		})
	}

}

func refFile(t *testing.T, filename string) string {
	t.Helper()
	return path.Join("testdata", filename)
}

func md5sum(t *testing.T, filename string) []byte {
	t.Helper()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}
