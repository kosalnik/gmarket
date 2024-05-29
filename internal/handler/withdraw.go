package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/kosalnik/gmarket/pkg/domain/service"
	"github.com/shopspring/decimal"
)

type WithdrawRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}

func NewWithdrawHandler(ctx context.Context, orderService *service.OrderService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewWithdrawHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Warn("NewWithdrawHandler: unable to close body")
			}
		}()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Debug("NewWithdrawHandler close body error", "err", err)
			http.Error(w, "500: no body", http.StatusInternalServerError)
			return
		}
		var req WithdrawRequest
		if err := json.Unmarshal(b, &req); err != nil {
			logger.Debug("NewWithdrawHandler unmarshal error", "err", err)
			http.Error(w, "500: fail unmarshal", http.StatusInternalServerError)
		}

		err = orderService.Withdraw(ctx, userID, entity.OrderNumber(req.Order), req.Sum)
		if err == nil {
			// 200 — успешная обработка запроса;
			return
		}
		if errors.Is(err, service.ErrMoneyNotEnough) {
			http.Error(w, "402: money not enough", http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, service.ErrWrongOrderNumber) {
			http.Error(w, "422: wrong order number", http.StatusUnprocessableEntity)
			return
		}
		// 401 — пользователь не авторизован;
		// 402 — на счету недостаточно средств;
		// 422 — неверный номер заказа;
		// 500 — внутренняя ошибка сервера.
		logger.Error("withdraw error", "err", err)
		http.Error(w, "500: internal error", http.StatusInternalServerError)
	}
}
