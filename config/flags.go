package config

import (
	"flag"
	"strconv"
	"time"
)

// -- bool Value
type boolValue bool

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*b = boolValue(v)
	return nil
}

func (b *boolValue) Get() any { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

// -- int32 Value
type int32Value int32

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return err
	}

	*i = int32Value(v)
	return err
}

func (i *int32Value) Get() any { return int32(*i) }

func (i *int32Value) String() string { return strconv.Itoa(int(*i)) }

// -- int64 Value
type int64Value int64

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() any { return int64(*i) }

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// -- uint Value
type uint32Value uint

func (i *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return err
	}

	*i = uint32Value(v)
	return err
}

func (i *uint32Value) Get() any { return uint32(*i) }

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- uint64 Value
type uint64Value uint64

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return err
	}

	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() any { return uint64(*i) }

func (i *uint64Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// -- string Value
type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() any { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

// -- float32 Value
type float32Value float32

func (f *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}

	*f = float32Value(v)
	return nil
}

func (f *float32Value) Get() any { return float32(*f) }

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

// -- float64 Value
type float64Value float64

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*f = float64Value(v)
	return nil
}

func (f *float64Value) Get() any { return float64(*f) }

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// -- time.Duration Value
type durationValue time.Duration

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = durationValue(v)
	return nil
}

func (d *durationValue) Get() any { return time.Duration(*d) }

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// -- time.Time Value
type timeValue time.Time

func (d *timeValue) Set(s string) error {
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*d = timeValue(v)
	return nil
}

func (d *timeValue) Get() any { return time.Time(*d) }

func (d *timeValue) String() string { return (*time.Time)(d).String() }

func BoolFlag(name string, usage string) flag.Flag {
	var value boolValue
	return CustomFlag(name, usage, &value)
}

func StringFlag(name string, usage string) flag.Flag {
	var value stringValue
	return CustomFlag(name, usage, &value)
}

func Int32Flag(name string, usage string) flag.Flag {
	var value int32Value
	return CustomFlag(name, usage, &value)
}

func Int64Flag(name string, usage string) flag.Flag {
	var value int64Value
	return CustomFlag(name, usage, &value)
}

func Uint32Flag(name string, usage string) flag.Flag {
	var value uint32Value
	return CustomFlag(name, usage, &value)
}

func Uint64Flag(name string, usage string) flag.Flag {
	var value uint64Value
	return CustomFlag(name, usage, &value)
}

func Float32Flag(name string, usage string) flag.Flag {
	var value float32Value
	return CustomFlag(name, usage, &value)
}

func Float64Flag(name string, usage string) flag.Flag {
	var value float64Value
	return CustomFlag(name, usage, &value)
}

func DurationFlag(name string, usage string) flag.Flag {
	var value durationValue
	return CustomFlag(name, usage, &value)
}

func TimeFlag(name string, usage string) flag.Flag {
	var value timeValue
	return CustomFlag(name, usage, &value)
}

func CustomFlag(name string, usage string, value flag.Value) flag.Flag {
	return flag.Flag{
		Name:     name,
		Usage:    usage,
		Value:    value,
		DefValue: value.String(),
	}
}
