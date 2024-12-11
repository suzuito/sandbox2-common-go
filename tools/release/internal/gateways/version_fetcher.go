package gateways

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type VersionFetcher interface {
	GetLatestVersion(ctx context.Context) (*semver.Version, error)
}
