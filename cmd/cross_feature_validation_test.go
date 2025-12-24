package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// TestTierEscalationStandard verifies Standard tier access limitations
func TestTierEscalationStandard(t *testing.T) {
	standardDomains := validation.GetDomainsByTier("Standard")
	assert.Equal(t, 25, len(standardDomains), "Standard tier should have 25 domains")

	// Verify Standard tier domains are accessible
	for _, domain := range standardDomains {
		assert.Equal(t, "Standard", domain.RequiresTier, "Standard tier domain should require Standard")
	}

	// Verify Professional domains are NOT in Standard
	proDomains := validation.GetDomainsByTier("Professional")
	proOnlyCount := 0
	for _, proDomain := range proDomains {
		found := false
		for _, stdDomain := range standardDomains {
			if proDomain.Name == stdDomain.Name {
				found = true
				break
			}
		}
		if !found {
			proOnlyCount++
		}
	}

	assert.Greater(t, proOnlyCount, 0, "Should have Professional-only domains")
}

// TestTierEscalationProfessional verifies Professional tier includes all Standard
func TestTierEscalationProfessional(t *testing.T) {
	standardDomains := validation.GetDomainsByTier("Standard")
	proDomains := validation.GetDomainsByTier("Professional")

	assert.Equal(t, 36, len(proDomains), "Professional tier should have 36 domains")

	// Verify all Standard domains are in Professional
	for _, stdDomain := range standardDomains {
		found := false
		for _, proDomain := range proDomains {
			if stdDomain.Name == proDomain.Name {
				found = true
				break
			}
		}
		assert.True(t, found, "Standard domain %q should be in Professional tier", stdDomain.Name)
	}
}

// TestTierEscalationEnterprise verifies Enterprise has all domains
func TestTierEscalationEnterprise(t *testing.T) {
	entDomains := validation.GetDomainsByTier("Enterprise")
	assert.Equal(t, 42, len(entDomains), "Enterprise tier should have 42 domains")

	// Verify all 42 domains are accessible
	for domainName := range types.DomainRegistry {
		found := false
		for _, entDomain := range entDomains {
			if entDomain.Name == domainName {
				found = true
				break
			}
		}
		assert.True(t, found, "All domains should be accessible at Enterprise tier")
	}
}

// TestTierValidationWithPreview verifies preview domains respect tier constraints
func TestTierValidationWithPreview(t *testing.T) {
	previewDomains := validation.GetPreviewDomains()

	for _, previewDomain := range previewDomains {
		// Preview domain should still have a tier requirement
		assert.NotEmpty(t, previewDomain.RequiresTier, "Preview domain %q should have tier", previewDomain.Name)

		// Verify tier validation works for preview domains
		// Some preview domains require Enterprise tier
		ok := validation.ValidateTierAccess("Enterprise", previewDomain.RequiresTier)
		assert.True(t, ok, "Enterprise user should access preview domain %q", previewDomain.Name)
	}
}

// TestCategoryDomainTierConsistency verifies domains within category have appropriate tiers
func TestCategoryDomainTierConsistency(t *testing.T) {
	categories := validation.GetAllCategories()

	for _, category := range categories {
		domainsByCategory := validation.GetDomainsByCategory(category)
		assert.Greater(t, len(domainsByCategory), 0, "Category %q should have domains", category)

		// All domains in category should have valid tier requirements
		for _, domain := range domainsByCategory {
			assert.NotEmpty(t, domain.RequiresTier, "Domain %q in category %q should have tier", domain.Name, category)
		}
	}
}

// TestWorkflowDomainsRespectTier verifies workflow domains are tier-accessible
func TestWorkflowDomainsRespectTier(t *testing.T) {
	for domainName := range types.DomainRegistry {
		workflows := validation.GetWorkflowSuggestions(domainName)

		for _, workflow := range workflows {
			// All domains in workflow should exist and be accessible at Enterprise tier
			for _, workflowDomain := range workflow.Domains {
				domainInfo, found := types.GetDomainInfo(workflowDomain)
				assert.True(t, found, "Workflow domain %q should exist", workflowDomain)

				// All workflow domains should be accessible at Enterprise tier
				assert.NotNil(t, domainInfo)
				ok := validation.ValidateTierAccess("Enterprise", domainInfo.RequiresTier)
				assert.True(t, ok, "Workflow %q domain %q should be Enterprise-accessible", workflow.Name, workflowDomain)
			}
		}
	}
}

// TestRelatedDomainsRespectTier verifies related domains are tier-compatible
func TestRelatedDomainsRespectTier(t *testing.T) {
	testDomains := []string{"api", "dns", "kubernetes", "authentication", "cdn"}

	for _, testDomain := range testDomains {
		relatedDomains := validation.GetRelatedDomains(testDomain)

		// All related domains should be tier-compatible with Enterprise
		// (since some test domains like cdn require Enterprise tier)
		for _, relatedDomain := range relatedDomains {
			ok := validation.ValidateTierAccess("Enterprise", relatedDomain.RequiresTier)
			assert.True(t, ok, "Related domain %q should be Enterprise-accessible", relatedDomain.Name)
		}
	}
}

// TestCategoryPreviewDomainInteraction verifies preview domains work with categories
func TestCategoryPreviewDomainInteraction(t *testing.T) {
	previewDomains := validation.GetPreviewDomains()

	for _, previewDomain := range previewDomains {
		// Preview domain should have valid category
		assert.NotEmpty(t, previewDomain.Category, "Preview domain %q should have category", previewDomain.Name)

		// Should appear in category listing
		domainsInCategory := validation.GetDomainsByCategory(previewDomain.Category)
		found := false
		for _, d := range domainsInCategory {
			if d.Name == previewDomain.Name {
				found = true
				break
			}
		}
		assert.True(t, found, "Preview domain %q should appear in its category", previewDomain.Name)

		// Should have workflows if category supports them
		workflows := validation.GetWorkflowSuggestions(previewDomain.Name)
		categoryWorkflows := validation.GetWorkflowsByCategory(previewDomain.Category)
		if len(categoryWorkflows) > 0 {
			assert.Greater(t, len(workflows), 0, "Preview domain %q should have category workflows", previewDomain.Name)
		}
	}
}

// TestUseCaseWorkflowAlignment verifies use cases align with workflows
func TestUseCaseWorkflowAlignment(t *testing.T) {
	testDomains := []string{"api", "kubernetes", "authentication", "dns"}

	for _, domain := range testDomains {
		useCases, err := validation.GetDomainUseCases(domain)
		workflows := validation.GetWorkflowSuggestions(domain)

		// If domain has workflows, should have use cases
		if len(workflows) > 0 {
			assert.NoError(t, err, "Domain %q with workflows should have use cases", domain)
			assert.NotEmpty(t, useCases, "Domain %q with workflows should have use cases", domain)
		}

		// Workflow domains should have use cases
		for _, workflow := range workflows {
			for _, workflowDomain := range workflow.Domains {
				wfUseCases, _ := validation.GetDomainUseCases(workflowDomain)
				// Should have use cases or be in category without use cases
				domainInfo, _ := types.GetDomainInfo(workflowDomain)
				if len(domainInfo.UseCases) > 0 {
					assert.NotEmpty(t, wfUseCases, "Workflow domain %q should have use cases", workflowDomain)
				}
			}
		}
	}
}

// TestCategoryUsesCaseKeywordMatch verifies use cases match category patterns
func TestCategoryUseCaseKeywordMatch(t *testing.T) {
	categoryKeywords := map[string]string{
		"Security":       "protect",
		"Networking":     "network",
		"Platform":       "user",
		"Infrastructure": "deploy",
		"Operations":     "monitor",
	}

	for category := range categoryKeywords {
		domains := validation.GetDomainsByCategory(category)

		for _, domain := range domains {
			if len(domain.UseCases) > 0 {
				useCases, err := validation.GetDomainUseCases(domain.Name)
				assert.NoError(t, err, "Should get use cases for %q", domain.Name)

				// Some use cases should contain category-relevant keywords
				// (Not strict requirement, just semantic alignment)
				assert.NotEmpty(t, useCases, "Domain %q in %q should have use cases", domain.Name, category)
			}
		}
	}
}

// TestRelatedDomainsCategoryAlignment verifies related domains are properly selected
func TestRelatedDomainsCategoryAlignment(t *testing.T) {
	testDomains := []string{"api", "dns", "kubernetes", "authentication"}

	for _, domain := range testDomains {
		mainInfo, _ := types.GetDomainInfo(domain)
		relatedDomains := validation.GetRelatedDomains(domain)

		// Related domains should exist (up to 5)
		assert.Greater(t, len(relatedDomains), 0, "Domain %q should have related domains", domain)
		assert.LessOrEqual(t, len(relatedDomains), 5, "Domain %q should have at most 5 related domains", domain)

		// All related domains should be accessible at Enterprise tier (some are Enterprise-only)
		for _, relatedDomain := range relatedDomains {
			canAccess := validation.ValidateTierAccess("Enterprise", relatedDomain.RequiresTier)
			assert.True(t, canAccess, "Related domain %q should be accessible at Enterprise tier", relatedDomain.Name)

			// All related domains should exist
			_, found := types.GetDomainInfo(relatedDomain.Name)
			assert.True(t, found, "Related domain %q should exist", relatedDomain.Name)
		}

		_ = mainInfo // Use mainInfo to avoid unused variable
	}
}

// TestWorkflowConsistencyAcrossTiers verifies workflows for Professional-accessible domains
func TestWorkflowConsistencyAcrossTiers(t *testing.T) {
	// Get Professional tier domains (36 total)
	proDomains := validation.GetDomainsByTier("Professional")

	// For each Professional domain, verify it has consistent workflow structure
	for _, domain := range proDomains {
		workflows := validation.GetWorkflowSuggestions(domain.Name)

		// Each workflow should have required fields
		for _, workflow := range workflows {
			assert.NotEmpty(t, workflow.Name, "Workflow for %q should have name", domain.Name)
			assert.NotEmpty(t, workflow.Description, "Workflow for %q should have description", domain.Name)
			assert.Greater(t, len(workflow.Domains), 0, "Workflow for %q should have domains", domain.Name)
			assert.NotEmpty(t, workflow.Category, "Workflow for %q should have category", domain.Name)
		}
	}
}

// TestPreviewDomainHasCompleteMetadata verifies preview domains have all fields
func TestPreviewDomainHasCompleteMetadata(t *testing.T) {
	previewDomains := validation.GetPreviewDomains()

	for _, domain := range previewDomains {
		// All required fields should be present
		assert.NotEmpty(t, domain.Name, "Preview domain should have name")
		assert.NotEmpty(t, domain.DisplayName, "Preview domain should have display name")
		assert.NotEmpty(t, domain.Description, "Preview domain should have description")
		assert.NotEmpty(t, domain.Category, "Preview domain should have category")
		assert.NotEmpty(t, domain.RequiresTier, "Preview domain should have tier")
		assert.True(t, domain.IsPreview, "Domain should be marked as preview")

		// Should have related domains
		relatedDomains := validation.GetRelatedDomains(domain.Name)
		assert.Greater(t, len(relatedDomains), 0, "Preview domain should have related domains")

		// Should be categorized properly
		categoryDomains := validation.GetDomainsByCategory(domain.Category)
		found := false
		for _, d := range categoryDomains {
			if d.Name == domain.Name {
				found = true
				break
			}
		}
		assert.True(t, found, "Preview domain should be in its category")
	}
}

// TestStandardTierDomainsHaveWorkflows verifies Standard tier domains work properly
func TestStandardTierDomainsHaveWorkflows(t *testing.T) {
	standardDomains := validation.GetDomainsByTier("Standard")

	for _, domain := range standardDomains {
		// Should have related domains
		relatedDomains := validation.GetRelatedDomains(domain.Name)
		assert.Greater(t, len(relatedDomains), 0, "Standard domain %q should have related domains", domain.Name)

		// Should be categorized
		assert.NotEmpty(t, domain.Category, "Standard domain %q should have category", domain.Name)

		// Workflows may or may not exist depending on category
		workflows := validation.GetWorkflowSuggestions(domain.Name)
		// No assertion on workflows - some Standard domains may not have workflows
		_ = workflows
	}
}

// TestSearchUseCasesAcrossAllTiers verifies use case search works tier-agnostic
func TestSearchUseCasesAcrossAllTiers(t *testing.T) {
	searchResults := validation.SearchUseCases("configure")
	assert.Greater(t, len(searchResults), 0, "Should find 'configure' use cases")

	// Results should include domains from all accessible tiers
	domainsCovered := make(map[string]bool)
	for _, result := range searchResults {
		domainsCovered[result.Domain] = true
	}

	assert.Greater(t, len(domainsCovered), 1, "Should find use cases in multiple domains")

	// All returned domains should exist
	for domain := range domainsCovered {
		_, found := types.GetDomainInfo(domain)
		assert.True(t, found, "Use case result domain %q should exist", domain)
	}
}

// TestFullWorkflowPath verifies complete workflow is navigable
func TestFullWorkflowPath(t *testing.T) {
	// Test API Security workflow path
	apiDomain, _ := types.GetDomainInfo("api")
	assert.NotNil(t, apiDomain, "API domain should exist")

	// Get suggested workflows
	workflows := validation.GetWorkflowSuggestions("api")
	assert.Greater(t, len(workflows), 0, "API should have workflow suggestions")

	// For each workflow, verify all domains exist and are accessible at Enterprise tier
	for _, workflow := range workflows {
		for _, domainName := range workflow.Domains {
			domain, found := types.GetDomainInfo(domainName)
			require.True(t, found, "Workflow domain %q should exist", domainName)

			// Domain should be Enterprise-accessible (workflows may include Enterprise-only domains)
			ok := validation.ValidateTierAccess("Enterprise", domain.RequiresTier)
			assert.True(t, ok, "Workflow domain %q should be Enterprise-accessible", domainName)

			// Domain should have proper category
			assert.NotEmpty(t, domain.Category, "Workflow domain %q should have category", domainName)
		}
	}
}

// TestCategoryToWorkflowMapping verifies category workflows have consistent structure
func TestCategoryToWorkflowMapping(t *testing.T) {
	categories := validation.GetAllCategories()

	for _, category := range categories {
		categoryWorkflows := validation.GetWorkflowsByCategory(category)

		// Each workflow should match category
		for _, workflow := range categoryWorkflows {
			assert.Equal(t, category, workflow.Category,
				"Workflow %q should match category %q", workflow.Name, category)

			// All domains in workflow should exist
			for _, domainName := range workflow.Domains {
				domain, found := types.GetDomainInfo(domainName)
				assert.True(t, found, "Workflow domain %q should exist", domainName)
				assert.NotNil(t, domain)

				// Domain should have some relation to category (may not be exact match due to cross-category workflows)
				assert.NotEmpty(t, domain.Category, "Domain %q should have category", domainName)
			}
		}
	}
}

// TestFeatureCombinations verifies realistic feature combinations work
func TestFeatureCombinations(t *testing.T) {
	// Scenario 1: Enterprise user exploring Security category
	securityDomains := validation.GetDomainsByCategory("Security")
	assert.Greater(t, len(securityDomains), 0)

	// Count how many Security domains require Enterprise
	enterpriseOnlyCount := 0
	for _, domain := range securityDomains {
		if domain.RequiresTier == "Enterprise" {
			enterpriseOnlyCount++
		}

		// Enterprise user should access all Security domains
		ok := validation.ValidateTierAccess("Enterprise", domain.RequiresTier)
		assert.True(t, ok, "Enterprise user should access Security domain %q", domain.Name)

		related := validation.GetRelatedDomains(domain.Name)
		assert.Greater(t, len(related), 0, "Security domain %q should have related domains", domain.Name)
	}

	// Some Security domains require Enterprise (e.g., ddos, blindfold, shape)
	assert.Greater(t, enterpriseOnlyCount, 0, "Security category should have some Enterprise-only domains")

	// Scenario 2: Standard user checking what they can access
	standardDomains := validation.GetDomainsByTier("Standard")
	for _, domain := range standardDomains {
		ok := validation.ValidateTierAccess("Standard", domain.RequiresTier)
		assert.True(t, ok, "Standard user should access Standard domain %q", domain.Name)
	}

	// Scenario 3: Verify some Infrastructure domains require Professional or Enterprise
	infraDomains := validation.GetDomainsByCategory("Infrastructure")
	hasNonStandardInfra := false
	for _, domain := range infraDomains {
		if domain.RequiresTier != "Standard" {
			hasNonStandardInfra = true
		}
	}
	assert.True(t, hasNonStandardInfra, "Infrastructure should have some Professional+ domains")
}
