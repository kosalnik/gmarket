package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type OrderService struct {
	repo           domain.Repository
	accrualService AccrualService
}

var (
	ErrMoneyNotEnough   = errors.New("money are not enough")
	ErrWrongOrderNumber = errors.New("money are not enough")
)

func NewOrderService(repo domain.Repository, accrualSvc AccrualService) (*OrderService, error) {
	return &OrderService{
		repo:           repo,
		accrualService: accrualSvc,
	}, nil
}

func (s *OrderService) RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber *entity.OrderNumber) (*entity.Order, error) {
	var err error
	o, err := s.repo.RegisterOrder(ctx, userID, *orderNumber)
	if err != nil {
		return nil, err
	}
	go func() {
		for i := 0; i < 5; i++ {
			var result *Result
			result, err = s.accrualService.RegisterOrder(ctx, *orderNumber)
			if err != nil {
				logger.Error("RegisterOrder accrualError", "err", err)
				<-time.After(time.Second)
				continue
			}
			logger.Info("RegisterOrder accrual response", "response", result)
			switch result.Status {
			case "REGISTERED":
				if err = s.repo.MarkOrderProcessing(ctx, *orderNumber); err == nil {
					return
				}
			case "INVALID":
				if err = s.repo.MarkOrderInvalid(ctx, *orderNumber); err == nil {
					return
				}
			case "PROCESSING":
				return
			case "PROCESSED":
				if err = s.repo.MarkOrderProcessedAndDepositAccount(ctx, userID, *orderNumber, result.Amount); err == nil {
					return
				}
			}
			if err != nil {
				logger.Error("RegisterOrder result error", "err", err)
			}
		}
	}()
	if errors.Is(err, ErrNotRegistered) {
		return o, s.repo.MarkOrderInvalid(ctx, *orderNumber)
	}
	return o, err
}

func (s *OrderService) GetOrders(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error) {
	return s.repo.GetOrders(ctx, userID)
}

func (s *OrderService) Withdraw(ctx context.Context, userID uuid.UUID, orderNumber entity.OrderNumber, sum decimal.Decimal) error {
	return s.repo.Withdraw(ctx, userID, orderNumber, sum)
}

func (s *OrderService) Withdrawals(ctx context.Context, userID uuid.UUID) ([]*entity.Withdraw, error) {
	return s.repo.Withdrawals(ctx, userID)
}
