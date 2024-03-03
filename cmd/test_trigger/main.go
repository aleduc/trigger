package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"test_trigger/internal"
	"test_trigger/internal/call"
	"test_trigger/internal/call/pool"
	"test_trigger/internal/call/worker"
	"test_trigger/internal/http_wrapper"
	"test_trigger/internal/limiter"
	"test_trigger/internal/logger"
	"test_trigger/internal/realtime"
)

const (
	defaultTimeout         = 10 * time.Minute // depends on real call duration.
	workerStepTime         = time.Millisecond * 500
	maxWorkers             = 30
	poolDefaultRecheckTime = time.Second * 3
	poolCloseTimeout       = time.Minute * 10
	defaultPort            = ":8328"
	readHeaderTimeout      = 20 * time.Second
	readTimeout            = 1 * time.Minute
	writeTimeout           = 2 * time.Minute
	shutdownTimeout        = 30 * time.Second
	limiterSecondsSize     = 10
	limiterMaxRequests     = 25
	originateTriggerURL    = "https://google.com"
)

func main() {
	mainCtx, mainCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer mainCtxCancel()

	poolCtx, poolCancel := context.WithCancel(context.Background())
	defer poolCancel()

	l := logger.NewSimple()
	storage := call.NewStorage()
	lim := limiter.NewSlidingWindow(limiterSecondsSize, limiterMaxRequests, realtime.NewRealTime(time.Now))
	httpClient := http_wrapper.NewClient(defaultTimeout)
	externalAPIClient := call.NewClient(originateTriggerURL, httpClient)

	workerCreator := worker.NewCreate(lim, storage, storage, l, externalAPIClient, workerStepTime)
	p := pool.NewPool(workerCreator, storage, l)
	err := p.Start(poolCtx, maxWorkers)
	//defer pool.Close(poolCtx, poolCancel, poolDefaultRecheckTime)
	if err != nil {
		l.Error(err)
		return
	}

	handler := internal.NewServer(storage, func() string { return uuid.New().String() }, l)
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/trigger", handler.Trigger)
	server := &http.Server{
		Addr:              defaultPort,
		Handler:           serverMux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
	}

	serverStopped := make(chan struct{}, 1)
	go func() {
		l.Info("HTTP server is started")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			l.Error(err)
		}
		serverStopped <- struct{}{}
	}()

	select {
	case <-mainCtx.Done():
		l.Info("graceful shutting down…")
		// stop http server first
		serverShutdown(l, server)
		<-serverStopped
		l.Info("http server is stopped")
		// stop workers, but process all remaining calls(with deadline). Since I have memory storage and don't want to lose calls.
		p.Close(poolCtx, poolCancel, poolDefaultRecheckTime, poolCloseTimeout)
	case <-serverStopped:
		serverShutdown(l, server)
		l.Info("http server is stopped")
		p.Close(poolCtx, poolCancel, poolDefaultRecheckTime, poolCloseTimeout)
	}

	l.Info("done")
}

func serverShutdown(l *logger.Simple, server *http.Server) {
	l.Info("stop http server")
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelTimeout()
	if err := server.Shutdown(timeoutCtx); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, http.ErrServerClosed) {
		l.Error(err)
	}
}
