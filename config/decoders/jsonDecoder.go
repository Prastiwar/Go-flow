package decoders

import (
	"encoding/json"
	"io"
)

type jsonDecoder struct{}

func NewJson() *jsonDecoder {
	return &jsonDecoder{}
}

func (d *jsonDecoder) Decode(r io.Reader, v any) error {
	parser := json.NewDecoder(r)
	return parser.Decode(v)
}
