package config

import (
	"errors"
	"goflow/exception"
	"os"
	"reflect"
)

type envProvider struct {
	prefix string
}

func NewEnvProvider() *envProvider {
	return &envProvider{}
}

func NewEnvProviderWith(prefix string) *envProvider {
	return &envProvider{
		prefix: prefix,
	}
}

func (p *envProvider) Load(v any) (err error) {
	defer exception.HandlePanicError(func(er error) {
		err = er
	})

	toVal := reflect.ValueOf(v)
	for i := 0; i < toVal.NumField(); i++ {
		field := toVal.Field(i)
		if !field.CanSet() {
			continue
		}

		sf := toVal.Type().Field(i)
		envKey := p.prefix + sf.Name
		s, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		// TODO: parse common types
		envValue := reflect.ValueOf(s)
		field.Set(envValue)
	}
	return errors.New("not implemented")
}
