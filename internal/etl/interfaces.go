package etl

import (
	"context"

	"github.com/nisemenov/etl_service/internal/domain"
)

type Producer[D any, ID comparable] interface {
	Fetch(ctx context.Context) ([]D, error)
	Ack(ctx context.Context, ids []ID) error
}

type Repository[D any, ID comparable] interface {
	SaveBatch(ctx context.Context, data []D) error
	FetchForProcessing(ctx context.Context, limit int) ([]D, error)
	MarkSent(ctx context.Context, ids []ID) error
}

type Consumer[D any] interface {
	InsertBatch(ctx context.Context, data []D) error
}
