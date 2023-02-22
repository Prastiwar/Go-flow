package datas

import (
	"io"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestUnmarshalerFunc_Unmarshal(t *testing.T) {
	callCounter := assert.Count(t, 1)

	f := UnmarshalerFunc(func(data []byte, v any) error {
		callCounter.Inc()
		return nil
	})

	_ = f.Unmarshal(nil, nil)

	callCounter.Assert(t)
}

func TestReaderUnmarshalerFunc_UnmarshalFrom(t *testing.T) {
	callCounter := assert.Count(t, 1)

	f := ReaderUnmarshalerFunc(func(r io.Reader, v any) error {
		callCounter.Inc()
		return nil
	})

	_ = f.UnmarshalFrom(nil, nil)

	callCounter.Assert(t)
}
