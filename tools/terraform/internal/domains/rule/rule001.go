package rule

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type Rule001 struct{}

func (t *Rule001) Name() string {
	return "rule001"
}

func (t *Rule001) Check(
	ctx context.Context,
	dirPathBaes string,
	modules module.Modules,
	reporter reporter.Reporter,
) (bool, error) {
	var result bool = true
	for _, module := range modules {
		if !module.IsRoot {
			continue
		}

		hasTerraformBackendGCS := false
		terraformBackendBucket := ""
		terraformBackendPrefix := ""
		hasProviderGoogle := false
		providerGoogleProject := ""

		for _, file := range module.Files {
			for _, terraform := range file.Terraforms {

				if terraform.Backend != nil && terraform.Backend.Name == "gcs" {
					hasTerraformBackendGCS = true
					terraformBackendBucket = terraform.Backend.Bucket
					terraformBackendPrefix = terraform.Backend.Prefix
				}
			}

			for _, provider := range file.Providers {
				if provider.Name == "google" {
					hasProviderGoogle = true
					providerGoogleProject = provider.Project
				}
			}
		}

		dirAbsPathBase, err := filepath.Abs(dirPathBaes)
		if err != nil {
			return false, terrors.Errorf("invalid filepath.Abs: %w", err)
		}
		dirPathRel, err := filepath.Rel(dirAbsPathBase, module.AbsPath.String())
		if err != nil {
			return false, terrors.Errorf("invalid filepath.Rel: %w", err)
		}

		if !reporter.AssertTruef(
			module.AbsPath.String(),
			hasTerraformBackendGCS,
			`resource terraform.backend."gcs" not found`,
		) {
			result = false
		}

		if hasProviderGoogle && hasTerraformBackendGCS {
			if !reporter.AssertEqualf(
				module.AbsPath.String(),
				fmt.Sprintf("%s-terraform", providerGoogleProject),
				terraformBackendBucket,
				"invalid terraform.backend.\"gcs\".bucket",
			) {
				result = false
			}

			if !reporter.AssertEqualf(
				module.AbsPath.String(),
				dirPathRel,
				terraformBackendPrefix,
				"invalid terraform.backend.\"gcs\".prefix",
			) {
				result = false
			}
		}

		if !reporter.AssertTruef(
			module.AbsPath.String(),
			hasProviderGoogle,
			`resource provider."google" not found`,
		) {
			result = false
		}
	}

	return result, nil
}
