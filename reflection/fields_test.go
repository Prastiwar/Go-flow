package reflection_test

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/reflection"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestSetFieldValue(t *testing.T) {
	var nilBool *bool
	tests := []struct {
		name     string
		field    func() reflect.Value
		rawValue any
		want     reflect.Value
		wantErr  bool
	}{
		{
			name: "success-bool",
			field: func() reflect.Value {
				v := struct {
					Num bool
				}{
					Num: true,
				}
				return reflect.ValueOf(&v).Elem().Field(0)
			},
			rawValue: "false",
			want:     reflect.ValueOf(false),
			wantErr:  false,
		},
		{
			name: "success-nil-pointer-skip",
			field: func() reflect.Value {
				v := struct {
					Num *bool
				}{
					Num: ptr(true),
				}
				return reflect.ValueOf(&v).Elem().Field(0)
			},
			rawValue: nil,
			want:     reflect.ValueOf(ptr(true)),
			wantErr:  false,
		},
		{
			name: "success-nil-pointer-override",
			field: func() reflect.Value {
				v := struct {
					Num *bool
				}{
					Num: ptr(true),
				}
				return reflect.ValueOf(&v).Elem().Field(0)
			},
			rawValue: reflect.ValueOf(nil),
			want:     reflect.ValueOf(nilBool),
			wantErr:  false,
		},
		{
			name: "success-not-pointer-convertible-pointer",
			field: func() reflect.Value {
				v := struct {
					Num *bool
				}{
					Num: ptr(false),
				}
				return reflect.ValueOf(&v).Elem().Field(0)
			},
			rawValue: "true",
			want:     reflect.ValueOf(ptr(true)),
			wantErr:  false,
		},
		{
			name: "invalid-not-pointer-struct-bool",
			field: func() reflect.Value {
				v := struct {
					Num bool
				}{
					Num: true,
				}
				return reflect.ValueOf(v).Field(0)
			},
			wantErr: true,
		},
		{
			name:    "invalid-unaddressable-bool",
			field:   func() reflect.Value { return reflect.ValueOf(1) },
			wantErr: true,
		},
		{
			name: "invalid-not-convertible",
			field: func() reflect.Value {
				v := struct {
					Num bool
				}{
					Num: true,
				}
				return reflect.ValueOf(&v).Elem().Field(0)
			},
			rawValue: 1,
			want:     reflect.ValueOf(true),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.field()
			err := reflection.SetFieldValue(field, tt.rawValue)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)
			if t.Failed() {
				t.FailNow()
			}

			if tt.want.Kind() == reflect.Pointer {
				if tt.want.IsValid() && tt.want.IsNil() {
					assert.Equal(t, true, field.IsNil())
					return
				}

				assert.Equal(t, tt.want.Elem().Interface(), field.Elem().Interface())
				return
			}

			if !tt.want.IsValid() {
				assert.Equal(t, false, field.IsValid(), "field should not be valid")
				return
			}

			assert.Equal(t, tt.want.Type(), field.Type())
			assert.Equal(t, tt.want.Interface(), field.Interface())
		})
	}
}

func TestGetFieldValueFor(t *testing.T) {
	tests := []struct {
		name      string
		fieldType reflect.Type
		rawValue  any
		want      reflect.Value
		wantErr   bool
	}{
		{
			name:      "success-parsed",
			fieldType: reflect.TypeOf(int64(1)),
			rawValue:  "1",
			want:      reflect.ValueOf(int64(1)),
			wantErr:   false,
		},
		{
			name:      "success-convertible",
			fieldType: reflect.TypeOf(int64(1)),
			rawValue:  int32(1),
			want:      reflect.ValueOf(int64(1)),
			wantErr:   false,
		},
		{
			name:      "success-non-pointer-to-pointer",
			fieldType: reflect.TypeOf(ptr(int64(1))),
			rawValue:  int64(1),
			want:      reflect.ValueOf(ptr(int64(1))),
			wantErr:   false,
		},
		{
			name:      "success-non-pointer-convertible-to-pointer",
			fieldType: reflect.TypeOf(ptr(int64(1))),
			rawValue:  "1",
			want:      reflect.ValueOf(ptr(int64(1))),
			wantErr:   false,
		},
		{
			name:      "success-pointer-convertible-to-pointer",
			fieldType: reflect.TypeOf(ptr(int64(1))),
			rawValue:  ptr("1"),
			want:      reflect.ValueOf(ptr(int64(1))),
			wantErr:   false,
		},
		{
			name:      "success-pointer-convertible-to-non-pointer",
			fieldType: reflect.TypeOf(int64(1)),
			rawValue:  ptr("1"),
			want:      reflect.ValueOf(int64(1)),
			wantErr:   false,
		},
		{
			name:      "success-zero-value",
			fieldType: reflect.TypeOf(1),
			rawValue:  reflect.Zero(reflect.TypeOf(2)),
			want:      reflect.ValueOf(0),
			wantErr:   false,
		},
		{
			name:      "success-nil",
			fieldType: reflect.TypeOf(1),
			rawValue:  nil,
			want:      reflect.ValueOf(0),
			wantErr:   false,
		},
		{
			name:      "invalid-pointer-not-parsable",
			fieldType: reflect.TypeOf(ptr(struct{}{})),
			rawValue:  ptr("1"),
			want:      reflect.Value{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reflection.GetFieldValueFor(tt.fieldType, tt.rawValue)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NilError(t, err)
			if t.Failed() {
				t.FailNow()
			}

			if tt.want.Kind() == reflect.Pointer {
				assert.Equal(t, tt.want.Elem().Interface(), got.Elem().Interface())
				return
			}

			assert.Equal(t, tt.want.Interface(), got.Interface())
		})
	}
}
