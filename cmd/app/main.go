package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/nisemenov/etl_service/internal/config"
	"github.com/nisemenov/etl_service/internal/producer"
	"github.com/nisemenov/etl_service/internal/repository"
	"github.com/nisemenov/etl_service/internal/storage/sqlite"
)

var logger *slog.Logger

func InitLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, nil)
	logger = slog.New(handler)
	return logger
}

func main() {
	logger := InitLogger()
	logger.Info("service started")

	cfg := config.Load()

	db := sqlite.OpenSQLite(cfg.DBPath)
	sqlite.Migrate(db)
	defer db.Close()

	httpClient := &http.Client{Timeout: 15 * time.Second}

	paymentsProducer := producer.NewHTTPProducer(
		httpClient,
		cfg.APIBaseURL,
		config.FetchPaymentsForCH,
	)

	paymentsRepo := repository.NewSQLitePaymentRepo(db)

	// clickhouse := consumer.NewClickHouseClient(cfg.ClickhouseDSN)

	// etlService := etl.NewService(
	//     paymentsProducer,
	//     paymentsRepo,
	//     clickhouse,
	// )

	// worker := worker.New(etlService)
	//
	// worker.Run(context.Background())
}
