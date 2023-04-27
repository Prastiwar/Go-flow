package config_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/datas"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

// createTempContentFile creates temporary file in temporary directory and returns its filename
func createTempContentFile(t *testing.T, contents string) string {
	f, _ := os.CreateTemp(os.TempDir(), ".json")

	_, err := f.WriteString("{ \"title\": \"header\"}")
	assert.NilError(t, err)

	err = f.Close()
	assert.NilError(t, err)

	return f.Name()
}

func TestReaderProviderLoad(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t *testing.T) (config.Provider, any, func())
		wantErr bool
	}{
		{
			name: "success-json-reader",
			init: func(t *testing.T) (config.Provider, any, func()) {
				provider := config.NewReaderProvider(strings.NewReader("{ \"title\": \"header\"}"), datas.Json())

				v := struct {
					Title string
				}{}

				return provider, &v, func() {
					assert.Equal(t, "header", v.Title)
				}
			},
			wantErr: false,
		},
		{
			name: "success-file-reader",
			init: func(t *testing.T) (config.Provider, any, func()) {
				filename := createTempContentFile(t, "{ \"title\": \"header\"}")
				provider := config.NewFileProvider(filename, datas.Json())

				v := struct {
					Title string
				}{}

				return provider, &v, func() {
					assert.Equal(t, "header", v.Title)
				}
			},
			wantErr: false,
		},
		{
			name: "inaccesible-file-reader",
			init: func(t *testing.T) (config.Provider, any, func()) {
				provider := config.NewFileProvider("unknown", datas.Json())

				v := struct {
					Title string
				}{
					Title: "unchanged",
				}

				return provider, &v, func() {
					assert.Equal(t, "unchanged", v.Title)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, v, asserts := tt.init(t)

			err := provider.Load(context.Background(), v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NilError(t, err)
			}

			asserts()
		})
	}
}
