package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	configContent := `server-url: https://test.console.ves.volterra.io/api
p12-bundle: /path/to/cert.p12
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.ServerURL != "https://test.console.ves.volterra.io/api" {
		t.Errorf("Expected URL https://test.console.ves.volterra.io/api, got %s", cfg.ServerURL)
	}

	if cfg.P12Bundle != "/path/to/cert.p12" {
		t.Errorf("Expected P12Bundle /path/to/cert.p12, got %s", cfg.P12Bundle)
	}
}

func TestLoad_CertKeyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	configContent := `server-url: https://test.console.ves.volterra.io/api
cert: /path/to/cert.pem
key: /path/to/key.pem
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Cert != "/path/to/cert.pem" {
		t.Errorf("Expected Cert /path/to/cert.pem, got %s", cfg.Cert)
	}

	if cfg.Key != "/path/to/key.pem" {
		t.Errorf("Expected Key /path/to/key.pem, got %s", cfg.Key)
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if err == nil {
		t.Error("Expected error for empty path")
	}
}

func TestLoad_NonexistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/.f5xcconfig")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestConfig_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		P12Bundle: "/path/to/cert.p12",
	}

	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify file permissions (skip on Windows as it doesn't support Unix permissions)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(configPath)
		if err != nil {
			t.Fatalf("Failed to stat config file: %v", err)
		}
		if info.Mode().Perm() != 0600 {
			t.Errorf("Expected permissions 0600, got %o", info.Mode().Perm())
		}
	}

	// Reload and verify
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if loaded.ServerURL != cfg.ServerURL {
		t.Error("Saved config doesn't match original")
	}
}

func TestConfig_Validate_Valid_P12(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		P12Bundle: "/path/to/cert.p12",
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestConfig_Validate_Valid_CertKey(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		Cert:      "/path/to/cert.pem",
		Key:       "/path/to/key.pem",
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestConfig_Validate_NoServerURL(t *testing.T) {
	cfg := &Config{
		P12Bundle: "/path/to/cert.p12",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error for missing server-url")
	}
}

func TestConfig_Validate_NoAuth(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error for missing authentication")
	}
}

func TestConfig_Validate_OnlyCert(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		Cert:      "/path/to/cert.pem",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error when only cert is provided without key")
	}
}

func TestConfig_Validate_OnlyKey(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		Key:       "/path/to/key.pem",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error when only key is provided without cert")
	}
}

func TestConfig_EmptyServerURL(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	configContent := `server-url: ""
p12-bundle: /path/to/cert.p12
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.ServerURL != "" {
		t.Errorf("Expected empty server URL, got %s", cfg.ServerURL)
	}
}

func TestLoad_APITokenConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	configContent := `server-url: https://test.console.ves.volterra.io/api
api-token: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !cfg.APIToken {
		t.Error("Expected APIToken to be true")
	}

	if cfg.P12Bundle != "" {
		t.Error("Expected P12Bundle to be empty when using API token")
	}

	if cfg.Cert != "" || cfg.Key != "" {
		t.Error("Expected Cert and Key to be empty when using API token")
	}
}

func TestConfig_Validate_Valid_APIToken(t *testing.T) {
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		APIToken:  true,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected valid config with API token, got error: %v", err)
	}
}

func TestConfig_Save_APIToken(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".f5xcconfig")

	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		APIToken:  true,
	}

	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Reload and verify
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if !loaded.APIToken {
		t.Error("Expected APIToken to be true after reload")
	}
}

func TestConfig_Validate_APIToken_NoOtherAuth(t *testing.T) {
	// API token should be valid without P12 or cert/key
	cfg := &Config{
		ServerURL: "https://test.console.ves.volterra.io/api",
		APIToken:  true,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("API token config should be valid: %v", err)
	}
}
