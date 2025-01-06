package terraformexe

import "github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"

type PlanResult struct {
	Module     *module.Module
	IsPlanDiff bool
}

type PlanResults []*PlanResult

type ApplyResult struct{}

type ApplyResults []*ApplyResult
