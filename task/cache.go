package task

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	ExpiresAt  time.Time
}

type TaskCache struct {
	mu       sync.RWMutex
	items    map[string]CacheItem
	defaultTTL time.Duration
}

func NewTaskCache(ttl time.Duration) *TaskCache {
	cache := &TaskCache{
		items:      make(map[string]CacheItem),
		defaultTTL: ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanupLoop()
	
	return cache
}

func (c *TaskCache) Set(key string, value interface{}, ttl ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	duration := c.defaultTTL
	if len(ttl) > 0 {
		duration = ttl[0]
	}
	
	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *TaskCache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil
	}
	
	if time.Now().After(item.ExpiresAt) {
		return nil
	}
	
	return item.Value
}

func (c *TaskCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *TaskCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

func (c *TaskCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if now.After(v.ExpiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}