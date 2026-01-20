// Package domain defines core business entities and types
// used across the ETL service.
package domain

import (
	"time"

	"github.com/nisemenov/etl_service/internal/validation"
)

type PaymentID int64

type PaymentStatus string

const (
	StatusNew        PaymentStatus = "new"
	StatusProcessing PaymentStatus = "processing"
	StatusExported   PaymentStatus = "exported"
	StatusFailed     PaymentStatus = "failed"
)

type Payment struct {
	// from payment_yookassa
	ID                    PaymentID `validate:"required"`
	CaseID                int64     `validate:"required"`
	DebtorID              int64     `validate:"required"`
	FullName              string    `validate:"required"`
	CreditNumber          string    `validate:"required"`
	CreditIssueDate       time.Time `validate:"required"`
	Amount                float64   `validate:"required"`
	DebtAmount            float64   `validate:"required"`
	ExecutionDateBySystem time.Time `validate:"required"`
	Channel               string    `validate:"required"`

	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Payment) Validate() error {
	return validation.Validate.Struct(p)
}
