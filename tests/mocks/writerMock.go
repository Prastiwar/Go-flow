package mocks

type WriterMock struct {
	writeFn func(p []byte) (n int, err error)
}

func NewWriterMock(writeFn func(p []byte) (n int, err error)) *WriterMock {
	return &WriterMock{writeFn: writeFn}
}

func (m *WriterMock) Write(p []byte) (n int, err error) {
	return m.writeFn(p)
}
