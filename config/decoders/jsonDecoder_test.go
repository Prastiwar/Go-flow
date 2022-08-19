package decoders

import (
	"goflow/tests/assert"
	"io"
	"strings"
	"testing"
)

func TestJsonDecoderDecode(t *testing.T) {
	type TitleObj struct {
		Title string
	}

	tests := []struct {
		name      string
		r         io.Reader
		v         any
		assertion assert.ResultErrorFunc[any]
	}{
		{
			name: "success-struct-pointer",
			r:    strings.NewReader("{ \"title\": \"header\"}"),
			v:    &TitleObj{},
			assertion: func(t *testing.T, result any, err error) {
				assert.NilError(t, err)
				v := result.(*TitleObj)
				assert.Equal(t, "header", v.Title)
			},
		},
		{
			name: "invalid-struct-non-pointer",
			r:    strings.NewReader("{ \"title\": \"header\"}"),
			v:    TitleObj{},
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorWith(t, err, "non-pointer")
				v := result.(TitleObj)
				assert.Equal(t, "", v.Title)
			},
		},
		{
			name: "invalid-bool",
			r:    strings.NewReader("{ \"title\": \"header\"}"),
			v:    new(bool),
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorWith(t, err, "cannot unmarshal")
			},
		},
		{
			name: "invalid-nil",
			r:    strings.NewReader("{ \"title\": \"header\"}"),
			v:    nil,
			assertion: func(t *testing.T, result any, err error) {
				assert.ErrorWith(t, err, "nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewJson()
			err := d.Decode(tt.r, tt.v)
			tt.assertion(t, tt.v, err)
		})
	}
}
