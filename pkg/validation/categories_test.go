package validation

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// TestGetAllCategories verifies all categories are returned
func TestGetAllCategories(t *testing.T) {
	categories := GetAllCategories()

	// Should return multiple categories
	assert.Greater(t, len(categories), 0, "Should have at least one category")

	// Should be sorted
	sorted := make([]string, len(categories))
	copy(sorted, categories)
	sort.Strings(sorted)
	assert.Equal(t, sorted, categories, "Categories should be sorted alphabetically")

	// All entries should be non-empty
	for _, cat := range categories {
		assert.NotEmpty(t, cat, "Category should not be empty")
	}
}

// TestGetDomainsByCategory verifies domains are correctly grouped by category
func TestGetDomainsByCategory(t *testing.T) {
	categories := GetAllCategories()

	for _, cat := range categories {
		domains := GetDomainsByCategory(cat)

		// Should have at least one domain per category
		assert.Greater(t, len(domains), 0, "Category %q should have at least one domain", cat)

		// All domains should belong to this category
		for _, domain := range domains {
			assert.Equal(t, cat, domain.Category, "Domain %q should belong to category %q", domain.Name, cat)
		}

		// Should be sorted by display name
		sorted := make([]*types.DomainInfo, len(domains))
		copy(sorted, domains)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].DisplayName < sorted[j].DisplayName
		})
		assert.Equal(t, sorted, domains, "Domains should be sorted by display name in category %q", cat)
	}
}

// TestGroupDomainsByCategory verifies domain grouping structure
func TestGroupDomainsByCategory(t *testing.T) {
	groupings := GroupDomainsByCategory()

	// Should have multiple groupings
	assert.Greater(t, len(groupings), 0, "Should have at least one category grouping")

	// Categories should be sorted
	categories := make([]string, 0, len(groupings))
	for _, g := range groupings {
		categories = append(categories, g.Category)
	}
	sorted := make([]string, len(categories))
	copy(sorted, categories)
	sort.Strings(sorted)
	assert.Equal(t, sorted, categories, "Categories in groupings should be sorted")

	// Each grouping should have at least one domain
	for _, g := range groupings {
		assert.Greater(t, len(g.Domains), 0, "Category %q should have at least one domain", g.Category)

		// Domains within grouping should be sorted
		sorted := make([]*types.DomainInfo, len(g.Domains))
		copy(sorted, g.Domains)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].DisplayName < sorted[j].DisplayName
		})
		assert.Equal(t, sorted, g.Domains, "Domains in category %q should be sorted", g.Category)
	}
}

// TestGetDomainsInCategories verifies filtering by multiple categories
func TestGetDomainsInCategories(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		expectNil  bool
	}{
		{
			name:       "Empty filter",
			categories: []string{},
			expectNil:  true,
		},
		{
			name:       "Single category",
			categories: []string{"Security"},
			expectNil:  false,
		},
		{
			name:       "Multiple categories",
			categories: []string{"Security", "Platform"},
			expectNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains := GetDomainsInCategories(tt.categories)

			if tt.expectNil {
				assert.Nil(t, domains, "Should return nil for empty filter")
			} else {
				assert.NotNil(t, domains, "Should return domains")
				assert.Greater(t, len(domains), 0, "Should have at least one domain")

				// All returned domains should be in one of the filter categories
				filterSet := make(map[string]bool)
				for _, cat := range tt.categories {
					filterSet[cat] = true
				}
				for _, domain := range domains {
					assert.True(t, filterSet[domain.Category], "Domain %q should be in filtered categories", domain.Name)
				}

				// Domains should be sorted by display name
				sorted := make([]*types.DomainInfo, len(domains))
				copy(sorted, domains)
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i].DisplayName < sorted[j].DisplayName
				})
				assert.Equal(t, sorted, domains, "Domains should be sorted by display name")
			}
		})
	}
}

// TestGetCategoryDistribution verifies category counts
func TestGetCategoryDistribution(t *testing.T) {
	distribution := GetCategoryDistribution()

	// Should have multiple categories
	assert.Greater(t, len(distribution), 0, "Should have at least one category")

	// Count should match actual domain count
	totalCount := 0
	for _, cd := range distribution {
		assert.Greater(t, cd.Count, 0, "Category %q should have at least one domain", cd.Category)
		totalCount += cd.Count
	}
	assert.Equal(t, len(types.DomainRegistry), totalCount, "Total count should match registry size")

	// Should be sorted by count descending
	for i := 0; i < len(distribution)-1; i++ {
		if distribution[i].Count != distribution[i+1].Count {
			assert.Greater(t, distribution[i].Count, distribution[i+1].Count,
				"Categories should be sorted by count descending")
		}
	}
}

// TestAllDomainsInCategories verifies all domains are in some category
func TestAllDomainsInCategories(t *testing.T) {
	groupings := GroupDomainsByCategory()

	// Collect all domains in groupings
	domainSet := make(map[string]bool)
	for _, g := range groupings {
		for _, domain := range g.Domains {
			domainSet[domain.Name] = true
		}
	}

	// All registry domains should be present
	for domainName := range types.DomainRegistry {
		assert.True(t, domainSet[domainName], "Domain %q should be in a category", domainName)
	}
}

// TestCategoryNames verifies expected category names exist
func TestCategoryNames(t *testing.T) {
	categories := GetAllCategories()
	categorySet := make(map[string]bool)
	for _, cat := range categories {
		categorySet[cat] = true
	}

	// Expected categories based on domain metadata
	expectedCategories := []string{"AI", "Infrastructure", "Networking", "Operations", "Other", "Platform", "Security"}

	for _, expected := range expectedCategories {
		assert.True(t, categorySet[expected], "Category %q should exist", expected)
	}
}

// TestSecurityCategoryDomains verifies specific security domains
func TestSecurityCategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("Security")
	assert.NotEmpty(t, domains, "Security category should have domains")

	// Check for known security domains
	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedSecurityDomains := []string{"api", "application_firewall", "certificates", "ddos", "network_security", "shape", "threat_campaign"}
	for _, expected := range expectedSecurityDomains {
		assert.True(t, domainSet[expected], "Security category should contain domain %q", expected)
	}
}

// TestPlatformCategoryDomains verifies specific platform domains
func TestPlatformCategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("Platform")
	assert.NotEmpty(t, domains, "Platform category should have domains")

	// Check for known platform domains
	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedPlatformDomains := []string{"authentication", "bigip", "marketplace", "users"}
	for _, expected := range expectedPlatformDomains {
		assert.True(t, domainSet[expected], "Platform category should contain domain %q", expected)
	}
}

// TestAICategoryDomains verifies AI category
func TestAICategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("AI")
	assert.NotEmpty(t, domains, "AI category should have domains")

	// AI should have exactly one domain: generative_ai
	assert.Equal(t, 1, len(domains), "AI category should have exactly one domain")
	assert.Equal(t, "generative_ai", domains[0].Name)
}

// TestOperationsCategoryDomains verifies specific operations domains
func TestOperationsCategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("Operations")
	assert.NotEmpty(t, domains, "Operations category should have domains")

	// Check for known operations domains
	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedOperationsDomains := []string{"data_intelligence", "observability", "statistics", "telemetry_and_insights"}
	for _, expected := range expectedOperationsDomains {
		assert.True(t, domainSet[expected], "Operations category should contain domain %q", expected)
	}
}

// TestNetworkingCategoryDomains verifies specific networking domains
func TestNetworkingCategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("Networking")
	assert.NotEmpty(t, domains, "Networking category should have domains")

	// Check for known networking domains
	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedNetworkingDomains := []string{"cdn", "dns", "network", "rate_limiting", "virtual"}
	for _, expected := range expectedNetworkingDomains {
		assert.True(t, domainSet[expected], "Networking category should contain domain %q", expected)
	}
}

// TestInfrastructureCategoryDomains verifies specific infrastructure domains
func TestInfrastructureCategoryDomains(t *testing.T) {
	domains := GetDomainsByCategory("Infrastructure")
	assert.NotEmpty(t, domains, "Infrastructure category should have domains")

	// Check for known infrastructure domains
	domainSet := make(map[string]bool)
	for _, domain := range domains {
		domainSet[domain.Name] = true
	}

	expectedInfrastructureDomains := []string{"cloud_infrastructure", "kubernetes", "service_mesh", "site"}
	for _, expected := range expectedInfrastructureDomains {
		assert.True(t, domainSet[expected], "Infrastructure category should contain domain %q", expected)
	}
}

// TestDistributionTotals verifies distribution totals
func TestDistributionTotals(t *testing.T) {
	distribution := GetCategoryDistribution()

	// Expected distribution based on domain count analysis
	expectedTotals := map[string]int{
		"Other":          10,
		"Security":       9,
		"Platform":       7,
		"Operations":     5,
		"Networking":     5,
		"Infrastructure": 4,
		"AI":             1,
	}

	actualTotals := make(map[string]int)
	for _, cd := range distribution {
		actualTotals[cd.Category] = cd.Count
	}

	for category, expectedCount := range expectedTotals {
		actualCount := actualTotals[category]
		assert.Equal(t, expectedCount, actualCount,
			"Category %q should have %d domains, but has %d", category, expectedCount, actualCount)
	}
}

// BenchmarkGetAllCategories benchmarks category retrieval
func BenchmarkGetAllCategories(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllCategories()
	}
}

// BenchmarkGetDomainsByCategory benchmarks domain retrieval by category
func BenchmarkGetDomainsByCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetDomainsByCategory("Security")
	}
}

// BenchmarkGroupDomainsByCategory benchmarks full grouping operation
func BenchmarkGroupDomainsByCategory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GroupDomainsByCategory()
	}
}

// BenchmarkGetCategoryDistribution benchmarks distribution calculation
func BenchmarkGetCategoryDistribution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetCategoryDistribution()
	}
}
