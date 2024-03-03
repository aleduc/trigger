// Code generated by MockGen. DO NOT EDIT.
// Source: external.go

// Package call is a generated GoMock package.
package call

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHTTPWrapper is a mock of HTTPWrapper interface.
type MockHTTPWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockHTTPWrapperMockRecorder
}

// MockHTTPWrapperMockRecorder is the mock recorder for MockHTTPWrapper.
type MockHTTPWrapperMockRecorder struct {
	mock *MockHTTPWrapper
}

// NewMockHTTPWrapper creates a new mock instance.
func NewMockHTTPWrapper(ctrl *gomock.Controller) *MockHTTPWrapper {
	mock := &MockHTTPWrapper{ctrl: ctrl}
	mock.recorder = &MockHTTPWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHTTPWrapper) EXPECT() *MockHTTPWrapperMockRecorder {
	return m.recorder
}

// MakePostRequest mocks base method.
func (m *MockHTTPWrapper) MakePostRequest(ctx context.Context, url string, body []byte) ([]byte, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakePostRequest", ctx, url, body)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// MakePostRequest indicates an expected call of MakePostRequest.
func (mr *MockHTTPWrapperMockRecorder) MakePostRequest(ctx, url, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakePostRequest", reflect.TypeOf((*MockHTTPWrapper)(nil).MakePostRequest), ctx, url, body)
}