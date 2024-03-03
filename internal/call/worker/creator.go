package worker

import (
	"context"
	"sync"
	"time"

	"test_trigger/internal/logger"
)

//go:generate go run github.com/golang/mock/mockgen --source=creator.go --destination=creator_mock.go --package=worker

// Worker processes calls.
type Worker interface {
	ProcessCalls(ctx context.Context, wg *sync.WaitGroup)
}

// Create is responsible for Worker creation.
type Create struct {
	Limiter        Limiter
	Storage        ProcessStorage
	StatusStorage  StatusStorage
	Logger         logger.Logger
	ExternalCaller ExternalCaller
	StepTime       time.Duration
}

func NewCreate(limiter Limiter, storage ProcessStorage, statusStorage StatusStorage, logger logger.Logger, externalCaller ExternalCaller, stepTime time.Duration) *Create {
	return &Create{Limiter: limiter, Storage: storage, StatusStorage: statusStorage, Logger: logger, ExternalCaller: externalCaller, StepTime: stepTime}
}

// NewWorker returns Worker interface(not structure), since it should return only specific implementation.
func (c *Create) NewWorker() Worker {
	return NewWorker(c.Limiter, c.Storage, c.StatusStorage, c.Logger, c.ExternalCaller, c.StepTime)
}
