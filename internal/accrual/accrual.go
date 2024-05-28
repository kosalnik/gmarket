package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/kosalnik/gmarket/internal/config"
	"github.com/kosalnik/gmarket/internal/infra/logger"
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

var (
	ErrToManyRequests = errors.New(`too many requests`)
	ErrInternalError  = errors.New(`internal error`)
	ErrNotFound       = errors.New(`order is not registered`)
	ErrUnknown        = errors.New(`unknown error`)
)

const (
	StatusSuccess = "ok"
)

func (a *Accrual) RegisterOrder(ctx context.Context, orderNumber big.Int) (*service.Result, error) {
	resp, err := a.client.Get(fmt.Sprintf("%s/api/orders/%s", a.cfg.Address, orderNumber.String()))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warn("Accrual.RegisterOrder: Failed to close body")
		}
	}()
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, ErrInternalError
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrToManyRequests
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnknown
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r struct {
		Order   big.Int
		Status  string
		Accrual decimal.Decimal
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &service.Result{
		Amount: r.Accrual,
		Status: r.Status,
	}, nil
}
