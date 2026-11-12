package limiter

import (
	"sync"
	"time"
)

// 固定时间窗口限流器
type FixedWindowLimiter struct {
	WindowSize  time.Duration
	WindowStart time.Time
	Capacity    int
	Count       int
	mu          sync.Mutex
}

func NewFixedWinowLimiter(windowSize time.Duration, capacity int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		WindowSize:  windowSize,
		WindowStart: time.Now(),
		Capacity:    capacity,
	}
}

func (limiter *FixedWindowLimiter) Allow() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := time.Now()
	if now.Sub(limiter.WindowStart) > limiter.WindowSize {
		limiter.WindowStart = now
		limiter.Capacity = 0
	}
	limiter.Count++
	return limiter.Count <= limiter.Capacity
}
