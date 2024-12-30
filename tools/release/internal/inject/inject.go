package inject

import (
	"net/http"

	"github.com/google/go-github/v67/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/gh/repositories"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/local/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/usecases"
)

func NewUsecase(
	filePathGit string,
	githubAppToken string,
) (usecases.Usecase, error) {
	var env Environment
	if err := envconfig.Process("", &env); err != nil {
		return nil, terrors.Errorf("failed to load environment variable: %w", err)
	}

	githubHTTPClient := http.DefaultClient
	if env.E2ETestID != "" {
		githubHTTPClient = &http.Client{
			Transport: e2ehelpers.NewRoundTripperForE2E(
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
