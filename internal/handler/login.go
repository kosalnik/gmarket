package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

func NewLoginHandler(
	ctx context.Context,
	userService *service.UserService,
	authService auth.TokenEncoder,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest LoginRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &loginRequest); err != nil {
			http.Error(w, "Wrong request", http.StatusBadRequest)
			return
		}
		u, err := userService.FindByLoginAndPassword(ctx, loginRequest.Login, loginRequest.Password)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
