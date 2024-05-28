package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	UserID  uuid.UUID       `db:"user_id"`
	Balance decimal.Decimal `db:"balance"`
}
