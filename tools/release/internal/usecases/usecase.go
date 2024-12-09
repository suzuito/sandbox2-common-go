package usecases

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
)

type impl struct {
}

func (t *impl) ValidateReleaseVersion(
	ctx context.Context,
	versionString string,
) error {
	semver.NewVersion(versionString)
	return fmt.Errorf("not impl")
}

func New() *impl {
	return &impl{}
}
