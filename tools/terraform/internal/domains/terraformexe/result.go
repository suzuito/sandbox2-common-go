package terraformexe

import "github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"

type PlanResult struct {
	Module     *module.Module
	IsPlanDiff bool
}

type PlanResults []*PlanResult

func (t *PlanResults) ModulesApplied() module.Modules {
	mods := module.Modules{}
	for _, a := range *t {
		mods = append(mods, a.Module)
	}
	return mods
}

type ApplyResult struct{}

type ApplyResults []*ApplyResult
