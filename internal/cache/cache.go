package cache

import (
	"sync"
	"time"
)

// Cache simple generic cache using sync.Map
type Cache[T any] struct {
	store sync.Map
}

type Item[T any] struct {
	Value     T
	ExpiresAt int64
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		store: sync.Map{},
	}
}

// FetchAll retrieves all items from cache
func (c *Cache[T]) FetchAll() []T {
	var result []T
	c.store.Range(func(k, v any) bool {
		if item, ok := v.(Item[T]); ok {
			// Check if the item has expired
			if time.Now().Unix() <= item.ExpiresAt {
				result = append(result, item.Value)
			} else {
				c.store.Delete(k)
			}
		}
		return true
	})
	return result
}

// Get retrieves an item from cache
func (c *Cache[T]) Get(key string) (T, bool) {
	var t T
	data, exists := c.store.Load(key)
	if !exists {
		return t, false
	}

	item := data.(Item[T])
	// Check if expired
	if time.Now().Unix() > item.ExpiresAt {
		c.store.Delete(key) // Remove expired item
		return t, false
	}

	return item.Value, true
}

// Set stores an item in cache with TTL
func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	expiration := time.Now().Add(ttl).Unix()
	c.store.Store(key, Item[T]{Value: value, ExpiresAt: expiration})
}

// Delete deletes an item in cache
func (c *Cache[T]) Delete(key string) {
	c.store.Delete(key)
}
