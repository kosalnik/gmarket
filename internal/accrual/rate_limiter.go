package accrual

import (
	"context"
	"sync"
	"time"
)

//go:generate mockgen -source=rate_limiter.go -destination=./mock/rate_limiter.go -package=mock
type RateLimiter interface {
	Fire() int
	Wait(ctx context.Context)
}

type RateLimit struct {
	mu            sync.Mutex
	lastPoint     time.Time
	count         int
	limit         int
	limitDuration time.Duration
}

func NewRateLimiter(limit int, limitDuration time.Duration) *RateLimit {
	return &RateLimit{
		mu:            sync.Mutex{},
		limit:         limit,
		limitDuration: limitDuration,
	}
}

// Возвращает сколько попыток до исчерпания лимита осталось
func (r *RateLimit) Fire() int {
	now := time.Now().Round(r.limitDuration)
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.lastPoint.Equal(now) {
		r.count = 0
		r.lastPoint = now
		return r.limit
	}
	r.count++
	return r.limit - r.count
}

func (r *RateLimit) Wait(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(r.limitDuration):
		return
	}
}
