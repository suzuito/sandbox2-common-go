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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	testCases := []struct {
		desc             string
		args             []string
		expectedExitCode int
		expectedStdout   string
		expectedStderr   string
	}{
		{
			desc:             "",
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

			stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			fmt.Println("aho")
			fmt.Println(filePathBin)

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
