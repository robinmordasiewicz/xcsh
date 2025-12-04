package testutil

import (
	"os"
	"testing"
)

// TestConfig holds test configuration from environment variables
type TestConfig struct {
	APIURL      string
	P12File     string
	P12Password string
}

// LoadTestConfig loads test configuration from environment variables
func LoadTestConfig(t *testing.T) *TestConfig {
	cfg := &TestConfig{
		APIURL:      os.Getenv("F5XC_API_URL"),
		P12File:     os.Getenv("F5XC_API_P12_FILE"),
		P12Password: os.Getenv("F5XC_P12_PASSWORD"),
	}

	return cfg
}

// RequireIntegrationEnv skips the test if integration test environment is not configured
func RequireIntegrationEnv(t *testing.T) *TestConfig {
	cfg := LoadTestConfig(t)

	if cfg.APIURL == "" {
		t.Skip("F5XC_API_URL not set, skipping integration test")
	}
	if cfg.P12File == "" {
		t.Skip("F5XC_API_P12_FILE not set, skipping integration test")
	}
	if cfg.P12Password == "" {
		t.Skip("F5XC_P12_PASSWORD not set, skipping integration test")
	}

	// Verify P12 file exists
	if _, err := os.Stat(cfg.P12File); os.IsNotExist(err) {
		t.Skipf("P12 file not found at %s, skipping integration test", cfg.P12File)
	}

	return cfg
}

// SetupTestEnv sets up environment for integration tests
func SetupTestEnv(cfg *TestConfig) func() {
	// Store original values
	origURL := os.Getenv("F5XC_API_URL")
	origP12 := os.Getenv("F5XC_API_P12_FILE")
	origPass := os.Getenv("VES_P12_PASSWORD")

	// Set test values
	_ = os.Setenv("F5XC_API_URL", cfg.APIURL)
	_ = os.Setenv("F5XC_API_P12_FILE", cfg.P12File)
	_ = os.Setenv("VES_P12_PASSWORD", cfg.P12Password)

	// Return cleanup function
	return func() {
		if origURL != "" {
			_ = os.Setenv("F5XC_API_URL", origURL)
		} else {
			_ = os.Unsetenv("F5XC_API_URL")
		}
		if origP12 != "" {
			_ = os.Setenv("F5XC_API_P12_FILE", origP12)
		} else {
			_ = os.Unsetenv("F5XC_API_P12_FILE")
		}
		if origPass != "" {
			_ = os.Setenv("VES_P12_PASSWORD", origPass)
		} else {
			_ = os.Unsetenv("VES_P12_PASSWORD")
		}
	}
}

// UniqueResourceName generates a unique name for test resources
func UniqueResourceName(prefix string) string {
	return prefix + "-test-" + randomSuffix()
}

func randomSuffix() string {
	// Use timestamp-based suffix for uniqueness
	return os.Getenv("TEST_RUN_ID")
}
