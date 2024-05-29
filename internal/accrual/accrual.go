package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
	"github.com/kosalnik/gmarket/pkg/domain/entity"
	"github.com/kosalnik/gmarket/pkg/domain/service"
	"github.com/shopspring/decimal"
)

type Accrual struct {
	cfg    config.AccrualSystem
	client *http.Client
	pool   *Pool
}

func NewAccrual(cfg config.AccrualSystem) (*Accrual, error) {
	return &Accrual{
		cfg: cfg,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

var _ service.AccrualService = &Accrual{}

//const (
//	StatusSuccess = "ok"
//)

func (a *Accrual) RegisterOrder(_ context.Context, orderNumber entity.OrderNumber) (*service.Result, error) {
	uri := fmt.Sprintf("%s/api/orders/%v", a.cfg.Address, orderNumber)
	logger.Debug("Accrual.RegisterOrder: Request", "uri", uri)
	resp, err := a.client.Get(uri)
	if err != nil {
		logger.Debug("Accrual.RegisterOrder: Request failed", "uri", uri, "err", err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warn("Accrual.RegisterOrder: Failed to close body")
		}
	}()
	logger.Debug("Accrual.RegisterOrder: Response", "uri", uri, "status", resp.StatusCode)
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, service.ErrInternalError
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, service.ErrToManyRequests
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, service.ErrNotRegistered
	}
	if resp.StatusCode != http.StatusOK {
		return nil, service.ErrUnknown
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("Accrual.RegisterOrder: Response body fail", "uri", uri, "body", b)
		return nil, err
	}
	logger.Debug("Accrual.RegisterOrder: Response", "uri", uri, "body", b)
	var r struct {
		Order   entity.OrderNumber
		Status  string
		Accrual decimal.Decimal
	}
	if err := json.Unmarshal(b, &r); err != nil {
		logger.Debug("Accrual.RegisterOrder: Unmarshal fail", "uri", uri, "err", err)
		return nil, err
	}
	ret := service.Result{
		Amount: r.Accrual,
		Status: r.Status,
	}
	logger.Debug("Accrual.RegisterOrder: Result", "uri", uri, "ret", ret)
	return &ret, nil
}
