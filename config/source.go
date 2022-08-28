// Package config provides single source of configuration management. The default providers are
// file, environment and command line (flag) configuration. You can set default values for specified key
// which will be loaded at first place and can be overriden by one of providers during loading process.
// The package contains also helpers to pass one struct fields to another struct to easily bind values to it.
package config

import (
	"encoding/json"
	"reflect"
)

// Provider is implemented by any value that has a Load method, which loads
// configuration and overrides if applicable matching field values for v
type Provider interface {
	Load(v any, opts ...LoadOption) error
}

// Source stores shared options, default values and configured providers to manage
// multi-source configuration loading.
type Source struct {
	providers []Provider
	defaults  []byte
	options   []LoadOption
}

type opt struct {
	key   string
	value any
}

// Opts creates an instance used for initializing default value for named key
func Opt(key string, value any) opt {
	return opt{
		key:   key,
		value: value,
	}
}

// Provide returns a Source of configuration. The order of passed providers matters in terms of
// overriding field value since each provider will be Load'ed in the same order as they were
// passed to this function
func Provide(providers ...Provider) *Source {
	return &Source{providers: providers}
}

// ShareOptions shares provided options to be used across each call to Load
func (s *Source) ShareOptions(options ...LoadOption) {
	s.options = options
}

// SetDefault sets default values in json format to be easily unmarshaled by Default method
func (s *Source) SetDefault(defaults ...opt) error {
	opts := make(map[string]interface{}, len(defaults))
	for _, opt := range defaults {
		_, ok := opts[opt.key]
		if ok {
			return wrapErrDuplicateKey(opt.key)
		}
		opts[opt.key] = opt.value
	}

	bytes, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	s.defaults = bytes
	return nil
}

// Default parses the JSON-encoded default values and stores the result
// in the value pointed to by v. Read json.Unmarshal for more information.
func (s *Source) Default(v any) error {
	if len(s.defaults) > 0 {
		return json.Unmarshal(s.defaults, v)
	}
	return nil
}

// Load calls LoadWithOptions with LoadOptions stored by ShareOptions method.
func (s *Source) Load(v any) error {
	return s.LoadWithOptions(v, s.options...)
}

// LoadWithOptions calls Load method on each provider which binds matching v fields by
// corresponding key value. LoadWithOptions can return ErrNonPointer or ErrNonStruct if v is not valid.
// If field was not found in provider - it will not override the value. But it can be overriden by
// provider which will be called as next in order if the value can be found.
func (s *Source) LoadWithOptions(v any, opts ...LoadOption) error {
	if _, err := valueLoadOf(v); err != nil {
		return err
	}

	if err := s.Default(v); err != nil {
		return err
	}

	for _, p := range s.providers {
		if err := p.Load(v, opts...); err != nil {
			return err
		}
	}

	return nil
}

// Bind sets each 'to' field value from corresponding field from 'from'.
// It will not return an error if will not find matching field.
func Bind(from any, to any) error {
	if reflect.ValueOf(from).Kind() != reflect.Pointer {
		return ErrNonPointer
	}

	fromVal := reflect.ValueOf(from).Elem()
	if fromVal.Kind() != reflect.Struct {
		return ErrNonStruct
	}

	setter := NewFieldSetter("bind", *NewLoadOptions())
	return setter.SetFields(to, func(key string) (any, error) {
		return fromVal.FieldByName(key), nil
	})
}
