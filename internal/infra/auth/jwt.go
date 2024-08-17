package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
)

type JwtEncoder struct {
	privateKey *rsa.PrivateKey
	parser     *jwt.Parser
	keyFunc    jwt.Keyfunc
}

var ErrInvalidToken = errors.New("invalid token")
var ErrEmptyPublicKey = errors.New("empty public key")
var ErrEmptyPrivateKey = errors.New("empty private key")

func NewJwtEncoder(cfg config.JWT, opts ...jwt.ParserOption) (*JwtEncoder, error) {
	if cfg.PublicKey == "" {
		return nil, ErrEmptyPublicKey
	}
	if cfg.PrivateKey == "" {
		return nil, ErrEmptyPrivateKey
	}
	pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		return nil, err
	}
	opts = append(DefaultOpts(), opts...)
	return &JwtEncoder{
		privateKey: pk,
		parser:     jwt.NewParser(opts...),
		keyFunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
			return key, err
		},
	}, nil
}

func (e JwtEncoder) Encode(userID uuid.UUID) (string, error) {
	p := JwtClaims{
		jwt.RegisteredClaims{
			ID:        userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, p)
	return t.SignedString(e.privateKey)
}

func (e JwtEncoder) Decode(tokenString string) (*JwtClaims, error) {
	payload := &JwtClaims{}
	token, err := e.parser.ParseWithClaims(tokenString, payload, e.keyFunc)
	if err != nil {
		logger.Info("failed to parse token", err)
		return nil, err
	}
	if !token.Valid {
		logger.Info("invalid token")
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func DefaultOpts() []jwt.ParserOption {
	return []jwt.ParserOption{
		jwt.WithIssuedAt(),
		jwt.WithLeeway(time.Minute),
	}
}
