package businesslogics

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v68/github"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/rule"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformexe"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/file"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/gateways"
)

type BusinessLogic interface {
	ParseDir(
		ctx context.Context,
		path string,
	) (*module.Module, bool, error)
	ParseBaseDir(
		ctx context.Context,
		path string,
	) (module.Modules, error)
	CheckRules(
		ctx context.Context,
		dirPathBase string,
		modules module.Modules,
		rules rule.Rules,
	) (bool, error)
	FetchPathsChangedInPR(
		ctx context.Context,
		owner string,
		repo string,
		pr int,
	) ([]string, error)
	IsPRMergeable(
		ctx context.Context,
		owner string,
		repo string,
		pr int,
	) (bool, error)
	CommentResults(
		ctx context.Context,
		owner string,
		repo string,
		issueNumber int,
		results []string,
	) error
	TerraformInit(
		ctx context.Context,
		module *module.Module,
	) error
	TerraformPlan(
		ctx context.Context,
		module *module.Module,
	) (*terraformexe.PlanResult, error)
	TerraformApply(
		ctx context.Context,
		module *module.Module,
	) (*terraformexe.ApplyResult, error)
}

type impl struct {
	Reporter                  reporter.Reporter
	GithubPullRequestsService gateways.GithubPullRequestsService
	GithubIssuesService       gateways.GithubIssuesService
	Terraform                 gateways.TerraformGateway
}

func (t *impl) ParseDir(
	ctx context.Context,
	path string,
) (*module.Module, bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, false, terrors.Errorf("failed to os.ReadDir: %s: %w", path, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, false, terrors.Errorf("failed to filepath.Abs: %s: %w", path, err)
	}

	module := module.Module{
		AbsPath: module.ModulePath(absPath),
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
			return nil, false, terrors.Errorf("failed to os.ReadFile: %s: %w", filePath, err)
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			return nil, false, terrors.Errorf("failed to filepath.Abs: %s: %w", filePath, err)
		}

		tffile := file.File{
			AbsPath: absFilePath,
		}
		if err := hclsimple.Decode(filePath+".hcl", content, nil, &tffile); err != nil {
			_, ok := err.(hcl.Diagnostics)
			if !ok {
				return nil, false, terrors.Errorf("failed to hclsimple.Decode: %w", err)
			}
		}

		module.Files = append(module.Files, &tffile)
	}

	if len(module.Files) <= 0 {
		return nil, false, nil
	}

	return &module, true, nil
}

func (t *impl) ParseBaseDir(
	ctx context.Context,
	basePath string,
) (module.Modules, error) {
	modules := module.Modules{}

	if err := filepath.Walk(basePath, func(path string, info fs.FileInfo, _ error) error {
		if !info.IsDir() {
			return nil
		}

		module, ok, err := t.ParseDir(ctx, path)
		if err != nil {
			return terrors.Wrap(err)
		} else if !ok {
			return nil
		}

		modules = append(modules, module)

		return nil
	}); err != nil {
		return modules, terrors.Wrap(err)
	}

	return modules, nil
}

func (t *impl) CheckRules(
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

func (t *impl) FetchPathsChangedInPR(
	ctx context.Context,
	owner string,
	repo string,
	pr int,
) ([]string, error) {
	returned := []string{}

	perPage := 100
	for page := 1; ; page++ {
		commitFiles, _, err := t.GithubPullRequestsService.ListFiles(
			ctx,
			owner,
			repo,
			pr,
			&github.ListOptions{Page: page, PerPage: perPage},
		)
		if err != nil {
			return []string{}, terrors.Errorf("failed to t.GithubPullRequestsService.ListFiles: %w", err)
		}

		for _, commitFile := range commitFiles {
			returned = append(returned, *commitFile.Filename)
		}

		if len(commitFiles) < perPage {
			break
		}
	}

	return returned, nil
}

func (t *impl) IsPRMergeable(
	ctx context.Context,
	owner string,
	repo string,
	prNumber int,
) (bool, error) {
	pr, _, err := t.GithubPullRequestsService.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return false, terrors.Errorf("failed to GithubPullRequestsService.Get: %w", err)
	}

	if pr.Mergeable == nil {
		return false, nil
	}

	return *pr.Mergeable, nil
}

func (t *impl) CommentResults(
	ctx context.Context,
	owner string,
	repo string,
	issueNumber int,
	results []string,
) error {
	bodyString := fmt.Sprintf(
		"```\n%s```\n",
		strings.Join(
			results,
			"\n----------------------------------------\n"+
				"----------------------------------------\n"+
				"----------------------------------------\n",
		),
	)

	if _, _, err := t.GithubIssuesService.CreateComment(
		ctx,
		owner,
		repo,
		issueNumber,
		&github.IssueComment{
			Body: &bodyString,
		},
	); err != nil {
		return terrors.Errorf("failed to GithubIssuesService.CreateComment: %w", err)
	}

	return nil
}

func (t *impl) TerraformInit(
	ctx context.Context,
	module *module.Module,
) error {
	if err := t.Terraform.Init(ctx, module); err != nil {
		return terrors.Errorf("failed to terraform.Init: %w", err)
	}
	return nil
}

func (t *impl) TerraformPlan(
	ctx context.Context,
	module *module.Module,
) (*terraformexe.PlanResult, error) {
	r, err := t.Terraform.Plan(ctx, module)
	if err != nil {
		return nil, terrors.Errorf("failed to terraform.Plan: %w", err)
	}
	return r, nil
}

func (t *impl) TerraformApply(
	ctx context.Context,
	module *module.Module,
) (*terraformexe.ApplyResult, error) {
	r, err := t.Terraform.Apply(ctx, module)
	if err != nil {
		return nil, terrors.Errorf("failed to terraform.Apply: %w", err)
	}
	return r, nil
}

func New(
	reporter reporter.Reporter,
	githubPullRequestsService gateways.GithubPullRequestsService,
	githubIssuesService gateways.GithubIssuesService,
	terraform gateways.TerraformGateway,
) *impl {
	return &impl{
		Reporter:                  reporter,
		GithubPullRequestsService: githubPullRequestsService,
		GithubIssuesService:       githubIssuesService,
		Terraform:                 terraform,
	}
}
