package mocks

import (
	"io"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

var _ io.Writer = &Writer{}

type Writer struct {
	OnWrite func(p []byte) (n int, err error)
}

func NewWriterMock(writeFn func(p []byte) (n int, err error)) *Writer {
	return &Writer{OnWrite: writeFn}
}

func (m *Writer) Write(p []byte) (n int, err error) {
	assert.ExpectCall(m.OnWrite)
	return m.OnWrite(p)
}
