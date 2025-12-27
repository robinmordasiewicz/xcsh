package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// TestCheckAndWarnPreviewDomain verifies preview domain detection
func TestCheckAndWarnPreviewDomain(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		expectWarning bool
	}{
		{
			name:          "Preview domain",
			domain:        "generative_ai",
			expectWarning: true,
		},
		{
			name:          "Stable domain",
			domain:        "dns",
			expectWarning: false,
		},
		{
			name:          "Another stable domain",
			domain:        "kubernetes_and_orchestration",
			expectWarning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := CheckAndWarnPreviewDomain(tt.domain)

			if tt.expectWarning {
				assert.NotNil(t, warning, "Should return warning for preview domain")
				assert.Equal(t, tt.domain, warning.Domain)
			} else {
				assert.Nil(t, warning, "Should return nil for stable domain")
			}
		})
	}
}

// TestCheckAndWarnPreviewDomainUnknown tests handling of non-existent domains
func TestCheckAndWarnPreviewDomainUnknown(t *testing.T) {
	warning := CheckAndWarnPreviewDomain("nonexistent_domain")
	assert.Nil(t, warning, "Should return nil for unknown domain")
}

// TestPreviewWarningContent verifies preview warning message content
func TestPreviewWarningContent(t *testing.T) {
	warning := CheckAndWarnPreviewDomain("generative_ai")

	require.NotNil(t, warning)
	warningMsg := warning.Error()

	// Verify warning contains required information
	assert.Contains(t, warningMsg, "PREVIEW", "Should indicate preview status")
	assert.Contains(t, warningMsg, "Generative Ai", "Should mention domain name")
	assert.Contains(t, warningMsg, "beta", "Should mention beta nature")
	assert.Contains(t, warningMsg, "experimental", "Should mention experimental nature")
	assert.Contains(t, warningMsg, "breaking changes", "Should warn about breaking changes")
	assert.Contains(t, warningMsg, "support@f5.com", "Should include support contact")
	assert.Contains(t, warningMsg, "console.volterra.io", "Should include status URL")
}

// TestAllPreviewDomainsDetected verifies all preview domains are properly identified
func TestAllPreviewDomainsDetected(t *testing.T) {
	previewCount := 0
	stableCount := 0

	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist", domain)

		if info.IsPreview {
			previewCount++

			// Verify CheckAndWarnPreviewDomain correctly identifies it
			warning := CheckAndWarnPreviewDomain(domain)
			assert.NotNil(t, warning, "Domain %q is marked preview but CheckAndWarnPreviewDomain didn't warn", domain)
		} else {
			stableCount++

			// Verify CheckAndWarnPreviewDomain doesn't warn for stable domains
			warning := CheckAndWarnPreviewDomain(domain)
			assert.Nil(t, warning, "Domain %q is stable but CheckAndWarnPreviewDomain warned", domain)
		}
	}

	assert.Equal(t, 1, previewCount, "Should have 1 preview domain")
	assert.Equal(t, 38, stableCount, "Should have 38 stable domains")
}

// TestPreviewBadgeFormatting verifies preview badges are formatted correctly
func TestPreviewBadgeFormatting(t *testing.T) {
	tests := []struct {
		name          string
		isPreview     bool
		expectedRegex string
	}{
		{
			name:      "Preview domain short description",
			isPreview: true,
		},
		{
			name:      "Stable domain short description",
			isPreview: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDesc := "Manage resources"
			result := validation.AppendPreviewToShortDescription(originalDesc, tt.isPreview)

			if tt.isPreview {
				assert.Contains(t, result, "[PREVIEW]", "Should include preview badge")
				assert.Contains(t, result, originalDesc, "Should preserve original description")
			} else {
				assert.Equal(t, originalDesc, result, "Should not modify stable domain descriptions")
			}
		})
	}
}

// TestPreviewWarningIndicator verifies preview indicator symbol
func TestPreviewWarningIndicator(t *testing.T) {
	indicator := validation.GetPreviewIndicator()
	assert.Equal(t, "⚠️", indicator)

	// Verify it appears in warnings
	warning := CheckAndWarnPreviewDomain("generative_ai")
	require.NotNil(t, warning)
	assert.Contains(t, warning.Error(), indicator)
}

// TestPreviewAndTierCombination verifies preview domain with tier requirements
func TestPreviewAndTierCombination(t *testing.T) {
	info, found := types.GetDomainInfo("generative_ai")
	require.True(t, found)

	// Generative AI should be both preview and require Advanced tier
	assert.True(t, info.IsPreview, "generative_ai should be in preview")
	assert.Equal(t, validation.TierAdvanced, info.RequiresTier,
		"generative_ai should require Advanced tier")
}

// TestPreviewWarningNonBlocking verifies preview warnings don't block domain access
func TestPreviewWarningNonBlocking(t *testing.T) {
	// Preview warnings should be nil for stable domains (allowing access)
	for domain := range types.DomainRegistry {
		info, _ := types.GetDomainInfo(domain)
		if !info.IsPreview {
			warning := CheckAndWarnPreviewDomain(domain)
			assert.Nil(t, warning, "Stable domain %q should not have warning", domain)
		}
	}
}

// TestPreviewWarningStructure verifies PreviewWarning struct initialization
func TestPreviewWarningStructure(t *testing.T) {
	domain := "generative_ai"
	displayName := "Generative AI"

	warning := validation.GetPreviewWarning(domain, displayName)

	assert.Equal(t, domain, warning.Domain)
	assert.Equal(t, displayName, warning.DomainDisplay)
	assert.Equal(t, "beta", warning.Status)
}

// TestPreviewWarningText verifies preview warning text formatting
func TestPreviewWarningText(t *testing.T) {
	text := validation.GetPreviewWarningText("Generative AI")

	assert.NotEmpty(t, text)
	assert.Contains(t, text, "⚠️")
	assert.Contains(t, text, "PREVIEW")
	assert.Contains(t, text, "Generative AI")
	assert.Contains(t, text, "beta")
	assert.Contains(t, text, "support@f5.com")
	assert.Contains(t, text, "console.volterra.io/status")
}

// TestPreviewFormatBadge verifies preview badge format
func TestPreviewFormatBadge(t *testing.T) {
	badge := validation.FormatPreviewBadge()
	assert.Equal(t, "[PREVIEW]", badge)
}

// TestGenerativeAIPreviewWarning specifically tests generative_ai domain
func TestGenerativeAIPreviewWarning(t *testing.T) {
	warning := CheckAndWarnPreviewDomain("generative_ai")

	require.NotNil(t, warning)
	assert.Equal(t, "generative_ai", warning.Domain)
	assert.Equal(t, "Generative Ai", warning.DomainDisplay) // Check actual display name from registry

	errMsg := warning.Error()

	// Verify comprehensive warning message
	assert.Contains(t, errMsg, "PREVIEW")
	assert.Contains(t, errMsg, "Generative Ai")
	assert.Contains(t, errMsg, "beta")
	assert.Contains(t, errMsg, "experimental")
	assert.Contains(t, errMsg, "support@f5.com")
}

// TestStableDomainNoWarning verifies stable domains don't generate warnings
func TestStableDomainNoWarning(t *testing.T) {
	stableDomains := []string{
		"dns",
		"kubernetes_and_orchestration",
		"authentication",
		"api",
		"network_security",
	}

	for _, domain := range stableDomains {
		t.Run("No_warning_for_"+domain, func(t *testing.T) {
			warning := CheckAndWarnPreviewDomain(domain)
			assert.Nil(t, warning, "Domain %q should not generate warning", domain)
		})
	}
}

// TestPreviewWarningMultiline verifies warning is multi-line for readability
func TestPreviewWarningMultiline(t *testing.T) {
	warning := CheckAndWarnPreviewDomain("generative_ai")
	require.NotNil(t, warning)

	errMsg := warning.Error()
	lines := strings.Split(errMsg, "\n")

	// Should have multiple lines for readability
	assert.Greater(t, len(lines), 2, "Warning should be multi-line for readability")
}

// TestIsPreviewFunction verifies IsPreview helper function
func TestIsPreviewFunction(t *testing.T) {
	tests := []struct {
		name      string
		isPreview bool
		expected  bool
	}{
		{"True", true, true},
		{"False", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validation.IsPreview(tt.isPreview)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPreviewConsistency verifies preview metadata is consistent
func TestPreviewConsistency(t *testing.T) {
	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found)

		// Verify consistency between domain info and warning generation
		if info.IsPreview {
			warning := CheckAndWarnPreviewDomain(domain)
			assert.NotNil(t, warning, "Preview domain %q should generate warning", domain)
			assert.Equal(t, domain, warning.Domain)
		} else {
			warning := CheckAndWarnPreviewDomain(domain)
			assert.Nil(t, warning, "Stable domain %q should not generate warning", domain)
		}
	}
}

// Benchmark tests for preview functions

func BenchmarkCheckAndWarnPreviewDomain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CheckAndWarnPreviewDomain("generative_ai")
	}
}

func BenchmarkCheckAndWarnPreviewDomainStable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CheckAndWarnPreviewDomain("dns")
	}
}

func BenchmarkPreviewFormattingWithBadge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validation.AppendPreviewToShortDescription("Manage resources", true)
	}
}
