package repositories

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type ReleaseRepository interface {
	CreateDraft(ctx context.Context, version *semver.Version) error
}
