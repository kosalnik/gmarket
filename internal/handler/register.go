package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewRegisterHandler(
	ctx context.Context,
	userService *service.UserService,
	authService auth.TokenEncoder,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var loginRequest LoginRequest
		body, err := io.ReadAll(req.Body)
		logger.Debug("Handle RegisterUser", body)
		if err != nil {
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &loginRequest); err != nil {
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}
		logger.Debug("RegisterUser", loginRequest)
		u, err := userService.Register(ctx, loginRequest.Login, loginRequest.Password)
		if err != nil {
			logger.Info("RegisterUser failed", err)
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}
		t, err := authService.Encode(u.ID)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Authorization", t)
	}
}
