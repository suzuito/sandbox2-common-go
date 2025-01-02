package fakecmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

func TestFakeCMDFakeCMD(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	dirPath := domains.DirPathFakeCommand(filepath.Dir(filePathBin))

	testCases := []e2ehelpers.CLITestCase{
		{
			Desc: "ng - command's behaviors file does not exist",
			ExpectedStderr: fmt.Sprintf(
				"FAKE_CMD_ERROR failed to read behavior file: %s: open %s: no such file or directory",
				dirPath.FilePathBehaviors(),
				dirPath.FilePathBehaviors(),
			),
			ExpectedExitCode: 127,
		},
		{
			Desc: "ng - command's behaviors file is not json",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteFile(
					dirPath.FilePathBehaviors(),
					[]byte("a"),
				)

				return nil
			},
			ExpectedStderr: fmt.Sprintf(
				"FAKE_CMD_ERROR failed to unmarshal behavior file: %s: invalid character 'a' looking for beginning of value",
				dirPath.FilePathBehaviors(),
			),
			ExpectedExitCode: 127,
		},
		{
			Desc: "ng - fake no command's behaviors",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{},
				)

				return nil
			},
			ExpectedStderr:   "FAKE_CMD_ERROR no behaviors",
			ExpectedExitCode: 127,
		},
		{
			Desc: "ok - fake a command's behavior",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout",
								Stderr:   "this is a test stderr",
								ExitCode: 10,
							},
						},
					},
				)

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStdout:   "this is a test stdout",
			ExpectedStderr:   "this is a test stderr",
			ExpectedExitCode: 10,
		},
		{
			Desc: "ok - fake a command's behaviors and no executions are done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout1",
								ExitCode: 11,
							},
						},
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout2",
								ExitCode: 12,
							},
						},
					},
				)

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStdout:   "this is a test stdout1",
			ExpectedExitCode: 11,
		},
		{
			Desc: "ok - fake a command's behaviors and already first execution is done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout1",
								ExitCode: 11,
							},
						},
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout2",
								ExitCode: 12,
							},
						},
					},
				)

				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathState(),
					&domains.State{
						ExecutedHistories: domains.ExecutedHistories{
							{},
						},
					},
				)

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStdout:   "this is a test stdout2",
			ExpectedExitCode: 12,
		},
		{
			Desc: "ng - fake a command's behaviors and already all executions are done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout1",
								ExitCode: 11,
							},
						},
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout2",
								ExitCode: 12,
							},
						},
					},
				)

				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathState(),
					&domains.State{
						ExecutedHistories: domains.ExecutedHistories{
							{}, {},
						},
					},
				)

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStderr:   "FAKE_CMD_ERROR all expected executions are done: expected=2 histories=2",
			ExpectedExitCode: 127,
		},
		{
			Desc: "ng - cannot get lock",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout1",
								ExitCode: 11,
							},
						},
					},
				)

				e2ehelpers.MustWriteJSONFile(dirPath.FilePathProcessing(), struct{}{})

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStderr: fmt.Sprintf(
				"FAKE_CMD_ERROR failed to get lock: %s: processing file already exists",
				dirPath.FilePathProcessing(),
			),
			ExpectedExitCode: 127,
		},
		{
			Desc: "ng - state file is not json",
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{
						{
							Type: domains.BehaviorTypeStdoutStderrExitCode,
							BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
								Stdout:   "this is a test stdout",
								Stderr:   "this is a test stderr",
								ExitCode: 10,
							},
						},
					},
				)

				e2ehelpers.MustWriteFile(
					dirPath.FilePathState(),
					[]byte("a"),
				)

				return nil
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) error {
				return errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				)
			},
			ExpectedStderr: fmt.Sprintf(
				"FAKE_CMD_ERROR failed to unmarshal state file: %s: invalid character 'a' looking for beginning of value",
				dirPath.FilePathState(),
			),
			ExpectedExitCode: 127,
		},
		{
			Desc:             "ng - no error logs",
			Envs:             []string{"FAKECMD_ERROR_LOG=discard"},
			ExpectedExitCode: 127,
		},
	}
	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
