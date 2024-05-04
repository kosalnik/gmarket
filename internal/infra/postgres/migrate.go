package postgres

import (
	"context"
	"database/sql"
	"embed"

	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func RunMigrations(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	logger.Info("DB Migration: start")
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return err
	}
	logger.Info("DB Migration: success")

	return nil
}
