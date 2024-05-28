package service

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
)

type OrderService struct {
	repo           domain.Repository
	accrualService AccrualService
}

func NewOrderService(repo domain.Repository, accrualSvc AccrualService) (*OrderService, error) {
	return &OrderService{
		repo:           repo,
		accrualService: accrualSvc,
	}, nil
}

func (s *OrderService) RegisterOrder(ctx context.Context, userID uuid.UUID, orderNumber *big.Int) (*entity.Order, error) {
	o, err := s.repo.RegisterOrder(ctx, userID, *orderNumber)
	if err != nil {
		return nil, err
	}
	go func() {
		for i := 0; i < 5; i++ {
			result, err := s.accrualService.RegisterOrder(ctx, *orderNumber)
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
	return o, nil
}
