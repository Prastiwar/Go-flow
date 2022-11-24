package decoders

import (
	"encoding/json"
	"io"
)

type jsonDecoder struct{}

// NewJson returns json decoder which is implementation from encoding/json.
func NewJson() *jsonDecoder {
	return &jsonDecoder{}
}

func (d *jsonDecoder) Decode(r io.Reader, v any) error {
	parser := json.NewDecoder(r)
	return parser.Decode(v)
}
