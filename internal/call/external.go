package call

import (
	"context"
	"encoding/json"
	"fmt"
)

//go:generate go run github.com/golang/mock/mockgen --source=external.go --destination=external_mock.go --package=call

type HTTPWrapper interface {
	MakePostRequest(ctx context.Context, url string, body []byte) ([]byte, int, error)
}

// Client is responsible for interaction with external API.
type Client struct {
	URL         string
	HTTPWrapper HTTPWrapper
}

func NewClient(URL string, HTTPWrapper HTTPWrapper) *Client {
	return &Client{URL: URL, HTTPWrapper: HTTPWrapper}
}

func (c *Client) Call(ctx context.Context, phoneNumber, virtualAgentID string) (status int, err error) {
	b := Body{
		PhoneNumber:    phoneNumber,
		VirtualAgentID: virtualAgentID,
	}
	// options:
	// 1. pass json.Marshal as an interface/function for unit tests.
	// 2. add body builder and return, for example, channel/func instead of struct.
	body, err := json.Marshal(b)
	if err != nil {
		return 0, fmt.Errorf("client call: %v", err)
	}

	_, status, err = c.HTTPWrapper.MakePostRequest(ctx, c.URL, body)
	if err != nil {
		return 0, fmt.Errorf("client call: make request: %v", err)
	}

	return
}
