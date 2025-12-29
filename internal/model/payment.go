// Package model описывает основные модели
package model

import "time"

type PaymentID int64

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusExported  PaymentStatus = "expored"
	StatusCompleted PaymentStatus = "completed"
)

type Payment struct {
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
