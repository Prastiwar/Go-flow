package net

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Prastiwar/Go-flow/rest"
)

type client struct {
	client *http.Client
}

func (c *client) Send(ctx context.Context, method string, url string) (rest.HttpResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO: rest.HttpResponse from resp
	fmt.Print(resp)
	return nil, nil
}

// NewHttpClient returns rest.HttpClient adapter for http.Client.
func NewHttpClient(c *http.Client) rest.HttpClient {
	return &client{
		client: c,
	}
}
