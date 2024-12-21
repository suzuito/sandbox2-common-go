package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestXxx(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	dirPathTestdata, err := filepath.Abs("testdata")
	if err != nil {
		panic(err)
	}

	testCases := []e2ehelpers.CLITestCase{
		{
			Desc:             `ng - -d is required`,
			Args:             []string{},
			ExpectedExitCode: 1,
			ExpectedStderr:   "-d is required",
		},
		{
			Desc: `ok - some rule violations exist`,
			Args: []string{
				"-d", fmt.Sprintf("%s/case001", dirPathTestdata),
			},
			ExpectedExitCode: 5,
			ExpectedStdout: strings.Join(
				[]string{
					fmt.Sprintf("ok: %s/mods/ok001", dirPathTestdata),
					fmt.Sprintf(`ng: terraform > backend "gcs" is not found (%s/mods/ng001)`, dirPathTestdata),
					fmt.Sprintf(`ng: provider > google > project is not found (%s/mods/ng002)`, dirPathTestdata),
					fmt.Sprintf(`ng: terraform > backend "gcs" > bucket is invalid (%s/mods/ng003/hoge.tf:y)`, dirPathTestdata),
					fmt.Sprintf(`ng: terraform > backend "gcs" > prefix is invalid (%s/mods/ng004/fuga.tf:z)`, dirPathTestdata),
				},
				"\n",
			),
		},
		{
			Desc: `ok - no rule violations`,
			Args: []string{
				"-d", fmt.Sprintf("%s/case002", dirPathTestdata),
			},
			ExpectedStdout: strings.Join(
				[]string{
					fmt.Sprintf("ok: %s/mods/ok001", dirPathTestdata),
				},
				"\n",
			),
		},
	}
	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
