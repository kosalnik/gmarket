package service

import (
	"context"
	"math/big"

	"github.com/shopspring/decimal"
)

type Result struct {
	Amount decimal.Decimal
	Status string
}

type AccrualService interface {
	RegisterOrder(ctx context.Context, orderNumber big.Int) (*Result, error)
}
