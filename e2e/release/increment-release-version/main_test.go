package main

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/smocker-dev/smocker/server/types"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	smockerURL, _ := url.Parse("http://localhost:8081")
	smockerClient := e2ehelpers.NewSmockerClient(
		smockerURL,
		http.DefaultClient,
	)

	externalCommandFaker := e2ehelpers.ExternalCommandFaker{}
	defer externalCommandFaker.Cleanup() //nolint:errcheck

	testCases := []e2ehelpers.CLITestCase{
		{
			Desc: "ok - increment patch (implicit -increment option)",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 0,
			ExpectedStdout: e2ehelpers.NewLines(
				"created release draft v1.1.5",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e001.sh",
					Stdout: e2ehelpers.NewLines(
						"v1.1.2",
						"v1.1.3",
						"hoge", // no semver is skipped
						"v1.1.4",
					),
				})
				if err != nil {
					return err
				}

				if err := smockerClient.PostMocks(
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
				); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ok - increment major",
			Args: []string{
				"-git", "/tmp/e2e002.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
				"-increment", "major",
			},
			ExpectedExitCode: 0,
			ExpectedStdout: e2ehelpers.NewLines(
				"created release draft v2.0.0",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e002.sh",
					Stdout: e2ehelpers.NewLines(
						"v1.1.2",
						"v1.1.3",
						"v1.1.4",
					),
				})
				if err != nil {
					return err
				}

				if err := smockerClient.PostMocks(
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
				); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ok - increment minor",
			Args: []string{
				"-git", "/tmp/e2e003.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
				"-increment", "minor",
			},
			ExpectedExitCode: 0,
			ExpectedStdout: e2ehelpers.NewLines(
				"created release draft v1.2.0",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e003.sh",
					Stdout: e2ehelpers.NewLines(
						"v1.1.2",
						"v1.1.3",
						"v1.1.4",
					),
				})
				if err != nil {
					return err
				}

				if err := smockerClient.PostMocks(
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
				); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ng - no existing versions in git (missmatch prefix)",
			Args: []string{
				"-git", "/tmp/e2e005.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 2,
			ExpectedStderr: e2ehelpers.NewLines(
				"no existing git versions",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e005.sh",
					Stdout: e2ehelpers.NewLines(
						"1.1.2",
						"1.1.3",
						"1.1.4",
					),
				})
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ng - no existing versions in git",
			Args: []string{
				"-git", "/tmp/e2e006.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 2,
			ExpectedStderr: e2ehelpers.NewLines(
				"no existing git versions",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e006.sh",
				})
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ng - git command is failed",
			Args: []string{
				"-git", "/tmp/e2e007.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 3,
			ExpectedStderr: e2ehelpers.NewLines(
				"failed to git command with code 127",
			),
			Setup: func(t *testing.T, testID e2ehelpers.TestID) error {
				err := externalCommandFaker.Add(&e2ehelpers.ExternalCommandBehavior{
					FilePath: "/tmp/e2e007.sh",
					ExitCode: 127,
				})
				if err != nil {
					return err
				}

				return nil
			},
		},
		{
			Desc: "ng - option -git required",
			Args: []string{
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"-git is required",
			),
		},
		{
			Desc: "ng - option -owner required",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"-owner is required",
			),
		},
		{
			Desc: "ng - option -repo required",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"-repo is required",
			),
		},
		{
			Desc: "ng - option -branch required",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-owner", "owner01",
				"-token", "token01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"-branch is required",
			),
		},
		{
			Desc: "ng - option -token required",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-owner", "owner01",
				"-branch", "branch01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"-token is required",
			),
		},
		{
			Desc: "ng - unknown -increment-type",
			Args: []string{
				"-git", "/tmp/e2e001.sh",
				"-increment", "x",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			ExpectedExitCode: 1,
			ExpectedStderr: e2ehelpers.NewLines(
				"invalid increment type 'x'",
			),
		},
	}
	for _, tC := range testCases {
		tC.Run(t, filePathBin)
	}
}
