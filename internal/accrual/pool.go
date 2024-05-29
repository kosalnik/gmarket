package accrual

import (
	"context"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type AccrualHandler func(ctx context.Context, orderNumber entity.OrderNumber)
type Pool struct {
	ch          chan entity.OrderNumber
	handler     AccrualHandler
	rateLimiter RateLimiter
}

func NewPool(ctx context.Context, cfg config.AccrualSystem, rateLimiter RateLimiter, hdl AccrualHandler) *Pool {
	p := &Pool{
		ch:          make(chan entity.OrderNumber),
		rateLimiter: rateLimiter,
		handler:     hdl,
	}
	logger.Info("Start accrual workers", "count", cfg.Concurrence)
	for i := 0; i < cfg.Concurrence; i++ {
		go p.worker(ctx)
	}
	return p
}

func (p *Pool) Handle(orderID *entity.OrderNumber) {
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
