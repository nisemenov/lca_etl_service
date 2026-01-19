package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nisemenov/etl_service/internal/domain"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) SaveBatch(
	ctx context.Context,
	payments []domain.Payment,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO payments (id, status)
        VALUES (?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range payments {
		_, err := stmt.ExecContext(ctx, p.ID, p.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) FetchForProcessing(
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

func (r *SQLiteRepository) MarkSent(
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

func (r *SQLiteRepository) markStatusTx(
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
