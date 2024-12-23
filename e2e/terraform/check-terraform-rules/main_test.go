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
					fmt.Sprintf(`resource terraform.backend."gcs" not found (%s/case001/mods/ng001)`, dirPathTestdata),
					fmt.Sprintf(`resource provider."google" not found (%s/case001/mods/ng002)`, dirPathTestdata),
					fmt.Sprintf("invalid terraform.backend.\"gcs\".bucket (%s/case001/mods/ng003)", dirPathTestdata),
					"  expected: base-999-terraform",
					"  actual: hoge-terraform",
					fmt.Sprintf("invalid terraform.backend.\"gcs\".prefix (%s/case001/mods/ng004)\n", dirPathTestdata),
				},
				"\n",
			),
			ExpectedStderr: strings.Join(
				[]string{
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
		},
		{
			Desc: `ng - set no existing dir as -d opt value`,
			Args: []string{
				"-d", fmt.Sprintf("%s/caseXXX", dirPathTestdata),
			},
			ExpectedExitCode: 10,
			ExpectedStderr: strings.Join(
				[]string{
					fmt.Sprintf("%s/caseXXX does not exist", dirPathTestdata),
				},
				"\n",
			),
		},
		{
			Desc: `ok - broken hcl file is ignored`,
			Args: []string{
				"-d", fmt.Sprintf("%s/case003", dirPathTestdata),
			},
		},
	}
	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
