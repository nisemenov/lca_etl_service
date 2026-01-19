// package main
//
// import (
// 	"database/sql"
//
// 	sqlite "github.com/nisemenov/etl_service/internal/db"
// 	"github.com/nisemenov/etl_service/internal/repository"
// )
//
// func migr() error {
// 	db, err := sql.Open("sqlite3", dsn)
// 	if err != nil {
// 		return err
// 	}
//
// 	if err := sqlite.Migrate(db); err != nil {
// 		return err
// 	}
//
// 	// repo := repository.NewSQLiteRepository(db)
// }
