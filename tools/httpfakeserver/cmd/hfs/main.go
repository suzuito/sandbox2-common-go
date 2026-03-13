package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/domain/mock"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/handler/admin"
	"github.com/suzuito/sandbox2-common-go/tools/httpfakeserver/internal/handler/fakeserver"
)

func main() {
	port := 8080
	portString := os.Getenv("PORT")
	if portString != "" {
		var err error
		port, err = strconv.Atoi(portString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert port into int\n")
			os.Exit(1)
		}
	}

	basePathAdmin := os.Getenv("BASE_PATH_ADMIN")
	if basePathAdmin == "" {
		basePathAdmin = "/admin"
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

	os.Exit(exitCode.Int())
}
