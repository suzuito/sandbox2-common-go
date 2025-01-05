package main

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/fakecmd/domains"
)

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	smockerURL, _ := url.Parse("http://localhost:8081")
	smockerClient := e2ehelpers.NewSmockerClient(
		smockerURL,
		http.DefaultClient,
	)

	commandFaker := domains.MustByEnv()
	defer commandFaker.Cleanup() //nolint:errcheck

	envs := []string{
		"GITHUB_HTTP_CLIENT_FAKE_SCHEME=http",
		"GITHUB_HTTP_CLIENT_FAKE_HOST=localhost:8080",
	}

	testCases := []e2ehelpers.CLITestCaseV2{
		{
			Desc: "ok - increment patch (implicit -increment option)",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: e2ehelpers.NewLines(
								"v1.1.2",
								"v1.1.3",
								"hoge", // no semver is skipped
								"v1.1.4",
							),
						},
					},
				})

				require.NoError(t, smockerClient.PostMocks(
					types.Mocks{
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/releases",
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

				input.Envs = envs
				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}
				expected.Stdout = e2ehelpers.NewLines(
					"created release draft v1.1.5",
				)
			},
		},
		{
			Desc: "ok - increment major",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: e2ehelpers.NewLines(
								"v1.1.2",
								"v1.1.3",
								"v1.1.4",
							),
						},
					},
				})

				require.NoError(t, smockerClient.PostMocks(
					types.Mocks{
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/releases",
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

				input.Envs = envs
				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
					"-increment", "major",
				}

				expected.Stdout = e2ehelpers.NewLines(
					"created release draft v2.0.0",
				)
			},
		},
		{
			Desc: "ok - increment minor",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: e2ehelpers.NewLines(
								"v1.1.2",
								"v1.1.3",
								"v1.1.4",
							),
						},
					},
				})

				require.NoError(t, smockerClient.PostMocks(
					types.Mocks{
						{
							Request: types.MockRequest{
								Method: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "POST",
								},
								Path: types.StringMatcher{
									Matcher: "ShouldEqual",
									Value:   "/repos/owner01/repo01/releases",
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

				input.Envs = envs
				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
					"-increment", "minor",
				}

				expected.Stdout = e2ehelpers.NewLines(
					"created release draft v1.2.0",
				)
			},
		},
		{
			Desc: "ng - no existing versions in git (missmatch prefix)",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							Stdout: e2ehelpers.NewLines(
								"1.1.2",
								"1.1.3",
								"1.1.4",
							),
						},
					},
				})

				input.Envs = envs
				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.ExitCode = 2
				expected.Stderr = e2ehelpers.NewLines(
					"no existing git versions",
				)
			},
		},
		{
			Desc: "ng - no existing versions in git",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type:                         domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{},
					},
				})

				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.Stderr = e2ehelpers.NewLines(
					"no existing git versions",
				)
				expected.ExitCode = 2
			},
		},
		{
			Desc: "ng - git command is failed",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				fcmd := commandFaker.AddInTest(t, domains.Behaviors{
					{
						Type: domains.BehaviorTypeStdoutStderrExitCode,
						BehaviorStdoutStderrExitCode: &domains.BehaviorStdoutStderrExitCode{
							ExitCode: 127,
						},
					},
				})

				input.Envs = envs
				input.Args = []string{
					"-git", fcmd.DirPath().FilePathCommand(),
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.ExitCode = 3
				expected.Stderr = e2ehelpers.NewLines(
					"failed to git command with code 127",
				)
			},
		},
		{
			Desc: "ng - option -git required",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-prefix", "v",
					"-owner", "owner01",
					"-repo", "repo01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"-git is required",
				)
			},
		},
		{
			Desc: "ng - option -owner required",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-git", "/tmp/e2e001.sh",
					"-prefix", "v",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"-owner is required",
				)
			},
		},
		{
			Desc: "ng - option -repo required",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-git", "/tmp/e2e001.sh",
					"-prefix", "v",
					"-owner", "owner01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"-repo is required",
				)
			},
		},
		{
			Desc: "ng - option -branch required",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-git", "/tmp/e2e001.sh",
					"-prefix", "v",
					"-repo", "repo01",
					"-owner", "owner01",
					"-token", "token01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"-branch is required",
				)
			},
		},
		{
			Desc: "ng - option -token required",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-git", "/tmp/e2e001.sh",
					"-prefix", "v",
					"-repo", "repo01",
					"-owner", "owner01",
					"-branch", "branch01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"-token is required",
				)
			},
		},
		{
			Desc: "ng - unknown -increment-type",
			Setup: func(t *testing.T, testID e2ehelpers.TestID, input *e2ehelpers.CLITestCaseV2Input, expected *e2ehelpers.CLITestCaseV2Expected) {
				input.Args = []string{
					"-git", "/tmp/e2e001.sh",
					"-increment", "x",
					"-owner", "owner01",
					"-repo", "repo01",
					"-branch", "branch01",
					"-token", "token01",
				}

				expected.ExitCode = 1
				expected.Stderr = e2ehelpers.NewLines(
					"invalid increment type 'x'",
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
