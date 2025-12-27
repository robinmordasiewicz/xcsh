// Package validation provides subscription tier validation and access control.
package validation

import (
	"fmt"
)

// PreviewWarning represents a warning for accessing a preview/beta domain
type PreviewWarning struct {
	Domain        string // The domain name
	DomainDisplay string // Human-readable domain name
	Status        string // Preview status (e.g., "beta", "experimental")
}

// Error implements the error interface and returns a formatted warning message.
func (w *PreviewWarning) Error() string {
	return fmt.Sprintf(
		"⚠️  Domain '%s' is in PREVIEW\n\n"+
			"This is a beta/experimental feature and may have limited functionality,\n"+
			"performance issues, or breaking changes.\n\n"+
			"For feedback or issues: contact support@f5.com\n"+
			"Status: https://console.volterra.io/status",
		w.DomainDisplay,
	)
}

// IsPreview checks if a domain is marked as preview/beta
func IsPreview(isPreview bool) bool {
	return isPreview
}

// GetPreviewWarning creates a new PreviewWarning with the given parameters.
func GetPreviewWarning(domain, domainDisplay string) *PreviewWarning {
	return &PreviewWarning{
		Domain:        domain,
		DomainDisplay: domainDisplay,
		Status:        "beta",
	}
}

// FormatPreviewBadge returns a formatted badge for preview domains
// Used in help text to indicate a domain is in preview
func FormatPreviewBadge() string {
	return "[PREVIEW]"
}

// GetPreviewIndicator returns a short indicator string for preview domains
// Used in list output and completions
func GetPreviewIndicator() string {
	return "⚠️"
}

// AppendPreviewToShortDescription adds preview indicator to a domain's short description
func AppendPreviewToShortDescription(shortDesc string, isPreview bool) string {
	if !isPreview {
		return shortDesc
	}
	return fmt.Sprintf("%s %s", FormatPreviewBadge(), shortDesc)
}

// GetPreviewWarningText returns a multi-line warning text for preview domains
func GetPreviewWarningText(domainDisplay string) string {
	return fmt.Sprintf(
		"⚠️  PREVIEW: Domain '%s' is in beta and may have limited functionality or breaking changes.\n"+
			"For feedback: support@f5.com | Status: https://console.volterra.io/status",
		domainDisplay,
	)
}
