package main

import (
	"net/http"
	"time"

	"github.com/nisemenov/etl_service/internal/config"
	"github.com/nisemenov/etl_service/internal/producer"
	"github.com/nisemenov/etl_service/internal/repository"
	"github.com/nisemenov/etl_service/internal/storage/sqlite"
)

func main() {
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
