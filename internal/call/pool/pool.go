package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"test_trigger/internal/call/worker"
	"test_trigger/internal/logger"
)

//go:generate go run github.com/golang/mock/mockgen --source=pool.go --destination=pool_mock.go --package=pool

var ErrWorkersCount = errors.New("workers count should be greater than 0")

type WorkerCreator interface {
	NewWorker() worker.Worker
}

type QueueLengthGetter interface {
	QueueLength(_ context.Context) (int, error)
}

// Pool manages workers.
type Pool struct {
	wg                *sync.WaitGroup
	WorkerCreator     WorkerCreator
	QueueLengthGetter QueueLengthGetter
	Logger            logger.Logger
}

func NewPool(workerCreator WorkerCreator, queueLengthGetter QueueLengthGetter, logger logger.Logger) *Pool {
	return &Pool{WorkerCreator: workerCreator, QueueLengthGetter: queueLengthGetter, Logger: logger, wg: &sync.WaitGroup{}}
}

// Start runs workers.
func (p *Pool) Start(ctx context.Context, maxWorkers int) error {
	if maxWorkers <= 0 {
		return ErrWorkersCount
	}
	for i := 1; i <= maxWorkers; i++ {
		w := p.WorkerCreator.NewWorker()
		p.wg.Add(1)
		go w.ProcessCalls(ctx, p.wg)
	}

	return nil
}

// Close stops workers, gives them time to finish all calls in the queue.
func (p *Pool) Close(ctx context.Context, cancelFunc context.CancelFunc, recheckTime, closeTimeout time.Duration) {
	ticker := time.NewTicker(recheckTime)
	defer ticker.Stop()
	timeoutTicker := time.NewTicker(closeTimeout)
	defer timeoutTicker.Stop()
	p.Logger.Info("pool closure started")
mainLoop:
	for {
		queueLength, err := p.QueueLengthGetter.QueueLength(ctx)
		if err != nil {
			p.Logger.Error(fmt.Errorf("pool close: %v", err))
			continue
		}
		if queueLength == 0 {
			break
		}
		select {
		case <-ticker.C:
			p.Logger.Info(fmt.Sprintf("%v calls should be processed", queueLength))
		case <-timeoutTicker.C:
			p.Logger.Info(fmt.Sprintf("%v calls should have been processed, but they weren't", queueLength))
			break mainLoop
		}
	}

	cancelFunc()
	p.wg.Wait()
	p.Logger.Info("pool closure finished")
}
