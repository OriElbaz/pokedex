package pokecache

import (
	"sync"
	"time"
)

// CACHE STRUCTS //

type Cache struct {
	Cache map[string]cacheEntry
	mu sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

// CACHE FUNCTIONS //

func NewCache(duration time.Duration) *Cache{
	c := &Cache{
		Cache: make(map[string]cacheEntry),
		mu: sync.Mutex{},
	}

	go c.reapLoop(duration)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Cache[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.Cache[key]
	if ok == false {
		return nil, ok
	}

	val := data.val
	return val, ok
}

func (c *Cache) reapLoop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.reap(duration)
	}
}

func (c *Cache) reap(duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	timeCutOff := time.Now().Add(-duration)

	for key, value := range c.Cache {
		if value.createdAt.Before(timeCutOff) {
			delete(c.Cache, key)
		}

	}

}