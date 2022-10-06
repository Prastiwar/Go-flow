package config

import (
	"reflect"

	"github.com/Prastiwar/Go-flow/reflection"
)

// FieldValueFinder defines function which should return field value or error for named key
type FieldValueFinder func(key string) (any, error)

// FieldSetter is implemented by any value that has a SetFields method.
// The implementation controls how fields for v value are set.
type FieldSetter interface {
	SetFields(v any, findFn FieldValueFinder) error
}

type fieldSetter struct {
	name    string
	options LoadOptions
}

// NewFieldSetter returns a new FieldSetter implementation which can store found value for each v field.
func NewFieldSetter(name string, o LoadOptions) *fieldSetter {
	return &fieldSetter{
		name:    name,
		options: o,
	}
}

// SetFields stores found value from FieldValueFinder function. If found value is nil it will
// skip assignment. To override value and set nil value, FieldValueFinder should return reflect.ValueOf(nil).
// If field type and found value are not the same type but found value is convertible, it will try to
// convert it to matching type. It also supports converting between pointers and non-pointers.
func (s *fieldSetter) SetFields(v any, findFn FieldValueFinder) error {
	toVal, err := valueLoadOf(v)
	if err != nil {
		return err
	}

	count := toVal.NumField()
	toValType := toVal.Type()
	for i := 0; i < count; i++ {
		field := toVal.Field(i)

		if !field.CanSet() {
			continue
		}

		sf := toValType.Field(i)
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

		err = reflection.SetFieldValue(field, rawValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// valueLoadOf returns reflect.Value for struct pointer. If 'v' is not a pointer or struct it will return an error
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
