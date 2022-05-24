package mocks

type ErrorMock struct {
	s string
}

func (e ErrorMock) Error() string {
	return e.s
}

func NewErrorMock(s string) *ErrorMock {
	return &ErrorMock{s}
}
