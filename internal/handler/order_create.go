package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/pkg/domain"

	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/internal/infra/postgres"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

func NewOrderCreateHandler(ctx context.Context, orderService domain.OrderService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewOrderCreateHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}
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
		orderNumber, err := new(entity.OrderNumber).FromBytes(b)
		if err != nil || !orderNumber.Valid() {
			logger.Info("NewOrderCreateHandler: wrong number", "err", err)
			http.Error(w, "422: wrong number", http.StatusUnprocessableEntity)
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
				http.Error(w, "200: Accepted", http.StatusOK)
				return
			}
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}

		http.Error(w, "202: Accepted", http.StatusAccepted)
	}
}
