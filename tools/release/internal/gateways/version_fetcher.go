package gateways

import (
	"context"
	"errors"

	"github.com/Masterminds/semver/v3"
)

var ErrNoVersionExists = errors.New("no version exists")

type VersionFetcher interface {
	GetLatestVersion(ctx context.Context, prefix string) (*semver.Version, error)
}
