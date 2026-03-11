package hfs

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

var (
	filePathServerBin string
	targetURL         string
)

func defaultEnvs() []string {
	return []string{}
}

func healthCheck(ctx context.Context) func() error {
	return func() error {
		return e2ehelpers.CheckHTTPServerHealth(
			ctx,
			fmt.Sprintf("%s/admin/health", targetURL),
		)
	}
}

func testMain(m *testing.M) int {
	filePathServerBin = os.Getenv("FILE_PATH_SERVER_BIN")
	targetURL = fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))

	ctx := context.Background()

	shutdown, okHealthCheck := e2ehelpers.RunServer(
		ctx,
		filePathServerBin,
		&e2ehelpers.RunServerInput{
			Envs: defaultEnvs(),
		},
		healthCheck(ctx),
	)
	defer shutdown() // nolint:errcheck

	if !okHealthCheck {
		return 1
	}

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}
