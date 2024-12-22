package businesslogics

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/file"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/module"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/rule"
)

type BusinessLogic interface {
	ParseDir(
		ctx context.Context,
		path string,
	) (*module.Module, error)
	CheckRules(
		ctx context.Context,
		dirPathBase string,
		modules module.Modules,
		rules rule.Rules,
	) (bool, error)
}

type impl struct {
	Reporter reporter.Reporter
}

func (t *impl) ParseDir(
	ctx context.Context,
	path string,
) (*module.Module, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, terrors.Errorf("failed to os.ReadDir: %s: %w", path, err)
	}

	module := module.Module{
		Path: path,
	}
	for _, entry := range entries {
		if entry.Name() == ".terraform.lock.hcl" {
			module.IsRoot = true
			continue
		}

		if filepath.Ext(entry.Name()) != ".tf" {
			continue
		}

		filePath := filepath.Join(path, entry.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, terrors.Errorf("failed to os.ReadFile: %s: %w", filePath, err)
		}

		tffile := file.File{
			Path: filePath,
		}
		if err := hclsimple.Decode(filePath+".hcl", content, nil, &tffile); err != nil {
			_, ok := err.(hcl.Diagnostics)
			if !ok {
				return nil, terrors.Errorf("failed to hclsimple.Decode: %w", err)
			}
		}

		module.Files = append(module.Files, &tffile)
	}

	return &module, nil
}

func (t impl) CheckRules(
	ctx context.Context,
	dirPathBase string,
	modules module.Modules,
	rules rule.Rules,
) (bool, error) {
	var result bool = false
	for _, rule := range rules {
		resultEach, err := rule.Check(ctx, dirPathBase, modules, t.Reporter)
		if err != nil {
			return false, terrors.Errorf("failed to check rule: %w", err)
		}

		if resultEach {
			result = resultEach
		}
	}

	return result, nil
}

func New(
	reporter reporter.Reporter,
) *impl {
	return &impl{
		Reporter: reporter,
	}
}
