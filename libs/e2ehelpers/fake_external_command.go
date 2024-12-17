package e2ehelpers

import (
	"errors"
	"fmt"
	"os"
	"text/template"
)

var fakeExternalCommandTemplate = template.Must(template.New("fakeCommand").Parse(`#!/bin/sh

{{if gt (len .Stdout) 0}}
cat << EOF
{{.Stdout}}
EOF
{{end}}

{{if gt (len .Stderr) 0}}
cat << EOF >&2
{{.Stderr}}
EOF
{{end}}

exit {{.ExitCode}}

`))

type FakeExternalCommand struct {
	filePath string
}

func (t *FakeExternalCommand) FilePath() string {
	return t.filePath
}

func (t *FakeExternalCommand) Cleanup() error {
	return os.Remove(t.filePath)
}

type NewFakeExternalCommandArg struct {
	FilePath string
	ExitCode int
	Stdout   string
	Stderr   string
}

func NewFakeExternalCommand(arg *NewFakeExternalCommandArg) (*FakeExternalCommand, error) {
	f, err := os.Create(arg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	if err := f.Chmod(0755); err != nil {
		return nil, fmt.Errorf("failed to chmod: %w", err)
	}

	if err := fakeExternalCommandTemplate.Execute(f, arg); err != nil {
		return nil, fmt.Errorf("failed to template: %w", err)
	}

	return &FakeExternalCommand{
		filePath: f.Name(),
	}, nil
}

type ExternalCommandFaker struct {
	commands []*FakeExternalCommand
}

func (t *ExternalCommandFaker) Cleanup() error {
	var err error
	for _, c := range t.commands {
		err = errors.Join(c.Cleanup())
	}

	return err
}

func (t *ExternalCommandFaker) New(arg *NewFakeExternalCommandArg) (*FakeExternalCommand, error) {
	cmd, err := NewFakeExternalCommand(arg)
	if err != nil {
		return nil, err
	}

	t.commands = append(t.commands, cmd)

	return cmd, nil
}
