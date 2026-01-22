// Package producer provides clients for fetching data
// from external producer services over HTTP.
package producer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPProducer struct {
	client  *http.Client
	baseURL string
}

func NewHTTPProducer(client *http.Client, baseURL string) *HTTPProducer {
	return &HTTPProducer{
		client:  client,
		baseURL: baseURL,
	}
}

func (p *HTTPProducer) Get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("http %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (p *HTTPProducer) Post(ctx context.Context, path string, body any) error {
	b, _ := json.Marshal(body)
	return p.PostRaw(ctx, path, "application/json", bytes.NewReader(b))
}

func (p *HTTPProducer) PostRaw(ctx context.Context, path string, contentType string, body io.Reader) error {
	req, _ := http.NewRequestWithContext(ctx, "POST", p.baseURL+path, body)
	req.Header.Set("Content-Type", contentType)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
