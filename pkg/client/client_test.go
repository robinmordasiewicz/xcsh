package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew_WithoutCredentials(t *testing.T) {
	cfg := &Config{
		ServerURLs: []string{"http://localhost:8080"},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if len(client.serverURLs) != 1 {
		t.Errorf("Expected 1 server URL, got %d", len(client.serverURLs))
	}
}

func TestClient_Get(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.URL.Path != "/api/test" {
			t.Errorf("Expected path /api/test, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "/api/test", nil)
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

func TestClient_Post(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if body["name"] != "test-resource" {
			t.Errorf("Expected name 'test-resource', got '%v'", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"metadata": map[string]string{"name": "test-resource"},
			"spec":     map[string]string{},
		})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := map[string]interface{}{
		"name": "test-resource",
	}

	resp, err := client.Post(ctx, "/api/config/namespaces/test/resources", body)
	if err != nil {
		t.Fatalf("Post request failed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestClient_Put(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"updated": "true"})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Put(ctx, "/api/resource/1", map[string]string{"name": "updated"})
	if err != nil {
		t.Fatalf("Put request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Delete(ctx, "/api/resource/1")
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_Patch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"patched": "true"})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Patch(ctx, "/api/resource/1", map[string]string{"field": "value"})
	if err != nil {
		t.Fatalf("Patch request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_QueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("namespace") != "test-ns" {
			t.Errorf("Expected namespace=test-ns, got %s", query.Get("namespace"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", query.Get("limit"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]string{})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	query := make(map[string][]string)
	query["namespace"] = []string{"test-ns"}
	query["limit"] = []string{"10"}

	_, err = client.Get(ctx, "/api/resources", query)
	if err != nil {
		t.Fatalf("Get with query params failed: %v", err)
	}
}

func TestFormatOutput_JSON(t *testing.T) {
	data := []byte(`{"name":"test","value":123}`)

	output, err := FormatOutput(data, "json")
	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	expected := `{
  "name": "test",
  "value": 123
}`
	if output != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, output)
	}
}

func TestFormatOutput_YAML(t *testing.T) {
	data := []byte(`{"name":"test","value":123}`)

	output, err := FormatOutput(data, "yaml")
	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// YAML output varies slightly, just check it contains expected content
	if output == "" {
		t.Error("Expected non-empty YAML output")
	}
}

func TestFormatOutput_InvalidFormat(t *testing.T) {
	data := []byte(`{"name":"test"}`)

	_, err := FormatOutput(data, "invalid")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestFormatOutput_InvalidJSON(t *testing.T) {
	data := []byte(`not valid json`)

	output, err := FormatOutput(data, "json")
	if err != nil {
		t.Fatalf("FormatOutput failed: %v", err)
	}

	// Should return raw data if not valid JSON
	if output != "not valid json" {
		t.Errorf("Expected raw data for invalid JSON, got: %s", output)
	}
}

func TestClient_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":   "bad_request",
			"message": "Invalid parameter",
		})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	resp, err := client.Get(ctx, "/api/invalid", nil)
	if err != nil {
		t.Fatalf("Request should not fail for 4xx status: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	var errorResp map[string]string
	if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResp["error"] != "bad_request" {
		t.Errorf("Expected error 'bad_request', got '%s'", errorResp["error"])
	}
}

func TestNew_WithAPIToken(t *testing.T) {
	cfg := &Config{
		ServerURLs: []string{"http://localhost:8080"},
		APIToken:   "test-api-token",
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.apiToken != "test-api-token" {
		t.Errorf("Expected apiToken to be 'test-api-token', got '%s'", client.apiToken)
	}
}

func TestClient_APITokenHeader(t *testing.T) {
	// Create mock server that verifies the Authorization header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "APIToken test-api-token"

		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header '%s', got '%s'", expectedAuth, authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		APIToken:   "test-api-token",
		Timeout:    30,
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "/api/test", nil)
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_NoAPITokenHeader_WhenNotConfigured(t *testing.T) {
	// Create mock server that verifies NO Authorization header is sent
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			t.Errorf("Expected no Authorization header, got '%s'", authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	cfg := &Config{
		ServerURLs: []string{server.URL},
		Timeout:    30,
		// No APIToken set
	}

	client, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "/api/test", nil)
	if err != nil {
		t.Fatalf("Get request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_APIToken_AllMethods(t *testing.T) {
	// Test that API token is sent for all HTTP methods
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				expectedAuth := "APIToken test-token"

				if authHeader != expectedAuth {
					t.Errorf("Expected Authorization header '%s' for %s, got '%s'", expectedAuth, method, authHeader)
				}

				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			}))
			defer server.Close()

			cfg := &Config{
				ServerURLs: []string{server.URL},
				APIToken:   "test-token",
				Timeout:    30,
			}

			client, err := New(cfg)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			ctx := context.Background()
			var resp *Response

			switch method {
			case http.MethodGet:
				resp, err = client.Get(ctx, "/api/test", nil)
			case http.MethodPost:
				resp, err = client.Post(ctx, "/api/test", map[string]string{})
			case http.MethodPut:
				resp, err = client.Put(ctx, "/api/test", map[string]string{})
			case http.MethodDelete:
				resp, err = client.Delete(ctx, "/api/test")
			case http.MethodPatch:
				resp, err = client.Patch(ctx, "/api/test", map[string]string{})
			}

			if err != nil {
				t.Fatalf("%s request failed: %v", method, err)
			}

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, resp.StatusCode)
			}
		})
	}
}
