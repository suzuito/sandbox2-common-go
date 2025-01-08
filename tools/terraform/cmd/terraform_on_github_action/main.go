package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformexe"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/inject"
)

var usageString = `Terraform command wrapper for Github Action.

When this command is used?
- When a user comments 'terraform plan' on a GitHub Pull request, this command will run 'terraform plan' for changed files in PR.
- When a user comments 'terraform apply' on a GitHub Pull request, this command will run 'terraform apply' for changed files in PR.
- When a commit on GitHub Pull request is created, this command will run 'terraform plan' for changed files in PR.

What is exit code on this command?
- 0: terraform command is sucessed with empty diff
- 1: command line arg error
- 2: terraform command is sucessed with non-empty diff
- others: unknown errors

`

func usage() {
	fmt.Fprintln(os.Stderr, usageString)
	flag.PrintDefaults()
}

func main() {
	var env inject.Environment
	if err := envconfig.Process("", &env); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load environment variable: %v\n", err)
		os.Exit(1)
	}

	var eventName string
	var eventPath string
	var dirPathBase string
	var dirPathRootGit string
	var autoMerge bool

	flag.StringVar(&eventName, "event-name", "", "Event name of GitHub Action")
	flag.StringVar(&eventPath, "event-path", "", "Event path of GitHub Action")
	flag.StringVar(&dirPathBase, "d", "", "Base directory path")
	flag.StringVar(&dirPathRootGit, "git-rootdir", "", "Base directory path of git")
	flag.BoolVar(&autoMerge, "automerge", false, "Automerge PR after apply is succeeded")
	flag.Usage = usage

	flag.Parse()

	if eventName == "" {
		fmt.Fprintln(os.Stderr, "-event-name is required")
		os.Exit(1)
	}

	if eventPath == "" {
		fmt.Fprintln(os.Stderr, "-event-path is required")
		os.Exit(1)
	}

	if dirPathBase == "" {
		fmt.Fprintln(os.Stderr, "-d is required")
		os.Exit(1)
	}

	if dirPathRootGit == "" {
		fmt.Fprintln(os.Stderr, "-git-rootdir is required")
		os.Exit(1)
	}

	if _, err := os.Stat(dirPathBase); err != nil {
		fmt.Fprintf(os.Stderr, "invalid base dir: %s: %s\n", dirPathBase, err)
		os.Exit(1)
	}

	arg, ok, err := terraformexe.NewTerraformExecutionArg(
		dirPathBase,
		eventName,
		eventPath,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	} else if !ok {
		fmt.Println("skipped")
		os.Exit(0)
	}

	ctx := context.Background()

	uc := inject.NewUsecase(&env)
	switch arg.TargetType {
	case terraformexe.PlanAll:
		err = uc.TerraformPlanAllModules(
			ctx,
			dirPathBase,
			dirPathRootGit,
		)
	case terraformexe.InPR:
		err = uc.TerraformInPR(
			ctx,
			dirPathBase,
			dirPathRootGit,
			arg.GitHubOwner,
			arg.GitHubRepository,
			arg.GitHubPullRequestNumber,
			arg.PlanOnly,
			autoMerge,
		)
	default:
		err = fmt.Errorf("target type is not supported: %d", arg.TargetType)
	}

	if err != nil {
		if clierr, ok := errordefcli.AsCLIError(err); ok {
			fmt.Fprintln(os.Stderr, clierr.Error())
			os.Exit(clierr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(125)
	}
}
