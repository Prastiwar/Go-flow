package config

import "errors"

var (
	ErrNonPointer = errors.New("cannot pass non pointer values")
	ErrNonStruct  = errors.New("cannot load config to non struct value")
)
