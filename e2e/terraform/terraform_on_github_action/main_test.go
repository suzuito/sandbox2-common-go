package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

var googleCloudProjectID = "prj01"

func TestTerraformOnGithubAction(t *testing.T) {
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

	externalCommandFaker := e2ehelpers.ExternalCommandFaker{}
	defer externalCommandFaker.Cleanup()

	testCases := []e2ehelpers.CLITestCase{
		{
			Desc: "ng - -d is required",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-p", googleCloudProjectID,
				"-github-context", "{}",
				"-git-rootdir", "dummygithubroot",
			},
			ExpectedExitCode: 1,
			ExpectedStderr:   "-d is required",
		},
		{
			Desc: "ng - -p is required",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-github-context", "{}",
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			ExpectedExitCode: 1,
			ExpectedStderr:   "-p is required",
		},
		{
			Desc: "ng - -github-context is required",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			ExpectedExitCode: 1,
			ExpectedStderr:   "-github-context is required",
		},
		{
			Desc: "ng - -git-rootdir is required",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", "{}",
			},
			ExpectedExitCode: 1,
			ExpectedStderr:   "-git-rootdir is required",
		},
		{
			Desc: "ng - invalid base directory",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/caseXX", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", "{}",
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			ExpectedExitCode: 1,
			ExpectedStderr: fmt.Sprintf(
				"invalid base dir: %s/ghrepo01/caseXX: stat %s/ghrepo01/caseXX: no such file or directory",
				dirPathTestdata, dirPathTestdata,
			),
		},
		{
			Desc: "ng - github context is invalid json",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", "dummygithubcontext",
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			ExpectedExitCode: 1,
			ExpectedStderr:   "github context is invalid json: dummygithubcontext: invalid character 'd' looking for beginning of value",
		},
		{
			Desc: "ok - skipped github action event",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", "{}",
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			ExpectedStdout: "skipped",
		},
		// issue comment
		{
			Desc: "ok - [issue comment] - no file changed",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", e2ehelpers.MinifyJSONString(`{
				    "event_name":"issue_comment",
					"repository":"owner01/repo01",
					"repository_owner":"owner01",
					"issue":{"number":123,"pull_request":{}},
					"event":{"comment":{"body":"///terraform plan"}}
				}`),
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
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
								Body: `[]`,
							},
						},
					},
					false,
				); err != nil {
					return err
				}

				return nil
			},
			ExpectedStdout: "no file changed in PR: 123",
		},
		{
			Desc: "ok - [issue comment] - no file changed",
			Envs: append(
				envs,
				"GITHUB_TOKEN=foo",
				"FILE_PATH_TERRAFORM=/tmp/terraform001.sh",
			),
			Args: []string{
				"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				"-p", googleCloudProjectID,
				"-github-context", e2ehelpers.MinifyJSONString(`{
				    "event_name":"issue_comment",
					"repository":"owner01/repo01",
					"repository_owner":"owner01",
					"issue":{"number":123,"pull_request":{}},
					"event":{"comment":{"body":"///terraform plan"}}
				}`),
				"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
			},
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
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
								    {"filename":"hoge/fuga.go"},
								    {"filename":"case01/roots/r1/main.tf"},
								    {"filename":"case01/roots/r2/.terraform.lock.hcl"},
								    {"filename":"case01/commons/r3m1m1/main.tf"}
								]`,
							},
						},
					},
					false,
				); err != nil {
					return err
				}

				if err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/terraform001.sh",
					Stdout:   "this is terraform command stdout",
					Stderr:   "this is terraform command stderr",
				}); err != nil {
					return err
				}

				return nil
			},
			ExpectedStdout: e2ehelpers.NewLines(
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r1 init -no-color",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r1 plan -no-color -detailed-exitcode",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r2 init -no-color",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r2 plan -no-color -detailed-exitcode",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r3 init -no-color",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
				"==== CMD ====",
				fmt.Sprintf(
					"/tmp/terraform001.sh -chdir=%s/ghrepo01/case01/roots/r3 plan -no-color -detailed-exitcode",
					dirPathTestdata,
				),
				"==== OUT ====",
				"this is terraform command stdout",
				"==== END ====",
				"exit with 0",
				"",
				"",
			),
			ExpectedStderr: e2ehelpers.NewLines(
				"this is terraform command stderr",
				"this is terraform command stderr",
				"this is terraform command stderr",
				"this is terraform command stderr",
				"this is terraform command stderr",
				"this is terraform command stderr",
			),
		},
	}

	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
