package domain

import (
	"context"
	"math/big"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateUserWithAccount(ctx context.Context, login, passwordHash string) (*entity.User, error)
	RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber big.Int) (*entity.Order, error)
	FindUserByLoginAndPassword(ctx context.Context, login, passwordHash string) (*entity.User, error)
	MarkOrderInvalid(ctx context.Context, orderNumber big.Int) error
	MarkOrderProcessing(ctx context.Context, orderNumber big.Int) error
	MarkOrderProcessedAndDepositAccount(ctx context.Context, userID uuid.UUID, orderNumber big.Int, amount decimal.Decimal) error
}
