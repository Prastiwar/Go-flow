package retry

import (
	"errors"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestPolicy_Execute(t *testing.T) {
	attempt := 0

	tests := []struct {
		name     string
		p        func(t *testing.T) *Policy
		fn       func() error
		asserter assert.ErrorFunc
	}{
		{
			name: "success-after-single-retry",
			p: func(t *testing.T) *Policy {
				c := assert.Count(t, 1, "failed retry call")
				return NewPolicy(
					WithCount(1),
					WithCancelPredicate(func(attempt int, err error) bool {
						c.Inc()
						return false
					}),
				)
			},
			fn: func() error {
				if attempt == 0 {
					attempt++
					return errors.New("invalid")
				}
				return nil
			},
			asserter: func(t *testing.T, err error) {
				assert.NilError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.p(t)
			attempt = 0
			err := p.Execute(tt.fn)
			tt.asserter(t, err)
		})
	}
}
