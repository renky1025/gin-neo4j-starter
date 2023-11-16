package gocache

import "sync"

type InMemoryCache struct {
	m     sync.Mutex
	store map[string]interface{}
}

func NewMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		m:     sync.Mutex{},
		store: make(map[string]interface{}),
	}
}

func (c *InMemoryCache) Get(key string) interface{} {
	return c.store[key]
}

func (c *InMemoryCache) Set(key string, value interface{}, timeSec int) bool {
	c.m.Lock()
	defer c.m.Unlock()
	c.store[key] = value
	return true
}

func (c *InMemoryCache) Delete(key string) error {
	//if _, ok := c.store[key]; ok {
	delete(c.store, key)
	//}
	return nil
}
func (c *InMemoryCache) Clear() error {
	c.store = make(map[string]interface{})
	return nil
}
