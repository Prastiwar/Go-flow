package mocks

import (
	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ config.Provider = ProviderMock{}
)

type ProviderMock struct {
	OnLoad func(v any, opts ...config.LoadOption) error
}

func (m ProviderMock) Load(v any, opts ...config.LoadOption) error {
	assert.ExpectCall(m.OnLoad)
	return m.OnLoad(v, opts...)
}
