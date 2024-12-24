package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	smockerURL, _ := url.Parse("http://localhost:8081")
	smockerClient := e2ehelpers.NewSmockerClient(
		smockerURL,
		http.DefaultClient,
	)

	testCases := []e2ehelpers.CLITestCase{
		{
			Desc: "ok - no file changes (changeds doesn't include any .tf files on base dir)",
			Envs: append(
				envs,
				"FILE_PATH_TERRAFORM_CMD=/tmp/dummyterraform01",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/case01", dirPathTestdata),
				"-o", "/tmp/case01",
				"-e", "false",
				"-owner", "owner01",
				"-repo", "repo01",
				"-pr", "123",
				"-token", "token01",
			},
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				require.NoError(t, os.RemoveAll("/tmp/case01"))
				require.NoError(t, os.MkdirAll("/tmp/case01", 0755))

				if err := smockerClient.PostMocks(
					types.Mocks{
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123/files",
								},
								QueryParams: types.MultiMapMatcher{
									"page": types.StringMatcherSlice{
										{Matcher: "ShouldEqual", Value: "1"},
									},
									"per_page": types.StringMatcherSlice{
										{Matcher: "ShouldEqual", Value: "100"},
									},
								},
								Headers: types.MultiMapMatcher{
									"E2e-Testid": {
										{
											Matcher: "ShouldEqual",
											Value:   testID.String(),
										},
									},
								},
							},
							Response: &types.MockResponse{
								Status: http.StatusCreated,
								Headers: types.MapStringSlice{
									"Content-Type": types.StringSlice{"application/json"},
								},
								Body: `[
								]`,
							},
						},
					},
					false,
				); err != nil {
					return err
				}

				return nil
			},
			Assertions: func(t *testing.T) {
				result, err := os.ReadFile("/tmp/case01/result.txt")
				require.NoError(t, err)

				assert.Equal(t, ``, string(result))
			},
		},
		{
			Desc: "ok - no file changes (changeds doesn't include any .tf files on base dir)",
			Envs: append(
				envs,
				"FILE_PATH_TERRAFORM_CMD=/tmp/dummyterraform01",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/case02", dirPathTestdata),
				"-o", "/tmp/case02",
				"-e", "false",
				"-owner", "owner01",
				"-repo", "repo01",
				"-pr", "123",
				"-token", "token01",
			},
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				require.NoError(t, os.RemoveAll("/tmp/case02"))
				require.NoError(t, os.MkdirAll("/tmp/case02", 0755))

				if err := smockerClient.PostMocks(
					types.Mocks{
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123/files",
								},
								QueryParams: types.MultiMapMatcher{
									"page": types.StringMatcherSlice{
										{Matcher: "ShouldEqual", Value: "1"},
									},
									"per_page": types.StringMatcherSlice{
										{Matcher: "ShouldEqual", Value: "100"},
									},
								},
								Headers: types.MultiMapMatcher{
									"E2e-Testid": {
										{
											Matcher: "ShouldEqual",
											Value:   testID.String(),
										},
									},
								},
							},
							Response: &types.MockResponse{
								Status: http.StatusCreated,
								Headers: types.MapStringSlice{
									"Content-Type": types.StringSlice{"application/json"},
								},
								Body: `[
								    {"filename": "hoge.go"},
								    {"filename": "/tmp/fuga.tf"}
								]`,
							},
						},
					},
					false,
				); err != nil {
					return err
				}

				return nil
			},
			Assertions: func(t *testing.T) {
				result, err := os.ReadFile("/tmp/case02/result.txt")
				require.NoError(t, err)

				assert.Equal(t, ``, string(result))
			},
		},
	}

	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
