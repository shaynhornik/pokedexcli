package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry 
	mu sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		entries: make(map[string]cacheEntry),
	}
	c.reapLoop()
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry {time.Now(), val}
	c.mu.Lock()
	c.entries[key] = entry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	value, ok := c.entries[key]
	c.mu.Unlock()
	if ok {
		return value, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop() {
	
