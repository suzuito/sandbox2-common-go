package usecases

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/file"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

var modules = module.Modules{
	{
		AbsPath: module.ModulePath("/roots/r1"),
		IsRoot:  true,
		Files: []*file.File{
			{
				AbsPath: "/roots/r1/f1.tf",
			},
			{
				AbsPath: "/roots/r1/f2.tf",
			},
			{
				AbsPath: "/roots/r1/f3.tf",
			},
		},
	},
	{
		AbsPath: module.ModulePath("/roots/r2"),
		IsRoot:  true,
		Files: []*file.File{
			{
				AbsPath: "/roots/r2/f1.tf",
				Modules: []*file.ModuleRef{
					{
						Source: "../../commons/r2m1",
					},
				},
			},
		},
	},
	{
		AbsPath: module.ModulePath("/roots/r3"),
		IsRoot:  true,
		Files: []*file.File{
			{
				AbsPath: "/roots/r3/f1.tf",
				Modules: []*file.ModuleRef{
					{
						Source: "../../commons/r3m1",
					},
				},
			},
		},
	},
	{
		AbsPath: module.ModulePath("/commons/r2m1"),
		Files: []*file.File{
			{
				AbsPath: "/commons/r2m1/f1.tf",
			},
			{
				AbsPath: "/commons/r2m1/f2.tf",
				Modules: []*file.ModuleRef{
					{
						Source: "../../commons/r2m1m1",
					},
				},
			},
		},
	},
	{
		AbsPath: module.ModulePath("/commons/r3m1"),
		Files: []*file.File{
			{
				AbsPath: "/commons/r3m1/f1.tf",
				Modules: []*file.ModuleRef{
					{
						Source: "../../commons/r3m1m1",
					},
				},
			},
		},
	},
	{
		AbsPath: module.ModulePath("/commons/r3m1m1"),
		Files: []*file.File{
			{
				AbsPath: "/commons/r3m1m1/f1.tf",
			},
		},
	},
	{
		AbsPath: module.ModulePath("/commons/unused"),
		Files: []*file.File{
			{
				AbsPath: "/commons/unused/f1.tf",
			},
		},
	},
}

func Test_filterModulesByTargetAbsFilePaths(t *testing.T) {
	cases := []struct {
		name                    string
		inputModules            module.Modules
		inputTargetAbsFilePaths []string
		expectedModulePaths     []string
		wantErr                 bool
		errMsg                  string
	}{
		{
			name:                    "ok - no targets",
			inputModules:            modules,
			inputTargetAbsFilePaths: []string{},
		},
		{
			name:         "ok - all targets are ignored",
			inputModules: modules,
			inputTargetAbsFilePaths: []string{
				"/foo/bar",
				"/hoge/fuga",
			},
		},
		{
			name:         "ok",
			inputModules: modules,
			inputTargetAbsFilePaths: []string{
				"/roots/r1/f1.tf",
				"/roots/r1/f3.tf",
				"/commons/r2m1/f1.tf",
				"/commons/r3m1m1/f1.tf",
				"/commons/unused/f1.tf",
			},
			expectedModulePaths: []string{
				"/roots/r1",
				"/roots/r2",
				"/roots/r3",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := filterModulesByTargetAbsFilePaths(c.inputModules, c.inputTargetAbsFilePaths)
			if c.wantErr {
				require.NoError(t, err)
			} else {
				actualModulePaths := make([]string, 0, len(actual))
				for _, m := range actual {
					actualModulePaths = append(actualModulePaths, m.AbsPath.String())
				}
				require.ElementsMatch(t, c.expectedModulePaths, actualModulePaths)
			}
		})
	}
}
