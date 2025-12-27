package client

import (
	"fmt"
	"net/url"
	"strings"
)

// NormalizeAPIURL normalizes F5 Distributed Cloud API URLs to a consistent format.
// It handles various user input formats and returns a properly formatted API URL.
//
// Supported input formats:
//   - https://tenant.domain.com
//   - https://tenant.domain.com/api
//   - https://tenant.domain.com/api/api
//   - https://tenant.domain.com/some/other/path
//   - http://localhost:8001 (development)
//   - tenant.domain.com (adds https://)
//
// Output format: https://tenant.domain.com (without /api suffix)
//
// The /api prefix is added by individual API endpoints, so the base URL
// should NOT include it.
func NormalizeAPIURL(rawURL string) (string, error) {
	// Trim whitespace first
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return "", fmt.Errorf("API URL cannot be empty")
	}

	// Handle localhost development URLs specially
	if strings.Contains(rawURL, "localhost") || strings.Contains(rawURL, "127.0.0.1") {
		// For localhost, preserve the URL as-is (including port)
		if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
			rawURL = "http://" + rawURL
		}
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return "", fmt.Errorf("invalid localhost URL: %w", err)
		}
		// Remove /api suffix if present for consistency
		parsed.Path = strings.TrimSuffix(parsed.Path, "/api")
		parsed.Path = strings.TrimSuffix(parsed.Path, "/")
		return parsed.String(), nil
	}

	// Add https:// if no scheme provided
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Parse the URL
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Ensure https for production (not localhost)
	if parsed.Scheme == "http" && !strings.Contains(parsed.Host, "localhost") && !strings.Contains(parsed.Host, "127.0.0.1") {
		parsed.Scheme = "https"
	}

	// Normalize the path by removing /api patterns
	// F5 XC API endpoints add /api themselves, so the base URL shouldn't include it
	normalizedPath := normalizePath(parsed.Path)
	parsed.Path = normalizedPath

	// Remove trailing slash for consistency
	result := strings.TrimSuffix(parsed.String(), "/")

	return result, nil
}

// normalizePath removes /api path components and normalizes the URL path.
// It handles various patterns:
//   - /api -> /
//   - /api/api -> /
//   - /api/api/api -> /
//   - /some/path/api -> /some/path
//   - /api/some/other/path -> /some/other/path
func normalizePath(path string) string {
	if path == "" || path == "/" {
		return ""
	}

	// Split path into components
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Filter out all "api" components
	var filtered []string
	for _, part := range parts {
		if part != "" && part != "api" {
			filtered = append(filtered, part)
		}
	}

	// If no parts remain, return empty path
	if len(filtered) == 0 {
		return ""
	}

	// Rebuild path
	return "/" + strings.Join(filtered, "/")
}

// ExtractTenant extracts the tenant name from an F5 XC URL.
// For example:
//   - https://tenant1.console.ves.volterra.io -> "tenant1"
//   - https://tenant2.staging.volterra.us -> "tenant2"
//   - https://localhost:8001 -> "local"
func ExtractTenant(apiURL string) string {
	parsed, err := url.Parse(apiURL)
	if err != nil {
		return "unknown"
	}

	host := parsed.Hostname()

	// Handle localhost
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
		return "local"
	}

	// Extract first part of hostname as tenant
	// Example: tenant1.console.ves.volterra.io -> tenant1
	parts := strings.Split(host, ".")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}

	return "unknown"
}

// ValidateAPIURL performs basic validation on the API URL.
func ValidateAPIURL(apiURL string) error {
	if apiURL == "" {
		return fmt.Errorf("API URL cannot be empty")
	}

	parsed, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Host == "" {
		return fmt.Errorf("URL must include a hostname")
	}

	// Check for common mistakes
	if strings.Count(apiURL, "/api") > 1 {
		return fmt.Errorf("URL contains multiple /api paths - use NormalizeAPIURL to fix")
	}

	return nil
}
