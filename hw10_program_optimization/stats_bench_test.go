package hw10_program_optimization //nolint:golint,stylecheck
import (
	"archive/zip"
	"io"
	"log"
	"testing"
)

func BenchmarkGetUsers(b *testing.B) {
	data := readZip()

	for i := 0; i < b.N; i++ {
		_, _ = getUsers(data)
	}
}

func readZip() io.Reader {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()


	data, err := r.File[0].Open()
	return data
}
