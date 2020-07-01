package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type exitCode int

const (
	exitOK exitCode = iota
	exitErrGeneral
	exitErrCmdCannotExecute exitCode = iota + 124
	exitErrCmdNotFound
	exitErrInvalidArg
)

func TestRunCmd(t *testing.T) {
	type tc struct {
		name               string
		command            string
		args               []string
		env                Environment
		expectedReturnCode exitCode
	}

	tests := []tc{
		{
			name:               "General error",
			command:            "/bin/bash",
			args:               []string{"-c", "touch /etc"},
			expectedReturnCode: exitErrGeneral,
		},
		{
			name:               "Permission denied",
			command:            "/bin/bash",
			args:               []string{"-c", "/dev/null"},
			expectedReturnCode: exitErrCmdCannotExecute,
		},
		{
			name:               "Command not found",
			command:            "/bin/bash",
			args:               []string{"-c", "some_command_that_does_not_exists"},
			expectedReturnCode: exitErrCmdNotFound,
		},
		{
			name:               "Command OK",
			command:            "pwd",
			expectedReturnCode: exitOK,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			cmd := []string{tst.command}
			cmd = append(cmd, tst.args...)

			code := RunCmd(cmd, tst.env)

			assert.Equal(t, int(tst.expectedReturnCode), code)
		})
	}
}
