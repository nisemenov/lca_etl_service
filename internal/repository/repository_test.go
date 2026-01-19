package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/stretchr/testify/require"
)

func NewTestSQLiteRepo(t *testing.T) PaymentRepository {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
        PRAGMA journal_mode = WAL;
        CREATE TABLE payments (
            id INTEGER PRIMARY KEY,
            status TEXT
        );
    `)
	require.NoError(t, err)

	repo := NewSQLiteRepository(db)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return repo
}

func TestRepository_SaveBatch(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	err := repo.SaveBatch(ctx, []domain.Payment{{ID: 1, Status: domain.StatusNew}})
	require.NoError(t, err)
}

func TestRepository_FetchForProcessing(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1, Status: domain.StatusNew}})

	payments, err := repo.FetchForProcessing(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.Equal(t, domain.StatusProcessing, payments[0].Status)
}

func TestRepository_MarkSent(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1, Status: domain.StatusNew}})
	repo.MarkSent(ctx, []domain.PaymentID{1})

	payments, _ := repo.FetchForProcessing(ctx, 10)
	require.Len(t, payments, 0)
}
