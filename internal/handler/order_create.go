package handler

import (
	"context"
	"errors"
	"io"
	"math/big"
	"net/http"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/internal/infra/postgres"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

func NewOrderCreateHandler(ctx context.Context, orderService *service.OrderService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Warn("NewOrderCreateHandler: unable to close body")
			}
		}()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "400: no body", http.StatusBadRequest)
			return
		}
		orderNumber, ok := new(big.Int).SetString(string(b), 10)
		if !ok {
			http.Error(w, "400: wrong number", http.StatusBadRequest)
			return
		}
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewOrderCreateHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}
		_, err = orderService.RegisterOrder(ctx, userID, orderNumber)
		if err != nil {
			logger.Info("NewOrderCreateHandler: fail create order", "err", err)
			if errors.Is(err, postgres.ErrAlien) {
				http.Error(w, "409: Conflict", http.StatusConflict)
				return
			}
			if errors.Is(err, postgres.ErrAlreadyExists) {
				http.Error(w, "202: Accepted", http.StatusAccepted)
				return
			}
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}

	}
}
