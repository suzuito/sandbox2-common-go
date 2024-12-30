package inject

import (
	"net/http"
	"os"

	"github.com/google/go-github/v68/github"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/infra/local/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/infra/local/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/usecases"
)

func NewUsecase(env *Environment) usecases.Usecase {
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
		WithAuthToken(env.GithubAppToken)

	terraform := gateways.NewTerraformGateway(
		env.FilePathTerraform,
		os.Stdout,
		os.Stderr,
	)

	return usecases.New(
		businesslogics.New(
			reporter.New(),
			githubClient.PullRequests,
			terraform,
		),
	)
}
