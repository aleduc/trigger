package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_trigger/internal/call"
	"test_trigger/internal/logger"
)

func TestAsync_ProcessCalls(t *testing.T) {
	type fields struct {
		StepTime time.Duration
	}
	type args struct {
		ctx context.Context
		wg  *sync.WaitGroup
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedFunc func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, logger *logger.MockLogger, caller *MockExternalCaller)
	}{
		{
			name: "exit after one successful process",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "777", "aaa").Return(200, nil).Times(1)
				statusStorage.EXPECT().SaveStatus(ctx, 200, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1).Do(func(_ context.Context, _ int, _ call.Meta) {
					cancelFunc()
				},
				)

			},
		},
		{
			name: "Next fail, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, errors.New("some err")).Times(1).Do(func(_ context.Context) {
					cancelFunc()
				},
				)
				l.EXPECT().Error(fmt.Errorf("processOneCall: %v", errors.New("some err")))
			},
		},
		{
			name: "empty queue, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, false, nil).Times(1).Do(func(_ context.Context) {
					cancelFunc()
				},
				)
			},
		},
		{
			name: "limit exceeded, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(false).Times(1)
				storage.EXPECT().AddToQueueFront(ctx, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1).Do(func(_ context.Context, _ call.Meta) {
					cancelFunc()
				})
			},
		},
		{
			name: "limit exceeded, fail save to the queue, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(false).Times(1)
				storage.EXPECT().AddToQueueFront(ctx, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(errors.New("some err")).Times(1).Do(func(_ context.Context, _ call.Meta) {
					cancelFunc()
				})
				l.EXPECT().Error(fmt.Errorf("processOneCall: %v", errors.New("some err")))
			},
		},
		{
			name: "Call fail, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "777", "aaa").Return(200, errors.New("some err")).Times(1)
				l.EXPECT().Error(errors.New("some err")).Times(1)
				storage.EXPECT().AddToQueueFront(ctx, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1).Do(func(_ context.Context, _ call.Meta) {
					cancelFunc()
				})

			},
		},
		{
			name: "SaveStatus fail, then exit",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "777", "aaa").Return(200, nil).Times(1)
				statusStorage.EXPECT().SaveStatus(ctx, 200, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(errors.New("some err")).Times(1).Do(func(_ context.Context, _ int, _ call.Meta) {
					cancelFunc()
				},
				)
				storage.EXPECT().AddToQueueFront(ctx, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1)
				l.EXPECT().Error(errors.New("some err")).Times(1)

			},
		},
		{
			name: "status != 200",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "777", "aaa").Return(429, nil).Times(1)
				l.EXPECT().Info("Status = 429 instead of 200").Times(1)
				statusStorage.EXPECT().SaveStatus(ctx, 429, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1).Do(func(_ context.Context, _ int, _ call.Meta) {
					cancelFunc()
				},
				)
				storage.EXPECT().AddToQueueFront(ctx, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1)

			},
		},
		{
			name: "exit after two successful process",
			fields: fields{
				StepTime: time.Millisecond,
			},
			args: args{
				wg: &sync.WaitGroup{},
			},
			expectedFunc: func(ctx context.Context, cancelFunc context.CancelFunc, limiter *MockLimiter, storage *MockProcessStorage, statusStorage *MockStatusStorage, l *logger.MockLogger, caller *MockExternalCaller) {
				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "777", "aaa").Return(200, nil).Times(1)
				statusStorage.EXPECT().SaveStatus(ctx, 200, call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil).Times(1)

				storage.EXPECT().Next(ctx).Return(call.Meta{
					PhoneNumber:    "888",
					VirtualAgentID: "bbb",
					ID:             "2",
				}, true, nil).Times(1)

				limiter.EXPECT().Allow().Return(true).Times(1)
				caller.EXPECT().Call(ctx, "888", "bbb").Return(200, nil).Times(1)
				statusStorage.EXPECT().SaveStatus(ctx, 200, call.Meta{
					PhoneNumber:    "888",
					VirtualAgentID: "bbb",
					ID:             "2",
				}).Return(nil).Times(1).Do(func(_ context.Context, _ int, _ call.Meta) {
					cancelFunc()
				},
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			limiter := NewMockLimiter(ctrl)
			storage := NewMockProcessStorage(ctrl)
			statusStorage := NewMockStatusStorage(ctrl)
			l := logger.NewMockLogger(ctrl)
			caller := NewMockExternalCaller(ctrl)
			a := &Async{
				Limiter:        limiter,
				Storage:        storage,
				StatusStorage:  statusStorage,
				Logger:         l,
				ExternalCaller: caller,
				StepTime:       tt.fields.StepTime,
			}
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()
			tt.args.ctx = ctx
			if tt.expectedFunc != nil {
				tt.expectedFunc(tt.args.ctx, cancelFunc, limiter, storage, statusStorage, l, caller)
			}
			tt.args.wg.Add(1)
			a.ProcessCalls(tt.args.ctx, tt.args.wg)
			tt.args.wg.Wait()
		})
	}
}

// TODO add tests.
func TestAsync_processFail(t *testing.T) {
	type fields struct {
		Limiter        Limiter
		Storage        ProcessStorage
		Logger         logger.Logger
		ExternalCaller ExternalCaller
		StepTime       time.Duration
	}
	type args struct {
		ctx context.Context
		val call.Meta
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Async{
				Limiter:        tt.fields.Limiter,
				Storage:        tt.fields.Storage,
				Logger:         tt.fields.Logger,
				ExternalCaller: tt.fields.ExternalCaller,
				StepTime:       tt.fields.StepTime,
			}
			a.processFail(tt.args.ctx, tt.args.val)
		})
	}
}

// TODO add tests or hide behind a new interface.
func TestAsync_processOneCall(t *testing.T) {
	type fields struct {
		Limiter        Limiter
		Storage        ProcessStorage
		Logger         logger.Logger
		ExternalCaller ExternalCaller
		StepTime       time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Async{
				Limiter:        tt.fields.Limiter,
				Storage:        tt.fields.Storage,
				Logger:         tt.fields.Logger,
				ExternalCaller: tt.fields.ExternalCaller,
				StepTime:       tt.fields.StepTime,
			}
			a.processOneCall(tt.args.ctx)
		})
	}
}

func TestNewWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	limiter := NewMockLimiter(ctrl)
	storage := NewMockProcessStorage(ctrl)
	statusStorage := NewMockStatusStorage(ctrl)
	l := logger.NewMockLogger(ctrl)
	caller := NewMockExternalCaller(ctrl)
	expected := &Async{
		Limiter:        limiter,
		Storage:        storage,
		StatusStorage:  statusStorage,
		Logger:         l,
		ExternalCaller: caller,
		StepTime:       time.Second,
	}

	assert.Equal(t, expected, NewWorker(limiter, storage, statusStorage, l, caller, time.Second))
}
