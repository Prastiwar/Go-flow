package httpf

import (
	"context"
	"errors"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
)

const (
	RateLimitLimitHeader     = "X-Rate-Limit-Limit"
	RateLimitRemainingHeader = "X-Rate-Limit-Remaining"
	RateLimitResetHeader     = "X-Rate-Limit-Reset"
	RetryAfterHeader         = "Retry-After"
)

var (
	ErrMissingRateStore      = errors.New("nil LimiterStore passed as parameter")
	ErrMissingRateKeyFactory = errors.New("nil RateHttpKeyFactory passed as parameter")
)

// RateHttpKeyFactory is factory func to create a string key based on http.Request.
type RateHttpKeyFactory func(r *http.Request) string

// IPRateKey returns RateHttpKeyFactory that extracts value from http.Request that's either remote IP or header value sent with
// request and are specified in 'headerNames'. If more than one header would match, only the first value will be returned.
func IPRateKey(headerNames ...string) RateHttpKeyFactory {
	headers := make([]string, len(headerNames))
	for i, v := range headerNames {
		headers[i] = textproto.CanonicalMIMEHeaderKey(v)
	}

	return func(r *http.Request) string {
		for _, v := range headers {
			val, ok := r.Header[v]
			if ok && len(val) > 0 {
				return val[0]
			}
		}
		return r.RemoteAddr
	}
}

// PathRateKey returns RateKeyFactory that extracts url path from http.Request. Note this extracts whole URL path without query and not
// the registered route pattern therefore should not be used together with endpoints which use path parameters.
// To distinguish between same pattern using different methods it appends http.Request Method as prefix with ':' separator.
func PathRateKey() RateHttpKeyFactory {
	return func(r *http.Request) string {
		return r.Method + ":" + r.URL.Path
	}
}

// ComposeRateKeyFactories aggregates many RateHttpKeyFactory into single which will invoke all factories and
// merge the keys into single string key using whitespace as separator.
func ComposeRateKeyFactories(factories ...RateHttpKeyFactory) RateHttpKeyFactory {
	if len(factories) == 0 {
		panic("factories cannot be empty")
	}

	if len(factories) == 1 {
		return factories[0]
	}

	return func(r *http.Request) string {
		str := strings.Builder{}
		lastIndex := len(factories) - 1
		for i := 0; i < len(factories); i++ {
			str.WriteString(factories[i](r))
			if i != lastIndex {
				str.WriteString(" ")
			}
		}
		return str.String()
	}
}

// RateLimitMiddleware returns httpf.Handler which used rate-limiting feature to decide if h Handler can be requested.
// If rate limit exceeds maximum value, the error is returned and should be handled by ErrorHandler to actually
// return 429 status code with appropiate body. This middleware writes all of
// "X-Rate-Limit-Limit", "X-Rate-Limit-Remaining" and "X-Rate-Limit-Reset" headers with correct values.
func RateLimitMiddleware(h Handler, store rate.LimiterStore, keyFactory RateHttpKeyFactory) Handler {
	if store == nil {
		panic(ErrMissingRateStore)
	}
	if keyFactory == nil {
		panic(ErrMissingRateKeyFactory)
	}
	return HandlerFunc(func(w ResponseWriter, r *http.Request) error {
		key := keyFactory(r)
		ctx := r.Context()

		limiter, err := store.Limit(ctx, key)
		if err != nil {
			return err
		}

		token, err := limiter.Take(ctx)
		if err != nil {
			return err
		}

		if err := token.Use(); err != nil {
			if errors.Is(err, rate.ErrRateLimitExceeded) {
				writeRateLimitHeaders(ctx, w.Header(), limiter, token, err)
			}
			return err
		}

		writeRateLimitHeaders(ctx, w.Header(), limiter, token, nil)
		return h.ServeHTTP(w, r)
	})
}

func writeRateLimitHeaders(ctx context.Context, headers http.Header, limiter rate.Limiter, token rate.Token, err error) {
	maxRate := strconv.FormatInt(int64(limiter.Limit()), 10)
	headers.Add(RateLimitLimitHeader, maxRate)

	if tokens, err := limiter.Tokens(ctx); err == nil {
		remaining := strconv.FormatInt(int64(tokens), 10)
		headers.Add(RateLimitRemainingHeader, remaining)
	}

	resetAt := token.ResetsAt()
	resetsAt := strconv.FormatInt(resetAt.Unix(), 10)
	headers.Add(RateLimitResetHeader, resetsAt)

	if err != nil {
		delta := time.Until(resetAt)
		retryAfter := strconv.FormatInt(int64(delta.Seconds()), 10)
		headers.Add(RetryAfterHeader, retryAfter)
	}
}
