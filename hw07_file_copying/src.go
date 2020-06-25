package main

import (
	"io"
	"os"
)

type src struct {
	file     *os.File
	name     string
	copySize int64
	offset   int64
}

func newSrc(filename string, offset, limit int64) (*src, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExists
		}
		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}

	copySize, err := getLength(fileInfo, offset, limit)
	if err != nil {
		return nil, err
	}

	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	_, err = fh.Seek(offset, io.SeekStart)
	if err != nil {
		fh.Close()
		return nil, err
	}

	return &src{
		file:     fh,
		name:     filename,
		copySize: copySize,
		offset:   offset,
	}, nil
}

func getLength(fi os.FileInfo, offset, limit int64) (int64, error) {
	if offset >= fi.Size() {
		return 0, ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		return fi.Size() - offset, nil
	}

	if limit+offset > fi.Size() {
		return fi.Size() - offset, nil
	}

	return limit - 1, nil
}
