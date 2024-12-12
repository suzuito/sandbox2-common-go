package repositories

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v67/github"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/repositories"
	"github.com/suzuito/sandbox2-common-go/utils"
)

type ReleaseRepository struct {
	githubClient *github.Client
}

var _ repositories.ReleaseRepository = &ReleaseRepository{}

func (t *ReleaseRepository) CreateDraft(
	ctx context.Context,
	githubOwner string,
	githubRepo string,
	branch string,
	prefix string,
	version *semver.Version,
) error {
	versionString := fmt.Sprintf("%s%s", prefix, version.String())
	_, _, err := t.githubClient.Repositories.CreateRelease(
		ctx,
		githubOwner,
		githubRepo,
		&github.RepositoryRelease{
			TagName:              utils.Ptr(versionString),
			TargetCommitish:      utils.Ptr(branch),
			Name:                 utils.Ptr(versionString),
			Draft:                utils.Ptr(true),
			GenerateReleaseNotes: utils.Ptr(true),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	return nil
}

func NewReleaseRepository(
	githubClient *github.Client,
) *ReleaseRepository {
	return &ReleaseRepository{
		githubClient: githubClient,
	}
}
