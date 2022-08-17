package config

import (
	"goflow/tests/assert"
	"os"
	"testing"
)

func TestEnvProviderLoad(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		init    func(t *testing.T) (any, func())
		wantErr bool
	}{
		{
			name:   "succes-with-prefix",
			prefix: "DEV_",
			init: func(t *testing.T) (any, func()) {
				os.Setenv("DEV_CI", "true")
				os.Setenv("DEV_PATH", "./tests")
				v := struct {
					CI      bool
					Path    string
					NotUsed string
				}{
					CI: false,
				}

				return &v, func() {
					assert.Equal(t, v.CI, true)
					assert.Equal(t, v.Path, "./tests")
					os.Unsetenv("DEV_CI")
					os.Unsetenv("DEV_PATH")
				}
			},
			wantErr: false,
		},
		{
			name: "success-unexported-field",
			init: func(t *testing.T) (any, func()) {
				os.Setenv("Obj", "true")
				v := struct {
					obj struct{}
				}{
					obj: struct{}{},
				}

				return &v, func() {
					os.Unsetenv("Obj")
				}
			},
			wantErr: false,
		},
		{
			name: "invalid-non-pointer",
			init: func(t *testing.T) (any, func()) {
				os.Setenv("CI", "true")
				v := struct {
					CI bool
				}{
					CI: false,
				}

				return v, func() {
					os.Unsetenv("CI")
				}
			},
			wantErr: true,
		},
		{
			name: "invalid-not-parsable",
			init: func(t *testing.T) (any, func()) {
				os.Setenv("Obj", "true")
				v := struct {
					Obj struct{}
				}{
					Obj: struct{}{},
				}

				return &v, func() {
					os.Unsetenv("CI")
				}
			},
			wantErr: true,
		},
		{
			name: "invalid-not-struct",
			init: func(t *testing.T) (any, func()) {
				os.Setenv("CI", "true")
				var v *bool

				return &v, func() {
					os.Unsetenv("CI")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewEnvProviderWith(tt.prefix)
			v, asserts := tt.init(t)
			err := p.Load(v)
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
