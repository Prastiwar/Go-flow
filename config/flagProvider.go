package config

import (
	"errors"
	"flag"
)

type flagProvider struct {
	set flag.FlagSe
}

func FlagProvider(set flag.FlagSet) *flagProvider {
	return &flagProvider{set: set}
}

func (p *flagProvider) Default(v any) error {
	return nil
}

func (p *flagProvider) Bind(v any) error {
	return errors.New("not implemented")
}
