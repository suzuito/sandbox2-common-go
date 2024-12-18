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

// FileExternalCommand is a object having a location of fake command file.
type fakeExternalCommand struct {
	filePath string
}

// Cleanup deletes fake command file
func (t *fakeExternalCommand) Cleanup() error {
	return os.Remove(t.filePath)
}

// newFakeExternalCommand returns new [fakeExternalCommand].
func newFakeExternalCommand(arg *ExternalCommandBehavior) (*fakeExternalCommand, error) {
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

	return &fakeExternalCommand{
		filePath: f.Name(),
	}, nil
}

type ExternalCommandBehavior struct {
	FilePath string
	ExitCode int
	Stdout   string
	Stderr   string
}

type ExternalCommandFaker struct {
	commands []*fakeExternalCommand
}

func (t *ExternalCommandFaker) Cleanup() error {
	var err error
	for _, c := range t.commands {
		err = errors.Join(c.Cleanup())
	}

	t.commands = []*fakeExternalCommand{}

	return err
}

func (t *ExternalCommandFaker) Add(arg *ExternalCommandBehavior) error {
	cmd, err := newFakeExternalCommand(arg)
	if err != nil {
		return err
	}

	for _, command := range t.commands {
		if command.filePath == arg.FilePath {
			return fmt.Errorf("fake command already exists: %s", arg.FilePath)
		}
	}

	t.commands = append(t.commands, cmd)

	return nil
}
