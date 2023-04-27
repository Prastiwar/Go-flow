package config

import (
	"context"
	"os"
)

type envProvider struct {
	prefix string
}

// NewEnvProvider returns a new environment provider used to load environments variables.
func NewEnvProvider() *envProvider {
	return &envProvider{}
}

// NewEnvProviderWith returns a new environment provider with prefix
// used to load environments variables. Prefix is used to distinguish variables in
// different environments. Mostly used ones are "DEV_", "PROD_".
func NewEnvProviderWith(prefix string) *envProvider {
	p := NewEnvProvider()
	p.prefix = prefix
	return p
}

// Load lookups os environment variables for each field name from v value and store result in matching v field.
// If there is no matching field it will be ignored and it's value will not be overridden.
func (p *envProvider) Load(ctx context.Context, v any, opts ...LoadOption) (err error) {
	options := NewLoadOptions(opts...)
	setter := NewFieldSetter(EnvProviderName, *options)

	return setter.SetFields(v, func(key string) (any, error) {
		envKey := p.prefix + key
		s, ok := os.LookupEnv(envKey)
		if !ok {
			return nil, nil
		}

		return s, nil
	})
}
