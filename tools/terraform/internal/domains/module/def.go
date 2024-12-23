package module

import "github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/file"

type Module struct {
	Path   string
	Files  []*file.File
	IsRoot bool
}

type Modules []*Module
