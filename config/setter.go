package config

import (
	"goflow/reflection"
	"reflect"
)

type FieldValueFinder func(key string) (any, error)

func setFields(v any, opts LoadOptions, findFn FieldValueFinder) error {
	toVal, err := valueLoadOf(v)
	if err != nil {
		return err
	}

	for i := 0; i < toVal.NumField(); i++ {
		field := toVal.Field(i)
		field.Type().Name()
		if !field.CanSet() {
			continue
		}

		sf := toVal.Type().Field(i)
		key := opts.Intercept(sf)
		rawValue, err := findFn(key)
		if err != nil {
			return err
		}

		if rawValue == nil {
			continue
		}

		val, ok := rawValue.(reflect.Value)
		if !ok {
			val = reflect.ValueOf(rawValue)
		}

		// TODO: duration case is invalid
		if val.Type() == field.Type() {
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
