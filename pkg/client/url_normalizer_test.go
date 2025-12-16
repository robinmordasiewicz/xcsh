package client

import (
	"testing"
)

func TestNormalizeAPIURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		// Basic formats
		{
			name:     "full URL with https",
			input:    "https://tenant.console.ves.volterra.io",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "domain only (adds https)",
			input:    "tenant.console.ves.volterra.io",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with /api suffix",
			input:    "https://tenant.console.ves.volterra.io/api",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with double /api/api",
			input:    "https://tenant.console.ves.volterra.io/api/api",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with triple /api/api/api",
			input:    "https://tenant.console.ves.volterra.io/api/api/api",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with trailing slash",
			input:    "https://tenant.console.ves.volterra.io/",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with /api and trailing slash",
			input:    "https://tenant.console.ves.volterra.io/api/",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},

		// Staging environment
		{
			name:     "staging environment",
			input:    "https://tenant.staging.volterra.us",
			expected: "https://tenant.staging.volterra.us",
			wantErr:  false,
		},
		{
			name:     "staging with /api",
			input:    "https://tenant.staging.volterra.us/api",
			expected: "https://tenant.staging.volterra.us",
			wantErr:  false,
		},

		// HTTP upgrade to HTTPS
		{
			name:     "http URL (upgrades to https)",
			input:    "http://tenant.console.ves.volterra.io",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "http with /api (upgrades and removes /api)",
			input:    "http://tenant.console.ves.volterra.io/api",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},

		// Localhost (development)
		{
			name:     "localhost with port",
			input:    "http://localhost:8001",
			expected: "http://localhost:8001",
			wantErr:  false,
		},
		{
			name:     "localhost with /api",
			input:    "http://localhost:8001/api",
			expected: "http://localhost:8001",
			wantErr:  false,
		},
		{
			name:     "localhost without scheme",
			input:    "localhost:8001",
			expected: "http://localhost:8001",
			wantErr:  false,
		},
		{
			name:     "127.0.0.1 with port",
			input:    "http://127.0.0.1:8001",
			expected: "http://127.0.0.1:8001",
			wantErr:  false,
		},
		{
			name:     "127.0.0.1 without scheme",
			input:    "127.0.0.1:8001",
			expected: "http://127.0.0.1:8001",
			wantErr:  false,
		},

		// Other paths that should be removed
		{
			name:     "URL with other path after /api",
			input:    "https://tenant.console.ves.volterra.io/api/config",
			expected: "https://tenant.console.ves.volterra.io/config",
			wantErr:  false,
		},
		{
			name:     "URL with /api in middle of path",
			input:    "https://tenant.console.ves.volterra.io/some/api/path",
			expected: "https://tenant.console.ves.volterra.io/some/path",
			wantErr:  false,
		},

		// Whitespace handling
		{
			name:     "URL with leading whitespace",
			input:    "  https://tenant.console.ves.volterra.io",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with trailing whitespace",
			input:    "https://tenant.console.ves.volterra.io  ",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},
		{
			name:     "URL with leading and trailing whitespace",
			input:    "  https://tenant.console.ves.volterra.io  ",
			expected: "https://tenant.console.ves.volterra.io",
			wantErr:  false,
		},

		// Error cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeAPIURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeAPIURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("NormalizeAPIURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "",
		},
		{
			name:     "single /api",
			input:    "/api",
			expected: "",
		},
		{
			name:     "double /api/api",
			input:    "/api/api",
			expected: "",
		},
		{
			name:     "triple /api/api/api",
			input:    "/api/api/api",
			expected: "",
		},
		{
			name:     "/api with path after",
			input:    "/api/config",
			expected: "/config",
		},
		{
			name:     "path with /api in middle",
			input:    "/some/api/path",
			expected: "/some/path",
		},
		{
			name:     "path with multiple /api",
			input:    "/api/some/api/path/api",
			expected: "/some/path",
		},
		{
			name:     "path without /api",
			input:    "/some/path",
			expected: "/some/path",
		},
		{
			name:     "trailing slash",
			input:    "/some/path/",
			expected: "/some/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePath(tt.input)
			if got != tt.expected {
				t.Errorf("normalizePath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestExtractTenant(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "production URL",
			input:    "https://tenant1.console.ves.volterra.io",
			expected: "tenant1",
		},
		{
			name:     "staging URL",
			input:    "https://tenant2.staging.volterra.us",
			expected: "tenant2",
		},
		{
			name:     "URL with /api",
			input:    "https://mytenant.console.ves.volterra.io/api",
			expected: "mytenant",
		},
		{
			name:     "localhost",
			input:    "http://localhost:8001",
			expected: "local",
		},
		{
			name:     "127.0.0.1",
			input:    "http://127.0.0.1:8001",
			expected: "local",
		},
		{
			name:     "invalid URL",
			input:    "not-a-url",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTenant(tt.input)
			if got != tt.expected {
				t.Errorf("ExtractTenant() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateAPIURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid URL",
			input:   "https://tenant.console.ves.volterra.io",
			wantErr: false,
		},
		{
			name:    "valid URL with /api",
			input:   "https://tenant.console.ves.volterra.io/api",
			wantErr: false,
		},
		{
			name:    "localhost",
			input:   "http://localhost:8001",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "URL with multiple /api (should be normalized first)",
			input:   "https://tenant.console.ves.volterra.io/api/api",
			wantErr: true,
		},
		{
			name:    "URL without hostname",
			input:   "https:///api",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAPIURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAPIURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNormalizeAPIURL(b *testing.B) {
	urls := []string{
		"https://tenant.console.ves.volterra.io",
		"https://tenant.console.ves.volterra.io/api",
		"https://tenant.console.ves.volterra.io/api/api",
		"tenant.console.ves.volterra.io",
		"http://localhost:8001",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, url := range urls {
			_, _ = NormalizeAPIURL(url)
		}
	}
}
