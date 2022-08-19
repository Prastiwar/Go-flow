package config

import (
	"errors"
	"flag"
	"fmt"
)

var (
	ErrNonPointer          = errors.New("cannot pass non pointer values")
	ErrNonStruct           = errors.New("cannot load config to non struct value")
	ErrMustImplementGetter = errors.New("must implement flag.Getter interface")
)

func wrapErrMustImplementGetter(f flag.Flag) error {
	return fmt.Errorf("flag with name '%v' is not valid: %w", f.Name, ErrMustImplementGetter)
}
