// Package httpclient provides clients for fetching data
// from external producer services over HTTP.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient struct {
	client     *http.Client
	baseURL    string
	headers    map[string]string
	middleware func(*http.Request)
}

func (p *HTTPClient) Get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.baseURL+path, nil)
	if err != nil {
		return err
	}

	p.applyHeadersAndMiddleware(req)

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

func (p *HTTPClient) Post(ctx context.Context, path string, body any) error {
	b, _ := json.Marshal(body)
	return p.PostRaw(ctx, path, "application/json", bytes.NewReader(b))
}

func (p *HTTPClient) PostRaw(ctx context.Context, path string, contentType string, body io.Reader) error {
	req, _ := http.NewRequestWithContext(ctx, "POST", p.baseURL+path, body)
	req.Header.Set("Content-Type", contentType)

	p.applyHeadersAndMiddleware(req)

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

func (p *HTTPClient) applyHeadersAndMiddleware(req *http.Request) {
	for k, v := range p.headers {
		req.Header.Set(k, v)
	}

	if p.middleware != nil {
		p.middleware(req)
	}
}

func NewHTTPClient(client *http.Client, baseURL string, opts ...Option) *HTTPClient {
	p := &HTTPClient{
		client:  client,
		baseURL: baseURL,
		headers: make(map[string]string),
	}

	// for modification HTTPClient from high levels
	for _, opt := range opts {
		opt(p)
	}
	return p
}
