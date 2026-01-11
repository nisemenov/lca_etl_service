// Package usecase описывает основные контракты по работе с Payment
package usecase

import (
	"context"

	"github.com/nisemenov/lca_etl_service/internal/model"
)

type PaymentFetcher interface {
	Fetch(ctx context.Context) ([]model.Payment, error)
}

type PaymentRepository interface {
	SaveBatch(ctx context.Context, payments []model.Payment) error
}

type FetchPaymentUseCase struct {
	fetcher PaymentFetcher
	repo    PaymentRepository
}
