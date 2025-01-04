package e2ehelpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CLITestCaseV2 struct {
	Desc       string
	Setup      func(t *testing.T, testID TestID, input *CLITestCaseV2Input, expected *CLITestCaseV2Expected)
	Teardown   func(t *testing.T, testID TestID)
	Assertions func(t *testing.T)
}

func (c *CLITestCaseV2) Run(
	t *testing.T,
	filePathBin string,
) {
	ctx := context.Background()

	if c.Setup == nil {
		panic(errors.New("setup function is required"))
	}

	testID := NewTestID()
	input := CLITestCaseV2Input{}
	expected := CLITestCaseV2Expected{}
	c.Setup(t, testID, &input, &expected)

	cmd := exec.CommandContext(
		ctx,
		filePathBin,
		input.Args...,
	)

	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("E2E_TEST_ID=%s", testID.String()),
	)
	cmd.Env = append(cmd.Env, input.Envs...)

	stdout, stderr := bytes.NewBufferString(""), bytes.NewBufferString("")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		var exiterr *exec.ExitError
		if !errors.As(err, &exiterr) {
			require.NoError(t, err, fmt.Sprintf("%s %s: %s", filePathBin, strings.Join(input.Args, " "), err.Error()))
		}
	}

	assert.Equal(t, expected.ExitCode, cmd.ProcessState.ExitCode())
	assert.Equal(t, strings.TrimRight(expected.Stdout, "\n"), strings.TrimRight(stdout.String(), "\n"), "unexpected stdout")
	assert.Equal(t, strings.TrimRight(expected.Stderr, "\n"), strings.TrimRight(stderr.String(), "\n"), "unexpected stderr")

	if c.Assertions != nil {
		c.Assertions(t)
	}

	if c.Teardown != nil {
		c.Teardown(t, testID)
	}
}

type CLITestCaseV2Input struct {
	Envs []string
	Args []string
}

type CLITestCaseV2Expected struct {
	Stdout   string
	Stderr   string
	ExitCode int
}
