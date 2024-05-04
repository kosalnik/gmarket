package application

import (
	"context"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kosalnik/gmarket/internal/infra/postgres"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

type Application struct {
	cfg *config.Config
	db  *postgres.DB
}

func New(cfg *config.Config) *Application {
	return &Application{cfg: cfg}
}

func (app *Application) Run(_ context.Context) (err error) {
	ctx := context.Background()

	app.db, err = postgres.NewDB(ctx, app.cfg.Database)
	if err != nil {
		return err
	}

	logger.Info("Listen " + app.cfg.Server.Address)

	return http.ListenAndServe(app.cfg.Server.Address, app.GetRoutes())
}
