package di

import (
	"reflect"
	"testing"

	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestCache(t *testing.T) {
	tests := []struct {
		name       string
		cache      Cache
		life       LifeTime
		typ        reflect.Type
		val        interface{}
		expectedOk bool
	}{
		{
			name:       "root-success-singleton",
			cache:      NewRootCache(),
			life:       Singleton,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: true,
		},
		{
			name:       "root-invalid-transient",
			cache:      NewRootCache(),
			life:       Transient,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: false,
		},
		{
			name:       "root-invalid-scoped",
			cache:      NewRootCache(),
			life:       Scoped,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: false,
		},
		{
			name:       "scope-success-singleton",
			cache:      NewScopeCache(NewRootCache()),
			life:       Singleton,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: true,
		},
		{
			name:       "scope-invalid-transient",
			cache:      NewScopeCache(NewRootCache()),
			life:       Transient,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: false,
		},
		{
			name:       "scope-success-scoped",
			cache:      NewScopeCache(NewRootCache()),
			life:       Scoped,
			typ:        reflect.TypeOf(""),
			val:        "",
			expectedOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.cache.Put(tt.life, tt.typ, tt.val)

			assert.Equal(t, tt.expectedOk, ok, "put assertion failed")
			if t.Failed() {
				t.FailNow()
			}

			v, got := tt.cache.Get(tt.life, tt.typ)
			assert.Equal(t, ok, got)
			if ok {
				assert.Equal(t, tt.val, v)
			}
		})
	}
}
