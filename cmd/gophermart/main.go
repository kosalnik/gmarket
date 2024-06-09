package main

import (
	"context"

	"github.com/kosalnik/gmarket/internal/application"
	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

func main() {
	cfg := config.NewConfig()

	if err := logger.InitLogger(cfg.Logger); err != nil {
		panic(err)
	}
	app := application.New(cfg)
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
