// rest package abstracts REST communication.
package rest

import (
	"context"
	"io"
	"net/url"
)

type HttpClient interface {
	Send(ctx context.Context, req HttpRequest) (HttpResponse, error)
}

type FluentClient struct {
	c HttpClient
}

func NewFluentClient(c HttpClient) *FluentClient {
	return &FluentClient{c: c}
}

func (c *FluentClient) Get(ctx context.Context, url string) (HttpResponse, error) {
	return c.c.Send(ctx, NewGetRequest(url))
}

func (c *FluentClient) Post(ctx context.Context, url string, body io.Reader) (HttpResponse, error) {
	return c.c.Send(ctx, NewPostRequest(url, io.NopCloser(body)))
}

func (c *FluentClient) PostForm(ctx context.Context, url string, form url.Values) (HttpResponse, error) {
	return c.c.Send(ctx, NewPostFormRequest(url, form))
}

func (c *FluentClient) Put(ctx context.Context, url string, body io.Reader) (HttpResponse, error) {
	return c.c.Send(ctx, NewPutRequest(url, io.NopCloser(body)))
}

func (c *FluentClient) Delete(ctx context.Context, url string) (HttpResponse, error) {
	return c.c.Send(ctx, NewDeleteRequest(url))
}
