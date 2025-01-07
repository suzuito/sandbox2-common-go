package terraformexe

import "github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"

type PlanResult struct {
	Module     *module.Module
	IsPlanDiff bool
}

func (t *PlanResult) String() string {
	return "not impl"
}

type PlanResults []*PlanResult

type ApplyResult struct{}

func (t *ApplyResult) String() string {
	return "not impl"
}

type ApplyResults []*ApplyResult
