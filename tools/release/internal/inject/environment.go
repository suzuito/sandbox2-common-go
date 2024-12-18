package inject

type Environment struct {
	E2ETestID                  string `envconfig:"E2E_TEST_ID"`
	GithubHTTPClientFakeScheme string `envconfig:"GITHUB_HTTP_CLIENT_FAKE_SCHEME"`
	GithubHTTPClientFakeHost   string `envconfig:"GITHUB_HTTP_CLIENT_FAKE_HOST"`
}
