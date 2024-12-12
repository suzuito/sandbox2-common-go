package repositories

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type ReleaseRepository interface {
	CreateDraft(
		ctx context.Context,
		githubOwner string,
		githubRepo string,
		branch string,
		prefix string,
		version *semver.Version,
	) error
}
