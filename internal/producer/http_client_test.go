package producer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testURL = "/test"

func TestHTTPProducer_Get_OK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, testURL, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"value": 42}`))
	}))
	defer server.Close()

	client := &http.Client{}
	p := NewHTTPProducer(client, server.URL)

	var resp struct {
		Value int `json:"value"`
	}

	err := p.Get(context.Background(), testURL, &resp)
	require.NoError(t, err)
	require.Equal(t, 42, resp.Value)
}

func TestHTTPProducer_Get_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	p := NewHTTPProducer(&http.Client{}, server.URL)

	var out any
	err := p.Get(context.Background(), testURL, &out)

	require.Error(t, err)
}

func TestHTTPProducer_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, testURL, r.URL.Path)

		var body map[string]int
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		require.Equal(t, 123, body["id"])

		w.WriteHeader(200)
	}))
	defer server.Close()

	p := NewHTTPProducer(&http.Client{}, server.URL)

	err := p.Post(context.Background(), testURL, map[string]int{"id": 123})
	require.NoError(t, err)
}

func TestHTTPProducer_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond}
	p := NewHTTPProducer(client, server.URL)

	var out any
	err := p.Get(context.Background(), testURL, &out)

	require.Error(t, err)
}
