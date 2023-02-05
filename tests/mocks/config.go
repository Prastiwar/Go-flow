package mocks

import (
	"io"

	"github.com/Prastiwar/Go-flow/config"
)

var (
	_ config.ReaderDecoder = ReaderDecoderMock{}
	_ config.Provider      = ProviderMock{}
)

type ReaderDecoderMock struct {
	OnDecode func(r io.Reader, v any) error
}

func (m ReaderDecoderMock) Decode(r io.Reader, v any) error {
	return m.OnDecode(r, v)
}

type ProviderMock struct {
	OnLoad func(v any, opts ...config.LoadOption) error
}

func (m ProviderMock) Load(v any, opts ...config.LoadOption) error {
	return m.OnLoad(v, opts...)
}
