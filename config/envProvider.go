package config

import (
	"errors"
	"flag"
)

type envProvider struct {
	prefix string
}

func EnvProvider() *envProvider {
	return &envProvider{}
}

func (p *envProvider) Default(v any) error {
	return errors.New("not implemented")
}

func (p *envProvider) Bind(v any) error {
	return errors.New("not implemented")
}
