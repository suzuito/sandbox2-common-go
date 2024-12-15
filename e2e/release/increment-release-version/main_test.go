package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type SmockerClient struct {
	baseURL *url.URL
	client  *http.Client
}

func (t *SmockerClient) PostMocks(
	body PostMocksRequest,
	reset bool,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return terrors.Wrapf("failed to json.Marshal: %w", err)
	}

	reqURL, _ := url.Parse(t.baseURL.String())
	reqURL.Path = "/mocks"
	query := reqURL.Query()
	query.Set("reset", strconv.FormatBool(reset))
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return terrors.Wrapf("failed to http.NewRequest: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := t.client.Do(req)
	if err != nil {
		return terrors.Wrapf("failed to http request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			resBodyBytes = []byte{}
		}
		return terrors.Wrapf(
			"http error: status=%d body=%s",
			res.StatusCode, string(resBodyBytes),
		)
	}

	return nil
}

type PostMocksRequest []Mock

type Mock struct {
	Request  *MockRequest  `json:"request"`
	Response *MockResponse `json:"response"`
}

type MockRequest struct {
	Method  StringMatcher              `json:"method"`
	Path    StringMatcher              `json:"path"`
	Headers map[string][]StringMatcher `json:"headers"`
}

type MockResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"string"`
}

type StringMatcher struct {
	MatchOp MatchOp `json:"matcher"`
	Value   string  `json:"value"`
}

type MatchOp string

const (
	MatchOpShouldEqual    MatchOp = "ShouldEqual"
	MatchOpShouldNotEqual MatchOp = "ShouldNotEqual"
)

func NewSmockerClient(
	baseURL *url.URL,
	client *http.Client,
) *SmockerClient {
	return &SmockerClient{
		baseURL: baseURL,
		client:  client,
	}
}

func TestA(t *testing.T) {
	filePathBin := os.Getenv("FILE_PATH_BIN")

	testCases := []struct {
		desc             string
		args             []string
		setup            func(uuid.UUID) error
		expectedExitCode int
		expectedStdout   string
		expectedStderr   string
	}{
		{
			desc: "ok",
			setup: func(testID uuid.UUID) error {
				// external command
				f, err := os.Create("/tmp/e2e001.sh")
				if err != nil {
					return err
				}
				defer f.Close()

				f.Chmod(0755)

				fmt.Fprintf(f, "#!/bin/sh\n")
				fmt.Fprintf(f, "echo 'v1.1.2'\n")
				fmt.Fprintf(f, "echo 'v1.1.3'\n")
				fmt.Fprintf(f, "echo 'v1.1.4'\n")

				// smocker mock
				smockerURL, _ := url.Parse("http://localhost:8081")
				smockerClient := NewSmockerClient(
					smockerURL,
					http.DefaultClient,
				)
				if err := smockerClient.PostMocks(
					[]Mock{
						{
							Request: &MockRequest{
								Method: StringMatcher{
									MatchOp: MatchOpShouldEqual,
									Value:   "POST",
								},
								Path: StringMatcher{
									MatchOp: MatchOpShouldEqual,
									Value:   "/repos/owner01/repo01/releases",
								},
								Headers: map[string][]StringMatcher{
									"E2e-Testid": {
										{
											MatchOp: MatchOpShouldEqual,
											Value:   testID.String(),
										},
									},
								},
							},
							Response: &MockResponse{
								Status: http.StatusCreated,
								Headers: map[string]string{
									"Content-Type": "application/json",
								},
								Body: `{}`,
							},
						},
					},
					true,
				); err != nil {
					return err
				}

				return nil
			},
			args: []string{
				"-git", "/tmp/e2e001.sh",
				"-prefix", "v",
				"-owner", "owner01",
				"-repo", "repo01",
			},
			expectedExitCode: 0,
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

			if err := tC.setup(testID); err != nil {
				require.Error(t, err, err.Error())
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
