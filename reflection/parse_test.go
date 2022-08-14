package reflection

import (
	"errors"
	"goflow/tests/assert"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	b := true
	validTime, _ := time.Parse(time.RFC3339, "2030-10-12T00:00:00.00Z")

	tests := []struct {
		name    string
		str     string
		target  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "success-pointer-value-bool",
			str:     "true",
			target:  reflect.ValueOf(&b),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success-pointer-type-bool",
			str:     "true",
			target:  reflect.TypeOf(&b),
			want:    true,
			wantErr: false,
		},
		{
			name:    "success-bool",
			str:     "true",
			target:  true,
			want:    true,
			wantErr: false,
		},
		{
			name:    "success-string",
			str:     "true",
			target:  "",
			want:    "true",
			wantErr: false,
		},
		{
			name:    "success-int",
			str:     "123",
			target:  0,
			want:    123,
			wantErr: false,
		},
		{
			name:    "success-int8",
			str:     "123",
			target:  int8(0),
			want:    int8(123),
			wantErr: false,
		},
		{
			name:    "success-int16",
			str:     "123",
			target:  int16(0),
			want:    int16(123),
			wantErr: false,
		},
		{
			name:    "success-int32",
			str:     "123",
			target:  int32(0),
			want:    int32(123),
			wantErr: false,
		},
		{
			name:    "success-int64",
			str:     "123",
			target:  int64(0),
			want:    int64(123),
			wantErr: false,
		},
		{
			name:    "success-uint",
			str:     "123",
			target:  uint(0),
			want:    uint(123),
			wantErr: false,
		},
		{
			name:    "success-uint8",
			str:     "123",
			target:  uint8(0),
			want:    uint8(123),
			wantErr: false,
		},
		{
			name:    "success-uint16",
			str:     "123",
			target:  uint16(0),
			want:    uint16(123),
			wantErr: false,
		},
		{
			name:    "success-uint32",
			str:     "123",
			target:  uint32(0),
			want:    uint32(123),
			wantErr: false,
		},
		{
			name:    "success-uint64",
			str:     "123",
			target:  uint64(0),
			want:    uint64(123),
			wantErr: false,
		},
		{
			name:    "success-float32",
			str:     "123",
			target:  float32(0),
			want:    float32(123),
			wantErr: false,
		},
		{
			name:    "success-float64",
			str:     "123",
			target:  float64(0),
			want:    float64(123),
			wantErr: false,
		},
		{
			name:    "success-complex64",
			str:     "123",
			target:  complex64(0),
			want:    complex64(123),
			wantErr: false,
		},
		{
			name:    "success-complex128",
			str:     "123",
			target:  complex128(0),
			want:    complex128(123),
			wantErr: false,
		},
		{
			name:    "success-duration",
			str:     "10s",
			target:  time.Duration(0),
			want:    time.Duration(time.Second * 10),
			wantErr: false,
		},
		{
			name:    "success-time",
			str:     "2030-10-12T00:00:00.00Z",
			target:  time.Time{},
			want:    validTime,
			wantErr: false,
		},
		{
			name:    "success-error",
			str:     "text",
			target:  errors.New(""),
			want:    errors.New("text"),
			wantErr: false,
		},
		{
			name:    "invalid-type",
			str:     "{}",
			target:  struct{}{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.str, tt.target)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NilError(t, err)
			}

			if t.Failed() {
				t.FailNow()
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
