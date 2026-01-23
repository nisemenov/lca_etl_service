package httpclient

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPClient_WithHeaders(t *testing.T) {
	var header string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header = r.Header.Get("testHeader")
	}))
	defer server.Close()

	p := NewHTTPClient(&http.Client{}, server.URL, WithHeaders(map[string]string{"testHeader": "test_header"}))

	var out any
	err := p.Get(context.Background(), testURL, &out)

	require.Error(t, err)
	require.Equal(t, "test_header", header)
}

func TestHTTPClient_BasicAuth(t *testing.T) {
	var auth bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, ok := r.BasicAuth()
		auth = ok
	}))
	defer server.Close()

	p := NewHTTPClient(
		&http.Client{},
		server.URL,
		WithMiddleware(func(r *http.Request) { r.SetBasicAuth("username", "pass") }),
	)

	err := p.PostRaw(context.Background(), testURL, "application/json", bytes.NewReader(nil))

	require.NoError(t, err)
	require.True(t, auth)
}
