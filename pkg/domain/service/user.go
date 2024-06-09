package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type PasswordHasher interface {
	Hash(pwd string) string
	IsEquals(pwd string, h string) bool
}

type UserService struct {
	repo   domain.Repository
	hasher PasswordHasher
}

func NewUserService(repo domain.Repository, h PasswordHasher) (*UserService, error) {
	return &UserService{
		repo:   repo,
		hasher: h,
	}, nil
}

func (s *UserService) Register(ctx context.Context, login, password string) (*entity.User, error) {
	pwdHash := s.hasher.Hash(password)
	return s.repo.CreateUserWithAccount(ctx, login, pwdHash)
}

func (s *UserService) FindByLoginAndPassword(ctx context.Context, login, password string) (*entity.User, error) {
	pwdHash := s.hasher.Hash(password)
	return s.repo.FindUserByLoginAndPassword(ctx, login, pwdHash)
}

func (s *UserService) GetAccount(ctx context.Context, userID uuid.UUID) (*entity.Account, error) {
	return s.repo.GetAccount(ctx, userID)
}

func (s *UserService) GetSumWithdraw(ctx context.Context, userID uuid.UUID) (*decimal.Decimal, error) {
	return s.repo.GetSumWithdraw(ctx, userID)
}
