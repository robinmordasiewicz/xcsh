package validation

import (
	"strings"
	"testing"
)

func TestTierLevel(t *testing.T) {
	tests := []struct {
		name     string
		tier     string
		expected int
	}{
		{"Standard tier", TierStandard, 1},
		{"Professional tier", TierProfessional, 2},
		{"Enterprise tier", TierEnterprise, 3},
		{"Unknown tier", "Unknown", 0},
		{"Empty tier", "", 0},
		{"Mixed case", "standard", 0}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TierLevel(tt.tier)
			if got != tt.expected {
				t.Errorf("TierLevel(%q) = %d, want %d", tt.tier, got, tt.expected)
			}
		})
	}
}

func TestIsSufficientTier(t *testing.T) {
	tests := []struct {
		name         string
		currentTier  string
		requiredTier string
		expected     bool
	}{
		// Standard user tests
		{"Standard accessing Standard", TierStandard, TierStandard, true},
		{"Standard accessing Professional", TierStandard, TierProfessional, false},
		{"Standard accessing Enterprise", TierStandard, TierEnterprise, false},

		// Professional user tests
		{"Professional accessing Standard", TierProfessional, TierStandard, true},
		{"Professional accessing Professional", TierProfessional, TierProfessional, true},
		{"Professional accessing Enterprise", TierProfessional, TierEnterprise, false},

		// Enterprise user tests
		{"Enterprise accessing Standard", TierEnterprise, TierStandard, true},
		{"Enterprise accessing Professional", TierEnterprise, TierProfessional, true},
		{"Enterprise accessing Enterprise", TierEnterprise, TierEnterprise, true},

		// Edge cases
		{"Empty requirement (defaults to accessible)", TierStandard, "", true},
		{"Unknown required tier", TierProfessional, "Unknown", false},
		{"Unknown current tier", "Unknown", TierStandard, false},
		{"Both unknown", "Unknown1", "Unknown2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSufficientTier(tt.currentTier, tt.requiredTier)
			if got != tt.expected {
				t.Errorf(
					"IsSufficientTier(%q, %q) = %v, want %v",
					tt.currentTier,
					tt.requiredTier,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestTierName(t *testing.T) {
	tests := []struct {
		name     string
		tier     string
		expected string
	}{
		{"Standard", TierStandard, "Standard"},
		{"Professional", TierProfessional, "Professional"},
		{"Enterprise", TierEnterprise, "Enterprise"},
		{"Unknown returns original", "Custom", "Custom"},
		{"Empty returns empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TierName(tt.tier)
			if got != tt.expected {
				t.Errorf("TierName(%q) = %q, want %q", tt.tier, got, tt.expected)
			}
		})
	}
}

func TestGetNextTier(t *testing.T) {
	tests := []struct {
		name     string
		tier     string
		expected string
	}{
		{"Standard → Professional", TierStandard, TierProfessional},
		{"Professional → Enterprise", TierProfessional, TierEnterprise},
		{"Enterprise → empty (highest)", TierEnterprise, ""},
		{"Unknown → Professional (default)", "Unknown", TierProfessional},
		{"Empty → Professional (default)", "", TierProfessional},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNextTier(tt.tier)
			if got != tt.expected {
				t.Errorf("GetNextTier(%q) = %q, want %q", tt.tier, got, tt.expected)
			}
		})
	}
}

func TestGetUpgradePath(t *testing.T) {
	tests := []struct {
		name          string
		currentTier   string
		requiredTier  string
		expectedPath  string
		shouldBeEmpty bool
	}{
		// Sufficient tier cases (should return empty)
		{
			"Standard for Standard (sufficient)",
			TierStandard,
			TierStandard,
			"",
			true,
		},
		{
			"Professional for Standard (sufficient)",
			TierProfessional,
			TierStandard,
			"",
			true,
		},
		{
			"Enterprise for Professional (sufficient)",
			TierEnterprise,
			TierProfessional,
			"",
			true,
		},

		// Insufficient tier cases
		{
			"Standard → Professional (insufficient)",
			TierStandard,
			TierProfessional,
			"Upgrade from Standard to Professional tier",
			false,
		},
		{
			"Professional → Enterprise (insufficient)",
			TierProfessional,
			TierEnterprise,
			"Upgrade from Professional to Enterprise tier",
			false,
		},
		{
			"Standard → Enterprise (insufficient)",
			TierStandard,
			TierEnterprise,
			"Upgrade from Standard to Enterprise tier",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUpgradePath(tt.currentTier, tt.requiredTier)
			if tt.shouldBeEmpty {
				if got != "" {
					t.Errorf(
						"GetUpgradePath(%q, %q) should be empty but got %q",
						tt.currentTier,
						tt.requiredTier,
						got,
					)
				}
			} else {
				if got != tt.expectedPath {
					t.Errorf(
						"GetUpgradePath(%q, %q) = %q, want %q",
						tt.currentTier,
						tt.requiredTier,
						got,
						tt.expectedPath,
					)
				}
			}
		})
	}
}

func TestTierAccessError(t *testing.T) {
	tests := []struct {
		name     string
		error    *TierAccessError
		wantText []string // Substrings that should appear in error message
	}{
		{
			"Standard tier insufficient for Professional domain",
			NewTierAccessError("api", "API", TierStandard, TierProfessional),
			[]string{
				"API",
				"Professional",
				"Standard",
				"Upgrade from Standard to Professional",
				"console.volterra.io",
			},
		},
		{
			"Professional insufficient for Enterprise domain",
			NewTierAccessError("generative_ai", "Generative AI", TierProfessional, TierEnterprise),
			[]string{
				"Generative AI",
				"Enterprise",
				"Professional",
				"Upgrade from Professional to Enterprise",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.error.Error()

			for _, want := range tt.wantText {
				if !strings.Contains(errMsg, want) {
					t.Errorf(
						"Error message missing %q. Got: %q",
						want,
						errMsg,
					)
				}
			}
		})
	}
}

func TestNewTierAccessError(t *testing.T) {
	err := NewTierAccessError("api", "API", TierStandard, TierProfessional)

	if err.Domain != "api" {
		t.Errorf("Domain = %q, want %q", err.Domain, "api")
	}
	if err.DomainDisplay != "API" {
		t.Errorf("DomainDisplay = %q, want %q", err.DomainDisplay, "API")
	}
	if err.CurrentTier != TierStandard {
		t.Errorf("CurrentTier = %q, want %q", err.CurrentTier, TierStandard)
	}
	if err.RequiredTier != TierProfessional {
		t.Errorf("RequiredTier = %q, want %q", err.RequiredTier, TierProfessional)
	}
}

// Benchmarks
func BenchmarkTierLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TierLevel(TierProfessional)
	}
}

func BenchmarkIsSufficientTier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsSufficientTier(TierProfessional, TierStandard)
	}
}

func BenchmarkGetUpgradePath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetUpgradePath(TierStandard, TierProfessional)
	}
}
