package repository

import (
	"database/sql"
	"testing"

	"github.com/nisemenov/etl_service/internal/storage/sqlite"
	"github.com/stretchr/testify/require"
)

func NewTestSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	require.NoError(t, err)

	err = sqlite.Migrate(db)
	require.NoError(t, err)

	t.Cleanup(func() { _ = db.Close() })
	return db
}

func NewTestSQLitePaymentRepo(t *testing.T) *SQLitePaymentRepo {
	db := NewTestSQLiteDB(t)
	repo := NewSQLitePaymentRepo(db)
	return repo
}
