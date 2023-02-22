package datas_test

import (
	"bytes"
	"testing"

	"github.com/Prastiwar/Go-flow/datas"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

type xmlStruct struct {
	Foo string `json:"foo"`
}

func TestXmlBinary(t *testing.T) {
	xml := datas.Xml()
	data := xmlStruct{Foo: "success"}

	b, err := xml.Marshal(data)

	assert.NilError(t, err, "xml.Marshal(..)")
	assert.Equal(t, `<xmlStruct><Foo>success</Foo></xmlStruct>`, string(b), "xml.Marshal(..)")

	data.Foo = "failure"
	err = xml.Unmarshal(b, &data)

	assert.NilError(t, err, "xml.Unmarshal(..)")
	assert.Equal(t, "success", data.Foo, "xml.Unmarshal(..)")
}

func TestXmlIO(t *testing.T) {
	xml := datas.Xml()
	data := xmlStruct{Foo: "success"}
	b := bytes.NewReader([]byte(`<xmlStruct><Foo>success</Foo></xmlStruct>`))

	err := xml.UnmarshalFrom(b, &data)

	assert.NilError(t, err, "xml.UnmarshalFrom(..)")
	assert.Equal(t, "success", data.Foo, "xml.UnmarshalFrom(..)")

	writerCallCounter := assert.Count(t, 1)
	w := &mocks.Writer{
		OnWrite: func(p []byte) (n int, err error) {
			writerCallCounter.Inc()
			assert.Equal(t, `<xmlStruct><Foo>success</Foo></xmlStruct>`, string(p))
			return len(p), nil
		},
	}

	err = xml.MarshalTo(w, data)

	assert.NilError(t, err, "xml.MarshalTo(..)")
	writerCallCounter.Assert(t, "xml.MarshalTo(..)")
}
