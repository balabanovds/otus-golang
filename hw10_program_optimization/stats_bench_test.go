package hw10_program_optimization //nolint:golint,stylecheck
import (
	"archive/zip"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStat(b *testing.B) {
	data, closer := readZip()
	defer closer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "com")
	}
}

func BenchmarkGetUsers(b *testing.B) {
	data, closer := readZip()
	defer closer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = getUsers(data)
	}
}

func BenchmarkCountDomains(b *testing.B) {
	data, closer := readZip()
	defer closer()

	users, err := getUsers(data)
	require.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = countDomains(users, "com")
	}
}

func readZip() (io.Reader, func()) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		log.Fatalln(err)
	}

	data, err := r.File[0].Open()
	return data, func() {
		r.Close()
	}
}
