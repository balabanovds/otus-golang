package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/balabanovds/otus-golang/hw07_file_copying/pkg/pb"
)

var (
	ErrFileNotExists         = errors.New("file not exists")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	src, err := newSrc(fromPath, offset, limit)
	if err != nil {
		return err
	}
	defer src.file.Close()

	dst, err := newDst(toPath)
	if err != nil {
		return err
	}
	defer dst.file.Close()

	position := offset

	fmt.Printf("Copying from '%s', to '%s'\n", fromPath, toPath)
	fmt.Printf("Offset %d bytes, limit %d bytes\n", offset, limit)
	fmt.Printf("Consider to copy %d bytes\n", src.copySize)

	p := pb.New(src.copySize)
	p.Start()

	for {
		n, err := io.CopyN(dst.file, src.file, chunkSize)

		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		position += n
		if position > offset+src.copySize {
			break
		}

		p.Add(n)
		time.Sleep(100 * time.Microsecond)
	}

	p.Finish()

	fmt.Println("Copy complete")

	return nil
}
