package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/pkg/domain"

	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

func NewBalanceHandler(ctx context.Context, userService domain.UserService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := auth.UserIDFromContext(r.Context())
		if userID == uuid.Nil {
			logger.Info("NewBalanceHandler: unauthorized")
			http.Error(w, "401: unauthorized", http.StatusUnauthorized)
			return
		}
		acc, err := userService.GetAccount(ctx, userID)
		if err != nil {
			logger.Info("NewBalanceHandler: error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
		withdrawn, err := userService.GetSumWithdraw(ctx, userID)
		if err != nil {
			logger.Info("NewBalanceHandler: error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}

		ret := struct {
			Current   float64 `json:"current"`
			Withdrawn float64 `json:"withdrawn"`
		}{
			Current:   acc.Balance.InexactFloat64(),
			Withdrawn: withdrawn.InexactFloat64(),
		}
		resp, err := json.Marshal(ret)
		if err != nil {
			logger.Info("NewBalanceHandler: marshal", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
		if _, err := fmt.Fprint(w, string(resp)); err != nil {
			logger.Error("NewBalanceHandler: write error", "err", err)
			http.Error(w, "500: internal error", http.StatusInternalServerError)
			return
		}
	}
}
