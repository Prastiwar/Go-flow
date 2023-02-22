package mocks

import (
	"io"

	"github.com/Prastiwar/Go-flow/datas"
)

var (
	_ datas.ByteIOFormatter = ByteIOFormatterMock{}
)

type ByteIOFormatterMock struct {
	OnMarshal       func(v any) ([]byte, error)
	OnUnmarshal     func(data []byte, v any) error
	OnMarshalTo     func(w io.Writer, v any) error
	OnUnmarshalFrom func(r io.Reader, v any) error
}

func (m ByteIOFormatterMock) Marshal(v any) ([]byte, error) {
	return m.OnMarshal(v)
}

func (m ByteIOFormatterMock) Unmarshal(data []byte, v any) error {
	return m.OnUnmarshal(data, v)
}

func (m ByteIOFormatterMock) MarshalTo(w io.Writer, v any) error {
	return m.OnMarshalTo(w, v)
}

func (m ByteIOFormatterMock) UnmarshalFrom(r io.Reader, v any) error {
	return m.OnUnmarshalFrom(r, v)
}
