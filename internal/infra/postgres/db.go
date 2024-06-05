package postgres

import (
	"context"
	"database/sql"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

func NewDB(ctx context.Context, cfg config.Database) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URI)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(ctx, db); err != nil {
		logger.Error("Run migration: failed", "err", err)
		return nil, err
	}

	return db, nil
}
