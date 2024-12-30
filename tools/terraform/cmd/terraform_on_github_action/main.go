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

	var dirPathBase string
	var projectID string
	var githubContextJSON string
	var dirPathRootGit string

	flag.StringVar(&dirPathBase, "d", "", "Base directory path")
	flag.StringVar(&projectID, "p", "", "Google Cloud Project ID")
	flag.StringVar(&githubContextJSON, "github-context", "", "Context object json of Github Action https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs#github-context")
	flag.StringVar(&dirPathRootGit, "git-rootdir", "", "Base directory path of git")

	flag.Parse()

	if dirPathBase == "" {
		fmt.Fprintln(os.Stderr, "-d is required")
		os.Exit(1)
	}

	if projectID == "" {
		fmt.Fprintln(os.Stderr, "-p is required")
		os.Exit(1)
	}

	if githubContextJSON == "" {
		fmt.Fprintln(os.Stderr, "-github-context is required")
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
		projectID,
		githubContextJSON,
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
	if err := uc.TerraformOnGithubAction(
		ctx,
		dirPathBase,
		dirPathRootGit,
		projectID,
		arg,
	); err != nil {
		if clierr, ok := errordefcli.AsCLIError(err); ok {
			fmt.Fprintln(os.Stderr, clierr.Error())
			os.Exit(clierr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(125)
	}
}
