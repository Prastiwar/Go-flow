package exception

import (
	"errors"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

type stringerMock struct {
	name string
}

func (s stringerMock) String() string {
	return s.name
}

var (
	errShared = errors.New("invalid")
)

func TestConvertToError(t *testing.T) {
	tests := []struct {
		name      string
		e         any
		assertErr assert.ErrorFunc
	}{
		{
			name: "success-from-error",
			e:    errShared,
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, errShared, err)
			},
		},
		{
			name: "success-from-string",
			e:    errShared.Error(),
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, errShared, err)
			},
		},
		{
			name: "success-from-custom-struct",
			e:    struct{ name string }{name: errShared.Error()},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, "{invalid}", err.Error())
			},
		},
		{
			name: "success-from-stringer",
			e:    stringerMock{name: "smth"},
			assertErr: func(t *testing.T, err error) {
				assert.Equal(t, "smth", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConvertToError(tt.e)
			tt.assertErr(t, err)
		})
	}
}

func TestHandlePanicError(t *testing.T) {
	tests := []struct {
		name     string
		panicArg any
		count    int
		onPanic  func(t *testing.T, counter *assert.Counter) func(err error)
	}{
		{
			name:     "success-error",
			panicArg: errShared,
			count:    1,
			onPanic: func(t *testing.T, counter *assert.Counter) func(err error) {
				return func(err error) {
					counter.Inc()
					assert.Equal(t, errShared, err)
				}
			},
		},
		{
			name:  "success-no-error",
			count: 0,
			onPanic: func(t *testing.T, counter *assert.Counter) func(err error) {
				return func(err error) {
					counter.Inc()
					t.Fail()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := assert.Count(tt.count)

			defer HandlePanicError(func(err error) {
				tt.onPanic(t, counter)(err)
				counter.Assert(t)
			})

			if tt.panicArg != nil {
				panic(tt.panicArg)
			}
		})
	}
}
