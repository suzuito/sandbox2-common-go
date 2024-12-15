package errordefcli

import "errors"

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

func AsCLIError(err error) (*CLIError, bool) {
	var cliError *CLIError
	if !errors.As(err, &cliError) {
		return nil, false
	}
	return cliError, true
}
