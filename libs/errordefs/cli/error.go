package errordefcli

import (
	"errors"
	"fmt"
)

type cliError struct {
	exitCode uint8
	origin   error
}

func (t *cliError) Error() string {
	return "cli error: " + t.origin.Error()
}

func (t *cliError) Unwrap() error {
	if t.origin == nil {
		return nil
	}
	return t.origin
}

func Errorf(exitCode uint8, message string, args ...any) error {
	return &cliError{
		exitCode: exitCode,
		origin:   fmt.Errorf(message, args...),
	}
}

func asCLIError(err error) (*cliError, bool) {
	var cliError *cliError
	if !errors.As(err, &cliError) {
		return nil, false
	}
	return cliError, true
}

func Code(err error, defaultCode uint8) (uint8, string) {
	cerr, ok := asCLIError(err)
	if !ok {
		return defaultCode, err.Error()
	}

	return cerr.exitCode, cerr.Error()
}
