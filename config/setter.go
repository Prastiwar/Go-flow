package config

import (
	"goflow/reflection"
	"reflect"
)

type FieldValueFinder func(key string) (any, error)

type FieldSetter interface {
	SetFields(v any, findFn FieldValueFinder) error
}

type fieldSetter struct {
	name    string
	options LoadOptions
}

func NewFieldSetter(name string, o LoadOptions) *fieldSetter {
	return &fieldSetter{
		name:    name,
		options: o,
	}
}

func (s *fieldSetter) SetFields(v any, findFn FieldValueFinder) error {
	toVal, err := valueLoadOf(v)
	if err != nil {
		return err
	}

	for i := 0; i < toVal.NumField(); i++ {
		field := toVal.Field(i)

		if !field.CanSet() {
			continue
		}

		sf := toVal.Type().Field(i)
		key := s.options.Intercept(s.name, sf)
		rawValue, err := findFn(key)
		if err != nil {
			return err
		}

		// nil value is treated as not existing (so skip)
		// return reflect.ValueOf(nil) to treat is as acual nil value
		if rawValue == nil {
			continue
		}

		val, ok := rawValue.(reflect.Value)
		if !ok {
			val = reflect.ValueOf(rawValue)
		}

		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		fieldType := field.Type()

		if field.Kind() == reflect.Pointer {
			fieldNonPointer := fieldType.Elem()
			if val.Type().ConvertibleTo(fieldNonPointer) {
				p := reflect.New(fieldNonPointer)
				val = val.Convert(fieldNonPointer)
				p.Elem().Set(val)
				field.Set(p)
				continue
			}
		}

		if val.Type().ConvertibleTo(fieldType) {
			val = val.Convert(fieldType)
			field.Set(val)
			continue
		}

		vv, err := reflection.Parse(val.String(), field.Interface())
		if err != nil {
			return err
		}

		field.Set(reflect.ValueOf(vv))
	}

	return nil
}

func valueLoadOf(v any) (reflect.Value, error) {
	if reflect.ValueOf(v).Kind() != reflect.Pointer {
		return reflect.Value{}, ErrNonPointer
	}

	toVal := reflect.ValueOf(v).Elem()
	if toVal.Kind() != reflect.Struct {
		return reflect.Value{}, ErrNonStruct
	}

	return toVal, nil
}
