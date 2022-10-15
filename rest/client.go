package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	// TODO: handle headers and body io.Reader
	Send(ctx context.Context, method string, url string) (HttpResponse, error)
}

// TODO: Improve name
type Convention struct {
	c HttpClient
}

func (c *Convention) Get(ctx context.Context, url string) (HttpResponse, error) {
	return c.c.Send(ctx, http.MethodGet, url)
}

func (c *Convention) Post(ctx context.Context, url string, body io.Reader) (HttpResponse, error) {
	return c.c.Send(ctx, http.MethodPost, url)
}

func (c *Convention) PostForm(ctx context.Context, url string, form url.Values) (HttpResponse, error) {
	// body := strings.NewReader(form.Encode())
	// contentType := "application/x-www-form-urlencoded"
	return c.c.Send(ctx, http.MethodPost, url)
}

func (c *Convention) Put(ctx context.Context, url string, body io.Reader) (HttpResponse, error) {
	return c.c.Send(ctx, http.MethodPut, url)
}

func (c *Convention) Delete(ctx context.Context, url string) (HttpResponse, error) {
	return c.c.Send(ctx, http.MethodDelete, url)
}
