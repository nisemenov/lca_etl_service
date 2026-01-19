// Package repository contains persistence interfaces and
// database-backed implementations for payment storage.
package repository

import (
	"context"

	"github.com/nisemenov/etl_service/internal/domain"
)

type PaymentRepository interface {
	SaveBatch(ctx context.Context, payments []domain.Payment) error
	FetchForProcessing(ctx context.Context, limit int) ([]domain.Payment, error)
	MarkSent(ctx context.Context, ids []domain.PaymentID) error
}
