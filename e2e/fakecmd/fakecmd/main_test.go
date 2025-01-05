package fakecmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

func TestFakeCMDFakeCMD(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	dirPath := domains.DirPathFakeCommand(filepath.Dir(filePathBin))

	testCases := []e2ehelpers.CLITestCaseV2{
		{
			Desc: "ng - command's behaviors file does not exist",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				expected.ExitCode = 127
				expected.Stderr = fmt.Sprintf(
					"FAKE_CMD_ERROR failed to read behavior file: %s: open %s: no such file or directory",
					dirPath.FilePathBehaviors(),
					dirPath.FilePathBehaviors(),
				)
			},
		},
		{
			Desc: "ng - command's behaviors file is not json",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				e2ehelpers.MustWriteFile(
					dirPath.FilePathBehaviors(),
					[]byte("a"),
				)

				expected.ExitCode = 127
				expected.Stderr = fmt.Sprintf(
					"FAKE_CMD_ERROR failed to unmarshal behavior file: %s: invalid character 'a' looking for beginning of value",
					dirPath.FilePathBehaviors(),
				)
			},
		},
		{
			Desc: "ng - fake no command's behaviors",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				e2ehelpers.MustWriteJSONFile(
					dirPath.FilePathBehaviors(),
					domains.Behaviors{},
				)

				expected.ExitCode = 127
				expected.Stderr = "FAKE_CMD_ERROR no behaviors"
			},
		},
		{
			Desc: "ok - fake a command's behavior",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.Stderr = "this is a test stderr"
				expected.Stdout = "this is a test stdout"
				expected.ExitCode = 10
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(
					t,
					errors.Join(
						os.RemoveAll(dirPath.FilePathBehaviors()),
						os.RemoveAll(dirPath.FilePathState()),
						os.RemoveAll(dirPath.FilePathProcessing()),
					),
				)
			},
		},
		{
			Desc: "ok - fake a command's behaviors and no executions are done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.ExitCode = 11
				expected.Stdout = "this is a test stdout1"
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(t, errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				))
			},
		},
		{
			Desc: "ok - fake a command's behaviors and already first execution is done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.ExitCode = 12
				expected.Stdout = "this is a test stdout2"
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(t, errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				))
			},
		},
		{
			Desc: "ng - fake a command's behaviors and already all executions are done",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.ExitCode = 127
				expected.Stderr = "FAKE_CMD_ERROR all expected executions are done: expected=2 histories=2"
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(t, errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				))
			},
		},
		{
			Desc: "ng - cannot get lock",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.ExitCode = 127
				expected.Stderr = fmt.Sprintf(
					"FAKE_CMD_ERROR failed to get lock: %s: processing file already exists",
					dirPath.FilePathProcessing(),
				)
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(t, errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				))
			},
		},
		{
			Desc: "ng - state file is not json",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
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

				expected.ExitCode = 127
				expected.Stderr = fmt.Sprintf(
					"FAKE_CMD_ERROR failed to unmarshal state file: %s: invalid character 'a' looking for beginning of value",
					dirPath.FilePathState(),
				)
			},
			Teardown: func(t *testing.T, testID e2ehelpers.TestID) {
				require.NoError(t, errors.Join(
					os.RemoveAll(dirPath.FilePathBehaviors()),
					os.RemoveAll(dirPath.FilePathState()),
					os.RemoveAll(dirPath.FilePathProcessing()),
				))
			},
		},
		{
			Desc: "ng - no error logs",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Envs = []string{"FAKECMD_ERROR_LOG=discard"}
				expected.ExitCode = 127
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Desc, func(t *testing.T) {
			tC.Run(t, filePathBin)
		})
	}
}
