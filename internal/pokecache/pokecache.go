package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.RWMutex
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cacheEntry, ok := c.entries[key]
	return cacheEntry.val, ok
}

func (c Cache) realLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for true {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.entries {
				if time.Since(entry.createdAt) > interval {
					delete(c.entries, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries: map[string]cacheEntry{},
		mu:      &sync.RWMutex{},
	}
	go cache.realLoop(interval)
	return cache
}
