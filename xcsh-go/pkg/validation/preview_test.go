package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPreview(t *testing.T) {
	tests := []struct {
		name      string
		isPreview bool
		expected  bool
	}{
		{"Preview domain", true, true},
		{"Stable domain", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPreview(tt.isPreview)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPreviewWarning(t *testing.T) {
	domain := "generative_ai"
	displayName := "Generative AI"

	warning := GetPreviewWarning(domain, displayName)

	assert.NotNil(t, warning)
	assert.Equal(t, domain, warning.Domain)
	assert.Equal(t, displayName, warning.DomainDisplay)
	assert.Equal(t, "beta", warning.Status)
}

func TestPreviewWarningError(t *testing.T) {
	warning := GetPreviewWarning("generative_ai", "Generative AI")

	errMsg := warning.Error()
	assert.NotEmpty(t, errMsg)

	// Verify warning message contains key information
	assert.Contains(t, errMsg, "PREVIEW", "Warning should indicate preview status")
	assert.Contains(t, errMsg, "Generative AI", "Warning should mention domain name")
	assert.Contains(t, errMsg, "beta/experimental", "Warning should explain beta nature")
	assert.Contains(t, errMsg, "support@f5.com", "Warning should include support contact")
	assert.Contains(t, errMsg, "console.volterra.io/status", "Warning should include status URL")
}

func TestFormatPreviewBadge(t *testing.T) {
	badge := FormatPreviewBadge()

	assert.Equal(t, "[PREVIEW]", badge)
	assert.NotEmpty(t, badge)
}

func TestGetPreviewIndicator(t *testing.T) {
	indicator := GetPreviewIndicator()

	assert.Equal(t, "⚠️", indicator)
	assert.NotEmpty(t, indicator)
}

func TestAppendPreviewToShortDescription(t *testing.T) {
	tests := []struct {
		name          string
		shortDesc     string
		isPreview     bool
		expectedRegex string
	}{
		{
			name:          "Preview domain",
			shortDesc:     "Manage AI resources",
			isPreview:     true,
			expectedRegex: "\\[PREVIEW\\].*Manage AI resources",
		},
		{
			name:          "Stable domain",
			shortDesc:     "Manage DNS resources",
			isPreview:     false,
			expectedRegex: "Manage DNS resources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendPreviewToShortDescription(tt.shortDesc, tt.isPreview)
			assert.NotEmpty(t, result)

			if tt.isPreview {
				assert.Contains(t, result, "[PREVIEW]")
				assert.Contains(t, result, tt.shortDesc)
			} else {
				assert.Equal(t, tt.shortDesc, result)
			}
		})
	}
}

func TestGetPreviewWarningText(t *testing.T) {
	text := GetPreviewWarningText("Generative AI")

	assert.NotEmpty(t, text)
	assert.Contains(t, text, "⚠️", "Should include warning indicator")
	assert.Contains(t, text, "PREVIEW", "Should indicate preview status")
	assert.Contains(t, text, "Generative AI", "Should mention domain name")
	assert.Contains(t, text, "beta", "Should mention beta status")
	assert.Contains(t, text, "support@f5.com", "Should include support contact")
}

func TestPreviewWarningMultiLine(t *testing.T) {
	warning := GetPreviewWarning("generative_ai", "Generative AI")
	errMsg := warning.Error()

	// Verify it's multi-line
	lines := strings.Split(errMsg, "\n")
	assert.Greater(t, len(lines), 2, "Warning should be multi-line")
}

func TestPreviewBadgeInDescription(t *testing.T) {
	originalDesc := "Manage generative AI resources"
	previewDesc := AppendPreviewToShortDescription(originalDesc, true)

	// Preview indicator should be at the start
	assert.True(t, strings.HasPrefix(previewDesc, "[PREVIEW]"),
		"Preview badge should be at start of description")

	// Original description should still be present
	assert.Contains(t, previewDesc, originalDesc,
		"Original description should be preserved")
}

func TestPreviewIndicatorFormat(t *testing.T) {
	indicator := GetPreviewIndicator()
	warning := GetPreviewWarning("test", "Test")
	warningMsg := warning.Error()

	// Indicator should be in warning message
	assert.Contains(t, warningMsg, indicator)
}

func TestPreviewWarningConsistency(t *testing.T) {
	domain := "generative_ai"
	displayName := "Generative AI"

	warning1 := GetPreviewWarning(domain, displayName)
	warning2 := GetPreviewWarning(domain, displayName)

	// Both warnings should produce the same message
	assert.Equal(t, warning1.Error(), warning2.Error())
}

// Benchmark tests for preview functions
func BenchmarkIsPreview(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPreview(true)
	}
}

func BenchmarkGetPreviewWarning(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetPreviewWarning("generative_ai", "Generative AI")
	}
}

func BenchmarkFormatPreviewBadge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatPreviewBadge()
	}
}

func BenchmarkGetPreviewIndicator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetPreviewIndicator()
	}
}

func BenchmarkAppendPreviewToShortDescription(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AppendPreviewToShortDescription("Manage resources", true)
	}
}
