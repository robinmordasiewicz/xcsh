package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// TestDomainsListCommandBasic verifies the domains list command works
func TestDomainsListCommandBasic(t *testing.T) {
	// Create a test command with output buffer
	cmd := domainsListCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	// Execute with no category filter
	cmd.Flags().Set("category", "")
	err := cmd.RunE(cmd, []string{})

	require.NoError(t, err)
	output := buf.String()

	// Verify output contains expected sections
	assert.Contains(t, output, "Available Domains by Category")
	assert.Contains(t, output, "Summary")
	assert.Contains(t, output, "Total")
}

// TestDomainsListByCategory verifies category filtering
func TestDomainsListByCategory(t *testing.T) {
	tests := []struct {
		name             string
		category         string
		expectedDomains  []string
		shouldContain    string
		shouldNotContain string
	}{
		{
			name:             "Security category",
			category:         "Security",
			expectedDomains:  []string{"api", "application_firewall", "certificates"},
			shouldContain:    "Security",
			shouldNotContain: "Networking",
		},
		{
			name:             "Networking category",
			category:         "Networking",
			expectedDomains:  []string{"cdn", "dns", "network"},
			shouldContain:    "Networking",
			shouldNotContain: "Security",
		},
		{
			name:             "AI category",
			category:         "AI",
			expectedDomains:  []string{"generative_ai"},
			shouldContain:    "AI",
			shouldNotContain: "Security",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := domainsListCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			// Set category filter
			cmd.Flags().Set("category", tt.category)
			err := cmd.RunE(cmd, []string{})

			require.NoError(t, err)
			output := buf.String()

			// Verify category is shown
			assert.Contains(t, output, tt.shouldContain)

			// Verify at least one expected domain is shown
			foundExpected := false
			for _, domain := range tt.expectedDomains {
				if bytes.Contains([]byte(output), []byte(domain)) {
					foundExpected = true
					break
				}
			}
			assert.True(t, foundExpected, "Should contain at least one expected domain from category")
		})
	}
}

// TestDomainsListByTierCommand verifies the by-tier subcommand
func TestDomainsListByTierCommand(t *testing.T) {
	cmd := domainsListByTierCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.RunE(cmd, []string{})

	require.NoError(t, err)
	output := buf.String()

	// Verify tier sections are shown
	assert.Contains(t, output, "Standard Tier")
	assert.Contains(t, output, "Professional Tier")
	assert.Contains(t, output, "Enterprise Tier")

	// Verify known standard tier domains
	assert.Contains(t, output, "dns")
	assert.Contains(t, output, "certificates")
}

// TestAllDomainsInCategories verifies all 42 domains are categorized
func TestAllDomainsInCategories(t *testing.T) {
	categories := validation.GetAllCategories()
	assert.Greater(t, len(categories), 0, "Should have at least one category")

	// Count domains across all categories
	totalDomains := 0
	domainSet := make(map[string]bool)

	for _, category := range categories {
		domains := validation.GetDomainsByCategory(category)
		assert.Greater(t, len(domains), 0, "Category %q should have at least one domain", category)

		for _, domain := range domains {
			domainSet[domain.Name] = true
			totalDomains++
		}
	}

	// Should have all 41 domains
	assert.Equal(t, 41, totalDomains, "Should have all 41 domains across categories")

	// All registry domains should be categorized
	for domainName := range types.DomainRegistry {
		assert.True(t, domainSet[domainName], "Domain %q should be in a category", domainName)
	}
}

// TestCategoryConsistency verifies category metadata consistency
func TestCategoryConsistency(t *testing.T) {
	for domainName, domainInfo := range types.DomainRegistry {
		// Every domain should have a category
		assert.NotEmpty(t, domainInfo.Category, "Domain %q should have a category", domainName)

		// Category should be a valid category
		categories := validation.GetAllCategories()
		categoryFound := false
		for _, cat := range categories {
			if domainInfo.Category == cat {
				categoryFound = true
				break
			}
		}
		assert.True(t, categoryFound, "Domain %q has invalid category %q", domainName, domainInfo.Category)
	}
}

// TestCategoryDomainCount verifies specific domain counts per category
func TestCategoryDomainCount(t *testing.T) {
	expectedCounts := map[string]int{
		"AI":             1,
		"Infrastructure": 4,
		"Networking":     5,
		"Operations":     5,
		"Other":          10,
		"Platform":       7,
		"Security":       9,
	}

	for category, expectedCount := range expectedCounts {
		domains := validation.GetDomainsByCategory(category)
		assert.Equal(t, expectedCount, len(domains),
			"Category %q should have %d domains, but has %d",
			category, expectedCount, len(domains))
	}
}

// TestSecurityCategoryDomains verifies specific security domains
func TestSecurityCategoryDomains(t *testing.T) {
	domains := validation.GetDomainsByCategory("Security")
	assert.Equal(t, 9, len(domains), "Security category should have 9 domains")

	domainMap := make(map[string]*types.DomainInfo)
	for _, domain := range domains {
		domainMap[domain.Name] = domain
	}

	// Verify key security domains
	expectedSecurityDomains := []string{
		"api", "application_firewall", "blindfold", "certificates",
		"ddos", "infrastructure_protection", "network_security", "shape", "threat_campaign",
	}

	for _, expected := range expectedSecurityDomains {
		_, found := domainMap[expected]
		assert.True(t, found, "Domain %q should be in Security category", expected)
	}
}

// TestNetworkingCategoryDomains verifies networking category
func TestNetworkingCategoryDomains(t *testing.T) {
	domains := validation.GetDomainsByCategory("Networking")
	assert.Equal(t, 5, len(domains), "Networking category should have 5 domains")

	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedDomains := []string{"cdn", "dns", "network", "rate_limiting", "virtual"}
	for _, expected := range expectedDomains {
		assert.True(t, domainSet[expected], "Domain %q should be in Networking category", expected)
	}
}

// TestPlatformCategoryDomains verifies platform category
func TestPlatformCategoryDomains(t *testing.T) {
	domains := validation.GetDomainsByCategory("Platform")
	assert.Equal(t, 7, len(domains), "Platform category should have 7 domains")

	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	// Platform should contain these domains
	expectedDomains := []string{"authentication", "bigip", "marketplace", "nginx_one", "object_storage", "users", "vpm_and_node_management"}
	for _, expected := range expectedDomains {
		assert.True(t, domainSet[expected], "Domain %q should be in Platform category", expected)
	}
}

// TestOperationsCategoryDomains verifies operations category
func TestOperationsCategoryDomains(t *testing.T) {
	domains := validation.GetDomainsByCategory("Operations")
	assert.Equal(t, 5, len(domains), "Operations category should have 5 domains")

	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedDomains := []string{"data_intelligence", "observability", "statistics", "support", "telemetry_and_insights"}
	for _, expected := range expectedDomains {
		assert.True(t, domainSet[expected], "Domain %q should be in Operations category", expected)
	}
}

// TestInfrastructureCategoryDomains verifies infrastructure category
func TestInfrastructureCategoryDomains(t *testing.T) {
	domains := validation.GetDomainsByCategory("Infrastructure")
	assert.Equal(t, 4, len(domains), "Infrastructure category should have 4 domains")

	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedDomains := []string{"cloud_infrastructure", "kubernetes", "service_mesh", "site"}
	for _, expected := range expectedDomains {
		assert.True(t, domainSet[expected], "Domain %q should be in Infrastructure category", expected)
	}
}

// TestDomainSortingWithinCategories verifies domains are sorted by display name
func TestDomainSortingWithinCategories(t *testing.T) {
	categories := validation.GetAllCategories()

	for _, category := range categories {
		domains := validation.GetDomainsByCategory(category)

		// Verify domains are sorted by display name
		for i := 0; i < len(domains)-1; i++ {
			assert.Less(t, domains[i].DisplayName, domains[i+1].DisplayName,
				"Domains in category %q should be sorted by display name", category)
		}
	}
}

// TestCategoryDistribution verifies distribution calculation
func TestCategoryDistribution(t *testing.T) {
	distribution := validation.GetCategoryDistribution()

	// Should have 7 categories
	assert.Equal(t, 7, len(distribution), "Should have 7 categories")

	// Should be sorted by count descending
	for i := 0; i < len(distribution)-1; i++ {
		assert.GreaterOrEqual(t, distribution[i].Count, distribution[i+1].Count,
			"Categories should be sorted by count descending")
	}

	// Verify total matches registry
	totalCount := 0
	for _, cd := range distribution {
		totalCount += cd.Count
	}
	assert.Equal(t, 41, totalCount, "Total domains should be 41")
}

// TestMultipleCategoryFilter verifies filtering by multiple categories
func TestMultipleCategoryFilter(t *testing.T) {
	categories := []string{"Security", "Platform"}
	domains := validation.GetDomainsInCategories(categories)

	// Should have domains from both categories
	assert.Greater(t, len(domains), 0, "Should have domains from specified categories")

	// All domains should be in one of the specified categories
	categorySet := make(map[string]bool)
	for _, cat := range categories {
		categorySet[cat] = true
	}

	for _, domain := range domains {
		assert.True(t, categorySet[domain.Category],
			"Domain %q should be in one of specified categories", domain.Name)
	}
}

// TestEmptyCategory Filter verifies empty filter returns nil
func TestEmptyCategoryFilter(t *testing.T) {
	domains := validation.GetDomainsInCategories([]string{})
	assert.Nil(t, domains, "Empty category filter should return nil")
}

// TestInvalidCategory verifies handling of non-existent categories
func TestInvalidCategory(t *testing.T) {
	domains := validation.GetDomainsInCategories([]string{"NonExistentCategory"})
	// Should return empty slice (no domains in non-existent category)
	// Note: Go distinguishes between nil slice and empty slice; both are "empty" but nil is returned for validation logic
	if domains != nil {
		assert.Empty(t, domains, "Invalid category should have no domains")
	}
}

// TestCategoryMetadataCompleteness verifies all domains have required metadata
func TestCategoryMetadataCompleteness(t *testing.T) {
	for domainName, domainInfo := range types.DomainRegistry {
		// All should have category
		assert.NotEmpty(t, domainInfo.Category,
			"Domain %q should have category", domainName)

		// All should have at least one of: complexity, display name, description
		assert.NotEmpty(t, domainInfo.DisplayName,
			"Domain %q should have display name", domainName)
		assert.NotEmpty(t, domainInfo.Description,
			"Domain %q should have description", domainName)
	}
}

// BenchmarkGetDomainsByCategory benchmarks category retrieval
func BenchmarkGetDomainsByCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validation.GetDomainsByCategory("Security")
	}
}

// BenchmarkGroupDomainsByCategory benchmarks grouping operation
func BenchmarkGroupDomainsByCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validation.GroupDomainsByCategory()
	}
}
