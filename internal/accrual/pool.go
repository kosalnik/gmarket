package accrual

import (
	"context"
	"math/big"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

type AccrualHandler func(ctx context.Context, orderNumber big.Int)
type Pool struct {
	ch          chan big.Int
	handler     AccrualHandler
	rateLimiter RateLimiter
}

func NewPool(ctx context.Context, cfg config.AccrualSystem, rateLimiter RateLimiter, hdl AccrualHandler) *Pool {
	p := &Pool{
		ch:          make(chan big.Int),
		rateLimiter: rateLimiter,
		handler:     hdl,
	}
	logger.Info("Start accrual workers", "count", cfg.Concurrence)
	for i := 0; i < cfg.Concurrence; i++ {
		go p.worker(ctx)
	}
	return p
}

func (p *Pool) Handle(orderID *big.Int) {
	p.ch <- *orderID
}

func (p *Pool) worker(ctx context.Context) {
	for {
		select {
		case orderNum, ok := <-p.ch:
			if !ok {
				return
			}
			p.handler(ctx, orderNum)
			if p.rateLimiter.Fire() < 1 {
				p.rateLimiter.Wait(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}
