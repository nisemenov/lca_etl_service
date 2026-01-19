package consumer

import (
	"context"

	"github.com/nisemenov/etl_service/internal/domain"
)

type ClickHouseLoader interface {
	Load(ctx context.Context, payments domain.Payment) error
	Insert(ctx context.Context, payment domain.Payment) error
}
