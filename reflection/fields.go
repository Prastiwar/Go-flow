package reflection

import (
	"errors"
	"reflect"
)

var (
	ErrNotAddresable = errors.New("field is not addressable")
)

// CastFieldValue returns a reflect.Value which matches fieldType and value of rawValue.
// If rawValue type is different than fieldType then it's converted to match the type.
// If value cannot be converted to expected type, it'll report false together with last visited reflect.Value.
func CastFieldValue(fieldType reflect.Type, rawValue any) (reflect.Value, bool) {
	if rawValue == nil {
		return reflect.Zero(fieldType), false
	}

	val, ok := rawValue.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(rawValue)
	}

	if !val.IsValid() {
		return reflect.Zero(fieldType), false
	}

	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	valType := val.Type()
	if valType.AssignableTo(fieldType) {
		return val, true
	}

	isFieldTypePointer := fieldType.Kind() == reflect.Pointer
	if isFieldTypePointer {
		fieldNonPointer := fieldType.Elem()
		if valType.ConvertibleTo(fieldNonPointer) {
			p := reflect.New(fieldNonPointer)
			val = val.Convert(fieldNonPointer)
			p.Elem().Set(val)
			return p, true
		}
	}

	if valType.ConvertibleTo(fieldType) {
		val = val.Convert(fieldType)
		return val, true
	}

	return val, false
}

// GetFieldValueFor returns a reflect.Value which matches fieldType and value of rawValue.
// If rawValue type is different than fieldType then it's converted or parsed to match the type.
// If value cannot be converted or parsed to expected type, it'll return an ErrNotSupportedType error.
func GetFieldValueFor(fieldType reflect.Type, rawValue any) (reflect.Value, error) {
	val, ok := CastFieldValue(fieldType, rawValue)
	if ok || val == reflect.Zero(fieldType) {
		return val, nil
	}

	isFieldTypePointer := fieldType.Kind() == reflect.Pointer
	if isFieldTypePointer {
		defaultValue := reflect.Zero(fieldType.Elem())
		vv, err := Parse(val.String(), defaultValue)
		if err != nil {
			return reflect.Value{}, err
		}
		vvValue := reflect.ValueOf(vv)
		p := reflect.New(vvValue.Type())
		p.Elem().Set(vvValue)
		return p, nil
	}

	defaultValue := reflect.Zero(fieldType)
	vv, err := Parse(val.String(), defaultValue)
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(vv), nil
}

// SetFieldValue calls GetFieldValueFor and sets got value directly to field or error if occurred.
// If rawValue is nil - it will not set nil value to field - use reflect.ValueOf(nil) in this case.
func SetFieldValue(field reflect.Value, rawValue any) error {
	if !field.CanSet() {
		return ErrNotSupportedType
	}

	// skip nil values - reflect.ValueOf(nil) is used for setting nil values
	if rawValue == nil {
		return nil
	}

	v, err := GetFieldValueFor(field.Type(), rawValue)
	if err != nil {
		return err
	}

	field.Set(v)
	return nil
}
