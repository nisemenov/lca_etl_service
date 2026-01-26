// Package etl coordinates the end-to-end ETL workflow,
// orchestrating producers, repositories, and workers.
package etl

import (
	"context"
	"log/slog"

	"github.com/nisemenov/etl_service/internal/consumer"
	"github.com/nisemenov/etl_service/internal/producer"
	"github.com/nisemenov/etl_service/internal/repository"
)


type etlPipline[D any, ID comparable] struct {
	producer Producer[D, ID]
	repo     Repository[D, ID]
	loader   Consumer[D]
	logger   *slog.Logger
}

func (p *etlPipline[D, ID]) FetchAndSave(ctx context.Context) error {
	return nil
}

func (etl *etlPipline[D,ID]) ProcessAndSend(ctx context.Context) error {
	return nil
}

func (etl *etlPipline[D,ID]) Acknowledge(ctx context.Context) error {
	return nil
}

func (etl *etlPipline) Run(ctx context.Context) error {
	if err := etl.FetchAndSave(ctx); err != nil {
		return err
	}

	if err := etl.ProcessAndSend(ctx); err != nil {
		return err
	}

	if err := etl.Acknowledge(ctx); err != nil {
		return err
	}

	return nil
}

func NewETLPipline(
	producer producer.PaymentProducer,
	repo repository.PaymentRepository,
	loader consumer.HTTPClickHouse,
	logger *slog.Logger,
) *etlPipline {
	return &etlPipline{producer: producer, repo: repo, loader: loader, logger: logger}
}
