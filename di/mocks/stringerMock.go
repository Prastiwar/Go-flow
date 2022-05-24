package mocks

type StringerMock struct {
	s string
}

func (s StringerMock) String() string {
	return s.s
}

func (s StringerMock) GoString() string {
	return s.s
}

func NewStringerMock(str string) *StringerMock {
	return &StringerMock{str}
}
