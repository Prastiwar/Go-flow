package di

import "reflect"

// Cache is implemented by any value that has a Get and Put method.
// The implementation controls where and how reflect.Type is stored for specified LifeTime.
type Cache interface {
	Get(LifeTime, reflect.Type) (interface{}, bool)
	Put(LifeTime, reflect.Type, interface{}) bool
}

type rootCache map[reflect.Type]interface{}

// NewRootCache returns a new Cache which is defined as root and will cache Singleton services.
func NewRootCache() *rootCache {
	return &rootCache{}
}

// Get returns a Singleton service if it exist in cache storage. Boolean defines if it was found.
// Other LifeTime's will always return nil and false.
func (c *rootCache) Get(life LifeTime, t reflect.Type) (interface{}, bool) {
	if life == Singleton {
		v, ok := (*c)[t]
		return v, ok
	}
	return nil, false
}

// Put returns true if LifeTime is Singleton and was successfully stored in cache.
func (c *rootCache) Put(life LifeTime, t reflect.Type, v interface{}) bool {
	if life == Singleton {
		(*c)[t] = v
		return true
	}
	return false
}

type scopeCache struct {
	root  Cache
	scope map[reflect.Type]interface{}
}

// NewScopeCache returns a new Cache which is defined as scope for specified root Cache.
// Can store both Singleton and Scoped LifeTime services.
func NewScopeCache(root Cache) *scopeCache {
	return &scopeCache{
		root:  root,
		scope: make(map[reflect.Type]interface{}),
	}
}

// Get returns service existing in cache. Singleton is retrieved from root Cache. Scoped is
// retrieved from internal storage.
func (c *scopeCache) Get(life LifeTime, t reflect.Type) (interface{}, bool) {
	v, ok := c.root.Get(life, t)
	if !ok && life == Scoped {
		v, ok = c.scope[t]
	}
	return v, ok
}

// Put returns true if LifeTime is Singleton or Scoped and was successfully stored in root or internal cache.
func (c *scopeCache) Put(life LifeTime, t reflect.Type, v interface{}) bool {
	ok := c.root.Put(life, t, v)
	if !ok && life == Scoped {
		c.scope[t] = v
		return true
	}
	return ok
}
