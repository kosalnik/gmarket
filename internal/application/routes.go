package application

import (
	"net/http"

	"github.com/kosalnik/gmarket/internal/handler"
)

func (app *Application) GetRoutes() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("POST /api/user/register", handler.NewRegisterHandler())
	m.HandleFunc("POST /api/user/login", handler.NewLoginHandler())
	m.HandleFunc("POST /api/user/orders", handler.NewOrderCreateHandler())
	m.HandleFunc("GET /api/user/orders", handler.NewOrderGetHandler())
	m.HandleFunc("GET /api/user/balance", handler.NewBalanceHandler())
	m.HandleFunc("POST /api/user/balance/withdraw", handler.NewWithdrawHandler())
	m.HandleFunc("GET /api/user/withdrawals", handler.NewWithdrawalsHandler())

	return m
}
