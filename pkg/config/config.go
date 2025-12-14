package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the f5xcctl configuration file structure
type Config struct {
	// ServerURL is the API server endpoint
	ServerURL string `yaml:"server-url"`

	// P12Bundle is the path to the P12 certificate bundle
	P12Bundle string `yaml:"p12-bundle,omitempty"`

	// Cert is the path to the client certificate file
	Cert string `yaml:"cert,omitempty"`

	// Key is the path to the client key file
	Key string `yaml:"key,omitempty"`

	// APIToken indicates API token authentication mode (actual token from F5XC_API_TOKEN env var)
	APIToken bool `yaml:"api-token,omitempty"`
}

// rawConfig is used for YAML parsing
type rawConfig struct {
	ServerURL string `yaml:"server-url"`
	P12Bundle string `yaml:"p12-bundle"`
	Cert      string `yaml:"cert"`
	Key       string `yaml:"key"`
	APIToken  bool   `yaml:"api-token"`
}

// Load reads and parses a f5xcctl config file
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg := &Config{
		ServerURL: raw.ServerURL,
		P12Bundle: raw.P12Bundle,
		Cert:      raw.Cert,
		Key:       raw.Key,
		APIToken:  raw.APIToken,
	}

	return cfg, nil
}

// Save writes the config to a file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the config has required fields
func (c *Config) Validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("server-url is required")
	}

	// Must have either P12 bundle, cert/key pair, or API token
	hasP12 := c.P12Bundle != ""
	hasCertKey := c.Cert != "" && c.Key != ""
	hasAPIToken := c.APIToken

	if !hasP12 && !hasCertKey && !hasAPIToken {
		return fmt.Errorf("authentication required: p12-bundle, cert/key pair, or api-token")
	}

	return nil
}
