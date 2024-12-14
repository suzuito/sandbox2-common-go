package inject

import (
	"net/http"

	"github.com/google/go-github/v67/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox2-common-go/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/gh/repositories"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/local/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/usecases"
)

type RoundTripperForE2E struct {
	e2eTestID  string
	origin     http.RoundTripper
	fakeScheme string
	fakeHost   string
}

func (t *RoundTripperForE2E) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("E2E-TestId", t.e2eTestID)
	originalURL := req.URL
	originalURL.Scheme = t.fakeScheme
	originalURL.Host = t.fakeHost
	return t.origin.RoundTrip(req)
}

func NewRoundTripperForE2E(
	e2eTestID string,
	origin http.RoundTripper,
	fakeScheme string,
	fakeHost string,
) *RoundTripperForE2E {
	return &RoundTripperForE2E{
		e2eTestID:  e2eTestID,
		origin:     origin,
		fakeScheme: fakeScheme,
		fakeHost:   fakeHost,
	}
}

func NewUsecase(
	filePathGit string,
	githubAppToken string,
) (usecases.Usecase, error) {
	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		return nil, terrors.Wrapf("failed to load environment variable: %w", err)
	}

	githubHTTPClient := http.DefaultClient
	if env.E2ETestID != "" {
		githubHTTPClient = &http.Client{
			Transport: NewRoundTripperForE2E(
				env.E2ETestID,
				http.DefaultTransport,
				env.GithubHTTPClientFakeScheme,
				env.GithubHTTPClientFakeHost,
			),
		}
	}

	githubClient := github.
		NewClient(githubHTTPClient).
		WithAuthToken(githubAppToken)
	releaseRepository := repositories.NewReleaseRepository(githubClient)

	versionFetcher := gateways.NewVersionFetcher(filePathGit)

	bl := businesslogics.New(
		versionFetcher,
		releaseRepository,
	)

	uc := usecases.New(
		bl,
	)
	return uc, nil
}
