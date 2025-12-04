package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the vesctl configuration file structure
type Config struct {
	// ServerURLs are the API server endpoints
	ServerURLs []string `yaml:"server-urls"`

	// P12Bundle is the path to the P12 certificate bundle
	P12Bundle string `yaml:"p12-bundle,omitempty"`

	// Cert is the path to the client certificate file
	Cert string `yaml:"cert,omitempty"`

	// Key is the path to the client key file
	Key string `yaml:"key,omitempty"`
}

// rawConfig is used for flexible YAML parsing (supports both single string and array)
type rawConfig struct {
	ServerURLs interface{} `yaml:"server-urls"`
	P12Bundle  string      `yaml:"p12-bundle"`
	Cert       string      `yaml:"cert"`
	Key        string      `yaml:"key"`
}

// Load reads and parses a vesctl config file
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
		P12Bundle: raw.P12Bundle,
		Cert:      raw.Cert,
		Key:       raw.Key,
	}

	// Handle server-urls as either single string or array
	switch v := raw.ServerURLs.(type) {
	case string:
		if v != "" {
			cfg.ServerURLs = []string{v}
		}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				cfg.ServerURLs = append(cfg.ServerURLs, s)
			}
		}
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
	if len(c.ServerURLs) == 0 {
		return fmt.Errorf("server-urls is required")
	}

	// Must have either P12 bundle or cert/key pair
	hasP12 := c.P12Bundle != ""
	hasCertKey := c.Cert != "" && c.Key != ""

	if !hasP12 && !hasCertKey {
		return fmt.Errorf("either p12-bundle or cert/key pair is required for authentication")
	}

	return nil
}
