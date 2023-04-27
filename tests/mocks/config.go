package mocks

import (
	"context"

	"github.com/Prastiwar/Go-flow/config"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ config.Provider = ProviderMock{}
)

type ProviderMock struct {
	OnLoad func(ctx context.Context, v any, opts ...config.LoadOption) error
}

func (m ProviderMock) Load(ctx context.Context, v any, opts ...config.LoadOption) error {
	assert.ExpectCall(m.OnLoad)
	return m.OnLoad(ctx, v, opts...)
}
