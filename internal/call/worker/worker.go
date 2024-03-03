package worker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"test_trigger/internal/call"
	"test_trigger/internal/logger"
)

//go:generate go run github.com/golang/mock/mockgen --source=worker.go --destination=worker_mock.go --package=worker

// ProcessStorage describes methods for interaction with queue.
type ProcessStorage interface {
	Next(_ context.Context) (call.Meta, bool, error)
	AddToQueueFront(_ context.Context, meta call.Meta) error
}

// StatusStorage describes methods for status storage.
type StatusStorage interface {
	SaveStatus(_ context.Context, status int, meta call.Meta) error
}

// ExternalCaller send request to external call API.
type ExternalCaller interface {
	Call(ctx context.Context, phoneNumber, virtualAgentID string) (status int, err error)
}

// Limiter describes limiter internal implementation.
type Limiter interface {
	Allow() bool
}

type Async struct {
	Limiter        Limiter
	Storage        ProcessStorage
	StatusStorage  StatusStorage
	Logger         logger.Logger
	ExternalCaller ExternalCaller
	StepTime       time.Duration
}

func NewWorker(limiter Limiter, storage ProcessStorage, statusStorage StatusStorage, logger logger.Logger, externalCaller ExternalCaller, stepTime time.Duration) *Async {
	return &Async{Limiter: limiter, Storage: storage, StatusStorage: statusStorage, Logger: logger, ExternalCaller: externalCaller, StepTime: stepTime}
}

// ProcessCalls process any available calls from ProcessStorage.
func (a *Async) ProcessCalls(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(a.StepTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.processOneCall(ctx)
		}
	}
}

func (a *Async) processOneCall(ctx context.Context) {
	val, ok, err := a.Storage.Next(ctx)
	if err != nil {
		a.Logger.Error(fmt.Errorf("processOneCall: %v", err))
		return
	}
	if !ok {
		return
	}

	if !a.Limiter.Allow() {
		a.processFail(ctx, val)
		return
	}

	status, err := a.ExternalCaller.Call(ctx, val.PhoneNumber, val.VirtualAgentID)
	if err != nil {
		a.Logger.Error(err)
		a.processFail(ctx, val)
		return
	}

	err = a.StatusStorage.SaveStatus(ctx, status, val)
	if err != nil {
		a.Logger.Error(err)
		a.processFail(ctx, val)
		return
	}

	if status != http.StatusOK {
		a.Logger.Info(fmt.Sprintf("Status = %v instead of 200", status))
		a.processFail(ctx, val)
	}

	return
}

func (a *Async) processFail(ctx context.Context, val call.Meta) {
	err := a.Storage.AddToQueueFront(ctx, val)
	// Weak place, since it is possible to lose call there.
	// 1. Backoff approach can help.
	// 2. (preferable) Can be solved with different Storage approach, like in Kafka. Read message, commit message(mark as processed) after processing.
	// Even with commit/rollback implementation, will be possible to receive an error in some implementations and [lost data]/[process twice].
	if err != nil {
		a.Logger.Error(fmt.Errorf("processOneCall: %v", err))
	}
}
