package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	errordefcli "github.com/suzuito/sandbox2-common-go/libs/errordefs/cli"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/rule"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/inject"
)

var usageString = `This command verifies whether terraform files is preventing rules. If violation is occured, exit not zero.
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

	ctx := context.Background()

	var dirPathBase string

	flag.StringVar(&dirPathBase, "d", "", "base directory path")
	flag.Usage = usage

	flag.Parse()

	if dirPathBase == "" {
		fmt.Fprint(os.Stderr, "-d is required")
		os.Exit(1)
	}

	uc := inject.NewUsecase(&env)
	if err := uc.CheckTerraformRules(
		ctx,
		dirPathBase,
		rule.Rules{
			&rule.Rule001{},
		},
	); err != nil {
		if clierr, ok := errordefcli.AsCLIError(err); ok {
			fmt.Fprintln(os.Stderr, clierr.Error())
			os.Exit(clierr.ExitCode())
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(125)
	}
}
