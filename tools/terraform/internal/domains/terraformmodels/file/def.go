package file

type File struct {
	AbsPath    string
	Terraforms []*Terraform `hcl:"terraform,block"`
	Providers  []*Provider  `hcl:"provider,block"`
	Modules    []*ModuleRef `hcl:"module,block"`
}

type Terraform struct {
	Backend *TerraformBackend `hcl:"backend,block"`
}

type TerraformBackend struct {
	Name   string `hcl:"name,label"`
	Bucket string `hcl:"bucket"`
	Prefix string `hcl:"prefix"`
}

type Provider struct {
	Name    string `hcl:"name,label"`
	Project string `hcl:"project"`
}

type ModuleRef struct {
	Name   string `hcl:"name,label"`
	Source string `hcl:"source"`
}
