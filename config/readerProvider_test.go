package config

import (
	"os"
	"strings"
	"testing"

	"github.com/Prastiwar/Go-flow/config/decoders"
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
		init    func(t *testing.T) (*readerProvider, any, func())
		wantErr bool
	}{
		{
			name: "success-json-reader",
			init: func(t *testing.T) (*readerProvider, any, func()) {
				provider := NewReaderProvider(strings.NewReader("{ \"title\": \"header\"}"), decoders.NewJson())

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
			init: func(t *testing.T) (*readerProvider, any, func()) {
				filename := createTempContentFile(t, "{ \"title\": \"header\"}")
				provider := NewFileProvider(filename, decoders.NewJson())

				v := struct {
					Title string
				}{}

				return provider, &v, func() {
					assert.Equal(t, "header", v.Title)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, v, asserts := tt.init(t)

			err := provider.Load(v)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NilError(t, err)
			}

			asserts()
		})
	}
}
