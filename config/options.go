package config

import "reflect"

type KeyInterceptor func(reflect.StructField) string

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
func (o *LoadOptions) Intercept(f reflect.StructField) string {
	if o.Interceptor == nil {
		return f.Name
	}

	return o.Interceptor(f)
}
