package limiter

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	Capacity   int
	TokenCount int
	RefillTime time.Time // 记录每次补充令牌的时间
	RefilRate  int       // 补充频率，每秒补充多少令牌
	mu         sync.Mutex
}

func NewTokenLimiter(capacity int, startQuantity int, refillRate int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		Capacity:   capacity,
		TokenCount: startQuantity, // 起始数量少时是冷启动，起始数量满了可以应对一波突发
		RefilRate:  refillRate,
		RefillTime: time.Now(),
	}
}

func (limiter *TokenBucketLimiter) Allow() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := time.Now()

	// 惰性补充令牌, 计算距离上次补充的时间，将这段时间的令牌一次性补充
	elapsed := now.Sub(limiter.RefillTime).Seconds()
	tokenAdd := int(elapsed) * limiter.RefilRate
	if tokenAdd > 0 {
		limiter.TokenCount = min(limiter.Capacity, limiter.TokenCount+tokenAdd)
		limiter.RefillTime = now
	}

	if limiter.TokenCount > 0 {
		limiter.TokenCount--
		return true
	}
	return false
}
