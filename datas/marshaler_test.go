package datas

import (
	"io"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestMarshalerFunc_Marshal(t *testing.T) {
	callCounter := assert.Count(t, 1)

	f := MarshalerFunc(func(v any) ([]byte, error) {
		callCounter.Inc()
		return nil, nil
	})

	_, _ = f.Marshal(nil)

	callCounter.Assert(t)
}

func TestWriterMarshalerFunc_MarshalTo(t *testing.T) {
	callCounter := assert.Count(t, 1)

	f := WriterMarshalerFunc(func(w io.Writer, v any) error {
		callCounter.Inc()
		return nil
	})

	_ = f.MarshalTo(nil, nil)

	callCounter.Assert(t)
}
