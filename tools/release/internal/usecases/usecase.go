package usecases

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/release/internal/domains"
)

type impl struct {
	businessLogic businesslogics.BusinessLogic
}

func (t *impl) IncrementVersion(
	ctx context.Context,
	prefix string,
	incrementType domains.IncrementType,
) error {
	if err := t.businessLogic.IncrementVersion(ctx, prefix, incrementType); err != nil {
		return terrors.Wrapf("failed to IncrementVersion: %w", err)
	}

	return nil
}

func New(
	businessLogic businesslogics.BusinessLogic,
) *impl {
	return &impl{
		businessLogic: businessLogic,
	}
}
