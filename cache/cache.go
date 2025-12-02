package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.expiresAt)
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

func NewCache(ttlSeconds int64) *Cache {
	cache := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     time.Duration(ttlSeconds) * time.Second,
	}
	log.Printf("Cache initialized with TTL: %d seconds", ttlSeconds)
	return cache
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &CacheEntry{
		data:      value,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.entries[key] = entry
	log.Printf("Cache SET: %s (expires in %d seconds)", key, int64(c.ttl.Seconds()))
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		log.Printf("Cache GET: %s (MISS)", key)
		return nil, false
	}

	if entry.IsExpired() {
		log.Printf("Cache GET: %s (EXPIRED)", key)
		return nil, false
	}

	log.Printf("Cache GET: %s (HIT)", key)
	return entry.data, true
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[key]; exists {
		delete(c.entries, key)
		log.Printf("Cache INVALIDATE: %s", key)
	}
}

func (c *Cache) InvalidatePattern(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	invalidatedCount := 0
	for key := range c.entries {
		// Simple pattern matching: if key starts with pattern or contains pattern
		if len(key) >= len(pattern) && key[:len(pattern)] == pattern {
			delete(c.entries, key)
			invalidatedCount++
		}
	}
	if invalidatedCount > 0 {
		log.Printf("Cache INVALIDATE PATTERN: %s (invalidated %d entries)", pattern, invalidatedCount)
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.entries)
	c.entries = make(map[string]*CacheEntry)
	log.Printf("Cache CLEAR: cleared %d entries", count)
}

// Helper functions to generate cache keys
func BookListKey() string {
	return "books:list"
}

func BookIDKey(id uint) string {
	return fmt.Sprintf("books:id:%d", id)
}

func ReaderListKey() string {
	return "readers:list"
}

func ReaderIDKey(id uint) string {
	return fmt.Sprintf("readers:id:%d", id)
}
