package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigDirName is the XDG-compliant config directory name
	ConfigDirName = "xcsh"
	// ProfilesDirName is the subdirectory for profile files
	ProfilesDirName = "profiles"
	// GlobalConfigName is the filename for global configuration
	GlobalConfigName = "config.yaml"
	// FilePermissions for profile files (owner read/write only)
	FilePermissions = 0600
	// DirPermissions for config directories
	DirPermissions = 0700
)

// Manager handles profile operations and storage.
type Manager struct {
	configDir    string        // ~/.config/xcsh
	profilesDir  string        // ~/.config/xcsh/profiles
	globalConfig *GlobalConfig // Global settings
	current      *Profile      // Currently active profile
}

// NewManager creates a new profile manager with XDG-compliant paths.
func NewManager() (*Manager, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config directory: %w", err)
	}

	m := &Manager{
		configDir:   configDir,
		profilesDir: filepath.Join(configDir, ProfilesDirName),
		globalConfig: &GlobalConfig{
			DefaultProfile: "",
		},
	}

	// Ensure directories exist
	if err := m.ensureDirectories(); err != nil {
		return nil, err
	}

	// Load global config if exists
	if err := m.loadGlobalConfig(); err != nil {
		// Ignore error if file doesn't exist
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return m, nil
}

// getConfigDir returns the XDG-compliant config directory.
func getConfigDir() (string, error) {
	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, ConfigDirName), nil
	}

	// Fall back to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", ConfigDirName), nil
}

// ensureDirectories creates necessary directories if they don't exist.
func (m *Manager) ensureDirectories() error {
	if err := os.MkdirAll(m.configDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.MkdirAll(m.profilesDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	return nil
}

// loadGlobalConfig loads the global configuration file.
func (m *Manager) loadGlobalConfig() error {
	configPath := filepath.Join(m.configDir, GlobalConfigName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config GlobalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse global config: %w", err)
	}

	m.globalConfig = &config
	return nil
}

// saveGlobalConfig saves the global configuration file.
func (m *Manager) saveGlobalConfig() error {
	data, err := yaml.Marshal(m.globalConfig)
	if err != nil {
		return fmt.Errorf("failed to serialize global config: %w", err)
	}

	configPath := filepath.Join(m.configDir, GlobalConfigName)
	if err := os.WriteFile(configPath, data, FilePermissions); err != nil {
		return fmt.Errorf("failed to write global config: %w", err)
	}

	return nil
}

// profilePath returns the file path for a profile.
func (m *Manager) profilePath(name string) string {
	return filepath.Join(m.profilesDir, name+".yaml")
}

// Create saves a new profile.
func (m *Manager) Create(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}

	// Check if profile already exists
	if m.Exists(p.Name) {
		return fmt.Errorf("profile %q already exists", p.Name)
	}

	return m.save(p)
}

// Update saves changes to an existing profile.
func (m *Manager) Update(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}

	if !m.Exists(p.Name) {
		return fmt.Errorf("profile %q does not exist", p.Name)
	}

	return m.save(p)
}

// save writes a profile to disk.
func (m *Manager) save(p *Profile) error {
	data, err := p.ToYAML()
	if err != nil {
		return fmt.Errorf("failed to serialize profile: %w", err)
	}

	path := m.profilePath(p.Name)
	if err := os.WriteFile(path, data, FilePermissions); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil
}

// Load reads a profile from disk.
func (m *Manager) Load(name string) (*Profile, error) {
	path := m.profilePath(name)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile %q not found", name)
		}
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	p, err := FromYAML(data)
	if err != nil {
		return nil, err
	}

	// Ensure name matches filename
	p.Name = name

	return p, nil
}

// Delete removes a profile.
func (m *Manager) Delete(name string) error {
	if !m.Exists(name) {
		return fmt.Errorf("profile %q not found", name)
	}

	// Don't allow deleting the current default profile without changing it first
	if m.globalConfig.DefaultProfile == name {
		return fmt.Errorf("cannot delete default profile %q; use 'login profile use <other>' first", name)
	}

	path := m.profilePath(name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// Exists checks if a profile exists.
func (m *Manager) Exists(name string) bool {
	path := m.profilePath(name)
	_, err := os.Stat(path)
	return err == nil
}

// List returns all available profile names.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read profiles directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			profileName := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
			names = append(names, profileName)
		}
	}

	sort.Strings(names)
	return names, nil
}

// ListProfiles returns all profiles with their details.
func (m *Manager) ListProfiles() ([]*Profile, error) {
	names, err := m.List()
	if err != nil {
		return nil, err
	}

	profiles := make([]*Profile, 0, len(names))
	for _, name := range names {
		p, err := m.Load(name)
		if err != nil {
			// Log but continue with other profiles
			continue
		}
		profiles = append(profiles, p)
	}

	return profiles, nil
}

// SetDefault sets the default profile.
func (m *Manager) SetDefault(name string) error {
	if !m.Exists(name) {
		return fmt.Errorf("profile %q not found", name)
	}

	m.globalConfig.DefaultProfile = name
	return m.saveGlobalConfig()
}

// GetDefault returns the default profile name.
func (m *Manager) GetDefault() string {
	return m.globalConfig.DefaultProfile
}

// GetCurrent returns the currently active profile.
// Priority: explicitly set current > environment variable > default profile
func (m *Manager) GetCurrent() (*Profile, error) {
	if m.current != nil {
		return m.current, nil
	}

	// Check environment variable
	if envProfile := os.Getenv("F5XC_PROFILE"); envProfile != "" {
		return m.Load(envProfile)
	}

	// Use default profile
	if m.globalConfig.DefaultProfile != "" {
		return m.Load(m.globalConfig.DefaultProfile)
	}

	// Return first available profile if no default set
	names, err := m.List()
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no profiles configured; use 'xcsh login profile create' to create one")
	}

	return m.Load(names[0])
}

// SetCurrent sets the current profile for this session.
func (m *Manager) SetCurrent(name string) error {
	p, err := m.Load(name)
	if err != nil {
		return err
	}

	m.current = p
	return nil
}

// GetCurrentName returns the name of the current profile.
func (m *Manager) GetCurrentName() string {
	if m.current != nil {
		return m.current.Name
	}

	if envProfile := os.Getenv("F5XC_PROFILE"); envProfile != "" {
		return envProfile
	}

	if m.globalConfig.DefaultProfile != "" {
		return m.globalConfig.DefaultProfile
	}

	names, _ := m.List()
	if len(names) > 0 {
		return names[0]
	}

	return ""
}

// IsDefault returns true if the given profile is the default.
func (m *Manager) IsDefault(name string) bool {
	return m.globalConfig.DefaultProfile == name
}

// ConfigDir returns the configuration directory path.
func (m *Manager) ConfigDir() string {
	return m.configDir
}

// ProfilesDir returns the profiles directory path.
func (m *Manager) ProfilesDir() string {
	return m.profilesDir
}
