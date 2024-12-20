package main

import (
	"os"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestXxx(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	testCases := []e2ehelpers.CLITestCase{
		{
			Args: []string{
				"-d", "./testdata/success001",
			},
			ExpectedExitCode: 0,
		},
	}
	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
