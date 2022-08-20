package config

import (
	"os"
)

type envProvider struct {
	prefix string
}

func NewEnvProvider() *envProvider {
	return &envProvider{}
}

func NewEnvProviderWith(prefix string) *envProvider {
	p := NewEnvProvider()
	p.prefix = prefix
	return p
}

func (p *envProvider) Load(v any, opts ...LoadOption) (err error) {
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
