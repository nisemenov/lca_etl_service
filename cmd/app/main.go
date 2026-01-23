package main

import (
	"log/slog"
	"os"

	"github.com/nisemenov/etl_service/internal/config"
	"github.com/nisemenov/etl_service/internal/storage/sqlite"
)

func main() {
	level := slog.LevelInfo
	cfg := config.Load()

	if cfg.Debug {
		level = slog.LevelDebug
	}

	// init Logger
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.Debug, // показываем source только при cnf.Debug == true
	})
	logger := slog.New(handler)
	slog.SetDefault(logger) // опционально, если хочешь использовать slog.Info() без переменной

	logger.Info("starting application",
		"version", "v0.0.1",
		"debug", cfg.Debug,
	)

	db := sqlite.OpenSQLite(cfg.DBPath)
	sqlite.Migrate(db)
	defer db.Close()

	// worker := worker.New(etlService)
	//
	// worker.Run(context.Background())
}
