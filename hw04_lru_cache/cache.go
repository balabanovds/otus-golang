package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	// Set new item to cache, returning true if item was just updated, or false for brand new item
	Set(key string, value interface{}) bool
	// Get item from cache, returning (nil, false) if 404, otherwise - (item, true)
	Get(key string) (interface{}, bool)
	// Clear cache
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[string]*cacheItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func newCacheItem(key string, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*cacheItem),
	}
}

func (c *lruCache) Set(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	i, ok := c.items[key]
	if ok {
		// if element exists in cache, than update it's value and move it at the beginning of queue
		i.value = value
		c.queue.MoveToFront(c.queue.Get(i))
		return true
	}

	// if there is no element in cache
	newItem := newCacheItem(key, value)

	// .. then ut it at the beginning of queue
	c.queue.PushFront(newItem)
	c.items[key] = newItem

	if c.queue.Len() > c.capacity {
		// remove last (less used) element from queue and from items
		last := c.queue.Back()
		c.queue.Remove(last)

		ci, ok := last.Value.(*cacheItem)
		if !ok {
			panic("failed to cast to *cacheItem")
		}

		delete(c.items, ci.key)
	}

	return false
}

func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// if element found than move it to the front
	c.queue.MoveToFront(c.queue.Get(val))

	return val.value, true
}

func (c *lruCache) Clear() {
	c.items = make(map[string]*cacheItem)
	c.queue = NewList()
}
