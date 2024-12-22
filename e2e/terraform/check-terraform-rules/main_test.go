package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestCheckTerraformRules(t *testing.T) {
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
					fmt.Sprintf(`resource terraform.backend."gcs" not found %s/case001/mods/ng001`, dirPathTestdata),
					fmt.Sprintf(`resource provider."google" not found %s/case001/mods/ng002`, dirPathTestdata),
					fmt.Sprintf(`invalid terraform.backend."gcs".bucket: hoge-terraform %s/case001/mods/ng003`, dirPathTestdata),
					fmt.Sprintf(`invalid terraform.backend."gcs".prefix: hoge %s/case001/mods/ng004`, dirPathTestdata),
					"not pass\n",
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
