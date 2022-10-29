package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type responseReaderWriter interface {
	http.ResponseWriter
	io.ReadCloser

	StatusCode() int
}

type responseWriterDecorator struct {
	w          http.ResponseWriter
	data       *bytes.Buffer
	statusCode int
}

func (r *responseWriterDecorator) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriterDecorator) Write(data []byte) (int, error) {
	n, err := r.w.Write(data)
	if err != nil {
		return r.data.Write(data)
	}
	return n, err
}

func (r *responseWriterDecorator) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *responseWriterDecorator) Read(p []byte) (n int, err error) {
	if r.data == nil {
		return 0, nil
	}
	return r.data.Read(p)
}

func (r *responseWriterDecorator) Close() error {
	return nil
}

func (r *responseWriterDecorator) StatusCode() int {
	return r.statusCode
}

type responseWriter struct {
	data       *bytes.Buffer
	headers    http.Header
	statusCode int
}

func (r *responseWriter) Header() http.Header {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	return r.headers
}

func (r *responseWriter) Write(data []byte) (int, error) {
	if r.data == nil {
		r.data = bytes.NewBuffer(data)
		return len(data), nil
	}
	return r.data.Write(data)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if statusCode < 100 || statusCode > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", statusCode))
	}
	r.statusCode = statusCode
}

func (r *responseWriter) Read(p []byte) (n int, err error) {
	return r.data.Read(p)
}

func (r *responseWriter) Close() error {
	return nil
}

func (r *responseWriter) StatusCode() int {
	return r.statusCode
}
