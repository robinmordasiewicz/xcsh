package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// TestGetRelatedDomains verifies related domains retrieval
func TestGetRelatedDomains(t *testing.T) {
	// Test API domain
	related := GetRelatedDomains("api")
	assert.Greater(t, len(related), 0, "API should have related domains")

	// Check that security-related domains are included
	relatedNames := make(map[string]bool)
	for _, domain := range related {
		relatedNames[domain.Name] = true
	}

	// API should be related to other security domains
	assert.True(t, len(relatedNames) > 0, "Should have at least one related domain")

	// Verify all returned domains are valid
	for _, domain := range related {
		assert.NotEmpty(t, domain.Name)
		assert.NotEmpty(t, domain.Category)
	}
}

// TestGetRelatedDomainsNonExistent tests non-existent domain
func TestGetRelatedDomainsNonExistent(t *testing.T) {
	related := GetRelatedDomains("nonexistent_domain")
	assert.Nil(t, related, "Non-existent domain should return nil")
}

// TestGetRelatedDomainsCategoryGrouping verifies same category domains are related
func TestGetRelatedDomainsCategoryGrouping(t *testing.T) {
	// Test security domain (has same category domains)
	related := GetRelatedDomains("waf")
	assert.Greater(t, len(related), 0, "Security domain should have related domains")

	// At least some should be in same category
	wafInfo, _ := types.GetDomainInfo("waf")
	sameCategoryFound := false
	for _, domain := range related {
		if domain.Category == wafInfo.Category {
			sameCategoryFound = true
			break
		}
	}
	assert.True(t, sameCategoryFound, "Should have same-category related domains")
}

// TestGetWorkflowSuggestions verifies workflow suggestions
func TestGetWorkflowSuggestions(t *testing.T) {
	tests := []struct {
		domain        string
		shouldHaveMin int
		expectedNames []string
	}{
		{
			domain:        "api",
			shouldHaveMin: 1,
			expectedNames: []string{"API Security Workflow", "Network Protection Workflow"},
		},
		{
			domain:        "authentication",
			shouldHaveMin: 1,
			expectedNames: []string{"Access Management Workflow"},
		},
		{
			domain:        "kubernetes",
			shouldHaveMin: 1,
			expectedNames: []string{"Kubernetes Management Workflow"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			suggestions := GetWorkflowSuggestions(tt.domain)
			assert.Greater(t, len(suggestions), tt.shouldHaveMin-1,
				"Domain %q should have workflow suggestions", tt.domain)

			// Verify structure of suggestions
			for _, suggestion := range suggestions {
				assert.NotEmpty(t, suggestion.Name)
				assert.NotEmpty(t, suggestion.Description)
				assert.Greater(t, len(suggestion.Domains), 0,
					"Workflow should include domains")
			}
		})
	}
}

// TestGetWorkflowSuggestionsNonExistent tests non-existent domain
func TestGetWorkflowSuggestionsNonExistent(t *testing.T) {
	suggestions := GetWorkflowSuggestions("nonexistent")
	assert.Empty(t, suggestions, "Non-existent domain should return empty suggestions")
}

// TestFormatRelatedDomains verifies formatting
func TestFormatRelatedDomains(t *testing.T) {
	related := GetRelatedDomains("api")
	require.Greater(t, len(related), 0)

	formatted := FormatRelatedDomains(related)
	assert.Contains(t, formatted, "RELATED DOMAINS:")
	assert.Contains(t, formatted, "•")
	assert.Greater(t, strings.Count(formatted, "\n"), 1)
}

// TestFormatRelatedDomainsEmpty verifies empty case
func TestFormatRelatedDomainsEmpty(t *testing.T) {
	formatted := FormatRelatedDomains(nil)
	assert.Empty(t, formatted)
}

// TestFormatWorkflowSuggestions verifies formatting
func TestFormatWorkflowSuggestions(t *testing.T) {
	suggestions := GetWorkflowSuggestions("api")
	require.Greater(t, len(suggestions), 0)

	formatted := FormatWorkflowSuggestions(suggestions)
	assert.Contains(t, formatted, "SUGGESTED WORKFLOWS:")
	assert.Contains(t, formatted, "Workflow")
	assert.Greater(t, strings.Count(formatted, "\n"), 1)
}

// TestFormatWorkflowSuggestionsEmpty verifies empty case
func TestFormatWorkflowSuggestionsEmpty(t *testing.T) {
	formatted := FormatWorkflowSuggestions(nil)
	assert.Empty(t, formatted)
}

// TestGetWorkflowsByCategory verifies category grouping
func TestGetWorkflowsByCategory(t *testing.T) {
	// Test security category
	workflows := GetWorkflowsByCategory("Security")
	assert.Greater(t, len(workflows), 0, "Security category should have workflows")

	// Verify all workflows are for that category
	for _, workflow := range workflows {
		assert.Equal(t, "Security", workflow.Category,
			"Workflow %q should be in Security category", workflow.Name)
	}
}

// TestGetWorkflowsByCategoryNonExistent tests non-existent category
func TestGetWorkflowsByCategoryNonExistent(t *testing.T) {
	workflows := GetWorkflowsByCategory("NonExistent")
	assert.Empty(t, workflows, "Non-existent category should have no workflows")
}

// TestWorkflowCoverage verifies workflows cover key domains
func TestWorkflowCoverage(t *testing.T) {
	keyDomains := []string{"api", "kubernetes", "authentication", "dns"}

	for _, domain := range keyDomains {
		suggestions := GetWorkflowSuggestions(domain)
		assert.Greater(t, len(suggestions), 0,
			"Key domain %q should have workflow suggestions", domain)
	}
}

// TestRelatedDomainsNotSelf verifies domain is not related to itself
func TestRelatedDomainsNotSelf(t *testing.T) {
	related := GetRelatedDomains("api")

	for _, domain := range related {
		assert.NotEqual(t, "api", domain.Name,
			"Domain should not be related to itself")
	}
}

// TestWorkflowSuggestionsStructure verifies workflow structure
func TestWorkflowSuggestionsStructure(t *testing.T) {
	suggestions := GetWorkflowSuggestions("api")

	for _, suggestion := range suggestions {
		// Verify all fields populated
		assert.NotEmpty(t, suggestion.Name, "Workflow name required")
		assert.NotEmpty(t, suggestion.Description, "Workflow description required")
		assert.Greater(t, len(suggestion.Domains), 0, "Workflow must include domains")
		assert.NotEmpty(t, suggestion.Category, "Workflow category required")

		// Verify domains exist
		for _, domainName := range suggestion.Domains {
			_, found := types.GetDomainInfo(domainName)
			assert.True(t, found, "Domain %q in workflow should exist", domainName)
		}
	}
}

// BenchmarkGetRelatedDomains benchmarks related domains retrieval
func BenchmarkGetRelatedDomains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRelatedDomains("api")
	}
}

// BenchmarkGetWorkflowSuggestions benchmarks workflow suggestions
func BenchmarkGetWorkflowSuggestions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetWorkflowSuggestions("api")
	}
}

// BenchmarkFormatRelatedDomains benchmarks formatting
func BenchmarkFormatRelatedDomains(b *testing.B) {
	related := GetRelatedDomains("api")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatRelatedDomains(related)
	}
}
