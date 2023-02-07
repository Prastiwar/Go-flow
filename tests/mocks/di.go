package mocks

import (
	"reflect"

	"github.com/Prastiwar/Go-flow/di"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

var (
	_ di.Cache = DiCacheMock{}
)

type DiCacheMock struct {
	OnGet func(l di.LifeTime, t reflect.Type) (interface{}, bool)
	OnPut func(l di.LifeTime, t reflect.Type, v interface{}) bool
}

func (m DiCacheMock) Get(l di.LifeTime, t reflect.Type) (interface{}, bool) {
	assert.ExpectCall(m.OnGet)
	return m.OnGet(l, t)
}

func (m DiCacheMock) Put(l di.LifeTime, t reflect.Type, v interface{}) bool {
	assert.ExpectCall(m.OnPut)
	return m.OnPut(l, t, v)
}
