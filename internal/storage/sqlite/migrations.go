package sqlite

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	// Если очень хочется видеть, какие миграции применились:
	// goose.SetVerbose(true)

	return goose.Up(db, "migrations")
}
