package config

import "reflect"

// KeyInterceptor allows to intercept field name before it's used to find it in provider.
// It's useful when you want to use different field names than they're defined in struct.
// For example you can use tag to define field name and intercept it there.
type KeyInterceptor func(providerName string, field reflect.StructField) string

type LoadOption func(*LoadOptions)

type LoadOptions struct {
	Interceptor KeyInterceptor
}

func NewLoadOptions(options ...LoadOption) *LoadOptions {
	opts := &LoadOptions{}
	for _, o := range options {
		o(opts)
	}

	return opts
}

func WithInterceptor(i KeyInterceptor) LoadOption {
	return func(s *LoadOptions) {
		s.Interceptor = i
	}
}

// Intercept provides default behaviour in case Interceptor is not set the exact field name will be used
func (o *LoadOptions) Intercept(providerName string, f reflect.StructField) string {
	if o.Interceptor == nil {
		return f.Name
	}

	return o.Interceptor(providerName, f)
}
