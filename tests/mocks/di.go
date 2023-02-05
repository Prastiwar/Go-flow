package mocks

import (
	"reflect"

	"github.com/Prastiwar/Go-flow/di"
)

var (
	_ di.Cache = DiCacheMock{}
)

type DiCacheMock struct {
	OnGet func(l di.LifeTime, t reflect.Type) (interface{}, bool)
	OnPut func(l di.LifeTime, t reflect.Type, v interface{}) bool
}

func (m DiCacheMock) Get(l di.LifeTime, t reflect.Type) (interface{}, bool) {
	return m.OnGet(l, t)
}

func (m DiCacheMock) Put(l di.LifeTime, t reflect.Type, v interface{}) bool {
	return m.OnPut(l, t, v)
}
