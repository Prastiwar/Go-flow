package reflection

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrNotSupportedType = errors.New("target type is not supported by parser")
)

func Parse(s string, target interface{}) (interface{}, error) {
	if val, ok := target.(reflect.Value); ok {
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		target = val.Interface()
	}

	if val, ok := target.(reflect.Type); ok {
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}

		valPtr := reflect.New(val)
		target = valPtr.Elem().Interface()
	}

	switch target.(type) {
	case string:
		return s, nil

	case bool:
		return strconv.ParseBool(s)

	case int:
		v, err := strconv.ParseInt(s, 0, 0)
		return int(v), err
	case int8:
		v, err := strconv.ParseInt(s, 0, 8)
		return int8(v), err
	case int16:
		v, err := strconv.ParseInt(s, 0, 16)
		return int16(v), err
	case int32:
		v, err := strconv.ParseInt(s, 0, 32)
		return int32(v), err
	case int64:
		v, err := strconv.ParseInt(s, 0, 64)
		return int64(v), err

	case uint:
		v, err := strconv.ParseUint(s, 0, 0)
		return uint(v), err
	case uint8:
		v, err := strconv.ParseUint(s, 0, 8)
		return uint8(v), err
	case uint16:
		v, err := strconv.ParseUint(s, 0, 16)
		return uint16(v), err
	case uint32:
		v, err := strconv.ParseUint(s, 0, 32)
		return uint32(v), err
	case uint64:
		v, err := strconv.ParseUint(s, 0, 64)
		return uint64(v), err

	case float32:
		v, err := strconv.ParseFloat(s, 32)
		return float32(v), err
	case float64:
		v, err := strconv.ParseFloat(s, 64)
		return float64(v), err

	case complex64:
		v, err := strconv.ParseComplex(s, 64)
		return complex64(v), err
	case complex128:
		v, err := strconv.ParseComplex(s, 128)
		return complex128(v), err

	case time.Duration:
		return time.ParseDuration(s)
	case time.Time:
		return time.Parse(time.RFC3339, s)

	case error:
		return errors.New(s), nil
	}

	return nil, ErrNotSupportedType
}
