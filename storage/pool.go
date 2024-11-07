package storage

import (
	"sync"
	"time"
)

type PoolItem struct {
	Value      interface{}
	LastUsed   time.Time
	InUse      bool
}

type ConnectionPool struct {
	mu       sync.Mutex
	items    map[string]*PoolItem
	maxSize  int
	ttl      time.Duration
	createFn func() (interface{}, error)
	closeFn  func(interface{}) error
}

func NewConnectionPool(maxSize int, ttl time.Duration, 
	createFn func() (interface{}, error), 
	closeFn func(interface{}) error) *ConnectionPool {
	
	pool := &ConnectionPool{
		items:    make(map[string]*PoolItem),
		maxSize:  maxSize,
		ttl:      ttl,
		createFn: createFn,
		closeFn:  closeFn,
	}
	
	go pool.cleanupLoop()
	return pool
}

func (p *ConnectionPool) Acquire() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// Find available item
	for id, item := range p.items {
		if !item.InUse && time.Since(item.LastUsed) < p.ttl {
			item.InUse = true
			item.LastUsed = time.Now()
			return item.Value, nil
		}
	}
	
	// Create new if under max size
	if len(p.items) < p.maxSize {
		value, err := p.createFn()
		if err != nil {
			return nil, err
		}
		
		id := generateID()
		p.items[id] = &PoolItem{
			Value:    value,
			LastUsed: time.Now(),
			InUse:    true,
		}
		return value, nil
	}
	
	return nil, fmt.Errorf("connection pool exhausted")
}

func (p *ConnectionPool) Release(value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	for _, item := range p.items {
		if item.Value == value {
			item.InUse = false
			item.LastUsed = time.Now()
			return
		}
	}
}

func (p *ConnectionPool) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		p.mu.Lock()
		for id, item := range p.items {
			if !item.InUse && time.Since(item.LastUsed) > p.ttl {
				p.closeFn(item.Value)
				delete(p.items, id)
			}
		}
		p.mu.Unlock()
	}
}

func generateID() string {
	return fmt.Sprintf("conn-%d", time.Now().UnixNano())
}