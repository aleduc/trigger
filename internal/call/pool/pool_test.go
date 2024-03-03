package pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_trigger/internal/call/worker"
	"test_trigger/internal/logger"
)

func TestNewPool(t *testing.T) {
	ctrl := gomock.NewController(t)
	workerCreator := NewMockWorkerCreator(ctrl)
	queueLengthGetter := NewMockQueueLengthGetter(ctrl)
	loggerMock := logger.NewMockLogger(ctrl)
	expected := &Pool{
		wg:                &sync.WaitGroup{},
		WorkerCreator:     workerCreator,
		QueueLengthGetter: queueLengthGetter,
		Logger:            loggerMock,
	}
	assert.Equal(t, expected, NewPool(workerCreator, queueLengthGetter, loggerMock))
}

func TestPool_Close(t *testing.T) {
	type fields struct {
		wg *sync.WaitGroup
	}
	type args struct {
		ctx          context.Context
		cancelFunc   context.CancelFunc
		recheckTime  time.Duration
		closeTimeout time.Duration
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedFunc func(ctx context.Context, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger)
	}{
		{
			name: "successful closure after second check with timer",
			fields: fields{
				wg: &sync.WaitGroup{},
			},
			args: args{
				ctx:          context.Background(),
				cancelFunc:   func() {},
				recheckTime:  time.Millisecond,
				closeTimeout: time.Minute,
			},
			expectedFunc: func(ctx context.Context, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info("pool closure started").Times(1)
				getter.EXPECT().QueueLength(ctx).Return(1, nil).Times(1)
				mockLogger.EXPECT().Info(fmt.Sprintf("%v calls should be processed", 1)).Times(1)
				getter.EXPECT().QueueLength(ctx).Return(0, nil).Times(1)
				mockLogger.EXPECT().Info("pool closure finished").Times(1)
			},
		},
		{
			name: "successful closure after second check, queueLength getter error",
			fields: fields{
				wg: &sync.WaitGroup{},
			},
			args: args{
				ctx:          context.Background(),
				cancelFunc:   func() {},
				recheckTime:  time.Millisecond,
				closeTimeout: time.Minute,
			},
			expectedFunc: func(ctx context.Context, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info("pool closure started").Times(1)
				getter.EXPECT().QueueLength(ctx).Return(1, errors.New("some err")).Times(1)
				mockLogger.EXPECT().Error(fmt.Errorf("pool close: %v", errors.New("some err")))
				getter.EXPECT().QueueLength(ctx).Return(0, nil).Times(1)
				mockLogger.EXPECT().Info("pool closure finished").Times(1)
			},
		},
		{
			name: "successful closure after timeout",
			fields: fields{
				wg: &sync.WaitGroup{},
			},
			args: args{
				ctx:          context.Background(),
				cancelFunc:   func() {},
				recheckTime:  time.Minute * 10,
				closeTimeout: time.Millisecond,
			},
			expectedFunc: func(ctx context.Context, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info("pool closure started").Times(1)
				getter.EXPECT().QueueLength(ctx).Return(1, nil).Times(1)
				mockLogger.EXPECT().Info("pool closure finished").Times(1)
				mockLogger.EXPECT().Info(fmt.Sprintf("%v calls should have been processed, but they weren't", 1)).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			workerCreator := NewMockWorkerCreator(ctrl)
			queueLengthGetter := NewMockQueueLengthGetter(ctrl)
			loggerMock := logger.NewMockLogger(ctrl)
			p := &Pool{
				wg:                tt.fields.wg,
				WorkerCreator:     workerCreator,
				QueueLengthGetter: queueLengthGetter,
				Logger:            loggerMock,
			}
			if tt.expectedFunc != nil {
				tt.expectedFunc(tt.args.ctx, workerCreator, queueLengthGetter, loggerMock)
			}
			p.Close(tt.args.ctx, tt.args.cancelFunc, tt.args.recheckTime, tt.args.closeTimeout)
		})
	}
}

func TestPool_Start(t *testing.T) {
	type fields struct {
		wg *sync.WaitGroup
	}
	type args struct {
		ctx        context.Context
		maxWorkers int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedErr  error
		expectedFunc func(ctx context.Context, controller *gomock.Controller, group *sync.WaitGroup, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger)
	}{
		{
			name: "success",
			fields: fields{
				wg: &sync.WaitGroup{},
			},
			args: args{
				ctx:        context.Background(),
				maxWorkers: 3,
			},
			expectedErr: nil,
			expectedFunc: func(ctx context.Context, controller *gomock.Controller, group *sync.WaitGroup, creator *MockWorkerCreator, getter *MockQueueLengthGetter, mockLogger *logger.MockLogger) {
				worker1 := worker.NewMockWorker(controller)
				worker2 := worker.NewMockWorker(controller)
				worker3 := worker.NewMockWorker(controller)
				creator.EXPECT().NewWorker().Times(1).Return(worker1)
				creator.EXPECT().NewWorker().Times(1).Return(worker2)
				creator.EXPECT().NewWorker().Times(1).Return(worker3)
				worker1.EXPECT().ProcessCalls(ctx, group).Times(1).Do(func(_ context.Context, wg *sync.WaitGroup) {
					wg.Done()
				})
				worker2.EXPECT().ProcessCalls(ctx, group).Times(1).Do(func(_ context.Context, wg *sync.WaitGroup) {
					wg.Done()
				})
				worker3.EXPECT().ProcessCalls(ctx, group).Times(1).Do(func(_ context.Context, wg *sync.WaitGroup) {
					wg.Done()
				})
			},
		},
		{
			name: "ErrWorkersCount",
			fields: fields{
				wg: &sync.WaitGroup{},
			},
			args: args{
				maxWorkers: -10,
			},
			expectedErr:  ErrWorkersCount,
			expectedFunc: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			workerCreator := NewMockWorkerCreator(ctrl)
			queueLengthGetter := NewMockQueueLengthGetter(ctrl)
			loggerMock := logger.NewMockLogger(ctrl)
			p := &Pool{
				wg:                tt.fields.wg,
				WorkerCreator:     workerCreator,
				QueueLengthGetter: queueLengthGetter,
				Logger:            loggerMock,
			}
			if tt.expectedFunc != nil {
				tt.expectedFunc(tt.args.ctx, ctrl, tt.fields.wg, workerCreator, queueLengthGetter, loggerMock)
			}
			err := p.Start(tt.args.ctx, tt.args.maxWorkers)
			p.wg.Wait()
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
