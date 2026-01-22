package consumer

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nisemenov/etl_service/internal/domain"
	"github.com/nisemenov/etl_service/internal/producer"
	"github.com/stretchr/testify/require"
)

func TestClickHouseLoader_InsertBatch_OK(t *testing.T) {
	called := false
	receivedBody := make([]byte, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/?query=INSERT+INTO+payments+FORMAT+JSONEachRow", r.URL.RequestURI())

		var err error
		receivedBody, err = io.ReadAll(r.Body)
		require.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{}
	prod := producer.NewHTTPProducer(client, server.URL)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	loader := NewHTTPClickHouseLoader(prod, "payments", logger)

	payments := []domain.Payment{
		{ID: 1, FullName: "Ivan", Amount: 10000},
		{ID: 2, FullName: "Petr", Amount: 20000},
	}

	err := loader.InsertBatch(context.Background(), payments)
	require.NoError(t, err)
	require.True(t, called)

	// Проверяем что реально ушло
	lines := strings.Split(strings.TrimSpace(string(receivedBody)), "\n")
	require.Len(t, lines, 2)

	var r1, r2 map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &r1))
	require.Equal(t, float64(1), r1["id"])
	require.Equal(t, "Ivan", r1["full_name"])
	require.Equal(t, float64(10000), r1["amount"])

	require.NoError(t, json.Unmarshal([]byte(lines[1]), &r2))
	require.Equal(t, float64(2), r2["id"])
	require.Equal(t, "Petr", r2["full_name"])
	require.Equal(t, float64(20000), r2["amount"])
}
