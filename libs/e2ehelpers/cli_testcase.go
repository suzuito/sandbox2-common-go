package e2ehelpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CLITestCase struct {
	Desc             string
	Args             []string
	Setup            func(t *testing.T, testID TestID) error
	ExpectedExitCode int
	ExpectedStdout   string
	ExpectedStderr   string
}

func (c *CLITestCase) Run(
	tt *testing.T,
	filePathBin string,
) bool {
	return tt.Run(c.Desc, func(t *testing.T) {
		ctx := context.Background()

		testID := NewTestID()

		cmd := exec.CommandContext(
			ctx,
			filePathBin,
			c.Args...,
		)

		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("E2E_TEST_ID=%s", testID.String()),
			"GITHUB_HTTP_CLIENT_FAKE_SCHEME=http",
			"GITHUB_HTTP_CLIENT_FAKE_HOST=localhost:8080",
		)

		stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		if c.Setup != nil {
			if err := c.Setup(t, testID); err != nil {
				require.Error(t, err, err.Error())
			}
		}

		err := cmd.Run()
		if err != nil {
			var exiterr *exec.ExitError
			if !errors.As(err, &exiterr) {
				require.NoError(t, err, fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(c.Args, " "), err.Error()))
			}
		}

		assert.Equal(t, c.ExpectedExitCode, cmd.ProcessState.ExitCode())
		assert.Equal(t, c.ExpectedStdout, stdout.String())
		assert.Equal(t, c.ExpectedStderr, stderr.String())
	})
}
