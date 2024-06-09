package crypt

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/kosalnik/gmarket/pkg/domain/service"
)

type PasswordHasher struct {
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (p PasswordHasher) Hash(pwd string) string {
	t := md5.Sum([]byte(pwd))
	return hex.EncodeToString(t[:])
}

func (p PasswordHasher) IsEquals(pwd string, h string) bool {
	return p.Hash(pwd) == h
}

var _ service.PasswordHasher = &PasswordHasher{}
