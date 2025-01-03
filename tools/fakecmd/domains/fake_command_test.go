package domains_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

func TestDirPathFakeCommand(t *testing.T) {
	t.Parallel()

	d := domains.DirPathFakeCommand("hoge")
	assert.Equal(t, "hoge/cmd", d.FilePathCommand())
	assert.Equal(t, "hoge/processing", d.FilePathProcessing())
	assert.Equal(t, "hoge/behaviors.json", d.FilePathBehaviors())
	assert.Equal(t, "hoge/state.json", d.FilePathState())
	assert.Equal(t, "hoge", d.String())
}

func TestFakeCommand_Init(t *testing.T) {
	t.Parallel()

	filePathFakeCMD := fmt.Sprintf("/tmp/%s", uuid.NewString())
	e2ehelpers.MustWriteFile(filePathFakeCMD, []byte{})

	dirPath := domains.DirPathFakeCommand(fmt.Sprintf("/tmp/%s", uuid.NewString()))

	t.Cleanup(func() {
		os.RemoveAll(filePathFakeCMD)
		os.RemoveAll(dirPath.String())
	})

	testCases := []struct {
		desc                 string
		inputFilePathFakeCMD string
		inputDirPath         domains.DirPathFakeCommand
		inputBehaviors       domains.Behaviors
		setup                func(t *testing.T)
		wantErr              bool
		errMsg               string
		assertionFunc        func(t *testing.T)
	}{
		{
			desc:                 "ok",
			inputFilePathFakeCMD: filePathFakeCMD,
			inputDirPath:         dirPath,
			inputBehaviors:       domains.Behaviors{},
			assertionFunc: func(t *testing.T) {
				require.DirExists(t, dirPath.String())
				require.FileExists(t, dirPath.FilePathBehaviors())
				require.FileExists(t, dirPath.FilePathCommand())
			},
		},
		{
			desc:                 "ng - file path of fakecmd is invalid",
			inputFilePathFakeCMD: "/aaa/bbb/ccc",
			inputDirPath:         dirPath,
			inputBehaviors:       domains.Behaviors{},
			wantErr:              true,
			errMsg:               "fakecmd is invalid: /aaa/bbb/ccc: open /aaa/bbb/ccc: no such file or directory",
		},
		{
			desc:                 "ng - cloned fakecmd already exists",
			inputFilePathFakeCMD: filePathFakeCMD,
			inputDirPath:         domains.DirPathFakeCommand("/tmp/case003"),
			inputBehaviors:       domains.Behaviors{},
			setup: func(t *testing.T) {
				dirPath := domains.DirPathFakeCommand("/tmp/case003")
				os.RemoveAll(dirPath.String())
				e2ehelpers.MustMkdir(dirPath.String())
				e2ehelpers.MustWriteFile(dirPath.FilePathCommand(), []byte{})
			},
			wantErr: true,
			errMsg:  "cloned fakecmd already exists: /tmp/case003/cmd",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			if tC.setup != nil {
				tC.setup(t)
			}

			fcmd := domains.NewFakeCommand(
				tC.inputFilePathFakeCMD,
				tC.inputDirPath,
				tC.inputBehaviors,
			)

			err := fcmd.Init(false)
			if tC.wantErr {
				require.EqualError(t, err, tC.errMsg)
			} else {
				require.NoError(t, err)

				if tC.assertionFunc != nil {
					tC.assertionFunc(t)
				}
			}
		})
	}
}
