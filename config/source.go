package config

import (
	"encoding/json"
	"goflow/exception"
	"reflect"
)

type Provider interface {
	Load(v any, opts ...LoadOption) error
}

type Source struct {
	providers []Provider
	defaults  []byte
	options   []LoadOption
}

type opt struct {
	key   string
	value any
}

func Opt(key string, value any) *opt {
	return &opt{
		key:   key,
		value: value,
	}
}

// Provide returns a pointer to new instance of Source.
// Field value will be overriden by each provider in the same order as they're passed in.
func Provide(providers ...Provider) *Source {
	return &Source{providers: providers}
}

func (s *Source) ShareOptions(options ...LoadOption) {
	s.options = options
}

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

func (s *Source) Default(v any) error {
	if len(s.defaults) > 0 {
		return json.Unmarshal(s.defaults, v)
	}
	return nil
}

// Load calls LoadWithOptions with shared LoadOptions
func (s *Source) Load(v any) error {
	return s.LoadWithOptions(v, s.options...)
}

// Load run loading on each provider in order as it was initialized and bind found properties
// to corresponding 'v' field by it's name. 'v' must be a Pointer.
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
func Bind(from any, to any) (err error) {
	defer exception.HandlePanicError(func(er error) {
		err = er
	})

	toVal := reflect.ValueOf(to)
	fromVal := reflect.ValueOf(from)
	for i := 0; i < toVal.NumField(); i++ {
		sf := toVal.Type().Field(i)
		fv := fromVal.FieldByName(sf.Name)
		if fv.IsZero() {
			continue
		}

		// TODO: copy field value if possible
	}

	return nil
}
