package datas

import (
	"encoding/xml"
	"io"
)

var (
	_ ByteIOFormatter = &xmlData{}
)

type xmlData struct{}

// Xml returns a ByteIOFormatter for encoding and decoding data in XML format.
// The returned ByteIOFormatter is implemented using the encoding/xml package from the Go standard library.
func Xml() ByteIOFormatter {
	return &xmlData{}
}

func (d *xmlData) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

func (*xmlData) MarshalTo(w io.Writer, v any) error {
	return xml.NewEncoder(w).Encode(v)
}

func (d *xmlData) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

func (*xmlData) UnmarshalFrom(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
