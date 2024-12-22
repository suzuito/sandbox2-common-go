package usecases

import (
	"context"
	"io/fs"
	"path/filepath"

	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/module"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/rule"
)

type Usecase interface {
	CheckTerraformRules(
		ctx context.Context,
		dirPathBase string,
		rules rule.Rules,
	) error
}

type impl struct {
	businessLogic businesslogics.BusinessLogic
}

func (t *impl) CheckTerraformRules(
	ctx context.Context,
	dirPathBase string,
	rules rule.Rules,
) error {
	modules := module.Modules{}

	if err := filepath.WalkDir(dirPathBase, func(path string, d fs.DirEntry, errInArg error) error {
		if errInArg != nil {
			return errInArg
		}

		if !d.IsDir() {
			return nil
		}

		module, err := t.businessLogic.ParseDir(ctx, path)
		if err != nil {
			return terrors.Errorf("failed to parse dir: %w", err)
		} else if len(module.Files) <= 0 {
			return nil
		}

		modules = append(modules, module)

		return nil
	}); err != nil {
		return err
	}

	if result, err := t.businessLogic.CheckRules(ctx, dirPathBase, modules, rules); err != nil {
		return terrors.Errorf("failed to check rules: %w", err)
	} else if result {
		return errordefcli.NewCLIError(5, "not pass")
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
