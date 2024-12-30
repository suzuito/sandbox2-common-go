package module

import "github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/file"

type ModulePath string

func (t *ModulePath) String() string {
	return string(*t)
}

type Module struct {
	AbsPath ModulePath
	Files   []*file.File
	IsRoot  bool
}

func (t *Module) GoogleProjectID() (string, bool) {
	for _, f := range t.Files {
		for _, p := range f.Providers {
			if p.Name == "google" {
				return p.Project, true
			}
		}
	}
	return "", false
}

type Modules []*Module

func (t Modules) Len() int           { return len(t) }
func (t Modules) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Modules) Less(i, j int) bool { return t[i].AbsPath < t[j].AbsPath }
