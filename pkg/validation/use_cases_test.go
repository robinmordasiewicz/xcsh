package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// TestFormatUseCases verifies use case formatting with bullets
func TestFormatUseCases(t *testing.T) {
	useCases := []string{
		"Configure firewall rules",
		"Manage security policies",
		"Enable threat detection",
	}

	formatted := FormatUseCases(useCases)

	// Should contain USE CASES header
	assert.Contains(t, formatted, "USE CASES:")

	// Should contain all use cases with bullet points
	assert.Contains(t, formatted, "• Configure firewall rules")
	assert.Contains(t, formatted, "• Manage security policies")
	assert.Contains(t, formatted, "• Enable threat detection")

	// Should have newlines for multi-line format
	assert.Greater(t, strings.Count(formatted, "\n"), 2)
}

// TestFormatUseCasesEmpty verifies empty use cases return empty string
func TestFormatUseCasesEmpty(t *testing.T) {
	formatted := FormatUseCases([]string{})
	assert.Empty(t, formatted)
}

// TestFormatUseCasesShort verifies compact use case formatting
func TestFormatUseCasesShort(t *testing.T) {
	useCases := []string{
		"Configure firewall rules",
		"Manage security policies",
		"Enable threat detection",
		"Monitor alerts",
	}

	// Request first 2 use cases
	formatted := FormatUseCasesShort(useCases, 2)

	assert.Contains(t, formatted, "Configure firewall rules")
	assert.Contains(t, formatted, "Manage security policies")
	assert.NotContains(t, formatted, "Enable threat detection")
	assert.NotContains(t, formatted, "Monitor alerts")
}

// TestFormatUseCasesInline verifies inline comma-separated format
func TestFormatUseCasesInline(t *testing.T) {
	useCases := []string{
		"Configure rules",
		"Manage policies",
		"Enable detection",
	}

	formatted := FormatUseCasesInline(useCases)

	// Should be comma-separated
	assert.Equal(t, "Configure rules, Manage policies, Enable detection", formatted)

	// Should not contain newlines
	assert.NotContains(t, formatted, "\n")
}

// TestGetDomainUseCases verifies fetching use cases for a domain
func TestGetDomainUseCases(t *testing.T) {
	// Test domain with use cases
	formatted, err := GetDomainUseCases("api")
	require.NoError(t, err)
	assert.Contains(t, formatted, "USE CASES:")
	assert.NotEmpty(t, formatted)

	// Test domain without use cases
	formatted, err = GetDomainUseCases("admin_console_and_ui")
	require.NoError(t, err)
	assert.Empty(t, formatted)

	// Test non-existent domain
	_, err = GetDomainUseCases("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestGetDomainsWithUseCases verifies retrieval of domains with use cases
func TestGetDomainsWithUseCases(t *testing.T) {
	domains := GetDomainsWithUseCases()

	// Should have 31 domains with use cases
	assert.Greater(t, len(domains), 0)

	// All should have at least one use case
	for _, domain := range domains {
		assert.Greater(t, len(domain.UseCases), 0,
			"Domain %q should have at least one use case", domain.Name)
	}
}

// TestGetDomainsWithoutUseCases verifies retrieval of domains without use cases
func TestGetDomainsWithoutUseCases(t *testing.T) {
	domains := GetDomainsWithoutUseCases()

	// Should have 11 domains without use cases (42 total - 31 with)
	assert.Greater(t, len(domains), 0)

	// None should have use cases
	for _, domain := range domains {
		assert.Equal(t, 0, len(domain.UseCases),
			"Domain %q should have no use cases", domain.Name)
	}
}

// TestCalculateUseCaseStatistics verifies statistics calculation
func TestCalculateUseCaseStatistics(t *testing.T) {
	stats := CalculateUseCaseStatistics()

	// Verify totals
	assert.Equal(t, 42, stats.TotalDomains)
	assert.Equal(t, 31, stats.DomainsWithUseCases, "Should have 31 domains with use cases")
	assert.Equal(t, 11, stats.DomainsWithoutUseCases, "Should have 11 domains without use cases")

	// Verify percentages are reasonable
	assert.Greater(t, stats.CoveragePercentage, 70.0, "Coverage should be over 70%")
	assert.Less(t, stats.CoveragePercentage, 80.0, "Coverage should be under 80%")

	// Verify average is positive
	assert.Greater(t, stats.AveragePerDomain, 0.0)

	// Verify total use cases
	assert.Greater(t, stats.TotalUseCases, 0)
}

// TestFormatUseCaseStatistics verifies statistics formatting
func TestFormatUseCaseStatistics(t *testing.T) {
	stats := CalculateUseCaseStatistics()
	formatted := FormatUseCaseStatistics(stats)

	assert.Contains(t, formatted, "Use Case Coverage Summary")
	assert.Contains(t, formatted, "Total Domains")
	assert.Contains(t, formatted, "Coverage")
	assert.Contains(t, formatted, "Total Use Cases")
}

// TestGetAllUseCases verifies retrieving all use cases
func TestGetAllUseCases(t *testing.T) {
	useCases := GetAllUseCases()

	// Should have multiple use cases
	assert.Greater(t, len(useCases), 100, "Should have over 100 total use cases")

	// All should have domain and category
	for _, uc := range useCases {
		assert.NotEmpty(t, uc.Domain)
		assert.NotEmpty(t, uc.Description)
		assert.NotEmpty(t, uc.Category)
	}
}

// TestSearchUseCases verifies use case searching
func TestSearchUseCases(t *testing.T) {
	tests := []struct {
		name           string
		keyword        string
		shouldContain  string
		minResultCount int
	}{
		{
			name:           "Search for firewall",
			keyword:        "firewall",
			shouldContain:  "application_firewall",
			minResultCount: 1,
		},
		{
			name:           "Search for configure",
			keyword:        "configure",
			shouldContain:  "authentication",
			minResultCount: 10,
		},
		{
			name:           "Search for manage",
			keyword:        "manage",
			shouldContain:  "dns",
			minResultCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SearchUseCases(tt.keyword)

			assert.GreaterOrEqual(t, len(results), tt.minResultCount,
				"Should find at least %d results for keyword %q", tt.minResultCount, tt.keyword)

			// Verify at least one result matches expected domain
			foundDomain := false
			for _, uc := range results {
				if uc.Domain == tt.shouldContain {
					foundDomain = true
					break
				}
			}
			assert.True(t, foundDomain, "Should find use cases from domain %q", tt.shouldContain)
		})
	}
}

// TestSearchUseCasesEmpty verifies empty keyword returns all use cases
func TestSearchUseCasesEmpty(t *testing.T) {
	all := GetAllUseCases()
	empty := SearchUseCases("")

	assert.Equal(t, len(all), len(empty),
		"Empty search should return all use cases")
}

// TestSearchUseCasesCaseInsensitive verifies case-insensitive searching
func TestSearchUseCasesCaseInsensitive(t *testing.T) {
	lowercase := SearchUseCases("configure")
	uppercase := SearchUseCases("CONFIGURE")
	mixed := SearchUseCases("CoNfIgUrE")

	assert.Equal(t, len(lowercase), len(uppercase),
		"Search should be case-insensitive")
	assert.Equal(t, len(lowercase), len(mixed),
		"Search should be case-insensitive")
}

// TestUseCaseCoverageRatio verifies coverage ratio
func TestUseCaseCoverageRatio(t *testing.T) {
	with := GetDomainsWithUseCases()
	without := GetDomainsWithoutUseCases()

	totalCovered := len(with)
	totalUncovered := len(without)

	assert.Equal(t, 31, totalCovered, "Should have 31 domains with use cases")
	assert.Equal(t, 11, totalUncovered, "Should have 11 domains without use cases")
	assert.Equal(t, 42, totalCovered+totalUncovered, "Total should be 42 domains")
}

// TestSpecificDomainUseCases verifies known domains have expected use cases
func TestSpecificDomainUseCases(t *testing.T) {
	tests := []struct {
		domain          string
		expectedKeyword string
	}{
		{"api", "Discover"},
		{"authentication", "OIDC"},
		{"dns", "load balancing"},
		{"kubernetes", "Kubernetes"},
		{"generative_ai", "AI"},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			info, found := types.GetDomainInfo(tt.domain)
			require.True(t, found, "Domain %q should exist", tt.domain)

			if len(info.UseCases) > 0 {
				found := false
				for _, uc := range info.UseCases {
					if strings.Contains(strings.ToLower(uc), strings.ToLower(tt.expectedKeyword)) {
						found = true
						break
					}
				}
				assert.True(t, found, "Domain %q should have use case containing %q",
					tt.domain, tt.expectedKeyword)
			}
		})
	}
}

// TestUseCaseFormat verifies use case format is consistent
func TestUseCaseFormat(t *testing.T) {
	allUseCases := GetAllUseCases()

	for _, uc := range allUseCases {
		// Each use case should be a complete sentence or phrase
		assert.NotEmpty(t, uc.Description)
		assert.Greater(t, len(uc.Description), 3, "Use case should be meaningful")

		// Should not contain leading/trailing whitespace
		assert.Equal(t, uc.Description, strings.TrimSpace(uc.Description),
			"Use case should not have leading/trailing whitespace")
	}
}

// BenchmarkFormatUseCases benchmarks use case formatting
func BenchmarkFormatUseCases(b *testing.B) {
	useCases := []string{"Configure rules", "Manage policies", "Enable detection"}
	for i := 0; i < b.N; i++ {
		FormatUseCases(useCases)
	}
}

// BenchmarkGetAllUseCases benchmarks retrieving all use cases
func BenchmarkGetAllUseCases(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllUseCases()
	}
}

// BenchmarkCalculateUseCaseStatistics benchmarks statistics calculation
func BenchmarkCalculateUseCaseStatistics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateUseCaseStatistics()
	}
}

// BenchmarkSearchUseCases benchmarks use case searching
func BenchmarkSearchUseCases(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SearchUseCases("configure")
	}
}
