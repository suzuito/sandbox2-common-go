package e2ehelpers_test

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func Test_ExternalCommandFaker_Add_cannot_dup_filepath(t *testing.T) {
	f := e2ehelpers.ExternalCommandFaker{}
	defer f.Cleanup()

	require.NoError(t, f.Add(&e2ehelpers.ExternalCommandBehavior{
		FilePath: "/tmp/hoge001",
	}))

	err := f.Add(&e2ehelpers.ExternalCommandBehavior{
		FilePath: "/tmp/hoge001",
	})
	require.Error(t, err)
	require.EqualError(t, err, "fake command already exists: /tmp/hoge001")
}

func Test_ExternalCommandFaker(t *testing.T) {
	testCases := []struct {
		desc             string
		input            e2ehelpers.ExternalCommandBehavior
		expectedExitCode int
		expectedStdout   string
		expectedStderr   string
	}{
		{
			desc: "ok - stdout and stderr",
			input: e2ehelpers.ExternalCommandBehavior{
				FilePath: "/tmp/hoge001",
				Stdout:   "stdout",
				Stderr:   "stderr",
			},
			expectedStdout: "stdout\n",
			expectedStderr: "stderr\n",
		},
		{
			desc: "ok - stdout only",
			input: e2ehelpers.ExternalCommandBehavior{
				FilePath: "/tmp/hoge001",
				Stdout:   "stdout",
			},
			expectedStdout: "stdout\n",
		},
		{
			desc: "ok - stderr only",
			input: e2ehelpers.ExternalCommandBehavior{
				FilePath: "/tmp/hoge001",
				Stderr:   "stderr",
			},
			expectedStderr: "stderr\n",
		},
		{
			desc: "ok - failed command",
			input: e2ehelpers.ExternalCommandBehavior{
				FilePath: "/tmp/hoge001",
				ExitCode: 1,
				Stdout:   "stdout",
				Stderr:   "stderr",
			},
			expectedExitCode: 1,
			expectedStdout:   "stdout\n",
			expectedStderr:   "stderr\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			f := e2ehelpers.ExternalCommandFaker{}
			defer f.Cleanup()

			require.NoError(t, f.Add(&tC.input))

			actualStdout := bytes.NewBufferString("")
			actualStderr := bytes.NewBufferString("")

			fakeCmd := exec.Command("/tmp/hoge001")
			fakeCmd.Stdout = actualStdout
			fakeCmd.Stderr = actualStderr

			err := fakeCmd.Run()
			var exiterr *exec.ExitError
			if !errors.As(err, &exiterr) {
				require.NoError(t, err)
			}
			require.Equal(t, tC.expectedExitCode, fakeCmd.ProcessState.ExitCode())
			require.Equal(t, tC.expectedStdout, actualStdout.String())
			require.Equal(t, tC.expectedStderr, actualStderr.String())
		})
	}
}
