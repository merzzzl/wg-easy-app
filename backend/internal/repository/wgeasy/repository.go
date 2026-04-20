package wgeasy

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Repository struct {
	baseURL    *url.URL
	username   string
	password   string
	httpClient *http.Client
}

var (
	ErrUnexpectedStatus = errors.New("wgeasy: unexpected status")
	ErrDefaultTransport = errors.New("wgeasy: default transport is not *http.Transport")
)

func New(baseURL, username, password string, insecureTLS bool) (*Repository, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	baseTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, ErrDefaultTransport
	}

	transport := baseTransport.Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureTLS}

	return &Repository{
		baseURL:  parsed,
		username: username,
		password: password,
		httpClient: &http.Client{
			Transport: transport,
		},
	}, nil
}

func (r *Repository) doJSON(ctx context.Context, method, endpoint string, requestBody, responseBody any) error {
	body, err := r.doRaw(ctx, method, endpoint, requestBody)
	if err != nil {
		return err
	}

	if responseBody == nil || len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, responseBody); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (r *Repository) doRaw(ctx context.Context, method, endpoint string, requestBody any) ([]byte, error) {
	var bodyReader io.Reader

	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}

		bodyReader = strings.NewReader(string(payload))
	}

	reqURL := *r.baseURL
	reqURL.Path = path.Join(r.baseURL.Path, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(r.username, r.password)

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%w: status %d: %s", ErrUnexpectedStatus, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return body, nil
}
