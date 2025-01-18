package utils

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// RunHTTPServerWithGracefulShutdown shutdowns server with graceful shutdown.
// When first argument ctx's Done channel is closed, graceful shutdown starts.
// Returned integer is exit code.
func RunHTTPServerWithGracefulShutdown(
	ctx context.Context,
	server *http.Server,
	logger *slog.Logger,
) int {
	chGracefulShutdown := make(chan error)
	defer close(chGracefulShutdown)
	go func() {
		// Signalのハンドラー
		// SIGINT,SIGTERMをキャッチした後、ctx.Doneが制御を返す
		<-ctx.Done()
		logger.Info("start graceful shut down")
		ctxSignalHandler, cancel := context.WithTimeout(context.Background(), time.Second*100) // 100秒待ってもserver.Shutdown(ctx)が返ってこない場合、強制的にシャットダウンする
		defer cancel()
		// Graceful Shutdownをスタートする。
		// Graceful Shutdownが成功したら、server.Shutdown(ctxSignalHandler)はnilを返す。
		// Graceful Shutdownが失敗したら、server.Shutdown(ctxSignalHandler)は非nilを返す。
		chGracefulShutdown <- server.Shutdown(ctxSignalHandler)
	}()

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server.ListenAndServe() is failed", "err", err)
			return 1
		}
	}
	if err := <-chGracefulShutdown; err != nil {
		logger.Error("graceful shut down is failed (server.Shutdown(ctx) is failed)", "err", err)
		return 2
	}

	logger.Info("graceful shut down is complete")
	return 0
}
