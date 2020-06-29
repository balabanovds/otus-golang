package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		key := file.Name()
		if strings.Contains(key, "=") {
			continue
		}

		// try to read file
		filename := path.Join(dir, key)
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("failed to open file %s\n", filename)
			continue
		}

		value, err := readValue(file)
		if err != nil {
			continue
		}

		env[key] = value
	}

	return env, nil
}

func readValue(r io.Reader) (string, error) {
	br := bufio.NewReader(r)

	// read first line until '\n'
	line, _, err := br.ReadLine()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", nil
		}
		return "", err
	}

	line = bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))

	return string(line), nil
}

func prepareCmdEnv(env Environment) []string {
	for k, v := range env {
		if v == "" {
			os.Unsetenv(k)
			continue
		}

		os.Setenv(k, v)
	}

	return os.Environ()
}
