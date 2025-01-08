package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

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

	commandFaker := domains.MustByEnv()
	defer commandFaker.Cleanup() //nolint:errcheck

	eventPath1 := fmt.Sprintf("/tmp/%s", uuid.NewString())
	e2ehelpers.MustWriteFile(eventPath1, []byte(`{}`))
	defer os.RemoveAll(eventPath1) //nolint:errcheck

	testCases := []e2ehelpers.CLITestCaseV2{
		{
			Desc: "ng - -event-name is required",
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
					"-event-path", "foo.json",
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-git-rootdir", "dummygithubroot",
				}
				expected.ExitCode = 1
				expected.Stderr = "-event-name is required"
			},
		},
		{
			Desc: "ng - -event-path is required",
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
					"-event-name", "issue_comment",
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-git-rootdir", "dummygithubroot",
				}
				expected.ExitCode = 1
				expected.Stderr = "-event-path is required"
			},
		},
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
					"-event-name", "issue_comment",
					"-event-path", eventPath1,
					"-git-rootdir", "dummygithubroot",
				}
				expected.ExitCode = 1
				expected.Stderr = "-d is required"
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
					"-event-name", "issue_comment",
					"-event-path", eventPath1,
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = "-git-rootdir is required"
			},
		},
		{
			Desc: "ok - skipped - event name",
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
					"-event-name", "ev",
					"-event-path", eventPath1,
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-git-rootdir", "dummygithubroot",
				}
				expected.Stdout = "skipped"
			},
		},
		{
			Desc: "ng - invalid event path",
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
					"-event-name", "issue_comment",
					"-event-path", "foo.json",
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-git-rootdir", "dummygithubroot",
				}
				expected.ExitCode = 1
				expected.Stderr = "open foo.json: no such file or directory"
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
					"-event-name", "issue_comment",
					"-event-path", eventPath1,
					"-d", fmt.Sprintf("%s/ghrepo01/caseXX", dirPathTestdata),
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
				}
				expected.ExitCode = 1
				expected.Stderr = fmt.Sprintf(
					"invalid base dir: %s/ghrepo01/caseXX: stat %s/ghrepo01/caseXX: no such file or directory",
					dirPathTestdata, dirPathTestdata,
				)
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform plan"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform plan"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
			Desc: "ok - [issue comment] - terraform plan - with diff",
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 2,
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform plan"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 2",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"diff at `terraform plan`",
				)

				expected.ExitCode = 2
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123",
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
								Body: `{"mergeable":true}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
		{
			Desc: "ok - [issue comment] - terraform apply - with diff",
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123",
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
								Body: `{"mergeable":true}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 2,
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 2",
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
		{
			Desc: "ng - [issue comment] - terraform apply - PR is not mergeable",
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123",
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
								Body: `{"mergeable":false}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 2,
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 2",
					"",
					"pr is not mergeable",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"cannot exec `terraform apply` because PR is not mergeable",
				)

				expected.ExitCode = 3
			},
		},
		{
			Desc: "ng - [issue comment] - terraform apply - failed to `terraform init`",
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 99,
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 99",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"failed to terraform.Init: failed to init",
				)

				expected.ExitCode = 125
			},
		},
		{
			Desc: "ng - [issue comment] - terraform apply - failed to `terraform plan`",
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 99,
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 99",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"failed to terraform.Plan: failed to plan",
				)

				expected.ExitCode = 125
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123",
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
								Body: `{"mergeable":true}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
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
							Stdout:   "this is terraform command stdout\n",
							Stderr:   "this is terraform command stderr\n",
							ExitCode: 99,
						},
					},
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
					"exit with 99",
					"",
					"",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
					"failed to terraform.Apply: failed to apply",
				)

				expected.ExitCode = 125
			},
		},
		{
			Desc: "ok - [issue comment] - terraform apply - auto merged",
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
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "GET",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123",
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
								Body: `{"mergeable":true}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/issues/123/comments",
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
								Body: `{}`,
							},
						},
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "PUT",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/pulls/123/merge",
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
								Status: http.StatusOK,
								Headers: types.MapStringSlice{
									"Content-Type": types.StringSlice{"application/json"},
								},
								Body: `{}`,
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
					"-event-name", "issue_comment",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
					    "comment": {"body": "///terraform apply"},
						"issue": {"number":123, "pull_request":{}},
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
					"-git-rootdir", fmt.Sprintf("%s/ghrepo01", dirPathTestdata),
					"-automerge",
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
					"PR is merged",
				)

				expected.Stderr = e2ehelpers.NewLines(
					"this is terraform command stderr",
					"this is terraform command stderr",
					"this is terraform command stderr",
				)
			},
		},
		// schedule
		{
			Desc: "ok - [schedule] - with empty diff",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				behavior := domains.Behavior{
					Type: domains.BehaviorTypeStdoutStderrExitCode,
					BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
						Stdout: "this is terraform command stdout\n",
						Stderr: "this is terraform command stderr\n",
					},
				}
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "schedule",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
						"%s -chdir=%s/ghrepo01/case01/roots/r4 init -no-color",
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
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r4 plan -no-color -detailed-exitcode",
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
					"this is terraform command stderr",
					"this is terraform command stderr",
				)
			},
		},
		// workflow_dispatch
		{
			Desc: "ok - [workflow_dispatch] - with empty diff",
			Setup: func(
				t *testing.T,
				testID e2ehelpers.TestID,
				input *e2ehelpers.CLITestCaseV2Input,
				expected *e2ehelpers.CLITestCaseV2Expected,
			) {
				behavior := domains.Behavior{
					Type: domains.BehaviorTypeStdoutStderrExitCode,
					BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
						Stdout: "this is terraform command stdout\n",
						Stderr: "this is terraform command stderr\n",
					},
				}
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
					behavior,
				})

				input.Envs = append(
					envs,
					"GITHUB_TOKEN=foo",
					fmt.Sprintf("FILE_PATH_TERRAFORM=%s", fcmd.DirPath().FilePathCommand()),
				)
				input.Args = []string{
					"-event-name", "workflow_dispatch",
					"-event-path", e2ehelpers.MustWriteFileAtRandomPath("/tmp", []byte(`{
						"repository":{
						  "name": "repo01",
						  "owner": {
						    "login": "owner01"
						  }
						}
					}`)),
					"-d", fmt.Sprintf("%s/ghrepo01/case01", dirPathTestdata),
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
						"%s -chdir=%s/ghrepo01/case01/roots/r4 init -no-color",
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
					"*************",
					"*************",
					"*************",
					"==== CMD ====",
					fmt.Sprintf(
						"%s -chdir=%s/ghrepo01/case01/roots/r4 plan -no-color -detailed-exitcode",
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
