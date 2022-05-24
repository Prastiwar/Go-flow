package reflection

import (
	"testing"

	"goflow/tests/assert"
)

type Foo struct{}

func TestTypeOf(t *testing.T) {
	stringType := TypeOf[string]()
	assert.Equal(t, "string", stringType.Name())

	fooType := TypeOf[Foo]()
	assert.Equal(t, "Foo", fooType.Name())
}
