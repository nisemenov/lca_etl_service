// Package domain описывает основные модели
package domain

import "time"

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
	ID                    PaymentID
	CaseID                int64
	DebtorID              int64
	FullName              string
	CreditNumber          string
	CreditIssueDate       time.Time
	Amount                float64
	DebtAmount            float64
	ExecutionDateBySystem time.Time
	Channel               string

	Status PaymentStatus
}
