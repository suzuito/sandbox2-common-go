package gateways

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformexe"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type terraformGateway struct {
	filePathBinTerraform string
	stdout               io.Writer
	stderr               io.Writer
}

func (t *terraformGateway) Init(
	ctx context.Context,
	module *module.Module,
) error {
	result, err := t.run(
		ctx,
		t.filePathBinTerraform,
		[]string{
			fmt.Sprintf("-chdir=%s", module.AbsPath),
			"init",
			"-no-color",
		},
	)
	if err != nil {
		return terrors.Wrap(err)
	} else if result.ExitCode != 0 {
		return fmt.Errorf("failed to init")
	}

	return nil
}

func (t *terraformGateway) Plan(
	ctx context.Context,
	module *module.Module,
) (*terraformexe.PlanResult, error) {
	result, err := t.run(
		ctx,
		t.filePathBinTerraform,
		[]string{
			fmt.Sprintf("-chdir=%s", module.AbsPath),
			"plan",
			"-no-color",
			"-detailed-exitcode",
		},
	)
	if err != nil {
		return nil, terrors.Wrap(err)
	} else if result.ExitCode != 0 && result.ExitCode != 2 {
		return nil, fmt.Errorf("failed to init")
	}

	return &terraformexe.PlanResult{
		IsPlanDiff: result.ExitCode == 2,
	}, nil
}

func (t *terraformGateway) Apply(
	ctx context.Context,
	module *module.Module,
) (*terraformexe.ApplyResult, error) {
	result, err := t.run(
		ctx,
		t.filePathBinTerraform,
		[]string{
			fmt.Sprintf("-chdir=%s", module.AbsPath),
			"apply",
			"-no-color",
			"-auto-approve",
		},
	)
	if err != nil {
		return nil, terrors.Wrap(err)
	} else if result.ExitCode != 0 {
		return nil, fmt.Errorf("failed to init")
	}

	return &terraformexe.ApplyResult{}, nil
}

type runResult struct {
	Cmd      string
	ExitCode int
}

func (t *terraformGateway) run(
	ctx context.Context,
	commandName string,
	args []string,
) (*runResult, error) {
	commandline := fmt.Sprintf("%s %s", commandName, strings.Join(args, " "))
	fmt.Fprintf(t.stdout, "\n")
	fmt.Fprintf(t.stdout, "==== CMD ====\n")
	fmt.Fprintf(t.stdout, "%s\n", commandline)
	fmt.Fprintf(t.stdout, "==== OUT ====\n")
	cmd := exec.CommandContext(
		ctx,
		commandName,
		args...,
	)
	envs := []string{}
	cmd.Stderr = t.stderr
	cmd.Stdout = t.stdout
	cmd.Env = append(cmd.Environ(), envs...)
	cmd.Run()
	fmt.Fprintf(t.stdout, "==== END ====\n")
	fmt.Fprintf(t.stdout, "exit with %d\n", cmd.ProcessState.ExitCode())
	fmt.Fprintf(t.stdout, "\n")
	return &runResult{
		Cmd:      commandline,
		ExitCode: cmd.ProcessState.ExitCode(),
	}, nil
}

func NewTerraformGateway(
	filePathBinTerraform string,
	stdout io.Writer,
	stderr io.Writer,
) *terraformGateway {
	return &terraformGateway{
		filePathBinTerraform: filePathBinTerraform,
		stdout:               stdout,
		stderr:               stderr,
	}
}
