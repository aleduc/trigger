// Code generated by MockGen. DO NOT EDIT.
// Source: pool.go

// Package pool is a generated GoMock package.
package pool

import (
	context "context"
	reflect "reflect"
	worker "test_trigger/internal/call/worker"

	gomock "github.com/golang/mock/gomock"
)

// MockWorkerCreator is a mock of WorkerCreator interface.
type MockWorkerCreator struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerCreatorMockRecorder
}

// MockWorkerCreatorMockRecorder is the mock recorder for MockWorkerCreator.
type MockWorkerCreatorMockRecorder struct {
	mock *MockWorkerCreator
}

// NewMockWorkerCreator creates a new mock instance.
func NewMockWorkerCreator(ctrl *gomock.Controller) *MockWorkerCreator {
	mock := &MockWorkerCreator{ctrl: ctrl}
	mock.recorder = &MockWorkerCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkerCreator) EXPECT() *MockWorkerCreatorMockRecorder {
	return m.recorder
}

// NewWorker mocks base method.
func (m *MockWorkerCreator) NewWorker() worker.Worker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewWorker")
	ret0, _ := ret[0].(worker.Worker)
	return ret0
}

// NewWorker indicates an expected call of NewWorker.
func (mr *MockWorkerCreatorMockRecorder) NewWorker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewWorker", reflect.TypeOf((*MockWorkerCreator)(nil).NewWorker))
}

// MockQueueLengthGetter is a mock of QueueLengthGetter interface.
type MockQueueLengthGetter struct {
	ctrl     *gomock.Controller
	recorder *MockQueueLengthGetterMockRecorder
}

// MockQueueLengthGetterMockRecorder is the mock recorder for MockQueueLengthGetter.
type MockQueueLengthGetterMockRecorder struct {
	mock *MockQueueLengthGetter
}

// NewMockQueueLengthGetter creates a new mock instance.
func NewMockQueueLengthGetter(ctrl *gomock.Controller) *MockQueueLengthGetter {
	mock := &MockQueueLengthGetter{ctrl: ctrl}
	mock.recorder = &MockQueueLengthGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueueLengthGetter) EXPECT() *MockQueueLengthGetterMockRecorder {
	return m.recorder
}

// QueueLength mocks base method.
func (m *MockQueueLengthGetter) QueueLength(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueLength", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueueLength indicates an expected call of QueueLength.
func (mr *MockQueueLengthGetterMockRecorder) QueueLength(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueLength", reflect.TypeOf((*MockQueueLengthGetter)(nil).QueueLength), arg0)
}
