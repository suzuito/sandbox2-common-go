package gateways

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformexe"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type TerraformGateway interface {
	Init(
		ctx context.Context,
		module *module.Module,
	) error
	Plan(
		ctx context.Context,
		modules *module.Module,
	) (*terraformexe.PlanResult, error)
	Apply(
		ctx context.Context,
		modules *module.Module,
	) (*terraformexe.ApplyResult, error)
}
