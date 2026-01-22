// Package producer provides clients for fetching data
// from external producer services over HTTP.
package producer

import (
	"context"
	"errors"
	"log/slog"
	"math"

	"github.com/nisemenov/etl_service/internal/config"
	"github.com/nisemenov/etl_service/internal/domain"
)

type PaymentProducer interface {
	FetchPayments(ctx context.Context) ([]domain.Payment, error)
	AckPayments(ctx context.Context, ids []domain.PaymentID) error
}

type paymentProducer struct {
	http   *HTTPProducer
	logger *slog.Logger
}

func NewPaymentProducer(http *HTTPProducer, logger *slog.Logger) PaymentProducer {
	return &paymentProducer{http: http, logger: logger}
}

func (p *paymentProducer) FetchPayments(ctx context.Context) ([]domain.Payment, error) {
	var resp fetchPaymentsResponse
	err := p.http.Get(ctx, config.FetchPaymentsPath, &resp)

	// fill output in with validated data
	out := make([]domain.Payment, 0, len(resp.Data))
	for _, dto := range resp.Data {
		if err := dto.Validate(); err != nil {
			p.logger.Warn(err.Error(), "invalid paymentDTO with ID", dto.ID)
			continue
		}

		amount := domain.Money(math.Round(dto.Amount * 100))
		debt := domain.Money(math.Round(dto.DebtAmount * 100))

		out = append(out, domain.Payment{
			ID:                    dto.ID,
			CaseID:                dto.CaseID,
			DebtorID:              dto.DebtorID,
			FullName:              dto.FullName,
			CreditNumber:          dto.CreditNumber,
			CreditIssueDate:       dto.CreditIssueDate,
			Amount:                amount,
			DebtAmount:            debt,
			ExecutionDateBySystem: dto.ExecutionDateBySystem,
			Channel:               dto.Channel,
		})
	}
	if len(out) == 0 {
		return nil, errors.New("all payments invalid")
	}

	return out, err
}

func (p *paymentProducer) AckPayments(ctx context.Context, ids []domain.PaymentID) error {
	if len(ids) == 0 {
		p.logger.Warn("empty ids batch for AckPayments")
		return nil
	}
	payload := struct {
		IDs []domain.PaymentID `json:"ids"`
	}{
		IDs: ids,
	}
	return p.http.Post(ctx, config.AckPaymentsPath, payload)
}
