package embedding

import (
	"sync"
	"time"
)

type embeddingCacheEntry struct {
	vector    []float32
	createdAt time.Time
}

type embeddingCache struct {
	mu      sync.RWMutex
	entries map[string]embeddingCacheEntry
	ttl     time.Duration
	maxSize int
}

func newEmbeddingCache(ttl time.Duration, maxSize int) *embeddingCache {
	return &embeddingCache{
		entries: make(map[string]embeddingCacheEntry, maxSize),
		ttl:     ttl,
		maxSize: maxSize,
	}
}

func (c *embeddingCache) get(key string) ([]float32, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Since(entry.createdAt) > c.ttl {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return nil, false
	}
	return entry.vector, true
}

func (c *embeddingCache) set(key string, vector []float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict expired entries if at capacity.
	if len(c.entries) >= c.maxSize {
		now := time.Now()
		for k, e := range c.entries {
			if now.Sub(e.createdAt) > c.ttl {
				delete(c.entries, k)
			}
		}
		// If still at capacity, evict oldest.
		if len(c.entries) >= c.maxSize {
			var oldestKey string
			var oldestTime time.Time
			for k, e := range c.entries {
				if oldestKey == "" || e.createdAt.Before(oldestTime) {
					oldestKey = k
					oldestTime = e.createdAt
				}
			}
			delete(c.entries, oldestKey)
		}
	}

	c.entries[key] = embeddingCacheEntry{
		vector:    vector,
		createdAt: time.Now(),
	}
}
