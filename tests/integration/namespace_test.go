package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/robinmordasiewicz/vesctl/pkg/client"
)

// getTestClient creates a client for integration tests
func getTestClient(t *testing.T) *client.Client {
	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_API_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured (VES_API_URL, VES_API_P12_FILE, VES_P12_PASSWORD)")
	}

	if _, err := os.Stat(p12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s", p12File)
	}

	_ = os.Setenv("VES_P12_PASSWORD", p12Password)

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          p12File,
		Timeout:            60,
		Debug:              testing.Verbose(),
		InsecureSkipVerify: true, // Skip TLS verification for staging environments
	}

	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return c
}

// TestNamespace_List tests listing namespaces
func TestNamespace_List(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/web/namespaces", nil)
	if err != nil {
		t.Fatalf("List namespaces failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check if we got items array
	items, ok := result["items"].([]interface{})
	if !ok {
		// Try alternative response format
		t.Logf("Response format: %+v", result)
	} else {
		t.Logf("Found %d namespaces", len(items))
		for i, item := range items {
			if i >= 5 { // Only log first 5
				t.Logf("... and %d more", len(items)-5)
				break
			}
			if ns, ok := item.(map[string]interface{}); ok {
				if name, ok := ns["name"].(string); ok {
					t.Logf("  - %s", name)
				}
			}
		}
	}
}

// TestNamespace_GetSystem tests getting the system namespace
func TestNamespace_GetSystem(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/web/namespaces/system", nil)
	if err != nil {
		t.Fatalf("Get system namespace failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify we got namespace data
	if name, ok := result["name"].(string); ok {
		if name != "system" {
			t.Errorf("Expected namespace name 'system', got '%s'", name)
		}
	}

	t.Logf("System namespace: %+v", result)
}

// TestNamespace_GetShared tests getting the shared namespace
func TestNamespace_GetShared(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/web/namespaces/shared", nil)
	if err != nil {
		t.Fatalf("Get shared namespace failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	t.Logf("Shared namespace: %+v", result)
}

// TestNamespace_GetNonexistent tests getting a nonexistent namespace
func TestNamespace_GetNonexistent(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Get(ctx, "/api/web/namespaces/nonexistent-test-namespace-12345", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Should get 404
	if resp.StatusCode != 404 {
		t.Logf("Expected 404, got %d: %s", resp.StatusCode, string(resp.Body))
	} else {
		t.Log("Got expected 404 for nonexistent namespace")
	}
}

// TestNamespace_WebAPI tests the web API for namespaces
func TestNamespace_WebAPI(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// List namespaces via web API (note: /api/config/namespaces doesn't exist,
	// config API requires namespace in path like /api/config/namespaces/{ns}/resources)
	resp, err := c.Get(ctx, "/api/web/namespaces", nil)
	if err != nil {
		t.Fatalf("List namespaces via web API failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API error %d: %s", resp.StatusCode, string(resp.Body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if items, ok := result["items"].([]interface{}); ok {
		t.Logf("Web API found %d namespaces", len(items))
	}
}

// TestNamespace_SystemNamespaceContents tests listing contents of system namespace
func TestNamespace_SystemNamespaceContents(t *testing.T) {
	c := getTestClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to list something in system namespace (e.g., sites)
	resp, err := c.Get(ctx, "/api/config/namespaces/system/sites", nil)
	if err != nil {
		t.Fatalf("List sites in system namespace failed: %v", err)
	}

	// Log the result regardless of status
	t.Logf("Sites in system namespace - status: %d", resp.StatusCode)
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			if items, ok := result["items"].([]interface{}); ok {
				t.Logf("Found %d sites in system namespace", len(items))
			}
		}
	}
}
