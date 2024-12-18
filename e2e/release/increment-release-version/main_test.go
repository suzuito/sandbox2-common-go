package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/smocker-dev/smocker/server/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	testCases := []struct {
		desc             string
		args             []string
		setup            func(*testing.T, uuid.UUID) error
		expectedExitCode int
		expectedStdout   string
		expectedStderr   string
	}{
		{
			desc: "ok - increment patch (implicit -increment option)",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 0,
			expectedStdout: e2ehelpers.NewLines(
				"created release draft v1.1.5",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ok - increment major",
			args: []string{
				"-git", "/tmp/e2e002.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
				"-increment", "major",
			},
			expectedExitCode: 0,
			expectedStdout: e2ehelpers.NewLines(
				"created release draft v2.0.0",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ok - increment minor",
			args: []string{
				"-git", "/tmp/e2e003.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
				"-increment", "minor",
			},
			expectedExitCode: 0,
			expectedStdout: e2ehelpers.NewLines(
				"created release draft v1.2.0",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ng - no existing versions in git (missmatch prefix)",
			args: []string{
				"-git", "/tmp/e2e005.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 2,
			expectedStderr: e2ehelpers.NewLines(
				"no existing git versions",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ng - no existing versions in git",
			args: []string{
				"-git", "/tmp/e2e006.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 2,
			expectedStderr: e2ehelpers.NewLines(
				"no existing git versions",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ng - git command is failed",
			args: []string{
				"-git", "/tmp/e2e007.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 3,
			expectedStderr: e2ehelpers.NewLines(
				"failed to git command with code 127",
			),
			setup: func(t *testing.T, testID uuid.UUID) error {
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
			desc: "ng - option -git required",
			args: []string{
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"-git is required",
			),
		},
		{
			desc: "ng - option -owner required",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"-owner is required",
			),
		},
		{
			desc: "ng - option -repo required",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"-repo is required",
			),
		},
		{
			desc: "ng - option -branch required",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-owner", "owner01",
				"-token", "token01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"-branch is required",
			),
		},
		{
			desc: "ng - option -token required",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-repo", "repo01",
				"-owner", "owner01",
				"-branch", "branch01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"-token is required",
			),
		},
		{
			desc: "ng - unknown -increment-type",
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-increment", "x",
				"-owner", "owner01",
				"-repo", "repo01",
				"-branch", "branch01",
				"-token", "token01",
			},
			expectedExitCode: 1,
			expectedStderr: e2ehelpers.NewLines(
				"invalid increment type 'x'",
			),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()

			testID := uuid.New()

			cmd := exec.CommandContext(
				ctx,
				filePathBin,
				tC.args...,
			)

			cmd.Env = append(
				os.Environ(),
				fmt.Sprintf("E2E_TEST_ID=%s", testID),
				"GITHUB_HTTP_CLIENT_FAKE_SCHEME=http",
				"GITHUB_HTTP_CLIENT_FAKE_HOST=localhost:8080",
			)

			stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			if tC.setup != nil {
				if err := tC.setup(t, testID); err != nil {
					require.Error(t, err, err.Error())
				}
			}

			err := cmd.Run()
			if err != nil {
				var exiterr *exec.ExitError
				if !errors.As(err, &exiterr) {
					require.NoError(t, err, fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(tC.args, " "), err.Error()))
				}
			}

			assert.Equal(t, tC.expectedExitCode, cmd.ProcessState.ExitCode())
			assert.Equal(t, tC.expectedStdout, stdout.String())
			assert.Equal(t, tC.expectedStderr, stderr.String())
		})
	}
}
