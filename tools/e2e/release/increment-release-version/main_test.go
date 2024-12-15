package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	testCases := []struct {
		desc             string
		args             []string
		setup            func() error
		expectedExitCode int
		expectedStdout   string
		expectedStderr   string
	}{
		{
			desc: "ok",
			setup: func() error {
				f, err := os.Create("/tmp/e2e001.sh")
				if err != nil {
					return err
				}
				defer f.Close()

				f.Chmod(0755)

				fmt.Fprintf(f, "#!/bin/sh\n")
				fmt.Fprintf(f, "echo 'v1.1.2'\n")
				fmt.Fprintf(f, "echo 'v1.1.3'\n")
				fmt.Fprintf(f, "echo 'v1.1.4'\n")

				return nil
			},
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
			},
			expectedExitCode: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()

			cmd := exec.CommandContext(
				ctx,
				filePathBin,
				tC.args...,
			)

			cmd.Env = append(
				os.Environ(),
				fmt.Sprintf("E2E_TEST_ID=%s", uuid.New()),
				"GITHUB_HTTP_CLIENT_FAKE_SCHEME=http",
				"GITHUB_HTTP_CLIENT_FAKE_HOST=localhost:8080",
			)

			stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			if err := tC.setup(); err != nil {
				require.Error(t, err, err.Error())
			}

			err := cmd.Run()
			var exiterr *exec.ExitError
			if !errors.As(err, &exiterr) {
				require.NoError(t, err, fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(tC.args, " "), err.Error()))
			}

			assert.Equal(t, tC.expectedExitCode, cmd.ProcessState.ExitCode())
			assert.Equal(t, tC.expectedStdout, stdout.String())
			assert.Equal(t, tC.expectedStderr, stderr.String())
		})
	}
}
