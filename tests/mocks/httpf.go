package mocks

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/Prastiwar/Go-flow/httpf"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ httpf.Client         = ClientMock{}
	_ httpf.ResponseWriter = HttpfResponseWriterMock{}
	_ httpf.ParamsParser   = ParamsParserMock{}
	_ httpf.Router         = RouterMock{}
	_ httpf.RouteBuilder   = RouteBuilderMock{}
	_ httpf.Server         = ServerMock{}
)

type ClientMock struct {
	OnClose    func()
	OnDelete   func(ctx context.Context, url string) (*http.Response, error)
	OnGet      func(ctx context.Context, url string) (*http.Response, error)
	OnPost     func(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	OnPostForm func(ctx context.Context, url string, form url.Values) (*http.Response, error)
	OnPut      func(ctx context.Context, url string, body io.Reader) (*http.Response, error)
	OnSend     func(ctx context.Context, req *http.Request) (*http.Response, error)
}

func (m ClientMock) Close() {
	assert.ExpectCall(m.OnClose)
	m.OnClose()
}

func (m ClientMock) Delete(ctx context.Context, url string) (*http.Response, error) {
	assert.ExpectCall(m.OnDelete)
	return m.OnDelete(ctx, url)
}

func (m ClientMock) Get(ctx context.Context, url string) (*http.Response, error) {
	assert.ExpectCall(m.OnGet)
	return m.OnGet(ctx, url)
}

func (m ClientMock) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	assert.ExpectCall(m.OnPost)
	return m.OnPost(ctx, url, body)
}

func (m ClientMock) PostForm(ctx context.Context, url string, form url.Values) (*http.Response, error) {
	assert.ExpectCall(m.OnPostForm)
	return m.OnPostForm(ctx, url, form)
}

func (m ClientMock) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	assert.ExpectCall(m.OnPut)
	return m.OnPut(ctx, url, body)
}

func (m ClientMock) Send(ctx context.Context, req *http.Request) (*http.Response, error) {
	assert.ExpectCall(m.OnSend)
	return m.OnSend(ctx, req)
}

type HttpfResponseWriterMock struct {
	OnHeader      func() http.Header
	OnWrite       func(data []byte) (int, error)
	OnWriteHeader func(statusCode int)
	OnResponse    func(code int, data interface{}) error
}

func (m HttpfResponseWriterMock) Header() http.Header {
	assert.ExpectCall(m.OnHeader)
	return m.OnHeader()
}

func (m HttpfResponseWriterMock) Write(data []byte) (int, error) {
	assert.ExpectCall(m.OnWrite)
	return m.OnWrite(data)
}

func (m HttpfResponseWriterMock) WriteHeader(statusCode int) {
	assert.ExpectCall(m.OnWriteHeader)
	m.OnWriteHeader(statusCode)
}

func (m HttpfResponseWriterMock) Response(code int, data interface{}) error {
	assert.ExpectCall(m.OnResponse)
	return m.OnResponse(code, data)
}

type ParamsParserMock struct {
	OnParseParams func(r *http.Request) map[string]string
}

func (m ParamsParserMock) ParseParams(r *http.Request) map[string]string {
	assert.ExpectCall(m.OnParseParams)
	return m.OnParseParams(r)
}

type RouterMock struct {
	OnServeHTTP func(http.ResponseWriter, *http.Request)
	OnHandle    func(pattern string, handler http.Handler)
}

func (m RouterMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	assert.ExpectCall(m.OnServeHTTP)
	m.OnServeHTTP(w, r)
}

func (m RouterMock) Handle(pattern string, handler http.Handler) {
	assert.ExpectCall(m.OnHandle)
	m.OnHandle(pattern, handler)
}

type RouteBuilderMock struct {
	OnBuild               func() httpf.Router
	OnDelete              func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnGet                 func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnOptions             func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnPatch               func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnPost                func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnPut                 func(pattern string, handler httpf.Handler) httpf.RouteBuilder
	OnWithErrorHandler    func(handler httpf.ErrorHandler) httpf.RouteBuilder
	OnWithParamsParser    func(parser httpf.ParamsParser) httpf.RouteBuilder
	OnWithWriterDecorator func(decorator func(http.ResponseWriter) httpf.ResponseWriter) httpf.RouteBuilder
}

func (m RouteBuilderMock) Build() httpf.Router {
	assert.ExpectCall(m.OnBuild)
	return m.OnBuild()
}

func (m RouteBuilderMock) Delete(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnDelete)
	return m.OnDelete(pattern, handler)
}

func (m RouteBuilderMock) Get(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnGet)
	return m.OnGet(pattern, handler)
}

func (m RouteBuilderMock) Options(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnOptions)
	return m.OnOptions(pattern, handler)
}

func (m RouteBuilderMock) Patch(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnPatch)
	return m.OnPatch(pattern, handler)
}

func (m RouteBuilderMock) Post(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnPost)
	return m.OnPost(pattern, handler)
}

func (m RouteBuilderMock) Put(pattern string, handler httpf.Handler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnPut)
	return m.OnPut(pattern, handler)
}

func (m RouteBuilderMock) WithErrorHandler(handler httpf.ErrorHandler) httpf.RouteBuilder {
	assert.ExpectCall(m.OnWithErrorHandler)
	return m.OnWithErrorHandler(handler)
}

func (m RouteBuilderMock) WithParamsParser(parser httpf.ParamsParser) httpf.RouteBuilder {
	assert.ExpectCall(m.OnWithParamsParser)
	return m.OnWithParamsParser(parser)
}

func (m RouteBuilderMock) WithWriterDecorator(decorator func(http.ResponseWriter) httpf.ResponseWriter) httpf.RouteBuilder {
	assert.ExpectCall(m.OnWithWriterDecorator)
	return m.OnWithWriterDecorator(decorator)
}

type ServerMock struct {
	OnClose              func() error
	OnListenAndServe     func() error
	OnListenAndServeTLS  func(certFile string, keyFile string) error
	OnRegisterOnShutdown func(f func())
	OnServe              func(l net.Listener) error
	OnServeTLS           func(l net.Listener, certFile string, keyFile string) error
	OnShutdown           func(ctx context.Context) error
}

func (m ServerMock) Close() error {
	assert.ExpectCall(m.OnClose)
	return m.OnClose()
}

func (m ServerMock) ListenAndServe() error {
	assert.ExpectCall(m.OnListenAndServe)
	return m.OnListenAndServe()
}

func (m ServerMock) ListenAndServeTLS(certFile string, keyFile string) error {
	assert.ExpectCall(m.OnListenAndServeTLS)
	return m.OnListenAndServeTLS(certFile, keyFile)
}

func (m ServerMock) RegisterOnShutdown(f func()) {
	assert.ExpectCall(m.OnRegisterOnShutdown)
	m.OnRegisterOnShutdown(f)
}

func (m ServerMock) Serve(l net.Listener) error {
	assert.ExpectCall(m.OnServe)
	return m.OnServe(l)
}

func (m ServerMock) ServeTLS(l net.Listener, certFile string, keyFile string) error {
	assert.ExpectCall(m.OnServeTLS)
	return m.OnServeTLS(l, certFile, keyFile)
}

func (m ServerMock) Shutdown(ctx context.Context) error {
	assert.ExpectCall(m.OnShutdown)
	return m.OnShutdown(ctx)
}
