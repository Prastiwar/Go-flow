// http package implements rest abstraction to provide standard net/http implementation over it.
package http

import (
	"context"
	"net/http"

	"github.com/Prastiwar/Go-flow/rest"
)

type client struct {
	client *http.Client
}

func (c *client) Send(ctx context.Context, req rest.HttpRequest) (rest.HttpResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, req.Method(), req.Url(), req.Body())
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return rest.NewResponse(resp.StatusCode, rest.WithBody(resp.Body), rest.WithHeaders(resp.Header)), nil
}

// NewHttpClient returns rest.HttpClient adapter for http.Client.
func NewHttpClient(c *http.Client) rest.HttpClient {
	return &client{
		client: c,
	}
}
