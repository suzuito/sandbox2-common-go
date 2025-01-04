package usecases

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"

	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/rule"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformexe"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type Usecase interface {
	CheckTerraformRules(
		ctx context.Context,
		dirPathBase string,
		rules rule.Rules,
	) error
	TerraformOnGithubAction(
		ctx context.Context,
		dirPathBase string,
		dirPathRootGit string,
		projectID string,
		arg *terraformexe.Arg,
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

		module, ok, err := t.businessLogic.ParseDir(ctx, path)
		if err != nil {
			return terrors.Errorf("failed to parse dir: %w", err)
		} else if !ok {
			return nil
		}

		modules = append(modules, module)

		return nil
	}); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errordefcli.NewCLIErrorf(
				10,
				"%s does not exist",
				dirPathBase,
			)
		}

		return err
	}

	if result, err := t.businessLogic.CheckRules(ctx, dirPathBase, modules, rules); err != nil {
		return terrors.Errorf("failed to check rules: %w", err)
	} else if !result {
		return errordefcli.NewCLIError(5, "not pass")
	}

	return nil
}

func (t *impl) TerraformOnGithubAction(
	ctx context.Context,
	dirPathBase string,
	dirPathRootGit string,
	projectID string,
	arg *terraformexe.Arg,
) error {
	modules, err := t.businessLogic.ParseBaseDir(ctx, dirPathBase)
	if err != nil {
		return terrors.Wrap(err)
	}

	modules = slices.Collect(utils.Filter(
		func(m *module.Module) bool {
			if !m.IsRoot {
				return true
			}

			pid, exists := m.GoogleProjectID()
			if !exists {
				return false
			}

			return pid == projectID
		},
		slices.Values(modules),
	))

	if arg.TargetType == terraformexe.ForOnlyChageFiles {
		paths, err := t.businessLogic.FetchPathsChangedInPR(
			ctx,
			arg.GitHubOwner,
			arg.GitHubRepository,
			arg.GitHubPullRequestNumber,
		)
		if err != nil {
			return terrors.Wrap(err)
		}

		absPaths := make([]string, 0, len(paths))
		for _, p := range paths {
			absPaths = append(absPaths, filepath.Join(dirPathRootGit, p))
		}

		modules, err = filterModulesByTargetAbsFilePaths(modules, absPaths)
		if err != nil {
			return terrors.Wrap(err)
		}
	}

	if len(modules) <= 0 {
		fmt.Printf("no file changed in PR: %d\n", arg.GitHubPullRequestNumber)
		return nil
	}

	diff := false
	for _, module := range modules {
		if err := t.businessLogic.TerraformInit(ctx, module); err != nil {
			return terrors.Wrap(err)
		}

		planResult, err := t.businessLogic.TerraformPlan(ctx, module)
		if err != nil {
			return terrors.Wrap(err)
		}

		if planResult.IsPlanDiff {
			diff = true
		}
	}

	if arg.PlanOnly {
		if diff {
			return errordefcli.NewCLIError(2, "diff at `terraform plan`")
		}
		return nil
	}

	for _, module := range modules {
		if _, err := t.businessLogic.TerraformApply(ctx, module); err != nil {
			return terrors.Wrap(err)
		}
	}

	return nil
}

// TODO UT書く
func filterModulesByTargetAbsFilePaths(modules module.Modules, targetAbsFilePaths []string) (module.Modules, error) {
	for _, f := range targetAbsFilePaths {
		if !filepath.IsAbs(f) {
			return nil, terrors.Errorf("target file '%s' is not abs path", f)
		}
	}

	modulesByAbsPath := map[module.ModulePath]*module.Module{}
	for _, m := range modules {
		modulesByAbsPath[m.AbsPath] = m
	}

	moduleParentAbsPaths := map[module.ModulePath][]module.ModulePath{}
	for _, mod := range modules {
		for _, file := range mod.Files {
			for _, moduleRef := range file.Modules {
				source := filepath.Join(mod.AbsPath.String(), moduleRef.Source)
				absSourceString, err := filepath.Abs(source)
				if err != nil {
					return nil, terrors.Errorf("failed to convert source.path to abs path: %s: %w", source, err)
				}
				absSource := module.ModulePath(absSourceString)

				if _, exists := moduleParentAbsPaths[absSource]; !exists {
					moduleParentAbsPaths[absSource] = []module.ModulePath{}
				}
				moduleParentAbsPaths[absSource] = append(moduleParentAbsPaths[absSource], mod.AbsPath)
			}
		}
	}

	filtered := module.Modules{}

	for _, targetAbsFilePath := range targetAbsFilePaths {
		targetAbsPath := module.ModulePath(filepath.Dir(targetAbsFilePath))

		filtered = append(
			filtered,
			search(
				modulesByAbsPath,
				moduleParentAbsPaths,
				targetAbsPath,
			)...,
		)
	}

	return filtered, nil
}

func search(
	modulesByAbsPath map[module.ModulePath]*module.Module,
	moduleParentAbsPaths map[module.ModulePath][]module.ModulePath,
	target module.ModulePath,
) module.Modules {
	targetMod, exists := modulesByAbsPath[target]
	if !exists {
		// targetAbsPath is not terraform module
		return module.Modules{}
	}

	if targetMod.IsRoot {
		// targetAbsPath is root moudle
		return module.Modules{targetMod}
	}

	// targetAbsPaths is not root module
	parentPaths, exists := moduleParentAbsPaths[target]
	if !exists {
		// targetAbsPath is unused module
		return module.Modules{}
	}

	r := module.Modules{}
	for _, parentPath := range parentPaths {
		mods := search(modulesByAbsPath, moduleParentAbsPaths, parentPath)
		r = append(r, mods...)
	}

	sort.Sort(r)
	r = slices.CompactFunc(
		r,
		func(l *module.Module, r *module.Module) bool {
			return l.AbsPath == r.AbsPath
		},
	)

	return r
}

func New(
	businessLogic businesslogics.BusinessLogic,
) *impl {
	return &impl{
		businessLogic: businessLogic,
	}
}
