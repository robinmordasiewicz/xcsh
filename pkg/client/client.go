package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the client configuration
type Config struct {
	ServerURLs         []string
	Cert               string
	Key                string
	CACert             string
	P12Bundle          string
	Debug              bool
	Timeout            int
	InsecureSkipVerify bool // Skip TLS certificate verification (for staging/testing)
}

// Client is the HTTP client for F5 XC API
type Client struct {
	httpClient *http.Client
	serverURLs []string
	debug      bool
	tlsConfig  *tls.Config
}

// New creates a new API client
func New(cfg *Config) (*Client, error) {
	client := &Client{
		serverURLs: cfg.ServerURLs,
		debug:      cfg.Debug,
	}

	// Create TLS config if credentials are provided
	if cfg.P12Bundle != "" || (cfg.Cert != "" && cfg.Key != "") || cfg.CACert != "" {
		tlsConfig, err := createTLSConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		client.tlsConfig = tlsConfig
	}

	// Create HTTP transport with TLS
	transport := &http.Transport{
		TLSClientConfig: client.tlsConfig,
	}

	client.httpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}

	return client, nil
}

// Request represents an API request
type Request struct {
	Method   string
	Path     string
	Body     interface{}
	Headers  map[string]string
	Query    url.Values
}

// Response represents an API response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Do executes an HTTP request
func (c *Client) Do(ctx context.Context, req *Request) (*Response, error) {
	// Build URL
	baseURL := c.serverURLs[0] // Use first server URL
	baseURL = strings.TrimRight(baseURL, "/")

	// Handle case where base URL ends with /api and path starts with /api
	// This allows server URLs like "https://tenant.console.ves.volterra.io/api"
	// to work correctly with paths like "/api/web/namespaces"
	if strings.HasSuffix(baseURL, "/api") && strings.HasPrefix(req.Path, "/api") {
		baseURL = strings.TrimSuffix(baseURL, "/api")
	}

	fullURL := fmt.Sprintf("%s%s", baseURL, req.Path)

	// Add query parameters
	if len(req.Query) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, req.Query.Encode())
	}

	// Prepare body
	var bodyReader io.Reader
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Debug logging
	if c.debug {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", req.Method, fullURL)
		if len(bodyBytes) > 0 {
			fmt.Fprintf(os.Stderr, "DEBUG: Request body: %s\n", string(bodyBytes))
		}
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug response logging
	if c.debug {
		fmt.Fprintf(os.Stderr, "DEBUG: Response status: %d\n", resp.StatusCode)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, query url.Values) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	})
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: http.MethodDelete,
		Path:   path,
	})
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, &Request{
		Method: http.MethodPatch,
		Path:   path,
		Body:   body,
	})
}

// FormatOutput formats the response body according to the specified format
func FormatOutput(data []byte, format string) (string, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		// If not valid JSON, return as-is
		return string(data), nil
	}

	switch strings.ToLower(format) {
	case "json":
		formatted, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to format as JSON: %w", err)
		}
		return string(formatted), nil
	case "yaml", "":
		formatted, err := yaml.Marshal(obj)
		if err != nil {
			return "", fmt.Errorf("failed to format as YAML: %w", err)
		}
		return string(formatted), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}
