package mocks

import (
	"context"

	"github.com/Prastiwar/Go-flow/policy/retry"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ retry.Policy = RetryPolicyMock{}

type RetryPolicyMock struct {
	OnExecute func(ctx context.Context, fn func() error) error
}

func (m RetryPolicyMock) Execute(ctx context.Context, fn func() error) error {
	assert.ExpectCall(m.OnExecute)
	return m.OnExecute(ctx, fn)
}
