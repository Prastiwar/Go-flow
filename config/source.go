package config

import (
	"errors"
)

type Provider interface {
	Default(v any) error
	Bind(v any) error
}

type Source struct {
	providers []Provider
}

type opt struct {
	key string
	value string
}

func Opt(key, value string) *opt {
	return &opt{}
}

// Provide returns a pointer to new instance of Source.
// providers order defines values-overriding order
func Provide(providers ...Provider) *Source {
	return &Source{providers: providers}
}

func (s *Source) Default(v any) error {
	for _, p in providers {
		err := p.Default(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Source) Bind(v any) error {
	for _, p in providers {
		err := p.Bind(v)
		if err != nil {
			return err
		}
	}

	return nil
}
