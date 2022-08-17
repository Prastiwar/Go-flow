package config

import (
	"encoding/json"
	"goflow/exception"
	"reflect"
)

type Provider interface {
	Load(v any) error
}

type Source struct {
	providers []Provider
	options   []byte
}

type opt struct {
	key   string
	value string
}

func Opt(key, value string) *opt {
	return &opt{}
}

// Provide returns a pointer to new instance of Source.
// Field value will be overriden by each provider in the same order as they're passed in.
func Provide(providers ...Provider) *Source {
	return &Source{providers: providers}
}

func (s *Source) SetDefault(options ...opt) {
	opts := make(map[string]interface{}, len(options))
	for _, opt := range options {
		opts[opt.key] = opt.value
	}

	bytes, err := json.Marshal(opts)
	if err != nil {
		panic(err)
	}
	s.options = bytes
}

func (s *Source) Default(v any) error {
	return json.Unmarshal(s.options, v)
}

//
func (s *Source) Load(v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		return ErrNonPointer
	}

	err := s.Default(v)
	if err != nil {
		return err
	}

	for _, p := range s.providers {
		err := p.Load(v)
		if err != nil {
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
