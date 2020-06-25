package main

import (
	"io"
	"os"
)

type dst struct {
	file   *os.File
	offset int64
}

func newDst(filename string) (*dst, error) {
	fileInfo, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			fh, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			return &dst{
				file: fh,
			}, nil
		}
		return nil, err
	}

	if !fileInfo.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}

	fh, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	if !append {
		_ = fh.Truncate(0)
	}

	offset := fileInfo.Size()

	_, err = fh.Seek(offset, io.SeekCurrent)
	if err != nil {
		fh.Close()
		return nil, err
	}

	return &dst{
		file:   fh,
		offset: offset,
	}, nil
}
