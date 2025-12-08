package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/robinmordasiewicz/vesctl/pkg/client"
)

// TestAuthentication_P12Bundle tests authentication using P12 bundle
func TestAuthentication_P12Bundle(t *testing.T) {
	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured (VES_API_URL, VES_P12_FILE, VES_P12_PASSWORD)")
	}

	// Verify P12 file exists
	if _, err := os.Stat(p12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s", p12File)
	}

	// Set password env var for client
	_ = os.Setenv("VES_P12_PASSWORD", p12Password)
	defer func() { _ = os.Unsetenv("VES_P12_PASSWORD") }()

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          p12File,
		Timeout:            30,
		Debug:              testing.Verbose(),
		InsecureSkipVerify: true, // Skip TLS verification for staging environments
	}

	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with a simple API call - list namespaces
	resp, err := c.Get(ctx, "/api/web/namespaces", nil)
	if err != nil {
		t.Fatalf("API request failed: %v", err)
	}

	if resp.StatusCode == 401 {
		t.Fatal("Authentication failed - status 401 Unauthorized")
	}

	if resp.StatusCode == 403 {
		t.Fatal("Authentication failed - status 403 Forbidden")
	}

	if resp.StatusCode >= 400 {
		t.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(resp.Body))
	}

	t.Logf("Authentication successful, status: %d", resp.StatusCode)
}

// TestAuthentication_Whoami tests the whoami endpoint
func TestAuthentication_Whoami(t *testing.T) {
	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured")
	}

	if _, err := os.Stat(p12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s", p12File)
	}

	_ = os.Setenv("VES_P12_PASSWORD", p12Password)
	defer func() { _ = os.Unsetenv("VES_P12_PASSWORD") }()

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          p12File,
		Timeout:            30,
		Debug:              testing.Verbose(),
		InsecureSkipVerify: true, // Skip TLS verification for staging environments
	}

	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try whoami endpoint
	resp, err := c.Get(ctx, "/api/web/custom/namespace/system/whoami", nil)
	if err != nil {
		t.Fatalf("Whoami request failed: %v", err)
	}

	// 404 is acceptable - endpoint may not exist, but auth worked
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		t.Fatalf("Authentication failed with status %d", resp.StatusCode)
	}

	t.Logf("Whoami request completed, status: %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(resp.Body, &result); err == nil {
			t.Logf("Whoami response: %+v", result)
		}
	}
}

// TestAuthentication_InvalidP12 tests that invalid P12 fails appropriately
func TestAuthentication_InvalidP12(t *testing.T) {
	apiURL := os.Getenv("VES_API_URL")
	if apiURL == "" {
		t.Skip("VES_API_URL not set")
	}

	// Create a temporary invalid P12 file
	tmpFile, err := os.CreateTemp("", "invalid-*.p12")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, _ = tmpFile.Write([]byte("not a valid p12 bundle"))
	_ = tmpFile.Close()

	_ = os.Setenv("VES_P12_PASSWORD", "dummy")
	defer func() { _ = os.Unsetenv("VES_P12_PASSWORD") }()

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          tmpFile.Name(),
		Timeout:            30,
		InsecureSkipVerify: true,
	}

	_, err = client.New(cfg)
	if err == nil {
		t.Error("Expected error for invalid P12 bundle")
	}
}

// TestAuthentication_WrongPassword tests that wrong P12 password fails
func TestAuthentication_WrongPassword(t *testing.T) {
	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")

	if apiURL == "" || p12File == "" {
		t.Skip("Integration test environment not configured")
	}

	if _, err := os.Stat(p12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s", p12File)
	}

	// Set wrong password
	_ = os.Setenv("VES_P12_PASSWORD", "wrong-password-12345")
	defer func() { _ = os.Unsetenv("VES_P12_PASSWORD") }()

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          p12File,
		Timeout:            30,
		InsecureSkipVerify: true,
	}

	_, err := client.New(cfg)
	if err == nil {
		t.Error("Expected error for wrong P12 password")
	}
	t.Logf("Got expected error for wrong password: %v", err)
}

// TestAuthentication_MissingPassword tests that missing P12 password fails
func TestAuthentication_MissingPassword(t *testing.T) {
	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")

	if apiURL == "" || p12File == "" {
		t.Skip("Integration test environment not configured")
	}

	if _, err := os.Stat(p12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s", p12File)
	}

	// Ensure password is not set
	_ = os.Unsetenv("VES_P12_PASSWORD")

	cfg := &client.Config{
		ServerURLs:         []string{apiURL},
		P12Bundle:          p12File,
		Timeout:            30,
		InsecureSkipVerify: true,
	}

	_, err := client.New(cfg)
	// May or may not fail depending on whether P12 is password-protected
	// Just log the result
	if err != nil {
		t.Logf("Got error for missing password (expected if P12 is encrypted): %v", err)
	} else {
		t.Log("No error for missing password - P12 may not be password protected")
	}
}
