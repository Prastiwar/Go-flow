package datas

import (
	"encoding/json"
	"io"
)

var (
	_ ByteIOFormatter = &jsonData{}
)

type jsonData struct{}

// Json returns a ByteIOFormatter for encoding and decoding data in JSON format.
// The returned ByteIOFormatter is implemented using the encoding/json package from the Go standard library.
func Json() ByteIOFormatter {
	return &jsonData{}
}

func (d *jsonData) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (*jsonData) MarshalTo(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (d *jsonData) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (*jsonData) UnmarshalFrom(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
