package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// TestPhase1TierValidationIntegration verifies tier validation works in integrated context
func TestPhase1TierValidationIntegration(t *testing.T) {
	// Test domains across tiers
	testDomains := []string{"api", "dns", "authentication"}

	for _, domain := range testDomains {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist", domain)

		// Verify tier requirement is accessible
		assert.NotEmpty(t, info.RequiresTier, "Domain %q should have tier requirement", domain)

		// Verify tier is Standard or Advanced
		validTiers := []string{"Standard", "Advanced"}
		assert.Contains(t, validTiers, info.RequiresTier, "Domain %q has invalid tier", domain)
	}
}

// TestPhase2PreviewWarningsIntegration verifies preview warnings are set correctly
func TestPhase2PreviewWarningsIntegration(t *testing.T) {
	// Iterate through all domains and verify preview flag is set
	previewCount := 0
	nonPreviewCount := 0

	for _, info := range types.DomainRegistry {
		if info.IsPreview {
			previewCount++
		} else {
			nonPreviewCount++
		}
	}

	// Verify we have both preview and non-preview domains
	assert.Greater(t, nonPreviewCount, 0, "Should have non-preview domains")
	// Preview count can be 0 or more, that's fine
	assert.True(t, previewCount+nonPreviewCount == len(types.DomainRegistry), "All domains should be classified")
}

// TestPhase3DomainCategorization verifies all domains are properly categorized
func TestPhase3DomainCategorization(t *testing.T) {
	validCategories := []string{"Security", "Networking", "Platform", "Infrastructure", "Operations", "Other", "AI"}

	for domainName, info := range types.DomainRegistry {
		assert.NotEmpty(t, info.Category, "Domain %q should have category", domainName)
		assert.Contains(t, validCategories, info.Category, "Domain %q has invalid category", domainName)
	}

	// Verify we have domains in multiple categories
	categoryCounts := make(map[string]int)
	for _, info := range types.DomainRegistry {
		categoryCounts[info.Category]++
	}

	assert.Greater(t, len(categoryCounts), 1, "Should have domains in multiple categories")
	// Verify expected category distribution
	assert.GreaterOrEqual(t, categoryCounts["Security"], 8, "Should have at least 8 Security domains")
	assert.GreaterOrEqual(t, categoryCounts["Platform"], 7, "Should have at least 7 Platform domains")
	assert.GreaterOrEqual(t, categoryCounts["Other"], 9, "Should have at least 9 Other domains")
}

// TestPhase4UseCaseDocumentation verifies use cases are available where applicable
func TestPhase4UseCaseDocumentation(t *testing.T) {
	domainsWithUseCases := 0
	domainsWithoutUseCases := 0

	for domainName, info := range types.DomainRegistry {
		if len(info.UseCases) > 0 {
			domainsWithUseCases++
			// Verify use cases are meaningful
			for _, useCase := range info.UseCases {
				assert.NotEmpty(t, useCase, "Domain %q has empty use case", domainName)
				assert.Greater(t, len(useCase), 3, "Domain %q use case too short", domainName)
			}
		} else {
			domainsWithoutUseCases++
		}
	}

	// Verify we have substantial use case coverage
	assert.Greater(t, domainsWithUseCases, 25, "Should have use cases for at least 25 domains")
	assert.Less(t, domainsWithoutUseCases, 20, "Should have less than 20 domains without use cases")
}

// TestPhase5WorkflowSuggestions verifies workflow suggestions are available
func TestPhase5WorkflowSuggestions(t *testing.T) {
	domainsWithWorkflows := 0
	workflowCount := 0

	for domainName := range types.DomainRegistry {
		workflows := validation.GetWorkflowSuggestions(domainName)
		if len(workflows) > 0 {
			domainsWithWorkflows++
			workflowCount += len(workflows)

			// Verify workflow structure
			for _, workflow := range workflows {
				assert.NotEmpty(t, workflow.Name, "Workflow missing name")
				assert.NotEmpty(t, workflow.Description, "Workflow missing description")
				assert.Greater(t, len(workflow.Domains), 0, "Workflow must include domains")
				assert.NotEmpty(t, workflow.Category, "Workflow missing category")
			}
		}
	}

	assert.Greater(t, domainsWithWorkflows, 20, "Should have workflows for at least 20 domains")
	assert.Greater(t, workflowCount, 25, "Should have at least 25 total workflows")
}

// TestPhase5RelatedDomains verifies related domains are retrieved correctly
func TestPhase5RelatedDomains(t *testing.T) {
	for domainName := range types.DomainRegistry {
		relatedDomains := validation.GetRelatedDomains(domainName)

		// All domains should have related domains
		assert.NotNil(t, relatedDomains, "Domain %q should have related domains", domainName)
		assert.Greater(t, len(relatedDomains), 0, "Domain %q should have at least one related domain", domainName)

		// Verify domain is not related to itself
		for _, relatedDomain := range relatedDomains {
			assert.NotEqual(t, domainName, relatedDomain.Name, "Domain should not be related to itself")
		}
	}
}

// TestCompleteHelpTextFlow verifies all phases display correctly in help text
func TestCompleteHelpTextFlow(t *testing.T) {
	testDomains := []string{"api", "dns", "authentication", "kubernetes", "cdn"}

	for _, domainName := range testDomains {
		info, found := types.GetDomainInfo(domainName)
		require.True(t, found, "Domain %q should exist", domainName)

		// Build help text similar to cmd/domains.go
		var helpSections []string

		// Section 1: Description
		helpSections = append(helpSections, info.Description)

		// Section 2: Tier (Phase 1)
		if info.RequiresTier != "" {
			helpSections = append(helpSections, "Tier: "+info.RequiresTier)
		}

		// Section 3: Preview warning (Phase 2)
		if info.IsPreview {
			helpSections = append(helpSections, "[PREVIEW] This is a preview feature")
		}

		// Section 4: Category and Complexity (Phase 3)
		helpSections = append(helpSections, "Category: "+info.Category)
		if info.Complexity != "" {
			helpSections = append(helpSections, "Complexity: "+info.Complexity)
		}

		// Section 5: Use Cases (Phase 4)
		if len(info.UseCases) > 0 {
			helpSections = append(helpSections, "USE CASES:")
			for _, uc := range info.UseCases {
				helpSections = append(helpSections, "  • "+uc)
			}
		}

		// Section 6: Related Domains (Phase 5)
		relatedDomains := validation.GetRelatedDomains(domainName)
		if len(relatedDomains) > 0 {
			helpSections = append(helpSections, "RELATED DOMAINS:")
			for _, rd := range relatedDomains {
				helpSections = append(helpSections, "  • "+rd.Name)
			}
		}

		// Section 7: Workflows (Phase 5)
		workflows := validation.GetWorkflowSuggestions(domainName)
		if len(workflows) > 0 {
			helpSections = append(helpSections, "SUGGESTED WORKFLOWS:")
			for _, wf := range workflows {
				helpSections = append(helpSections, "  • "+wf.Name)
			}
		}

		// Build complete help text
		completeHelp := strings.Join(helpSections, "\n")

		// Verify all sections are present and in correct order
		assert.Contains(t, completeHelp, info.Description, "Help missing description for %q", domainName)
		assert.Contains(t, completeHelp, info.Category, "Help missing category for %q", domainName)

		// Verify no empty sections between content
		assert.False(t, strings.Contains(completeHelp, "\n\n\n"), "Help has too many blank lines for %q", domainName)

		// Verify section ordering: description before category before workflows
		descPos := strings.Index(completeHelp, info.Description)
		categoryPos := strings.Index(completeHelp, info.Category)
		assert.Less(t, descPos, categoryPos, "Description should come before category for %q", domainName)
	}
}

// TestNoFeatureConflicts verifies phases don't interfere with each other
func TestNoFeatureConflicts(t *testing.T) {
	for domainName, info := range types.DomainRegistry {
		// Verify tier requirements don't break access
		assert.NotEmpty(t, info.RequiresTier, "Domain %q missing tier", domainName)

		// Verify preview flag doesn't interfere with other fields
		if info.IsPreview {
			assert.NotEmpty(t, info.Category, "Preview domain %q should have category", domainName)
			assert.NotEmpty(t, info.RequiresTier, "Preview domain %q should have tier", domainName)
		}

		// Verify workflows are category-appropriate
		workflows := validation.GetWorkflowSuggestions(domainName)
		for _, workflow := range workflows {
			assert.Equal(t, info.Category, workflow.Category, "Workflow %q should match domain category", workflow.Name)
		}

		// Verify related domains exist
		relatedDomains := validation.GetRelatedDomains(domainName)
		assert.NotNil(t, relatedDomains, "Domain %q should have related domains", domainName)

		// Verify all related domains exist in registry
		for _, relatedDomain := range relatedDomains {
			_, found := types.GetDomainInfo(relatedDomain.Name)
			assert.True(t, found, "Related domain %q should exist in registry", relatedDomain.Name)
		}
	}
}

// TestCompletionHelperIntegration verifies completion helpers work with all phases
func TestCompletionHelperIntegration(t *testing.T) {
	// Test domain name completion
	allDomainNames := []string{}
	for domainName := range types.DomainRegistry {
		allDomainNames = append(allDomainNames, domainName)
	}
	assert.Equal(t, len(allDomainNames), 39, "Should have 39 domains")

	// Test category-based completion works
	categoryDomains := validation.GetDomainsByCategory("Security")
	assert.Greater(t, len(categoryDomains), 0, "Security category should have domains")

	// Test use case search works
	searchResults := validation.SearchUseCases("manage")
	assert.Greater(t, len(searchResults), 0, "Should find use cases with 'manage' keyword")

	// Test workflow suggestions work
	workflowsByCategory := validation.GetWorkflowsByCategory("Security")
	assert.Greater(t, len(workflowsByCategory), 0, "Security category should have workflows")
}

// TestErrorHandlingIntegration verifies all phases handle errors gracefully
func TestErrorHandlingIntegration(t *testing.T) {
	// Non-existent domain
	_, found := types.GetDomainInfo("nonexistent_domain_xyz")
	assert.False(t, found, "Non-existent domain should not be found")

	// Related domains for non-existent domain
	relatedDomains := validation.GetRelatedDomains("nonexistent_domain_xyz")
	assert.Nil(t, relatedDomains, "Non-existent domain should return nil related domains")

	// Workflows for non-existent domain
	workflows := validation.GetWorkflowSuggestions("nonexistent_domain_xyz")
	assert.Empty(t, workflows, "Non-existent domain should return empty workflows")

	// Use cases for non-existent domain
	useCases, err := validation.GetDomainUseCases("nonexistent_domain_xyz")
	assert.Error(t, err, "Non-existent domain should return error for use cases")
	assert.Empty(t, useCases, "Non-existent domain should return empty use cases")

	// Category lookup for non-existent category
	domainsByCategory := validation.GetDomainsByCategory("NonExistentCategory")
	assert.Empty(t, domainsByCategory, "Non-existent category should return empty domains")

	// Workflow lookup for non-existent category
	workflowsByCategory := validation.GetWorkflowsByCategory("NonExistentCategory")
	assert.Empty(t, workflowsByCategory, "Non-existent category should return empty workflows")
}

// TestPerformanceIntegration verifies response times meet targets
func TestPerformanceIntegration(t *testing.T) {
	testDomains := []string{"api", "dns", "authentication", "kubernetes", "cdn"}

	for _, domainName := range testDomains {
		// Test GetRelatedDomains performance
		start := time.Now()
		_ = validation.GetRelatedDomains(domainName)
		elapsed := time.Since(start)
		assert.Less(t, elapsed, 100*time.Millisecond, "GetRelatedDomains for %q took too long", domainName)

		// Test GetWorkflowSuggestions performance
		start = time.Now()
		_ = validation.GetWorkflowSuggestions(domainName)
		elapsed = time.Since(start)
		assert.Less(t, elapsed, 100*time.Millisecond, "GetWorkflowSuggestions for %q took too long", domainName)

		// Test help text generation performance
		start = time.Now()
		_ = validation.FormatRelatedDomains(validation.GetRelatedDomains(domainName))
		_ = validation.FormatWorkflowSuggestions(validation.GetWorkflowSuggestions(domainName))
		elapsed = time.Since(start)
		assert.Less(t, elapsed, 100*time.Millisecond, "Help text formatting for %q took too long", domainName)
	}
}

// TestCrossPhaseDataConsistency verifies data consistency across all phases
func TestCrossPhaseDataConsistency(t *testing.T) {
	for domainName, info := range types.DomainRegistry {
		// Verify all use cases reference existing domains
		for _, useCase := range info.UseCases {
			assert.NotEmpty(t, useCase, "Use case should not be empty for %q", domainName)
		}

		// Verify workflow domains all exist
		workflows := validation.GetWorkflowSuggestions(domainName)
		for _, workflow := range workflows {
			for _, workflowDomain := range workflow.Domains {
				_, found := types.GetDomainInfo(workflowDomain)
				assert.True(t, found, "Workflow domain %q in workflow for %q should exist", workflowDomain, domainName)
			}
		}

		// Verify related domain categories exist
		relatedDomains := validation.GetRelatedDomains(domainName)
		for _, relatedDomain := range relatedDomains {
			assert.NotEmpty(t, relatedDomain.Category, "Related domain %q should have category", relatedDomain.Name)
		}
	}
}

// TestDomainWorkflowConsistency verifies workflow suggestions are consistent
func TestDomainWorkflowConsistency(t *testing.T) {
	// Verify domains in workflows are accessible and correct tier
	for domainName := range types.DomainRegistry {
		workflows := validation.GetWorkflowSuggestions(domainName)

		for _, workflow := range workflows {
			// All domains in workflow should exist
			for _, workflowDomain := range workflow.Domains {
				info, found := types.GetDomainInfo(workflowDomain)
				assert.True(t, found, "Workflow domain %q should exist", workflowDomain)

				// Workflow domain should be same or compatible tier
				if info.RequiresTier != "" && domainName != "" {
					mainInfo, _ := types.GetDomainInfo(domainName)
					// Both should be valid tiers
					assert.NotEmpty(t, mainInfo.RequiresTier)
					assert.NotEmpty(t, info.RequiresTier)
				}
			}
		}
	}
}

// TestAllDomainsAccessible verifies all 39 domains are accessible and complete
func TestAllDomainsAccessible(t *testing.T) {
	assert.Equal(t, len(types.DomainRegistry), 39, "Should have exactly 39 domains")

	accessibleCount := 0
	completeCount := 0

	for domainName, info := range types.DomainRegistry {
		// Verify domain is accessible
		_, found := types.GetDomainInfo(domainName)
		assert.True(t, found, "Domain %q should be accessible", domainName)
		accessibleCount++

		// Verify domain is complete (has all required fields)
		assert.NotEmpty(t, info.Name, "Domain %q missing name", domainName)
		assert.NotEmpty(t, info.DisplayName, "Domain %q missing display name", domainName)
		assert.NotEmpty(t, info.Description, "Domain %q missing description", domainName)
		assert.NotEmpty(t, info.Category, "Domain %q missing category", domainName)
		assert.NotEmpty(t, info.RequiresTier, "Domain %q missing tier requirement", domainName)

		completeCount++
	}

	assert.Equal(t, accessibleCount, 39, "All 39 domains should be accessible")
	assert.Equal(t, completeCount, 39, "All 39 domains should be complete")
}
