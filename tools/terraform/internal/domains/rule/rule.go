package rule

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/domains/terraformmodels/module"
)

type Rule interface {
	Name() string
	Check(ctx context.Context, dirPathBaes string, modules module.Modules, reporter reporter.Reporter) (bool, error)
}

type Rules []Rule
