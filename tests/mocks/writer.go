package mocks

type Writer struct {
	writeFn func(p []byte) (n int, err error)
}

func NewWriterMock(writeFn func(p []byte) (n int, err error)) *Writer {
	return &Writer{writeFn: writeFn}
}

func (m *Writer) Write(p []byte) (n int, err error) {
	return m.writeFn(p)
}
