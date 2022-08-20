package config

import (
	"testing"

	"goflow/tests/assert"
)

func TestSetFields(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) (any, func())
		opts    LoadOptions
		findFn  FieldValueFinder
		wantErr bool
	}{
		{
			name: "success-raw-value-pointer",
			init: func(t *testing.T) (any, func()) {
				v := struct {
					Field string
				}{}
				return &v, func() {
					assert.Equal(t, "str", v.Field)
				}
			},
			opts: *NewLoadOptions(),
			findFn: func(key string) (any, error) {
				field := "str"
				return &field, nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, asserts := tt.init(t)
			setter := NewFieldSetter("", tt.opts)

			err := setter.SetFields(v, tt.findFn)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NilError(t, err)
			}

			if t.Failed() {
				t.FailNow()
			}

			asserts()
		})
	}
}
