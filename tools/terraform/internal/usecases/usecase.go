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
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type Usecase interface {
	CheckTerraformRules(
		ctx context.Context,
		dirPathBase string,
		rules rule.Rules,
	) error
	TerraformInPR(
		ctx context.Context,
		dirPathBase string,
		dirPathRootGit string,
		githubOwner string,
		githubRepo string,
		githubPRNumber int,
		planOnly bool,
		autoMerge bool,
	) error
	TerraformPlanAllModules(
		ctx context.Context,
		dirPathBase string,
		dirPathRootGit string,
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
			return errordefcli.Errorf(
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
		return errordefcli.Errorf(5, "not pass")
	}

	return nil
}

func (t *impl) TerraformInPR(
	ctx context.Context,
	dirPathBase string,
	dirPathRootGit string,
	githubOwner string,
	githubRepo string,
	githubPRNumber int,
	planOnly bool,
	autoMerge bool,
) error {
	modules, err := t.businessLogic.ParseBaseDir(ctx, dirPathBase)
	if err != nil {
		return terrors.Wrap(err)
	}

	paths, err := t.businessLogic.FetchPathsChangedInPR(
		ctx,
		githubOwner,
		githubRepo,
		githubPRNumber,
	)
	if err != nil {
		return terrors.Wrap(err)
	}

	absPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		absPath, err := filepath.Abs(filepath.Join(dirPathRootGit, p))
		if err != nil {
			return terrors.Wrap(err)
		}
		absPaths = append(absPaths, absPath)
	}

	modules, err = filterModulesByTargetAbsFilePaths(modules, absPaths)
	if err != nil {
		return terrors.Wrap(err)
	}

	if len(modules) <= 0 {
		fmt.Printf("no file changed in PR: %d\n", githubPRNumber)
		return nil
	}

	for _, module := range modules {
		if err := t.businessLogic.TerraformInit(ctx, module); err != nil {
			return terrors.Wrap(err)
		}
	}

	diff := false
	results := []string{}
	for _, module := range modules {
		planResult, err := t.businessLogic.TerraformPlan(ctx, module)
		if err != nil {
			return terrors.Wrap(err)
		}

		if planResult.IsPlanDiff {
			diff = true
		}

		results = append(results, planResult.String())
	}

	isMergeable := false
	if !planOnly {
		isMergeable, err = t.businessLogic.IsPRMergeable(
			ctx,
			githubOwner,
			githubRepo,
			githubPRNumber,
		)
		if err != nil {
			return terrors.Wrap(err)
		}

		if !isMergeable {
			fmt.Println("pr is not mergeable")
			results = append(results, "pr is not mergeable")
		} else {
			for _, module := range modules {
				applyResult, err := t.businessLogic.TerraformApply(ctx, module)
				if err != nil {
					return terrors.Wrap(err)
				}

				results = append(results, applyResult.String())
			}
		}
	}

	if err := t.businessLogic.CommentResults(
		ctx,
		githubOwner,
		githubRepo,
		githubPRNumber,
		results,
	); err != nil {
		return terrors.Wrap(err)
	}

	if !planOnly && autoMerge {
		if err := t.businessLogic.MergePR(
			ctx,
			githubOwner,
			githubRepo,
			githubPRNumber,
		); err != nil {
			return terrors.Wrap(err)
		}
		fmt.Println("PR is merged")
	}

	switch {
	case planOnly && diff:
		return errordefcli.Errorf(2, "diff at `terraform plan`")
	case !planOnly && !isMergeable:
		return errordefcli.Errorf(3, "cannot exec `terraform apply` because PR is not mergeable")
	}

	return nil
}

func (t *impl) TerraformPlanAllModules(
	ctx context.Context,
	dirPathBase string,
	dirPathRootGit string,
) error {
	modules, err := t.businessLogic.ParseBaseDir(ctx, dirPathBase)
	if err != nil {
		return terrors.Wrap(err)
	}

	modules = slices.Collect(utils.Filter(
		func(m *module.Module) bool { return m.IsRoot },
		slices.Values(modules),
	))

	if len(modules) <= 0 {
		return nil
	}

	for _, module := range modules {
		if err := t.businessLogic.TerraformInit(ctx, module); err != nil {
			return terrors.Wrap(err)
		}
	}

	diff := false
	for _, module := range modules {
		planResult, err := t.businessLogic.TerraformPlan(ctx, module)
		if err != nil {
			return terrors.Wrap(err)
		}

		if planResult.IsPlanDiff {
			diff = true
		}
	}

	if diff {
		return errordefcli.Errorf(2, "diff at `terraform plan`")
	}

	return nil
}

func filterModulesByTargetAbsFilePaths(modules module.Modules, targetAbsFilePaths []string) (module.Modules, error) {
	modulesByAbsPath := map[module.ModulePath]*module.Module{}
	for _, m := range modules {
		modulesByAbsPath[m.AbsPath] = m
	}

	// source mod -> parent mod
	moduleParentAbsPaths := map[module.ModulePath][]module.ModulePath{}
	for _, mod := range modules {
		for _, file := range mod.Files {
			for _, moduleRef := range file.Modules {
				absSourceString := filepath.Join(mod.AbsPath.String(), moduleRef.Source)
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

	sort.Sort(filtered)
	filtered = slices.CompactFunc(
		filtered,
		func(l *module.Module, r *module.Module) bool {
			return l.AbsPath == r.AbsPath
		},
	)

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

	return r
}

func New(
	businessLogic businesslogics.BusinessLogic,
) *impl {
	return &impl{
		businessLogic: businessLogic,
	}
}
