package hw10_program_optimization_test //nolint:golint,stylecheck
import (
	"archive/zip"

	"io"
	"log"
	"testing"

	stats "github.com/balabanovds/otus-golang/hw10_program_optimization"
)

func BenchmarkGetDomainStat(b *testing.B) {
	data, closer := readZip()
	defer closer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = stats.GetDomainStat(data, "com")
	}
}

func readZip() (io.Reader, func()) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := r.File[0].Open()
	return data, func() {
		_ = r.Close()
	}
}
