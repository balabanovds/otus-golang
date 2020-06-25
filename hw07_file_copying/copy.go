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

	buf := make([]byte, chunkSize)

	var n int
	position := offset

	fmt.Printf("Copying from '%s', to '%s'\n", fromPath, toPath)
	fmt.Printf("Offset %d bytes, limit %d bytes\n", offset, limit)
	fmt.Printf("Consider to copy %d bytes\n", src.copySize)

	p := pb.New(src.copySize)
	if progress {
		p.Start()
	}

	for {
		position += int64(n)
		if position > offset+src.copySize {
			break
		}

		n, err = src.file.ReadAt(buf, position)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err = dst.file.Write(buf[:n]); err != nil {
			return err
		}

		if progress {
			p.Add(int64(n))
			time.Sleep(time.Millisecond)
		}
	}

	if progress {
		p.Finish()
	}

	fmt.Println("Copy complete")

	return nil
}
