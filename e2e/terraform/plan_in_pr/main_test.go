package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestCheckTerraformRules(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	dirPathTestdata, err := filepath.Abs("testdata")
	if err != nil {
		panic(err)
	}

	envs := []string{
		"GITHUB_HTTP_CLIENT_FAKE_SCHEME=http",
		"GITHUB_HTTP_CLIENT_FAKE_HOST=localhost:8080",
		"FILE_PATH_TERRAFORM_CMD=",
	}

	testCases := []e2ehelpers.CLITestCase{
		{
			Envs: append(
				envs,
				"FILE_PATH_TERRAFORM_CMD=/tmp/dummyterraform01",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/case01", dirPathTestdata),
				"-o", "/tmp",
				"-e", "false",
				"-owner", "owner01",
				"-repo", "repo01",
				"-token", "token01",
			},
		},
	}

	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
