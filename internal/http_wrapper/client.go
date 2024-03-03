package http_wrapper

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

// Client implements http request logic.
type Client struct {
	client *http.Client
}

func NewClient(requestTimeout time.Duration) *Client {
	client := &http.Client{
		Timeout: requestTimeout,
	}
	return &Client{client: client}
}

// MakePostRequest sends http POST request with body.
func (c *Client) MakePostRequest(ctx context.Context, url string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
