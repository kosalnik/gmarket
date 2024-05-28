package service

import (
	"context"

	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type PasswordHasher interface {
	Hash(pwd string) (string, error)
	IsEquals(pwd string, h string) bool
}

type UserService struct {
	repo   domain.Repository
	hasher PasswordHasher
}

func NewUserService(userRepo domain.Repository, h PasswordHasher) (*UserService, error) {
	return &UserService{
		repo:   userRepo,
		hasher: h,
	}, nil
}

func (s *UserService) Register(ctx context.Context, login, password string) (*entity.User, error) {
	pwdHash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateUserWithAccount(ctx, login, pwdHash)
}

func (s *UserService) FindByLoginAndPassword(ctx context.Context, login, password string) (*entity.User, error) {
	pwdHash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	return s.repo.FindUserByLoginAndPassword(ctx, login, pwdHash)
}
