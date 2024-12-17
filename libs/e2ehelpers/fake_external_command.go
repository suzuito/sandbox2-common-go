package e2ehelpers

import (
	"fmt"
	"os"
	"text/template"
)

var fakeExternalCommandTemplate = template.Must(template.New("fakeCommand").Parse(`
#!/bin/sh

{{if (len .Stdout) gt 0}}
cat << EOF
{{.Stdout}}
EOF
{{- end}}

{{if (len .Stderr) gt 0}}
cat << EOF >&2
{{.Stderr}}
EOF
{{- end}}

`))

type fakeExternalCommandTemplateVar struct {
	Stdout string
	Stderr string
}

type FakeExternalCommand struct {
	filePath string
}

func (t *FakeExternalCommand) FilePath() string {
	return t.filePath
}

func (t *FakeExternalCommand) Cleanup() error {
	return os.Remove(t.filePath)
}

func NewFakeExternalCommand(
	filePath string,
	exitCode int,
	stdout string,
	stderr string,
) (*FakeExternalCommand, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	if err := f.Chmod(0755); err != nil {
		return nil, fmt.Errorf("failed to chmod: %w", err)
	}

	if err := fakeExternalCommandTemplate.Execute(f, fakeExternalCommandTemplateVar{
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return nil, fmt.Errorf("failed to template: %w", err)
	}

	return &FakeExternalCommand{
		filePath: f.Name(),
	}, nil
}
