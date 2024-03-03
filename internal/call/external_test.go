package call

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestClient_Call(t *testing.T) {
	type fields struct {
		URL string
	}
	type args struct {
		ctx            context.Context
		phoneNumber    string
		virtualAgentID string
	}
	type expectedValues struct {
		err    error
		status int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedFunc func(wrapper *MockHTTPWrapper)
		expectedValues
	}{
		{
			name: "success",
			fields: fields{
				URL: "google.com",
			},
			args: args{
				ctx:            context.Background(),
				phoneNumber:    "777-77-77",
				virtualAgentID: "aaa-vvv-ddd",
			},
			expectedFunc: func(wrapper *MockHTTPWrapper) {
				val, _ := json.Marshal(Body{
					PhoneNumber:    "777-77-77",
					VirtualAgentID: "aaa-vvv-ddd",
				})
				wrapper.EXPECT().MakePostRequest(gomock.Any(), "google.com", val).Return([]byte{1}, 200, nil).Times(1)
			},
			expectedValues: expectedValues{
				status: 200,
				err:    nil,
			},
		},
		{
			name: "http err",
			fields: fields{
				URL: "google.com",
			},
			args: args{
				ctx:            context.Background(),
				phoneNumber:    "777-77-77",
				virtualAgentID: "aaa-vvv-ddd",
			},
			expectedFunc: func(wrapper *MockHTTPWrapper) {
				val, _ := json.Marshal(Body{
					PhoneNumber:    "777-77-77",
					VirtualAgentID: "aaa-vvv-ddd",
				})
				wrapper.EXPECT().MakePostRequest(gomock.Any(), "google.com", val).Return([]byte{1}, 403, errors.New("some err")).Times(1)
			},
			expectedValues: expectedValues{
				status: 0,
				err:    fmt.Errorf("client call: make request: %v", errors.New("some err")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := NewMockHTTPWrapper(ctrl)
			c := &Client{
				URL:         tt.fields.URL,
				HTTPWrapper: httpClient,
			}
			if tt.expectedFunc != nil {
				tt.expectedFunc(httpClient)
			}
			ao := assert.New(t)
			actualStatus, actualErr := c.Call(tt.args.ctx, tt.args.phoneNumber, tt.args.virtualAgentID)
			ao.Equal(tt.expectedValues.status, actualStatus)
			ao.Equal(tt.expectedValues.err, actualErr)
		})
	}
}

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := NewMockHTTPWrapper(ctrl)
	expected := &Client{
		URL:         "google.com",
		HTTPWrapper: httpClient,
	}
	assert.Equal(t, expected, NewClient("google.com", httpClient))
}
