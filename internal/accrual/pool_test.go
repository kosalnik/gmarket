package accrual

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kosalnik/gmarket/internal/accrual/mock"
	"github.com/kosalnik/gmarket/internal/config"
)

func TestPool_Worker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rateLimiterMock := mock.NewMockRateLimiter(ctrl)
	rateLimiterMock.EXPECT().Fire().Times(10).Return(1)
	w := NewPool(
		ctx,
		config.AccrualSystem{Concurrence: 10},
		rateLimiterMock,
		func(ctx context.Context, orderNumber big.Int) {

		},
	)
	for i := 0; i < 10; i++ {
		w.Handle(big.NewInt(12345))
	}
	<-time.After(time.Second)
}
