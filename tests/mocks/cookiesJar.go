package mocks

import (
	"net/http"
	"net/url"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

type CookiesJar struct {
	OnSetCookies func(u *url.URL, cookies []*http.Cookie)
	OnCookies    func(u *url.URL) []*http.Cookie
}

func (m *CookiesJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	assert.ExpectCall(m.OnSetCookies)
	m.OnSetCookies(u, cookies)
}

func (m *CookiesJar) Cookies(u *url.URL) []*http.Cookie {
	assert.ExpectCall(m.OnCookies)
	return m.OnCookies(u)
}
