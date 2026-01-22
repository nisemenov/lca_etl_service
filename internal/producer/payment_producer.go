// Package producer provides clients for fetching data
// from external producer services over HTTP.
package producer

import (
	"context"
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
	if err != nil {
		p.logger.Error(
			"failed to fetch payment data",
			"err", err,
		)
		return nil, err
	}

	// fill output in with validated data
	response := make([]domain.Payment, 0, len(resp.Data))
	for _, rawPayment := range resp.Data {
		if err := rawPayment.Validate(); err != nil {
			p.logger.Warn(
				"failed to validate raw payment data",
				"err", err,
			)
			continue
		}

		amount := domain.Money(math.Round(rawPayment.Amount * 100))
		debt := domain.Money(math.Round(rawPayment.DebtAmount * 100))

		response = append(response, domain.Payment{
			ID:                    rawPayment.ID,
			CaseID:                rawPayment.CaseID,
			DebtorID:              rawPayment.DebtorID,
			FullName:              rawPayment.FullName,
			CreditNumber:          rawPayment.CreditNumber,
			CreditIssueDate:       rawPayment.CreditIssueDate,
			Amount:                amount,
			DebtAmount:            debt,
			ExecutionDateBySystem: rawPayment.ExecutionDateBySystem,
			Channel:               rawPayment.Channel,
		})
	}
	if len(response) == 0 {
		p.logger.Error("all payments invalid")
	}

	return response, nil
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

	err := p.http.Post(ctx, config.AckPaymentsPath, payload)
	if err != nil {
		p.logger.Error(
			"failed to insert payment ids into prod service",
			"err", err,
		)
		return err
	}
	return nil
}
