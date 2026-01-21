package repository

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestPaymentRepo_SaveBatch(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLitePaymentRepo(t)

	err := repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})
	require.NoError(t, err)

	payments, err := repo.fetchNewPayments(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.Equal(t, payments[0].ID, domain.PaymentID(1))
	require.Equal(t, payments[0].Status, domain.StatusNew)
}

func TestPaymentRepo_FetchForProcessing(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLitePaymentRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})

	payments, err := repo.FetchForProcessing(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.Equal(t, payments[0].Status, domain.StatusProcessing)
}

func TestPaymentRepo_MarkSent(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLitePaymentRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})
	repo.MarkSent(ctx, []domain.PaymentID{1})

	payments, err := repo.fetchNewPayments(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 0)
}
