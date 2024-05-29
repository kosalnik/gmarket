package auth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

//type authMiddleware struct {
//	encoder TokenEncoder
//}

type TokenEncoder interface {
	Decode(tokenString string) (*JwtClaims, error)
	Encode(id uuid.UUID) (string, error)
}

var userIDKey = &struct{}{}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	t, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return uuid.Nil
	}
	if ret, err := uuid.Parse(t); err == nil {
		return ret
	}
	return uuid.Nil
}

func AuthMiddleware(encoder TokenEncoder) func(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				logger.Debug("authorization required")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token, err := encoder.Decode(tokenString)
			if err != nil {
				logger.Info("wrong jwt", "err", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next(w, r.WithContext(context.WithValue(r.Context(), userIDKey, token.ID)))
		}
	}
}
