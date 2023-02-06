package config

import "reflect"

// KeyInterceptor allows to intercept field name before it's used to find it in provider.
// It's useful when you want to use different field names than they're defined in struct.
// For example you can use tag to define field name and intercept it there.
type KeyInterceptor func(providerName string, field reflect.StructField) string

type LoadOption func(*LoadOptions)

// LoadOptions stores settings to control behaviour in configuration loading.
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

// WithIgnoreGlobalOptions returns empty LoadOption to indicate no shared options should be used and
// no additional configuration is provided. This behaviour applies to Source provider.
// See Source.Load() for more information.
func WithIgnoreGlobalOptions() LoadOption {
	return func(s *LoadOptions) {}
}

// Intercept provides default behaviour in case Interceptor is not set the exact field name will be used.
func (o *LoadOptions) Intercept(providerName string, f reflect.StructField) string {
	if o.Interceptor == nil {
		return f.Name
	}

	return o.Interceptor(providerName, f)
}
