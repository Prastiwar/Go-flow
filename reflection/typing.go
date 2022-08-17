package reflection

import "reflect"

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// InParamTypes returns a function input parameter types. It panics if the 'fn' Kind is not Func.
func InParamTypes(fn reflect.Type) []reflect.Type {
	if fn.Kind() != reflect.Func {
		return nil
	}

	paramLen := fn.NumIn()
	paramTypes := make([]reflect.Type, paramLen)
	for i := 0; i < paramLen; i++ {
		paramTypes[i] = fn.In(i)
	}

	return paramTypes
}

// OutParamTypesOf returns a function output parameter types. It panics if the 'fn' Kind is not Func.
func OutParamTypes(fn reflect.Type) []reflect.Type {
	if fn.Kind() != reflect.Func {
		return nil
	}

	paramLen := fn.NumOut()
	paramTypes := make([]reflect.Type, paramLen)
	for i := 0; i < paramLen; i++ {
		paramTypes[i] = fn.Out(i)
	}

	return paramTypes
}

// TogglePointer if u kind is pointer, returns its element else returns u as pointer
func TogglePointer(u reflect.Type) reflect.Type {
	if u.Kind() == reflect.Pointer {
		return u.Elem()
	}

	return reflect.PointerTo(u)
}
