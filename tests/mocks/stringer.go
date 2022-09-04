package mocks

type Stringer struct {
	Value string
}

func (s Stringer) String() string {
	return s.Value
}

func (s Stringer) GoString() string {
	return s.Value
}
