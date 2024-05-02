package application

import (
	"context"
	"net/http"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

type Application struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Application {
	return &Application{cfg: cfg}
}

func (app *Application) Run(_ context.Context) error {
	logger.Info("Listen " + app.cfg.Server.Address)

	return http.ListenAndServe(app.cfg.Server.Address, app.GetRoutes())
}
