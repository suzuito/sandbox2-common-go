package businesslogics

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/domains"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/repositories"
)

type BusinessLogic interface {
	IncrementVersion(
		ctx context.Context,
		githubOwner string,
		githubRepo string,
		branch string,
		prefix string,
		incrementType domains.IncrementType,
	) error
}

type Impl struct {
	versionFetcher    gateways.VersionFetcher
	releaseRepository repositories.ReleaseRepository
}

var _ BusinessLogic = &Impl{}

func (t *Impl) IncrementVersion(
	ctx context.Context,
	githubOwner string,
	githubRepo string,
	branch string,
	prefix string,
	incrementType domains.IncrementType,
) error {
	latestVersion, err := t.versionFetcher.GetLatestVersion(ctx, prefix)
	if err != nil {
		return terrors.Errorf("failed to GetLatestVersion: %w", err)
	}

	var nextVersion semver.Version
	switch incrementType {
	case domains.IncrementTypeMajor:
		nextVersion = latestVersion.IncMajor()
	case domains.IncrementTypeMinor:
		nextVersion = latestVersion.IncMinor()
	case domains.IncrementTypePatch:
		nextVersion = latestVersion.IncPatch()
	default:
		return terrors.Errorf("invalid latest version: %s", latestVersion)
	}

	if err := t.releaseRepository.CreateDraft(
		ctx,
		githubOwner,
		githubRepo,
		branch,
		prefix,
		&nextVersion,
	); err != nil {
		return terrors.Errorf("failed to CreateDraft: %w", err)
	}

	return nil
}

func New(
	versionFetcher gateways.VersionFetcher,
	releaseRepository repositories.ReleaseRepository,
) *Impl {
	return &Impl{
		versionFetcher:    versionFetcher,
		releaseRepository: releaseRepository,
	}
}