package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	urlpkg "net/url"
	"strconv"
	"strings"
	"time"
)

type HttpRequest interface {
	Context() context.Context
	Method() string
	Url() string
	Headers() http.Header
	Query() url.Values
	Body() io.ReadCloser
	ContentLength() int64
	Timestamp() time.Time
}

type request struct {
	ctx           context.Context
	method        string
	url           string
	contentLength int64
	body          io.ReadCloser
	headers       http.Header
	query         url.Values
	timestamp     time.Time
}

func (r request) Context() context.Context {
	return r.ctx
}

func (r request) Body() io.ReadCloser {
	return r.body
}

func (r request) ContentLength() int64 {
	return r.contentLength
}

func (r request) Headers() http.Header {
	return r.headers
}

func (r request) Method() string {
	return r.method
}

func (r request) Query() url.Values {
	return r.query
}

func (r request) Timestamp() time.Time {
	return r.timestamp
}

func (r request) Url() string {
	return r.url
}

func (r request) WithContext(ctx context.Context) request {
	r.ctx = ctx
	return r
}

func NewRequest(method string, url string, body io.ReadCloser, headers http.Header) *request {
	query := parseQuery(url)
	contentLength := calculateContentLength(body)
	headers.Set("Content-Length", strconv.FormatInt(contentLength, 10))

	return &request{
		method:        method,
		url:           url,
		body:          body,
		contentLength: contentLength,
		query:         query,
		timestamp:     time.Now(),
		headers:       headers,
	}
}

func NewGetRequest(url string) *request {
	return NewRequest(http.MethodGet, url, http.NoBody, make(http.Header))
}

func NewPostRequest(url string, body io.ReadCloser) *request {
	return NewRequest(http.MethodPost, url, body, make(http.Header))
}

func NewPostFormRequest(url string, form url.Values) *request {
	body := io.NopCloser(strings.NewReader(form.Encode()))
	headers := http.Header{
		ContentTypeHeader: {ApplicationFormEncodedType},
	}
	return NewRequest(http.MethodPost, url, body, headers)
}

func NewPutRequest(url string, body io.ReadCloser) *request {
	return NewRequest(http.MethodPut, url, body, make(http.Header))
}

func NewDeleteRequest(url string) *request {
	return NewRequest(http.MethodDelete, url, http.NoBody, make(http.Header))
}

func calculateContentLength(body io.Reader) int64 {
	if body == nil {
		return 0
	}

	switch v := body.(type) {
	case *bytes.Buffer:
		return int64(v.Len())
	case *bytes.Reader:
		return int64(v.Len())
	case *strings.Reader:
		return int64(v.Len())
	default:
		return 0
	}
}

func parseQuery(url string) url.Values {
	uri, err := urlpkg.Parse(url)
	if err != nil {
		return make(urlpkg.Values)
	}
	return uri.Query()
}
