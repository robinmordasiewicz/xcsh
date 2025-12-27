// Package profile provides multi-profile support for F5 XC tenant configurations.
// Profiles are stored in XDG-compliant locations (~/.config/xcsh/profiles/).
package profile

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Profile represents an F5 XC tenant configuration.
type Profile struct {
	// Name is the profile identifier (used as filename)
	Name string `yaml:"name"`
	// APIURL is the F5 XC console URL (e.g., https://tenant.console.ves.volterra.io)
	APIURL string `yaml:"api_url"`
	// APIToken is the authentication token (stored with 0600 permissions)
	APIToken string `yaml:"api_token,omitempty"`
	// P12Bundle is the path to a PKCS#12 certificate bundle
	P12Bundle string `yaml:"p12_bundle,omitempty"`
	// Cert is the path to a certificate file
	Cert string `yaml:"cert,omitempty"`
	// Key is the path to a private key file
	Key string `yaml:"key,omitempty"`
	// DefaultNamespace is the default namespace for operations
	DefaultNamespace string `yaml:"default_namespace,omitempty"`
}

// GlobalConfig holds application-wide settings.
type GlobalConfig struct {
	// DefaultProfile is the name of the profile to use when none is specified
	DefaultProfile string `yaml:"default_profile"`
}

// validNameRegex matches valid profile names (alphanumeric, hyphens, underscores)
var validNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// Validate checks if the profile configuration is valid.
func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile name is required")
	}

	if !validNameRegex.MatchString(p.Name) {
		return fmt.Errorf("profile name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}

	if p.APIURL == "" {
		return fmt.Errorf("api_url is required")
	}

	// Validate URL format
	parsedURL, err := url.Parse(p.APIURL)
	if err != nil {
		return fmt.Errorf("invalid api_url: %w", err)
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return fmt.Errorf("api_url must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("api_url must have a host")
	}

	// Check authentication method
	hasToken := p.APIToken != ""
	hasP12 := p.P12Bundle != ""
	hasCert := p.Cert != "" && p.Key != ""

	authMethods := 0
	if hasToken {
		authMethods++
	}
	if hasP12 {
		authMethods++
	}
	if hasCert {
		authMethods++
	}

	if authMethods == 0 {
		return fmt.Errorf("at least one authentication method is required (api_token, p12_bundle, or cert/key)")
	}

	if authMethods > 1 {
		return fmt.Errorf("only one authentication method should be configured (api_token, p12_bundle, or cert/key)")
	}

	// Validate file paths exist if specified
	if hasP12 {
		if _, err := os.Stat(expandPath(p.P12Bundle)); err != nil {
			return fmt.Errorf("p12_bundle file not found: %s", p.P12Bundle)
		}
	}

	if hasCert {
		if _, err := os.Stat(expandPath(p.Cert)); err != nil {
			return fmt.Errorf("cert file not found: %s", p.Cert)
		}
		if _, err := os.Stat(expandPath(p.Key)); err != nil {
			return fmt.Errorf("key file not found: %s", p.Key)
		}
	}

	return nil
}

// AuthMethod returns a string describing the authentication method.
func (p *Profile) AuthMethod() string {
	if p.APIToken != "" {
		return "API Token"
	}
	if p.P12Bundle != "" {
		return "P12 Bundle"
	}
	if p.Cert != "" && p.Key != "" {
		return "Certificate/Key"
	}
	return "None"
}

// TenantName extracts the tenant name from the API URL.
func (p *Profile) TenantName() string {
	parsedURL, err := url.Parse(p.APIURL)
	if err != nil {
		return ""
	}

	// Extract tenant from hostname (e.g., "tenant" from "tenant.console.ves.volterra.io")
	parts := strings.Split(parsedURL.Host, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return parsedURL.Host
}

// MaskedToken returns a masked version of the token for display.
func (p *Profile) MaskedToken() string {
	if p.APIToken == "" {
		return ""
	}
	if len(p.APIToken) <= 8 {
		return "****"
	}
	return p.APIToken[:4] + "****" + p.APIToken[len(p.APIToken)-4:]
}

// ToYAML serializes the profile to YAML format.
func (p *Profile) ToYAML() ([]byte, error) {
	return yaml.Marshal(p)
}

// FromYAML deserializes a profile from YAML format.
func FromYAML(data []byte) (*Profile, error) {
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}
	return &p, nil
}

// expandPath expands ~ to home directory.
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
