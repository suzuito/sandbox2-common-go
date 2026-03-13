package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// RunHTTPServerWithGracefulShutdown shutdowns server with graceful shutdown.
// When first argument ctx's Done channel is closed, graceful shutdown starts.
// Returned integer is exit code.
//
// Deprecated: use RunHandlerWithGracefulShutdown
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

type exitCode int

func (t exitCode) Int() int {
	return int(t)
}

const (
	shutodownGracefully exitCode = 0
	shutodownForcefully exitCode = 1
)

type Options struct {
	// シグナル受信後、http.Server.Shutdown メソッドが呼ばれるまでスリープする時間(秒) <- 待ち時間a
	WaitSecondsUntilGracefulShutdownIsStarted int

	// http.Server.Shutdown メソッドが呼ばれた後、全てのTCPコネクションをアイドル状態にする処理のタイムアウト時間(秒) <- 待ち時間b
	GracefulShutdownTimeoutSeconds int

	// HTTPリクエストのキャンセル発動後、サーバーが強制終了するまでスリープする時間(秒)
	// http.Server.Shutdown メソッドが呼ばれた後、GracefulShutdownTimeoutSeconds 秒だけ待ったにも関わらず
	// 全てのTCPコネクションをアイドル状態にできなかった場合、HTTPリクエストコンテキストのキャンセルが発動される。
	// HTTPハンドラーはHTTPリクエストコンテキストのキャンセルを受信した場合、
	// ForcefullyRequestCancellationTimeoutSeconds 秒以内に
	// リソースを解放し、処理を終了させ、レスポンスを返し、コネクションをアイドル状態にしなければならない。
	ForcefullyRequestCancellationTimeoutSeconds int

	IsSignalCatched *atomic.Bool

	Logger *slog.Logger
}

// グレースフルシャットダウン付HTTPサーバー
func RunHandlerWithGracefulShutdown(
	ctx context.Context,
	handler http.Handler,
	serverPort int,
	opts Options,
) exitCode {
	var logger *slog.Logger
	if opts.Logger != nil {
		logger = opts.Logger
	} else {
		logger = slog.Default()
	}

	var isSignalCatched *atomic.Bool
	if opts.IsSignalCatched != nil {
		isSignalCatched = opts.IsSignalCatched
	} else {
		isSignalCatched = &atomic.Bool{}
	}

	// HTTPリクエストコンテキストのキャンセルを発動させるために
	// BaseContextを設定したサーバーを作成する
	// 本サンプルでは、GracefulShutdownTimeoutSeconds 秒待ってもTCPコネクションがアイドル状態にならない場合
	// cancelCtxBaseRequest() を実行し、HTTPリクエストコンテキストのキャンセルを発動させる
	ctxBaseRequest, cancelCtxBaseRequest := context.WithCancel(context.Background())
	defer cancelCtxBaseRequest()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: handler,
		BaseContext: func(l net.Listener) context.Context {
			return ctxBaseRequest
		},
	}

	// シグナルハンドラーの登録
	// ctxSignal は、シグナルをキャッチしたらctxSignal.Done()チャンネルがクローズされる
	ctxSignal, stop := signal.NotifyContext(
		ctx,

		// キャッチするシグナルの種類を指定する

		// SIGINT はUnix互換OSだけなので、OSの違いが吸収できるos.Interruptを使った方が良い
		// os.Interruptはプログラムの中断シグナル。プログラム実行中に Ctrl+C を叩くと、SIGINT シグナルがプログラムへ送られる。
		os.Interrupt,

		// SIGTERM はUnix互換OSにおけるプログラムの強制終了シグナル
		// Cloud 上のプロセス終了時、このシグナルを送るケースがある
		// 例えば Cloud Run
		// https://docs.cloud.google.com/run/docs/container-contract#instance-shutdown
		syscall.SIGTERM,
	)
	defer stop()

	// サーバーの起動
	chServeIsDone := make(chan error)
	go func() {
		logger.Info("server started", slog.Int("port", serverPort))

		// リスン状態を開始する
		// 意図的なリスン状態の終了(http.Server.Shutdown または http.Server.Close が実行されたことによる終了)においては
		// ListenAndServeメソッドは ErrServerClosed エラーを返す
		// そうでない場合においては、そのエラー内容を返す
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Info("server finished with error", slog.Any("error", err))
			chServeIsDone <- err
		} else {
			// ErrServerClosed だった場合は異常終了ではない
			// なぜならグレースフルシャットダウンによる終了(http.Server.Shutdownメソッドが実行された)なため
			logger.Info("server finished")
		}

		close(chServeIsDone)
	}()

	// サーバーの終了、または、シグナルの受信、を待つ
	select {
	case err := <-chServeIsDone:
		// サーバーの終了

		if err != nil {
			// シグナルを受信していないけどなんらかの理由でサーバーがエラー終了した場合、このパスが実行される
			logger.Info("server listen is finished with error", slog.Any("error", err))
			return 1
		}

		// シグナルを受信していないけどサーバーが正常終了した場合、このパスが実行される
		// シグナルを受信するまでサーバーが正常終了することはありえないため
		// 理論上、このパスを通ることは考えられないが
		// もしこのパスを通るとしたら、意味としては
		// エラーなくリスン状態を終了したことを意味する
		logger.Info("server listen is finished")
		return 0
	case <-ctxSignal.Done():
		// シグナルを受信
	}
	// シグナル受信後の処理をここから下に書く

	// シグナルハンドラーを解除するために stop 関数を実行する
	stop()

	logger.Info("catch signal", slog.Any("error", context.Cause(ctxSignal)))
	if isSignalCatched != nil {
		isSignalCatched.Store(true)
	}

	logger.Info("sleeping before calling http.Server.Shutdown", slog.Int("seconds", opts.WaitSecondsUntilGracefulShutdownIsStarted))
	time.Sleep(time.Duration(opts.WaitSecondsUntilGracefulShutdownIsStarted) * time.Second)

	// グレースフルシャットダウンの実行が開始される
	logger.Info("starting graceful shutdown (calling http.Server.Shutdown)")
	logger.Info("waiting for all TCP connections to become idle", slog.Int("seconds", opts.GracefulShutdownTimeoutSeconds))
	ctxTimeout, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(opts.GracefulShutdownTimeoutSeconds)*time.Second,
	)
	defer cancel()
	err := server.Shutdown(ctxTimeout)
	if err != nil {
		// Shutdownメソッドがエラーを返した場合
		// GracefulShutdownTimeoutSeconds (秒)時間以内に全てのTCPコネクションをアイドル状態にできなかったことを意味する
		// アイドル状態に戻せなかったコネクションが残っているがもうこれ以上は待てないので、
		// cancelCtxBaseRequest を呼び、コンテキストを介してハンドラー側へキャンセル信号を伝播する
		// ハンドラー側はキャンセル信号を受信したら即座にリソースを解放し、処理を終了させ、レスポンスを返し、コネクションをアイドル状態にする
		cancelCtxBaseRequest()
		logger.Info("failed to idle all TCP connections", slog.Any("error", err))
		logger.Info("sending cancellation signal to handlers")
		logger.Info("server will be forcefully shut down, handlers must finish within timeout", slog.Int("seconds", opts.ForcefullyRequestCancellationTimeoutSeconds))
		time.Sleep(time.Duration(opts.ForcefullyRequestCancellationTimeoutSeconds) * time.Second)
		logger.Info("exit server forcefully")
		// 実はこれを呼ぶ必要があるらしい...
		// server.Close を呼ばない場合、TCPコネクションは生き続ける（サーバーは終了していない）
		err := server.Close()
		logger.Info("server is closed forcefully", slog.Any("error", err))
		return shutodownForcefully
	}

	logger.Info("exit server gracefully")
	return shutodownGracefully
}
