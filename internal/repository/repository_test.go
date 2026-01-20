package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/nisemenov/etl_service/internal/storage/sqlite"
	"github.com/stretchr/testify/require"
)

func NewTestSQLiteRepo(t *testing.T) *SQLiteRepository {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	require.NoError(t, err)

	err = sqlite.Migrate(db)
	require.NoError(t, err, "failed to apply migrations in test db")

	repo := NewSQLiteRepository(db)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return repo
}

func TestRepository_SaveBatch(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	err := repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})
	require.NoError(t, err)

	payments, err := repo.fetchNewPayments(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.Equal(t, payments[0].ID, domain.PaymentID(1))
	require.Equal(t, payments[0].Status, domain.StatusNew)
}

func TestRepository_FetchForProcessing(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})

	payments, err := repo.FetchForProcessing(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.Equal(t, payments[0].Status, domain.StatusProcessing)
}

func TestRepository_MarkSent(t *testing.T) {
	ctx := context.Background()
	repo := NewTestSQLiteRepo(t)

	repo.SaveBatch(ctx, []domain.Payment{{ID: 1}})
	repo.MarkSent(ctx, []domain.PaymentID{1})

	payments, err := repo.fetchNewPayments(ctx, 10)
	require.NoError(t, err)
	require.Len(t, payments, 0)
}
