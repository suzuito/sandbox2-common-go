package httpfakeserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/domain/mock"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/handler/admin"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/handler/fakeserver"
)

type Request = mock.Request
type Response = mock.Response
type Mock = mock.Mock
type Mocks = mock.Mocks

type Options struct {
	Port          int
	BasePathAdmin string
}

func Main(o Options) int {
	port := 8080
	if o.Port != 0 {
		port = o.Port
	}

	basePathAdmin := "/admin"
	if o.BasePathAdmin != "" {
		basePathAdmin = o.BasePathAdmin
	}

	caseRepository := mock.NewRepository()

	mux := http.NewServeMux()

	mux.HandleFunc(
		fmt.Sprintf("GET %s/health", basePathAdmin),
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)
	mux.HandleFunc(
		fmt.Sprintf("POST %s/cases", basePathAdmin),
		admin.PostAdminCase(caseRepository),
	)
	mux.HandleFunc(
		fmt.Sprintf("DELETE %s/cases", basePathAdmin),
		admin.DeleteAdminCase(caseRepository),
	)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/cases", basePathAdmin),
		admin.GetAdminCase(caseRepository),
	)
	mux.HandleFunc(
		"/",
		fakeserver.HandleFunc(caseRepository),
	)

	exitCode := utils.RunHandlerWithGracefulShutdown(
		context.Background(),
		mux,
		port,
		utils.Options{
			WaitSecondsUntilGracefulShutdownIsStarted:   1,
			GracefulShutdownTimeoutSeconds:              1,
			ForcefullyRequestCancellationTimeoutSeconds: 1,
		},
	)

	return exitCode.Int()
}
