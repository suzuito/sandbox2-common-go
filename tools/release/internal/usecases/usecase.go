package usecases

import (
	"context"

	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/domains"
)

type Usecase interface {
	IncrementVersion(
		ctx context.Context,
		githubOwner string,
		githubRepo string,
		branch string,
		prefix string,
		incrementTypeString string,
	) error
}

type Impl struct {
	businessLogic businesslogics.BusinessLogic
}

func (t *Impl) IncrementVersion(
	ctx context.Context,
	githubOwner string,
	githubRepo string,
	branch string,
	prefix string,
	incrementTypeString string,
) error {
	incrementType := domains.IncrementType(incrementTypeString)
	if err := incrementType.Validate(); err != nil {
		return errordefcli.Errorf(
			1,
			"invalid increment type '%s'",
			incrementType,
		)
	}

	if err := t.businessLogic.IncrementVersion(
		ctx,
		githubOwner,
		githubRepo,
		branch,
		prefix,
		incrementType,
	); err != nil {
		return terrors.Errorf("failed to IncrementVersion: %w", err)
	}

	return nil
}

func New(
	businessLogic businesslogics.BusinessLogic,
) *Impl {
	return &Impl{
		businessLogic: businessLogic,
	}
}
