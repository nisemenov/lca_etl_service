// Package producer provides clients for fetching data
// from external producer services over HTTP.
package producer

import (
	"context"

	"github.com/nisemenov/etl_service/internal/domain"
)

type Producer interface {
	FetchPayments(ctx context.Context) ([]domain.Payment, error)
	AckPayments(ctx context.Context, ids []int64) error
}
