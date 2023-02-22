package datas_test

import (
	"bytes"
	"testing"

	"github.com/Prastiwar/Go-flow/datas"
	"github.com/Prastiwar/Go-flow/tests/assert"
	"github.com/Prastiwar/Go-flow/tests/mocks"
)

type jsonStruct struct {
	Foo string `json:"foo"`
}

func TestJsonBinary(t *testing.T) {
	json := datas.Json()
	data := jsonStruct{Foo: "success"}

	b, err := json.Marshal(data)

	assert.NilError(t, err, "json.Marshal(..)")
	assert.Equal(t, `{"foo":"success"}`, string(b), "json.Marshal(..)")

	data.Foo = "failure"
	err = json.Unmarshal(b, &data)

	assert.NilError(t, err, "json.Unmarshal(..)")
	assert.Equal(t, "success", data.Foo, "json.Unmarshal(..)")
}

func TestJsonIO(t *testing.T) {
	json := datas.Json()
	data := jsonStruct{Foo: "success"}
	b := bytes.NewReader([]byte(`{"foo":"success"}`))

	err := json.UnmarshalFrom(b, &data)

	assert.NilError(t, err, "json.UnmarshalFrom(..)")
	assert.Equal(t, "success", data.Foo, "json.UnmarshalFrom(..)")

	writerCallCounter := assert.Count(t, 1)
	w := &mocks.Writer{
		OnWrite: func(p []byte) (n int, err error) {
			writerCallCounter.Inc()
			assert.Equal(t, `{"foo":"success"}`+"\n", string(p))
			return len(p), nil
		},
	}

	err = json.MarshalTo(w, data)

	assert.NilError(t, err, "json.MarshalTo(..)")
	writerCallCounter.Assert(t, "json.MarshalTo(..)")
}
