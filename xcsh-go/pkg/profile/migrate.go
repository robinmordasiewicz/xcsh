package profile

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// LegacyConfigFile is the old config file name
	LegacyConfigFile = ".f5xcconfig"
	// MigratedProfileName is the name for the migrated profile
	MigratedProfileName = "default"
)

// legacyConfig represents the old config file format
type legacyConfig struct {
	ServerURL string `yaml:"server-url"`
	P12Bundle string `yaml:"p12-bundle"`
	Cert      string `yaml:"cert"`
	Key       string `yaml:"key"`
	APIToken  bool   `yaml:"api-token"`
}

// MigrationResult contains information about a migration operation.
type MigrationResult struct {
	// Migrated indicates if migration was performed
	Migrated bool
	// ProfileName is the name of the created profile
	ProfileName string
	// LegacyPath is the path to the old config file
	LegacyPath string
	// Message describes what happened
	Message string
}

// NeedsMigration checks if there's a legacy config that should be migrated.
func (m *Manager) NeedsMigration() (bool, string) {
	// Check for legacy config file
	legacyPath := getLegacyConfigPath()
	if legacyPath == "" {
		return false, ""
	}

	if _, err := os.Stat(legacyPath); err != nil {
		return false, ""
	}

	// Check if we already have profiles
	names, err := m.List()
	if err != nil {
		return false, ""
	}

	// If we have profiles, no need to migrate
	if len(names) > 0 {
		return false, legacyPath
	}

	return true, legacyPath
}

// Migrate converts a legacy config file to the new profile format.
func (m *Manager) Migrate() (*MigrationResult, error) {
	result := &MigrationResult{}

	needsMigration, legacyPath := m.NeedsMigration()
	result.LegacyPath = legacyPath

	if !needsMigration {
		if legacyPath != "" {
			result.Message = "Legacy config exists but profiles already configured; no migration needed"
		} else {
			result.Message = "No legacy config found"
		}
		return result, nil
	}

	// Load legacy config
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return result, fmt.Errorf("failed to read legacy config: %w", err)
	}

	var legacy legacyConfig
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		return result, fmt.Errorf("failed to parse legacy config: %w", err)
	}

	// Create profile from legacy config
	profile, err := m.legacyToProfile(&legacy)
	if err != nil {
		return result, fmt.Errorf("failed to convert legacy config: %w", err)
	}

	// Save the new profile
	if err := m.Create(profile); err != nil {
		return result, fmt.Errorf("failed to save migrated profile: %w", err)
	}

	// Set as default
	if err := m.SetDefault(profile.Name); err != nil {
		return result, fmt.Errorf("failed to set default profile: %w", err)
	}

	result.Migrated = true
	result.ProfileName = profile.Name
	result.Message = fmt.Sprintf("Migrated legacy config to profile %q", profile.Name)

	return result, nil
}

// legacyToProfile converts a legacy config to a profile.
func (m *Manager) legacyToProfile(legacy *legacyConfig) (*Profile, error) {
	if legacy.ServerURL == "" {
		return nil, fmt.Errorf("legacy config missing server-url")
	}

	// Extract profile name from server URL
	profileName := extractProfileName(legacy.ServerURL)
	if profileName == "" {
		profileName = MigratedProfileName
	}

	profile := &Profile{
		Name:      profileName,
		APIURL:    legacy.ServerURL,
		P12Bundle: legacy.P12Bundle,
		Cert:      legacy.Cert,
		Key:       legacy.Key,
	}

	// Handle API token
	// Old config stores boolean, actual token in F5XC_API_TOKEN env var
	if legacy.APIToken {
		if token := os.Getenv("F5XC_API_TOKEN"); token != "" {
			profile.APIToken = token
		}
	}

	return profile, nil
}

// extractProfileName extracts a meaningful name from the server URL.
func extractProfileName(serverURL string) string {
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return ""
	}

	// Extract tenant from hostname
	// e.g., "tenant.console.ves.volterra.io" -> "tenant"
	// e.g., "tenant.staging.volterra.us" -> "tenant"
	parts := strings.Split(parsed.Host, ".")
	if len(parts) > 0 {
		name := parts[0]
		// Validate as profile name
		if validNameRegex.MatchString(name) {
			return name
		}
	}

	return ""
}

// getLegacyConfigPath returns the path to the legacy config file.
func getLegacyConfigPath() string {
	// Check F5XC_CONFIG environment variable first
	if configPath := os.Getenv("F5XC_CONFIG"); configPath != "" {
		return configPath
	}

	// Default to ~/.f5xcconfig
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, LegacyConfigFile)
}

// BackupLegacyConfig creates a backup of the legacy config file.
func (m *Manager) BackupLegacyConfig() error {
	legacyPath := getLegacyConfigPath()
	if legacyPath == "" {
		return nil
	}

	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		return nil
	}

	backupPath := legacyPath + ".backup"
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return fmt.Errorf("failed to read legacy config: %w", err)
	}

	if err := os.WriteFile(backupPath, data, FilePermissions); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}
