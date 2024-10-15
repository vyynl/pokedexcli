package pokecache

import (
	"sync"
	"time"
)

// Organization for our cache
type Cache struct {
	mu   sync.Mutex
	data map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Creating a new cache
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]CacheEntry),
	}

	go cache.reapLoop(interval)

	return cache
}

// Purging old cache data
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.data {
				if time.Since(entry.createdAt) > interval {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

// Interracting with the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}
