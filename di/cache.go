package di

import "reflect"

type DiCache interface {
	Get(LifeTime, reflect.Type) (interface{}, bool)
	Put(LifeTime, reflect.Type, interface{}) bool
}

type rootCache map[reflect.Type]interface{}

func NewRootCache() *rootCache {
	return &rootCache{}
}

func (c *rootCache) Get(life LifeTime, t reflect.Type) (interface{}, bool) {
	if life == Singleton {
		v, ok := (*c)[t]
		return v, ok
	}
	return nil, false
}

func (c *rootCache) Put(life LifeTime, t reflect.Type, v interface{}) bool {
	if life == Singleton {
		(*c)[t] = v
		return true
	}
	return false
}

type scopeCache struct {
	root  DiCache
	scope map[reflect.Type]interface{}
}

func NewScopeCache(root DiCache) *scopeCache {
	return &scopeCache{
		root:  root,
		scope: make(map[reflect.Type]interface{}),
	}
}

func (c *scopeCache) Get(life LifeTime, t reflect.Type) (interface{}, bool) {
	v, ok := c.root.Get(life, t)
	if !ok && life == Scoped {
		v, ok = c.scope[t]
	}
	return v, ok
}

func (c *scopeCache) Put(life LifeTime, t reflect.Type, v interface{}) bool {
	ok := c.root.Put(life, t, v)
	if !ok && life == Scoped {
		c.scope[t] = v
		return true
	}
	return ok
}
