// Code generated by MockGen. DO NOT EDIT.
// Source: worker.go

// Package worker is a generated GoMock package.
package worker

import (
	context "context"
	reflect "reflect"
	call "test_trigger/internal/call"

	gomock "github.com/golang/mock/gomock"
)

// MockProcessStorage is a mock of ProcessStorage interface.
type MockProcessStorage struct {
	ctrl     *gomock.Controller
	recorder *MockProcessStorageMockRecorder
}

// MockProcessStorageMockRecorder is the mock recorder for MockProcessStorage.
type MockProcessStorageMockRecorder struct {
	mock *MockProcessStorage
}

// NewMockProcessStorage creates a new mock instance.
func NewMockProcessStorage(ctrl *gomock.Controller) *MockProcessStorage {
	mock := &MockProcessStorage{ctrl: ctrl}
	mock.recorder = &MockProcessStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessStorage) EXPECT() *MockProcessStorageMockRecorder {
	return m.recorder
}

// AddToQueueFront mocks base method.
func (m *MockProcessStorage) AddToQueueFront(arg0 context.Context, meta call.Meta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToQueueFront", arg0, meta)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToQueueFront indicates an expected call of AddToQueueFront.
func (mr *MockProcessStorageMockRecorder) AddToQueueFront(arg0, meta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToQueueFront", reflect.TypeOf((*MockProcessStorage)(nil).AddToQueueFront), arg0, meta)
}

// Next mocks base method.
func (m *MockProcessStorage) Next(arg0 context.Context) (call.Meta, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", arg0)
	ret0, _ := ret[0].(call.Meta)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Next indicates an expected call of Next.
func (mr *MockProcessStorageMockRecorder) Next(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockProcessStorage)(nil).Next), arg0)
}

// MockStatusStorage is a mock of StatusStorage interface.
type MockStatusStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStatusStorageMockRecorder
}

// MockStatusStorageMockRecorder is the mock recorder for MockStatusStorage.
type MockStatusStorageMockRecorder struct {
	mock *MockStatusStorage
}

// NewMockStatusStorage creates a new mock instance.
func NewMockStatusStorage(ctrl *gomock.Controller) *MockStatusStorage {
	mock := &MockStatusStorage{ctrl: ctrl}
	mock.recorder = &MockStatusStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatusStorage) EXPECT() *MockStatusStorageMockRecorder {
	return m.recorder
}

// SaveStatus mocks base method.
func (m *MockStatusStorage) SaveStatus(arg0 context.Context, status int, meta call.Meta) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveStatus", arg0, status, meta)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveStatus indicates an expected call of SaveStatus.
func (mr *MockStatusStorageMockRecorder) SaveStatus(arg0, status, meta interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveStatus", reflect.TypeOf((*MockStatusStorage)(nil).SaveStatus), arg0, status, meta)
}

// MockExternalCaller is a mock of ExternalCaller interface.
type MockExternalCaller struct {
	ctrl     *gomock.Controller
	recorder *MockExternalCallerMockRecorder
}

// MockExternalCallerMockRecorder is the mock recorder for MockExternalCaller.
type MockExternalCallerMockRecorder struct {
	mock *MockExternalCaller
}

// NewMockExternalCaller creates a new mock instance.
func NewMockExternalCaller(ctrl *gomock.Controller) *MockExternalCaller {
	mock := &MockExternalCaller{ctrl: ctrl}
	mock.recorder = &MockExternalCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternalCaller) EXPECT() *MockExternalCallerMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockExternalCaller) Call(ctx context.Context, phoneNumber, virtualAgentID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, phoneNumber, virtualAgentID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Call indicates an expected call of Call.
func (mr *MockExternalCallerMockRecorder) Call(ctx, phoneNumber, virtualAgentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockExternalCaller)(nil).Call), ctx, phoneNumber, virtualAgentID)
}

// MockLimiter is a mock of Limiter interface.
type MockLimiter struct {
	ctrl     *gomock.Controller
	recorder *MockLimiterMockRecorder
}

// MockLimiterMockRecorder is the mock recorder for MockLimiter.
type MockLimiterMockRecorder struct {
	mock *MockLimiter
}

// NewMockLimiter creates a new mock instance.
func NewMockLimiter(ctrl *gomock.Controller) *MockLimiter {
	mock := &MockLimiter{ctrl: ctrl}
	mock.recorder = &MockLimiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimiter) EXPECT() *MockLimiterMockRecorder {
	return m.recorder
}

// Allow mocks base method.
func (m *MockLimiter) Allow() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Allow")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Allow indicates an expected call of Allow.
func (mr *MockLimiterMockRecorder) Allow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Allow", reflect.TypeOf((*MockLimiter)(nil).Allow))
}
