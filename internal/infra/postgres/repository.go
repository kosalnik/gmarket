package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/shopspring/decimal"
)

var (
	ErrAlien         = errors.New("alien order")
	ErrAlreadyExists = errors.New("already exists")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*Repository, error) {
	return &Repository{db: db}, nil
}

var _ domain.Repository = &Repository{}

func (r *Repository) inTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if er := tx.Rollback(); er != nil {
				err = fmt.Errorf("failed to rollback transaction %w: %w", er, err)
			}
			return
		}
		if recover() != nil {
			err = tx.Rollback()
			return
		}
	}()
	if err = fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) CreateUserWithAccount(ctx context.Context, login, passwordHash string) (u *entity.User, err error) {
	return u, r.inTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if u, err = r.createUser(ctx, tx, login, passwordHash); err != nil {
			return err
		}
		if _, err = r.createAccount(ctx, tx, u.ID); err != nil {
			return err
		}
		return nil
	})
}

func (*Repository) createUser(ctx context.Context, tx *sql.Tx, login, passwordHash string) (*entity.User, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}
	u := &entity.User{ID: id, Login: login, Password: passwordHash}
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO "user" (id, login, password) VALUES ($1, $2, $3) ON CONFLICT (login) DO NOTHING`,
		u.ID, u.Login, u.Password,
	)
	if err != nil {
		return nil, err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return nil, errors.New("already exists")
	}
	return u, nil
}

func (*Repository) createAccount(ctx context.Context, tx *sql.Tx, userID uuid.UUID) (*entity.Account, error) {
	a := &entity.Account{UserID: userID, Balance: decimal.NewFromInt(0)}
	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO "account" (user_id, balance) VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING`,
		a.UserID, a.Balance,
	)
	if err != nil {
		return nil, err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		return nil, errors.New("already exists")
	}
	return a, nil
}

func (r *Repository) RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber big.Int) (*entity.Order, error) {
	id, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a := &entity.Order{
		ID:          id,
		UserID:      userID,
		OrderNumber: orderNumber,
		Amount:      decimal.NewFromInt(0),
		Status:      entity.OrderStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO "order" (id, user_id, order_number, amount, status, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (order_number) DO NOTHING`,
		a.ID, a.UserID, a.OrderNumber.String(), a.Amount, a.Status, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if n, err := res.RowsAffected(); err != nil || n == 0 {
		o, err := r.getOrder(ctx, orderNumber)
		if err != nil {
			return nil, err
		}
		if o.UserID != userID {
			return nil, ErrAlien
		}
		return nil, ErrAlreadyExists
	}
	return a, nil
}

func (r *Repository) getOrder(ctx context.Context, orderNumber big.Int) (*entity.Order, error) {
	q := `SELECT id, user_id, order_number, amount, status, created_at, updated_at FROM "order" WHERE order_number = $2`
	res := r.db.QueryRowContext(ctx, q, orderNumber.String())
	var o entity.Order
	if err := res.Scan(&o.ID, &o.UserID, &o.OrderNumber, &o.Amount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) FindUserByLoginAndPassword(ctx context.Context, login, passwordHash string) (*entity.User, error) {
	q := `SELECT id, login, password FROM "user" WHERE login = $1 AND password = $2`
	logger.Debug("db: %s", login, passwordHash)
	res := r.db.QueryRowContext(ctx, q, login, passwordHash)
	var u entity.User
	if err := res.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		logger.Info("FindUserByLoginAndPassword failed", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) MarkOrderInvalid(ctx context.Context, orderNumber big.Int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "order" SET status = $1, updated_at = $2 WHERE order_number = $3`,
		entity.OrderStatusRejected, time.Now(), orderNumber,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) MarkOrderProcessing(ctx context.Context, orderNumber big.Int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "order" SET status = $1, updated_at = $2 WHERE order_number = $3`,
		entity.OrderStatusProcessing, time.Now(), orderNumber,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) MarkOrderProcessedAndDepositAccount(ctx context.Context, userID uuid.UUID, orderNumber big.Int, amount decimal.Decimal) error {
	return r.inTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		acc := entity.Account{UserID: userID}
		res := tx.QueryRowContext(
			ctx,
			`SELECT balance FROM "account" WHERE user_id = $1 FOR UPDATE`,
			acc.UserID,
		)
		if err := res.Scan(&acc.Balance); err != nil {
			logger.Info("MarkOrderProcessedAndDepositAccount failed", "err", err)
			return err
		}
		_, err := tx.ExecContext(
			ctx,
			`UPDATE "account" SET balance = balance + $1 WHERE user_id = $2`,
			amount, acc.UserID,
		)
		if err != nil {
			return err
		}
		if err := r.markOrderProcessed(ctx, tx, orderNumber, amount); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) markOrderProcessed(ctx context.Context, tx *sql.Tx, orderNumber big.Int, amount decimal.Decimal) error {
	_, err := tx.ExecContext(
		ctx,
		`UPDATE "order" SET status = $1, updated_at = $2 WHERE order_number = $3`,
		entity.OrderStatusProcessed, time.Now(), orderNumber,
	)
	if err != nil {
		return err
	}
	return nil
}
