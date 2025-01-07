package gateways

import (
	"bytes"
	"context"
	"errors"
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
		return nil, fmt.Errorf("failed to plan")
	}

	return &terraformexe.PlanResult{
		IsPlanDiff: result.ExitCode == 2,
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
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
		return nil, fmt.Errorf("failed to apply")
	}

	return &terraformexe.ApplyResult{
		Stdout: result.Stdout,
		Stderr: result.Stderr,
	}, nil
}

type runResult struct {
	Cmd      string
	ExitCode int
	Stdout   string
	Stderr   string
}

func (t *terraformGateway) run(
	ctx context.Context,
	commandName string,
	args []string,
) (*runResult, error) {
	stdoutBuffer := bytes.NewBufferString("")
	stderrBuffer := bytes.NewBufferString("")
	stdout := io.MultiWriter(stdoutBuffer, t.stdout)
	stderr := io.MultiWriter(stderrBuffer, t.stderr)

	commandline := fmt.Sprintf("%s %s", commandName, strings.Join(args, " "))
	fmt.Fprintf(stdout, "\n")
	fmt.Fprintln(stdout, "*************")
	fmt.Fprintln(stdout, "*************")
	fmt.Fprintln(stdout, "*************")
	fmt.Fprintf(stdout, "==== CMD ====\n")
	fmt.Fprintf(stdout, "%s\n", commandline)
	fmt.Fprintf(stdout, "==== OUT ====\n")
	cmd := exec.CommandContext(
		ctx,
		commandName,
		args...,
	)
	envs := []string{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	cmd.Env = append(cmd.Environ(), envs...)
	if err := cmd.Run(); err != nil {
		var exiterr *exec.ExitError
		if !errors.As(err, &exiterr) {
			return nil, terrors.Errorf("failed to cmd.Run: %w", err)
		}
	}
	fmt.Fprintf(stdout, "==== END ====\n")
	fmt.Fprintf(stdout, "exit with %d\n", cmd.ProcessState.ExitCode())
	fmt.Fprintf(stdout, "\n")
	return &runResult{
		Cmd:      commandline,
		ExitCode: cmd.ProcessState.ExitCode(),
		Stdout:   stdoutBuffer.String(),
		Stderr:   stderrBuffer.String(),
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
