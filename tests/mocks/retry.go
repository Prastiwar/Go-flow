package mocks

import (
	"github.com/Prastiwar/Go-flow/policy/retry"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ retry.Policy = RetryPolicyMock{}

type RetryPolicyMock struct {
	OnExecute func(fn func() error) error
}

func (m RetryPolicyMock) Execute(fn func() error) error {
	assert.ExpectCall(m.OnExecute)
	return m.OnExecute(fn)
}
