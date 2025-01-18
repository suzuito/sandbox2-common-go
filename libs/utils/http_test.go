package utils_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
)

func TestRunHTTPServerWithGracefulShutdown(t *testing.T) {
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)

	server := http.Server{Addr: ":8888"}
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})
	logger := slog.New(handler)

	var calledShutdown atomic.Bool
	server.RegisterOnShutdown(func() {
		calledShutdown.Store(true)
	})

	var exitCodeReturned atomic.Int64
	exitCodeReturned.Store(-1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		exitCodeReturned.Store(int64(utils.RunHTTPServerWithGracefulShutdown(ctxWithCancel, &server, logger)))
	}()

	// before graceful shutdown started
	require.False(t, calledShutdown.Load())
	require.Equal(t, int64(-1), exitCodeReturned.Load())

	cancel()
	wg.Wait()

	// after graceful shutdown completed
	require.True(t, calledShutdown.Load())
	require.Equal(t, int64(0), exitCodeReturned.Load())
}
