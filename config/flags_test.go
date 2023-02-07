package config_test

import (
	"flag"
	"strconv"
	"testing"
	"time"

	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

type customValue bool

func (b *customValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*b = customValue(v)
	return nil
}

func (b *customValue) Get() any { return bool(*b) }

func (b *customValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *customValue) IsBoolFlag() bool { return true }

func TestNewFlag(t *testing.T) {
	var val customValue
	got := config.CustomFlag("n", "u", &val)

	assert.Equal(t, "n", got.Name)
	assert.Equal(t, &val, got.Value)
	assert.Equal(t, val.String(), got.DefValue)
	assert.Equal(t, "u", got.Usage)
}

func TestFlags(t *testing.T) {
	const (
		name  = "test-name"
		usage = "test-usage"
	)
	validTime, _ := time.Parse(time.RFC3339, "2030-10-12T00:00:00.00Z")

	tests := []struct {
		name      string
		set       string
		flag      flag.Flag
		wantGet   any
		assertion assert.ErrorFunc //func(t *testing.T, err error, val flag.Value)
	}{
		{
			name:    "success-bool",
			set:     "true",
			flag:    config.BoolFlag(name, usage),
			wantGet: true,
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-bool",
			set:     "Maybe true",
			flag:    config.BoolFlag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-string",
			set:     "text",
			flag:    config.StringFlag(name, usage),
			wantGet: "text",
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "success-int32",
			set:     "10",
			flag:    config.Int32Flag(name, usage),
			wantGet: int32(10),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-int32",
			set:     "invalid 10",
			flag:    config.Int32Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-int64",
			set:     "10",
			flag:    config.Int64Flag(name, usage),
			wantGet: int64(10),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-int64",
			set:     "invalid 10",
			flag:    config.Int64Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-uint32",
			set:     "10",
			flag:    config.Uint32Flag(name, usage),
			wantGet: uint32(10),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-uint32",
			set:     "invalid 10",
			flag:    config.Uint32Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-uint64",
			set:     "10",
			flag:    config.Uint64Flag(name, usage),
			wantGet: uint64(10),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-uint64",
			set:     "invalid 10",
			flag:    config.Uint64Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-float32",
			set:     "10.1234",
			flag:    config.Float32Flag(name, usage),
			wantGet: float32(10.1234),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-float32",
			set:     "invalid 10",
			flag:    config.Float32Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-float64",
			set:     "10.1234",
			flag:    config.Float64Flag(name, usage),
			wantGet: float64(10.1234),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-float64",
			set:     "invalid 10",
			flag:    config.Float64Flag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-duration",
			set:     "10s",
			flag:    config.DurationFlag(name, usage),
			wantGet: time.Duration(time.Second * 10),
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-duration",
			set:     "invalid 10s",
			flag:    config.DurationFlag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:    "success-time",
			set:     "2030-10-12T00:00:00.00Z",
			flag:    config.TimeFlag(name, usage),
			wantGet: validTime,
			assertion: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
		{
			name:    "invalid-time",
			set:     "2030-10-12",
			flag:    config.TimeFlag(name, usage),
			wantGet: nil,
			assertion: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.flag.Value.Set(tt.set)
			tt.assertion(t, err)

			if tt.wantGet != nil {
				g, ok := tt.flag.Value.(flag.Getter)
				assert.Equal(t, true, ok, "flag is not getter")
				assert.Equal(t, tt.wantGet, g.Get())
			}
		})
	}
}
