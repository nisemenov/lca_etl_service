// Package domain defines core business entities and types
// used across the ETL service.
package domain

import (
	"math"
	"time"

	"github.com/nisemenov/etl_service/internal/validation"
)

const (
	StatusNew        PaymentStatus = "new"
	StatusProcessing PaymentStatus = "processing"
	StatusExported   PaymentStatus = "exported"
	StatusFailed     PaymentStatus = "failed"
)

type PaymentID int64

type PaymentStatus string

// Money type as float == int * 100
type Money int64

func FloatToMoney(f float64) Money {
	return Money(math.Round(f * 100))
}

func (m Money) Float64() float64 {
	return float64(m) / 100
}

type Payment struct {
	ID                    PaymentID `validate:"required"`
	CaseID                int64     `validate:"required"`
	DebtorID              int64     `validate:"required"`
	FullName              string    `validate:"required"`
	CreditNumber          string    `validate:"required"`
	CreditIssueDate       time.Time `validate:"required"`
	Amount                Money     `validate:"required"`
	DebtAmount            Money     `validate:"required"`
	ExecutionDateBySystem time.Time `validate:"required"`
	Channel               string    `validate:"required"`

	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Payment) Validate() error {
	return validation.Validate.Struct(p)
}
