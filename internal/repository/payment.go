// Package repository contains persistence interfaces and
// database-backed implementations for payment storage.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nisemenov/etl_service/internal/domain"
)

type PaymentRepository interface {
	SaveBatch(ctx context.Context, payments []domain.Payment) error
	FetchForProcessing(ctx context.Context, limit int) ([]domain.Payment, error)
	MarkSent(ctx context.Context, ids []domain.PaymentID) error
}

type SQLitePaymentRepo struct {
	db *sql.DB
}

func NewSQLitePaymentRepo(db *sql.DB) *SQLitePaymentRepo {
	return &SQLitePaymentRepo{db: db}
}

func (r *SQLitePaymentRepo) SaveBatch(
	ctx context.Context,
	payments []domain.Payment,
) error {
	if len(payments) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO payments (
            id,
			case_id,
			debtor_id,
			full_name,
			credit_number,
			credit_issue_date,
			amount,
			debt_amount,
			execution_date_by_system,
			channel,
			status
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            status = excluded.status,
            updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range payments {
		_, err := stmt.ExecContext(ctx,
			p.ID,
			p.CaseID,
			p.DebtorID,
			p.FullName,
			p.CreditNumber,
			p.CreditIssueDate,
			p.Amount,
			p.DebtAmount,
			p.ExecutionDateBySystem,
			p.Channel,
			domain.StatusNew,
		)
		if err != nil {
			return fmt.Errorf("insert payment %d: %w", p.ID, err)
		}
	}

	return tx.Commit()
}

func (r *SQLitePaymentRepo) FetchForProcessing(
	ctx context.Context,
	limit int,
) ([]domain.Payment, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT id
        FROM payments
        WHERE status = ?
        LIMIT ?
    `, domain.StatusNew, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []domain.PaymentID
	for rows.Next() {
		var id domain.PaymentID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, nil
	}

	if err := r.markStatusTx(ctx, tx, ids, domain.StatusProcessing); err != nil {
		return nil, err
	}

	payments := make([]domain.Payment, 0, len(ids))
	for _, id := range ids {
		payments = append(payments, domain.Payment{
			ID:     id,
			Status: domain.StatusProcessing,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *SQLitePaymentRepo) MarkSent(
	ctx context.Context,
	ids []domain.PaymentID,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := r.markStatusTx(ctx, tx, ids, domain.StatusExported); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLitePaymentRepo) markStatusTx(
	ctx context.Context,
	tx *sql.Tx,
	ids []domain.PaymentID,
	status domain.PaymentStatus,
) error {
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for _, id := range ids {
		args = append(args, id)
	}

	query := fmt.Sprintf(`
        UPDATE payments
        SET status = ?
        WHERE id IN (%s)
    `, placeholders)

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLitePaymentRepo) fetchNewPayments(ctx context.Context, limit int) ([]domain.Payment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			id,
			case_id,
			debtor_id,
			full_name,
			credit_number,
            credit_issue_date,
			amount,
			debt_amount,
			execution_date_by_system,
            channel,
			status,
			created_at,
			updated_at
        FROM payments
        WHERE status = ?
        ORDER BY created_at ASC
        LIMIT ?
	`, domain.StatusNew, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPayments(rows)
}

func scanPayments(rows *sql.Rows) ([]domain.Payment, error) {
	var payments []domain.Payment
	for rows.Next() {
		var p domain.Payment
		err := rows.Scan(
			&p.ID,
			&p.CaseID,
			&p.DebtorID,
			&p.FullName,
			&p.CreditNumber,
			&p.CreditIssueDate,
			&p.Amount,
			&p.DebtAmount,
			&p.ExecutionDateBySystem,
			&p.Channel,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, rows.Err()
}
