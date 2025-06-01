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

	testCases := []e2ehelpers.CLITestCaseV2{
		{
			Desc: `ng - -d is required`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				expected.ExitCode = 1
				expected.Stderr = "-d is required"
			},
		},
		{
			Desc: `ok - some rule violations exist`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-d", fmt.Sprintf("%s/case001", dirPathTestdata),
				}

				expected.ExitCode = 5
				expected.Stdout = strings.Join(
					[]string{
						fmt.Sprintf(`resource terraform.backend."gcs" not found (%s/case001/mods/ng001)`, dirPathTestdata),
						"  expected: true",
						"  actual: false",
						fmt.Sprintf(`resource provider."google" not found (%s/case001/mods/ng002)`, dirPathTestdata),
						"  expected: true",
						"  actual: false",
						fmt.Sprintf("invalid terraform.backend.\"gcs\".bucket (%s/case001/mods/ng003)", dirPathTestdata),
						"  expected: base-999-terraform",
						"  actual: hoge-terraform",
						fmt.Sprintf("invalid terraform.backend.\"gcs\".prefix (%s/case001/mods/ng004)", dirPathTestdata),
						"  expected: mods/ng004",
						"  actual: hoge\n",
					},
					"\n",
				)
				expected.Stderr = strings.Join(
					[]string{
						"cli error: not pass\n",
					},
					"\n",
				)
			},
		},
		{
			Desc: `ok - no rule violations`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-d", fmt.Sprintf("%s/case002", dirPathTestdata),
				}
			},
		},
		{
			Desc: `ng - set no existing dir as -d opt value`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-d", fmt.Sprintf("%s/caseXXX", dirPathTestdata),
				}
				expected.ExitCode = 10
				expected.Stderr = strings.Join(
					[]string{
						fmt.Sprintf("cli error: %s/caseXXX does not exist", dirPathTestdata),
					},
					"\n",
				)
			},
		},
		{
			Desc: `ok - broken hcl file is ignored`,
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-d", fmt.Sprintf("%s/case003", dirPathTestdata),
				}
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.Desc, func(t *testing.T) {
			tC.Run(t, filePathBin)
		})
	}
}
