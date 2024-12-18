package errordefcli

import (
	"errors"
	"fmt"
)

type CLIError struct {
	exitCode int
	message  string
}

func (t *CLIError) Error() string {
	return t.message
}

func (t *CLIError) ExitCode() int {
	return t.exitCode
}

func NewCLIError(exitCode int, message string) *CLIError {
	return &CLIError{
		exitCode: exitCode,
		message:  message,
	}
}

func NewCLIErrorf(exitCode int, message string, args ...any) *CLIError {
	return &CLIError{
		exitCode: exitCode,
		message:  fmt.Sprintf(message, args...),
	}
}

func AsCLIError(err error) (*CLIError, bool) {
	var cliError *CLIError
	if !errors.As(err, &cliError) {
		return nil, false
	}
	return cliError, true
}
