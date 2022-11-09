package mocks

import (
	"net/http"
	"net/url"
)

type CookiesJar struct {
	OnSetCookies func(u *url.URL, cookies []*http.Cookie)
	OnCookies    func(u *url.URL) []*http.Cookie
}

func (j *CookiesJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.OnSetCookies(u, cookies)
}

func (j *CookiesJar) Cookies(u *url.URL) []*http.Cookie {
	return j.OnCookies(u)
}
