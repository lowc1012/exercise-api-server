package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cache := NewCache[string]()

		// Test Set and Get
		cache.Set("key1", "value1", time.Minute)
		val, exists := cache.Get("key1")
		if !exists {
			t.Error("expected key1 to exist")
		}
		if val != "value1" {
			t.Errorf("expected value1, got %s", val)
		}

		_, exists = cache.Get("aaa")
		if exists {
			t.Error("expected `aaa` key to not exist")
		}

		cache.Delete("key1")
		_, exists = cache.Get("key1")
		if exists {
			t.Error("expected key1 to be deleted")
		}
	})

	t.Run("expiration", func(t *testing.T) {
		cache := NewCache[int]()

		// Set with short TTL
		cache.Set("exp", 123, time.Millisecond*50)

		// Verify it exists
		val, exists := cache.Get("exp")
		if !exists {
			t.Error("expected key to exist initially")
		}
		if val != 123 {
			t.Errorf("expected 123, got %d", val)
		}

		// Wait for expiration
		time.Sleep(time.Millisecond * 100)

		// Verify it's gone
		_, exists = cache.Get("exp")
		if exists {
			t.Error("expected key to be expired")
		}
	})

	t.Run("fetch all", func(t *testing.T) {
		cache := NewCache[int]()

		// Set multiple items
		cache.Set("a", 1, time.Minute)
		cache.Set("b", 2, time.Minute)
		cache.Set("c", 3, time.Millisecond*10)

		// Initial check
		items := cache.FetchAll()
		if len(items) != 3 {
			t.Errorf("expected 3 items, got %d", len(items))
		}

		// Wait for one item to expire
		time.Sleep(time.Millisecond * 200)

		// Check again
		items = cache.FetchAll()
		if len(items) != 2 {
			t.Errorf("expected 2 items after expiration, got %d", len(items))
		}
	})
}
