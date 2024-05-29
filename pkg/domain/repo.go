package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateUserWithAccount(ctx context.Context, login, passwordHash string) (*entity.User, error)
	RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber entity.OrderNumber) (*entity.Order, error)
	FindUserByLoginAndPassword(ctx context.Context, login, passwordHash string) (*entity.User, error)
	MarkOrderInvalid(ctx context.Context, orderNumber entity.OrderNumber) error
	MarkOrderProcessing(ctx context.Context, orderNumber entity.OrderNumber) error
	MarkOrderProcessedAndDepositAccount(ctx context.Context, userID uuid.UUID, orderNumber entity.OrderNumber, amount decimal.Decimal) error
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error)
	GetAccount(ctx context.Context, userID uuid.UUID) (*entity.Account, error)
	Withdraw(ctx context.Context, userID uuid.UUID, orderNumber entity.OrderNumber, sum decimal.Decimal) error
	Withdrawals(ctx context.Context, userID uuid.UUID) ([]*entity.Withdraw, error)
}
