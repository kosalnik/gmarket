package auth_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/stretchr/testify/require"
)

func TestHashCheckMiddleware(t *testing.T) {
	tests := map[string]struct {
		token      string
		want       int
		wantUserID uuid.UUID
	}{
		"success": {
			token:      goodToken,
			want:       200,
			wantUserID: uuid.Must(uuid.Parse("01ef18f3-9153-68d8-ab6b-74563c32efde")),
		},
		"empty token": {
			token:      "",
			want:       401,
			wantUserID: uuid.Nil,
		},
	}
	encoder, err := auth.NewJwtEncoder(
		config.JWT{PrivateKey: privateKey, PublicKey: publicKey},
		jwt.WithTimeFunc(func() time.Time { return time.Unix((goodExpiresAt.Unix()+goodIssuedAt.Unix())/2, 0) }),
	)
	require.NoError(t, err)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mw := auth.AuthMiddleware(encoder)
			mux := http.NewServeMux()
			h := mw(func(writer http.ResponseWriter, request *http.Request) {
				got := auth.UserIDFromContext(request.Context())
				require.Equal(t, tt.wantUserID, got)
			})
			mux.HandleFunc(`/`, h)
			r := httptest.NewRequest(http.MethodPost, `/`, strings.NewReader("test"))
			r.Header.Set("Authorization", tt.token)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			got := w.Code
			if got != tt.want {
				t.Errorf("AuthMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
