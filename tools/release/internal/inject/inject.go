package inject

import (
	"net/http"

	"github.com/google/go-github/v67/github"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/gh/repositories"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/infra/local/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/usecases"
)

func NewUsecase(
	filePathGit string,
	githubAppToken string,
) usecases.Usecase {
	githubClient := github.
		NewClient(http.DefaultClient).
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
	return uc
}
