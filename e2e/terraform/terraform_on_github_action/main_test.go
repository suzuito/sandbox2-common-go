package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
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

	commandFaker := domains.MustByEnv()
	defer commandFaker.Cleanup()

	testCases := []e2ehelpers.CLITestCaseV2{
		{
			Desc: "ng - -d is required",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-p", googleCloudProjectID,
					"-github-context", "{}",
					"-git-rootdir", "dummygithubroot",
				}
				expected.ExitCode = 1
				expected.Stderr = "-d is required"
			},
		},
		{
			Desc: "ng - -p is required",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-github-context", "{}",
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = "-p is required"
			},
		},
		{
			Desc: "ng - -github-context is required",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = "-github-context is required"
			},
		},
		{
			Desc: "ng - -git-rootdir is required",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-github-context", "{}",
				}
				expected.ExitCode = 1
				expected.Stderr = "-git-rootdir is required"
			},
		},
		{
			Desc: "ng - invalid base directory",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/caseXX", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-github-context", "{}",
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = fmt.Sprintf(
					"invalid base dir: %s/ghrepo01/caseXX: stat %s/ghrepo01/caseXX: no such file or directory",
					dirPathTestdata, dirPathTestdata,
				)
			},
		},
		{
			Desc: "ng - github context is invalid json",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-github-context", "dummygithubcontext",
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = "github context is invalid json: dummygithubcontext: invalid character 'd' looking for beginning of value"
			},
		},
		{
			Desc: "ok - skipped github action event",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-github-context", "{}",
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.Stdout = "skipped"
			},
		},
		// issue comment
		{
			Desc: "ok - [issue comment] - terraform plan - no file changed",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
				)
				input.Args = []string{
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
				}

				require.NoError(t, smockerClient.PostMocks(
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
				))

				expected.Stdout = "no file changed in PR: 123"
			},
		},
		{
			Desc: "ok - [issue comment] - terraform plan - with empty diff",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				require.NoError(t, smockerClient.PostMocks(
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
				))

				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
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
				}

				expected.Stdout = e2ehelpers.NewLines(
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r1 init -no-color",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r1 plan -no-color -detailed-exitcode",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r2 init -no-color",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r2 plan -no-color -detailed-exitcode",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r3 init -no-color",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r3 plan -no-color -detailed-exitcode",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
				)
			},
		},
		{
			Desc: "ok - [issue comment] - terraform apply - with empty diff",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				require.NoError(t, smockerClient.PostMocks(
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
								    {"filename":"case01/roots/r1/main.tf"}
								]`,
							},
						},
					},
					false,
				))

				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: "this is terraform command stdout\n",
							Stderr: "this is terraform command stderr\n",
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-p", googleCloudProjectID,
					"-github-context", e2ehelpers.MinifyJSONString(`{
				    "event_name":"issue_comment",
					"repository":"owner01/repo01",
					"repository_owner":"owner01",
					"issue":{"number":123,"pull_request":{}},
					"event":{"comment":{"body":"///terraform apply"}}
				}`),
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}

				expected.Stdout = e2ehelpers.NewLines(
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r1 init -no-color",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r1 plan -no-color -detailed-exitcode",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r1 apply -no-color -auto-approve",
						fcmd.DirPath().FilePathCommand(),
						dirPathTestdata,
					),
					"==== OUT ====",
					"this is terraform command stdout",
					"==== END ====",
					"exit with 0",
					"",
					"",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
				)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Desc, func(t *testing.T) {
			tC.Run(t, filePathBin)
		})
	}
}
