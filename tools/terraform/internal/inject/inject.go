package inject

import (
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/businesslogics"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/infra/local/domains/reporter"
	"github.com/suzuito/sandbox2-common-go/tools/terraform/internal/usecases"
)

func NewUsecase() usecases.Usecase {
	return usecases.New(
		businesslogics.New(reporter.New()),
	)
}
