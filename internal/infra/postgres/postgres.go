package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

type DB struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDB(ctx context.Context, cfg config.Database) (*DB, error) {
	db, err := sql.Open("pgx", cfg.URI)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(ctx, db); err != nil {
		logger.Error("Run migration: failed", "err", err)
		return nil, err
	}

	return &DB{db, sync.Mutex{}}, nil
}
