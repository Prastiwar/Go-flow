package config

import (
	"goflow/exception"
	"goflow/reflection"
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
	p := NewEnvProvider()
	p.prefix = prefix
	return p
}

func (p *envProvider) Load(v any) (err error) {
	defer exception.HandlePanicError(func(er error) {
		err = er
	})

	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		return ErrNonPointer
	}

	toVal := reflect.ValueOf(v).Elem()
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

		envValue := reflect.ValueOf(s)
		if envValue.Type() == field.Type() {
			field.Set(envValue)
			continue
		}

		vv, err := reflection.Parse(envValue.String(), field.Interface())
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(vv))
	}

	return nil
}
