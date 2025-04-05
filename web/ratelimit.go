package web

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.RWMutex
	visitors  map[string]*Visitor
	limit     int
	window    time.Duration
}

type Visitor struct {
	lastSeen  time.Time
	count     int
	mu        sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		limit:    limit,
		window:   window,
	}
	
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.RLock()
	visitor, exists := rl.visitors[ip]
	rl.mu.RUnlock()
	
	if !exists {
		rl.mu.Lock()
		visitor = &Visitor{
			lastSeen: time.Now(),
			count:    1,
		}
		rl.visitors[ip] = visitor
		rl.mu.Unlock()
		return true
	}
	
	visitor.mu.Lock()
	defer visitor.mu.Unlock()
	
	now := time.Now()
	if now.Sub(visitor.lastSeen) > rl.window {
		visitor.count = 1
		visitor.lastSeen = now
		return true
	}
	
	if visitor.count >= rl.limit {
		return false
	}
	
	visitor.count++
	visitor.lastSeen = now
	return true
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Backpressure mechanism
type Backpressure struct {
	queueSize    int
	maxQueueSize int
	mu           sync.Mutex
	cond         *sync.Cond
}

func NewBackpressure(maxSize int) *Backpressure {
	bp := &Backpressure{
		maxQueueSize: maxSize,
	}
	bp.cond = sync.NewCond(&bp.mu)
	return bp
}

func (bp *Backpressure) Acquire() bool {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	
	for bp.queueSize >= bp.maxQueueSize {
		bp.cond.Wait()
		return false
	}
	
	bp.queueSize++
	return true
}

func (bp *Backpressure) Release() {
	bp.mu.Lock()
	bp.queueSize--
	bp.mu.Unlock()
	bp.cond.Signal()
}

// Middleware for rate limiting
func RateLimitMiddleware(rl *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !rl.Allow(ip) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next(w, r)
		}
	}
}