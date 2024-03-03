package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_trigger/internal/call"
	"test_trigger/internal/logger"
)

func BuildTestReq(method, path string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var ioBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		ioBody = bytes.NewReader(jsonBody)
	}
	request := httptest.NewRequest(method, path, ioBody)
	response := httptest.NewRecorder()
	return request, response
}

func TestServer_Trigger(t *testing.T) {
	type fields struct {
		getUUID func() string
	}
	type args struct {
		method, path string
		body         interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedFunc   func(saver *MockCallSaver, l *logger.MockLogger)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "method not allowed",
			fields: fields{},
			args: args{
				method: http.MethodDelete,
				path:   "/",
			},
			expectedFunc:   nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
		{
			name:   "failed body",
			fields: fields{},
			args: args{
				method: http.MethodPost,
				path:   "/",
				body:   "123",
			},
			expectedFunc:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "json: cannot unmarshal string into Go value of type call.Body\n",
		},
		{
			name:   "failed, phone_number is empty",
			fields: fields{},
			args: args{
				method: http.MethodPost,
				path:   "/",
				body: call.Body{
					PhoneNumber:    "",
					VirtualAgentID: "aaa",
				},
			},
			expectedFunc:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "virtual_agent_id or phone_number can't be empty\n",
		},
		{
			name: "failed, save to storage",
			fields: fields{
				getUUID: func() string {
					return "1"
				},
			},
			args: args{
				method: http.MethodPost,
				path:   "/",
				body: call.Body{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
				},
			},
			expectedFunc: func(saver *MockCallSaver, l *logger.MockLogger) {
				saver.EXPECT().AddToQueueBack(gomock.Any(), call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(errors.New("some err"))
				l.EXPECT().Error(fmt.Errorf("trigger: AddToQueueBack: %v", errors.New("some err")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
		{
			name: "success",
			fields: fields{
				getUUID: func() string {
					return "1"
				},
			},
			args: args{
				method: http.MethodPost,
				path:   "/",
				body: call.Body{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
				},
			},
			expectedFunc: func(saver *MockCallSaver, l *logger.MockLogger) {
				saver.EXPECT().AddToQueueBack(gomock.Any(), call.Meta{
					PhoneNumber:    "777",
					VirtualAgentID: "aaa",
					ID:             "1",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"call_id":"1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			callSaver := NewMockCallSaver(ctrl)
			l := logger.NewMockLogger(ctrl)
			s := &Server{
				callSaver: callSaver,
				getUUID:   tt.fields.getUUID,
				logger:    l,
			}
			if tt.expectedFunc != nil {
				tt.expectedFunc(callSaver, l)
			}
			ao := assert.New(t)
			testReq, response := BuildTestReq(tt.args.method, tt.args.path, tt.args.body)
			s.Trigger(response, testReq)
			ao.Equal(tt.expectedStatus, response.Code)
			ao.Equal(tt.expectedBody, response.Body.String())
		})
	}
}
