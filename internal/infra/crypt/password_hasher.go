package crypt

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/service"
)

type PasswordHasher struct {
}

func NewPasswordHasher() (*PasswordHasher, error) {
	return &PasswordHasher{}, nil
}

func (p PasswordHasher) Hash(pwd string) (string, error) {
	t := md5.Sum([]byte(pwd))
	return hex.EncodeToString(t[:]), nil
}

func (p PasswordHasher) IsEquals(pwd string, h string) bool {
	n, err := p.Hash(pwd)
	if err != nil {
		logger.Info("Calculate password hash has been failed", err.Error())
		return false
	}
	return n == h
}

var _ service.PasswordHasher = &PasswordHasher{}
