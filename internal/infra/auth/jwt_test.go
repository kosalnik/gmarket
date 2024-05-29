package auth_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/auth"
	"github.com/stretchr/testify/require"
)

func TestJwtEncoder_Encode(t *testing.T) {
	cfg := config.JWT{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	e, err := auth.NewJwtEncoder(cfg)
	require.NoError(t, err)
	id, err := uuid.NewV6()
	require.NoError(t, err)
	got, err := e.Encode(id)
	require.NoError(t, err)
	require.NotEmpty(t, got)
}

func TestJwtEncoder_Decode(t *testing.T) {
	vasyaID, err := uuid.Parse(`01ef18f3-9153-68d8-ab6b-74563c32efde`)
	require.NoError(t, err)
	hackToken := "qqqq"
	tests := []struct {
		name        string
		tokenString string
		currentTime time.Time
		want        *auth.JwtClaims
		wantErr     bool
	}{
		{
			name:        "positive",
			tokenString: goodToken,
			currentTime: time.Unix((goodExpiresAt.Unix()+goodIssuedAt.Unix())/2, 0),
			want: &auth.JwtClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        vasyaID.String(),
					IssuedAt:  jwt.NewNumericDate(goodIssuedAt),
					ExpiresAt: jwt.NewNumericDate(goodExpiresAt),
				},
			},
			wantErr: false,
		},
		{
			name:        "expired jwt",
			tokenString: goodToken,
			currentTime: time.Unix(goodExpiresAt.Unix()+300, 0),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "not issued jwt",
			tokenString: goodToken,
			currentTime: time.Unix(goodIssuedAt.Unix()-100000, 0),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "malformed",
			tokenString: goodToken + "BAD",
			currentTime: time.Unix(goodIssuedAt.Unix()+1, 0),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "hack wrong sign method in the jwt",
			tokenString: hackToken,
			currentTime: time.Unix(goodIssuedAt.Unix()+1, 0),
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "empty",
			tokenString: "",
			currentTime: time.Unix(goodIssuedAt.Unix()+1, 0),
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.JWT{
				PublicKey:  publicKey,
				PrivateKey: privateKey,
			}
			e, err := auth.NewJwtEncoder(cfg, jwt.WithTimeFunc(func() time.Time { return tt.currentTime }))
			require.NoError(t, err)

			got, err := e.Decode(tt.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
