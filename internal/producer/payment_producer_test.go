package producer

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/nisemenov/etl_service/internal/httpclient"
	"github.com/stretchr/testify/require"
)

func TestPaymentProducer_Fetch_OK(t *testing.T) {
	var mthd string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mthd = r.Method

		require.Equal(t, FetchPaymentsPath, r.URL.Path)

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
	payments, err := payProducer.Fetch(context.Background())

	require.NoError(t, err)
	require.Equal(t, "GET", mthd)
	require.Len(t, payments, 1)

	first := payments[0]
	require.Equal(t, domain.PaymentID(1), first.ID)
	require.Equal(t, "Петров Петр Петрович", first.FullName)
	require.Equal(t, domain.Money(10000), first.Amount)
}

func TestPaymentProducer_Fetch_SkipsInvalid(t *testing.T) {
	var mthd string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mthd = r.Method

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
	payments, err := payProducer.Fetch(context.Background())

	require.NoError(t, err)
	require.Equal(t, "GET", mthd)
	require.Len(t, payments, 1)
	require.Equal(t, domain.PaymentID(2), payments[0].ID)
}

func TestPaymentProducer_Ack_OK(t *testing.T) {
	var mthd string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mthd = r.Method

		require.Equal(t, AckPaymentsPath, r.URL.Path)

		var body struct {
			IDs []int64 `json:"ids"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, []int64{1, 2}, body.IDs)
	}))
	defer server.Close()

	payProducer := getPayProducer(server)
	err := payProducer.Ack(context.Background(), []domain.PaymentID{domain.PaymentID(1), domain.PaymentID(2)})

	require.NoError(t, err)
	require.Equal(t, "POST", mthd)
}

func TestPaymentProducer_Ack_HTTPError(t *testing.T) {
	var mthd string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mthd = r.Method

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"db down"}`))
	}))
	defer server.Close()

	payProducer := getPayProducer(server)
	err := payProducer.Ack(context.Background(), []domain.PaymentID{1, 2})

	require.Error(t, err)
	require.Equal(t, "POST", mthd)
}

func TestPaymentProducer_Ack_Empty(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	payProducer := NewPaymentProducer(nil, logger)

	err := payProducer.Ack(context.Background(), nil)
	require.NoError(t, err)
}

func getPayProducer(server *httptest.Server) paymentProducer {
	prod := httpclient.NewHTTPClient(&http.Client{}, server.URL)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return *NewPaymentProducer(prod, logger)
}
