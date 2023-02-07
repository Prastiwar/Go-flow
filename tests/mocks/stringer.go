package mocks

type Stringer struct {
	Value string
}

func (m Stringer) String() string {
	return m.Value
}

func (m Stringer) GoString() string {
	return m.Value
}
