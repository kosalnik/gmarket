package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

type WithdrawalsResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func NewWithdrawalsHandler(ctx context.Context, orderService *service.OrderService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewOrderGetHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}

		res, err := orderService.Withdrawals(ctx, userID)
		if err != nil {
			logger.Error("GetWithdrawals error", "err", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
		if len(res) == 0 {
			http.Error(w, "204", http.StatusNoContent)
			return
		}
		ret := make([]WithdrawalsResponse, len(res))
		for i := range res {
			ret[i] = WithdrawalsResponse{
				Order:       res[i].OrderNumber.String(),
				Sum:         res[i].Amount.InexactFloat64(),
				ProcessedAt: res[i].UpdatedAt,
			}
		}

		b, err := json.Marshal(ret)
		if err != nil {
			logger.Error("GetWithdrawals marshal error", "err", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprint(w, string(b)); err != nil {
			logger.Error("GetWithdrawals write response error", "err", err)
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}
		// 200 — успешная обработка запроса.
		// 204 — нет ни одного списания.
		// 401 — пользователь не авторизован.
		// 500 — внутренняя ошибка сервера.
	}
}
