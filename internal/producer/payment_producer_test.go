package producer

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestPaymentProducer_FetchPayments_OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": [
        		{
            		"id": 1,
            		"case_id": 123,
            		"debtor_id": 123,
            		"full_name": "Петров Петр Петрович",
            		"credit_number": "XYZ789",
            		"credit_issue_date": "2023-05-05T00:00:00+00:00",
            		"amount": 100.00,
            		"debt_amount": 100.00,
            		"execution_date_by_system": "2024-07-01T12:00:00Z",
            		"channel": "email"
				}
			]
		}`))
	}))
	defer server.Close()

	payProducer := getPayProducer(server)

	payments, err := payProducer.FetchPayments(context.Background())

	require.NoError(t, err)
	require.Len(t, payments, 1)

	first := payments[0]
	require.Equal(t, domain.PaymentID(1), first.ID)
	require.Equal(t, "Петров Петр Петрович", first.FullName)
	require.Equal(t, domain.Money(10000), first.Amount)
}

func TestPaymentProducer_FetchPayments_SkipsInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"data": [
				{"id": 1},
        		{
            		"id": 2,
            		"case_id": 123,
            		"debtor_id": 123,
            		"full_name": "Петров Петр Петрович",
            		"credit_number": "XYZ789",
            		"credit_issue_date": "2023-05-05T00:00:00+00:00",
            		"amount": 100.00,
            		"debt_amount": 100.00,
            		"execution_date_by_system": "2024-07-01T12:00:00Z",
            		"channel": "email"
				},
				{"id": 3}
			]
		}`))
	}))
	defer server.Close()

	payProducer := getPayProducer(server)

	payments, err := payProducer.FetchPayments(context.Background())

	require.NoError(t, err)
	require.Len(t, payments, 1)

	first := payments[0]
	require.Equal(t, domain.PaymentID(2), first.ID)
}

func getPayProducer(server *httptest.Server) PaymentProducer {
	client := &http.Client{}
	prod := NewHTTPProducer(client, server.URL)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return NewPaymentProducer(prod, logger)
}
