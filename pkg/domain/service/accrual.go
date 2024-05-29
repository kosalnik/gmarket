package service

import (
	"context"
	"errors"

	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/shopspring/decimal"
)

var (
	ErrToManyRequests = errors.New(`too many requests`)
	ErrInternalError  = errors.New(`internal error`)
	ErrNotRegistered  = errors.New(`order is not registered`)
	ErrUnknown        = errors.New(`unknown error`)
)

type Result struct {
	Amount decimal.Decimal
	Status string
}

type AccrualService interface {
	RegisterOrder(ctx context.Context, orderNumber entity.OrderNumber) (*Result, error)
}