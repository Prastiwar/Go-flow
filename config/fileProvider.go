package config

import (
	"errors"
	"flag"
)

type fileProvider struct {
	filename string
}

func FileProvider(filename string) *fileProvider {
	return &fileProvider{filename: filename}
}

func (p *fileProvider) Default(v any) error {
	return errors.New("not implemented")
}

func (p *fileProvider) Bind(v any) error {
	return errors.New("not implemented")
}
