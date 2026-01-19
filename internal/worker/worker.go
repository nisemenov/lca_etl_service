// Package worker contains background workers responsible
// for processing and exporting payments.
package worker

import (
	"context"

	"github.com/nisemenov/etl_service/internal/consumer"
	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/nisemenov/etl_service/internal/repository"
)

func worker(
	ctx context.Context,
	jobs <-chan domain.Payment,
	repo repository.PaymentRepository,
	loader consumer.ClickHouseLoader,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case p, ok := <-jobs:
			if !ok {
				return
			}

			if err := loader.Insert(ctx, p); err != nil {
				// лог + retry позже
				continue
			}

			_ = repo.MarkSent(ctx, []int64{int64(p.ID)})
		}
	}
}
