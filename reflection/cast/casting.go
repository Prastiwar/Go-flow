// Package cast provides helpful functions for casting.
// It includes casting between different (castable)types for slice.
package cast

import (
	"reflect"

	"github.com/Prastiwar/Go-flow/reflection"
)

// As creates new slice instance and casts from elements to result slice.
// It reports true when from slice is empty.
func As[R any, T any](from []T) ([]R, bool) {
	count := len(from)
	result := make([]R, count)

	for i := 0; i < count; i++ {
		switch t := any(from[i]).(type) {
		case R:
			result[i] = R(t)
		default:
			var r R
			targetType := reflect.TypeOf(r)
			actualValue := reflect.ValueOf(from[i])

			castedValue, ok := reflection.CastFieldValue(targetType, actualValue)
			if !ok {
				return nil, false
			}

			resultVal, ok := castedValue.Interface().(R)
			if !ok {
				return nil, false
			}

			result[i] = resultVal
		}
	}

	return result, true
}

// Parse creates new slice instance and casts, converts or parses from elements to result slice.
// It reports true when from slice is empty or false if there was any error during parsing.
//
// NOTE: conversion from untyped int to string yields a string of one rune, not a string of digits.
func Parse[R any, T any](from []T) ([]R, bool) {
	result, ok := As[R](from)
	if ok {
		return result, true
	}

	count := len(from)
	result = make([]R, count)

	for i := 0; i < count; i++ {
		var r R
		targetType := reflect.TypeOf(r)
		actualValue := reflect.ValueOf(from[i])
		val, err := reflection.GetFieldValueFor(targetType, actualValue)
		if err != nil {
			return nil, false
		}

		resultVal, ok := val.Interface().(R)
		if !ok {
			return nil, false
		}

		result[i] = resultVal
	}
	return result, true
}
