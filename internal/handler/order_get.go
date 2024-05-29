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

func NewOrderGetHandler(ctx context.Context, orderService *service.OrderService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewOrderGetHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}
		orders, err := orderService.GetOrders(ctx, userID)
		if err != nil {
			logger.Info("NewOrderGetHandler: error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
		if len(orders) == 0 {
			http.Error(w, "204: Empty", http.StatusNoContent)
			return
		}
		ret := make([]struct {
			Number     string    `json:"number"`
			Status     string    `json:"status"`
			Accrual    float64   `json:"accrual"`
			UploadedAt time.Time `json:"uploaded_at"`
		}, len(orders))
		for i := range orders {
			ret[i] = struct {
				Number     string    `json:"number"`
				Status     string    `json:"status"`
				Accrual    float64   `json:"accrual"`
				UploadedAt time.Time `json:"uploaded_at"`
			}{
				Number:     orders[i].OrderNumber.String(),
				Status:     string(orders[i].Status),
				Accrual:    orders[i].Amount.InexactFloat64(),
				UploadedAt: orders[i].CreatedAt,
			}
		}
		resp, err := json.Marshal(ret)
		if err != nil {
			logger.Error("NewOrderGetHandler: marshal error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprint(w, string(resp)); err != nil {
			logger.Error("NewOrderGetHandler: write error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
	}
}
