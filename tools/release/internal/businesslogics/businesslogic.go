package businesslogics

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/suzuito/sandbox2-common-go/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/domains"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/gateways"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/repositories"
)

type BusinessLogic interface {
	IncrementVersion(
		ctx context.Context,
		prefix string,
		incrementType domains.IncrementType,
	) error
}

type impl struct {
	versionFetcher    gateways.VersionFetcher
	releaseRepository repositories.ReleaseRepository
}

func (t *impl) IncrementVersion(
	ctx context.Context,
	prefix string,
	incrementType domains.IncrementType,
) error {
	latestVersion, err := t.versionFetcher.GetLatestVersion(ctx)
	if err != nil {
		return terrors.Wrapf("failed to GetLatestVersion: %w", err)
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
		return terrors.Wrapf("invalid latest version: %s", latestVersion)
	}

	if err := t.releaseRepository.CreateDraft(ctx, &nextVersion); err != nil {
		return terrors.Wrapf("failed to CreateDraft: %w", err)
	}

	return nil
}

func New(
	versionFetcher gateways.VersionFetcher,
	releaseRepository repositories.ReleaseRepository,
) *impl {
	return &impl{
		versionFetcher:    versionFetcher,
		releaseRepository: releaseRepository,
	}
}
