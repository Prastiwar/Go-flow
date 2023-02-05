package mocks

import "github.com/Prastiwar/Go-flow/policy/retry"

var _ retry.Policy = RetryPolicyMock{}

type RetryPolicyMock struct {
	OnExecute func(fn func() error) error
}

func (m RetryPolicyMock) Execute(fn func() error) error {
	return m.OnExecute(fn)
}
