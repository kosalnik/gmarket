package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type UserService interface {
	Register(ctx context.Context, login, password string) (*entity.User, error)
	FindByLoginAndPassword(ctx context.Context, login, password string) (*entity.User, error)
	GetAccount(ctx context.Context, userID uuid.UUID) (*entity.Account, error)
	GetSumWithdraw(ctx context.Context, userID uuid.UUID) (*decimal.Decimal, error)
}

type OrderService interface {
	RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber *entity.OrderNumber) (*entity.Order, error)
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error)
	Withdraw(ctx context.Context, userID uuid.UUID, orderNumber entity.OrderNumber, sum decimal.Decimal) error
	Withdrawals(ctx context.Context, userID uuid.UUID) ([]*entity.Withdraw, error)
}
