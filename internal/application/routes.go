package application

import (
	"context"
	"net/http"

	"github.com/kosalnik/gmarket/internal/handler"
	"github.com/kosalnik/gmarket/internal/infra/auth"
)

func (app *Application) GetRoutes(ctx context.Context) http.Handler {
	m := http.NewServeMux()

	authMw := auth.AuthMiddleware(app.authService)
	m.HandleFunc("POST /api/user/register", handler.NewRegisterHandler(ctx, app.userService, app.authService))
	m.HandleFunc("POST /api/user/login", handler.NewLoginHandler(ctx, app.userService, app.authService))
	m.HandleFunc("POST /api/user/orders", authMw(handler.NewOrderCreateHandler(ctx, app.orderService)))
	m.HandleFunc("GET /api/user/orders", authMw(handler.NewOrderGetHandler(ctx, app.orderService)))
	m.HandleFunc("GET /api/user/balance", authMw(handler.NewBalanceHandler(ctx, app.userService)))
	m.HandleFunc("POST /api/user/balance/withdraw", authMw(handler.NewWithdrawHandler(ctx, app.orderService)))
	m.HandleFunc("GET /api/user/withdrawals", authMw(handler.NewWithdrawalsHandler(ctx, app.orderService)))

	return m
}
