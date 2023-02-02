package rate

var _ Limiter = LimiterMock{}

type LimiterMock struct {
}

// Limit implements Limiter
func (LimiterMock) Limit() uint64 {
	panic("unimplemented")
}

// Take implements Limiter
func (LimiterMock) Take() Token {
	panic("unimplemented")
}

// Tokens implements Limiter
func (LimiterMock) Tokens() uint64 {
	panic("unimplemented")
}
