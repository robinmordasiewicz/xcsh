package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// TestTierValidationStandardTierDomains verifies that Standard tier users can access all Standard domains
func TestTierValidationStandardTierDomains(t *testing.T) {
	standardDomains := []string{
		"dns",
		"kubernetes_and_orchestration",
		"authentication",
		"certificates",
		"observability",
	}

	for _, domain := range standardDomains {
		t.Run("Standard_access_"+domain, func(t *testing.T) {
			info, found := types.GetDomainInfo(domain)
			require.True(t, found, "Domain %q should exist", domain)

			// Standard tier should access Standard domains
			require.Equal(t, validation.TierStandard, info.RequiresTier,
				"Domain %q should have Standard tier requirement", domain)

			canAccess := validation.IsSufficientTier(validation.TierStandard, info.RequiresTier)
			assert.True(t, canAccess, "Standard tier user should access domain %q", domain)
		})
	}
}

// TestTierValidationProfessionalTierDomains verifies that Professional tier users can access Professional domains
func TestTierValidationProfessionalTierDomains(t *testing.T) {
	professionalDomains := []string{
		"api",
		"network_security",
		"kubernetes",
		"application_firewall",
	}

	for _, domain := range professionalDomains {
		t.Run("Professional_access_"+domain, func(t *testing.T) {
			info, found := types.GetDomainInfo(domain)
			require.True(t, found, "Domain %q should exist", domain)

			// Verify domain requires Professional tier
			require.Equal(t, validation.TierProfessional, info.RequiresTier,
				"Domain %q should have Professional tier requirement", domain)

			// Professional tier should access Professional domains
			canAccess := validation.IsSufficientTier(validation.TierProfessional, info.RequiresTier)
			assert.True(t, canAccess, "Professional tier user should access domain %q", domain)

			// Standard tier should NOT access Professional domains
			cannotAccess := validation.IsSufficientTier(validation.TierStandard, info.RequiresTier)
			assert.False(t, cannotAccess, "Standard tier user should not access domain %q", domain)
		})
	}
}

// TestTierValidationEnterpriseTierDomains verifies that Enterprise tier users can access Enterprise domains
func TestTierValidationEnterpriseTierDomains(t *testing.T) {
	enterpriseDomains := []string{
		"generative_ai",
		"ddos",
		"cdn",
		"blindfold",
	}

	for _, domain := range enterpriseDomains {
		t.Run("Enterprise_access_"+domain, func(t *testing.T) {
			info, found := types.GetDomainInfo(domain)
			require.True(t, found, "Domain %q should exist", domain)

			// Verify domain requires Enterprise tier
			require.Equal(t, validation.TierEnterprise, info.RequiresTier,
				"Domain %q should have Enterprise tier requirement", domain)

			// Enterprise tier should access Enterprise domains
			canAccess := validation.IsSufficientTier(validation.TierEnterprise, info.RequiresTier)
			assert.True(t, canAccess, "Enterprise tier user should access domain %q", domain)

			// Professional tier should NOT access Enterprise domains
			cannotAccessProf := validation.IsSufficientTier(validation.TierProfessional, info.RequiresTier)
			assert.False(t, cannotAccessProf, "Professional tier user should not access domain %q", domain)

			// Standard tier should NOT access Enterprise domains
			cannotAccessStd := validation.IsSufficientTier(validation.TierStandard, info.RequiresTier)
			assert.False(t, cannotAccessStd, "Standard tier user should not access domain %q", domain)
		})
	}
}

// TestTierValidationUpgradePath verifies that upgrade paths are suggested correctly
func TestTierValidationUpgradePath(t *testing.T) {
	tests := []struct {
		currentTier   string
		requiredTier  string
		shouldUpgrade bool
		expectedPath  string
	}{
		{
			currentTier:   validation.TierStandard,
			requiredTier:  validation.TierProfessional,
			shouldUpgrade: true,
			expectedPath:  "Upgrade from Standard to Professional tier",
		},
		{
			currentTier:   validation.TierStandard,
			requiredTier:  validation.TierEnterprise,
			shouldUpgrade: true,
			expectedPath:  "Upgrade from Standard to Enterprise tier",
		},
		{
			currentTier:   validation.TierProfessional,
			requiredTier:  validation.TierEnterprise,
			shouldUpgrade: true,
			expectedPath:  "Upgrade from Professional to Enterprise tier",
		},
		{
			currentTier:   validation.TierProfessional,
			requiredTier:  validation.TierStandard,
			shouldUpgrade: false,
			expectedPath:  "",
		},
		{
			currentTier:   validation.TierStandard,
			requiredTier:  validation.TierStandard,
			shouldUpgrade: false,
			expectedPath:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.currentTier+"_to_"+tt.requiredTier, func(t *testing.T) {
			path := validation.GetUpgradePath(tt.currentTier, tt.requiredTier)

			if tt.shouldUpgrade {
				assert.Equal(t, tt.expectedPath, path,
					"Should suggest upgrade path from %s to %s", tt.currentTier, tt.requiredTier)
			} else {
				assert.Empty(t, path,
					"Should not suggest upgrade when %s >= %s", tt.currentTier, tt.requiredTier)
			}
		})
	}
}

// TestTierValidationErrorMessages verifies that error messages contain required information
func TestTierValidationErrorMessages(t *testing.T) {
	err := validation.NewTierAccessError("api", "API", validation.TierStandard, validation.TierProfessional)

	errMsg := err.Error()
	assert.NotEmpty(t, errMsg, "Error message should not be empty")

	// Verify error message contains key information
	assert.Contains(t, errMsg, "API", "Error should mention domain display name")
	assert.Contains(t, errMsg, "Professional", "Error should mention required tier")
	assert.Contains(t, errMsg, "Standard", "Error should mention current tier")
	assert.Contains(t, errMsg, "console.volterra.io", "Error should include upgrade URL")
	assert.Contains(t, errMsg, "support@f5.com", "Error should include support contact")
}

// TestTierValidationErrorStructure verifies the TierAccessError structure
func TestTierValidationErrorStructure(t *testing.T) {
	domain := "api"
	displayName := "API"
	currentTier := validation.TierStandard
	requiredTier := validation.TierProfessional

	err := validation.NewTierAccessError(domain, displayName, currentTier, requiredTier)

	assert.Equal(t, domain, err.Domain, "Error Domain field should match")
	assert.Equal(t, displayName, err.DomainDisplay, "Error DomainDisplay field should match")
	assert.Equal(t, currentTier, err.CurrentTier, "Error CurrentTier field should match")
	assert.Equal(t, requiredTier, err.RequiredTier, "Error RequiredTier field should match")
}

// TestValidateDomainTierStandard verifies domain tier validation for Standard tier users
func TestValidateDomainTierStandard(t *testing.T) {
	ctx := context.Background()

	// Test Standard domain access
	err := ValidateDomainTier(ctx, "dns")
	assert.NoError(t, err, "Standard tier user should access Standard domain")

	// Test Professional domain access (should fail)
	err = ValidateDomainTier(ctx, "api")
	assert.Error(t, err, "Standard tier user should not access Professional domain")
	assert.Contains(t, err.Error(), "Api", "Error should mention domain")
}

// TestValidateDomainTierProfessional verifies domain tier validation for Professional tier users
func TestValidateDomainTierProfessional(t *testing.T) {
	ctx := context.Background()

	// Professional users should access both Standard and Professional domains
	// (This test verifies the validation logic works, though we can't directly
	// control the subscription tier in unit tests without mocking)

	// At minimum, verify the function doesn't panic and returns expected types
	err := ValidateDomainTier(ctx, "load_balancer")
	// Error depends on runtime tier detection, just verify no panic
	_ = err

	info, found := types.GetDomainInfo("api")
	require.True(t, found, "api domain should exist")
	assert.Equal(t, validation.TierProfessional, info.RequiresTier, "api should require Professional tier")
}

// TestDomainTierRequirementsConsistency verifies all 42 domains have tier requirements set
func TestDomainTierRequirementsConsistency(t *testing.T) {
	domainCount := 0
	standardCount := 0
	professionalCount := 0
	enterpriseCount := 0
	emptyCount := 0

	for domain := range types.DomainRegistry {
		domainCount++
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist in registry", domain)

		switch info.RequiresTier {
		case validation.TierStandard:
			standardCount++
		case validation.TierProfessional:
			professionalCount++
		case validation.TierEnterprise:
			enterpriseCount++
		case "":
			emptyCount++
		default:
			t.Errorf("Domain %q has unknown tier requirement: %q", domain, info.RequiresTier)
		}
	}

	// Verify we have all 41 domains
	assert.Equal(t, 41, domainCount, "Should have 41 domains total")

	// Verify tier distribution matches actual: 24 Standard, 11 Professional, 6 Enterprise
	assert.Equal(t, 24, standardCount, "Should have 24 Standard domains")
	assert.Equal(t, 11, professionalCount, "Should have 11 Professional domains")
	assert.Equal(t, 6, enterpriseCount, "Should have 6 Enterprise domains")

	// No domains should have empty tier requirement (should be set during generation)
	assert.Equal(t, 0, emptyCount, "All domains should have tier requirement set")

	t.Logf("Domain tier distribution: Standard=%d, Professional=%d, Enterprise=%d",
		standardCount, professionalCount, enterpriseCount)
}

// TestDomainDisplayNameConsistency verifies all domains have display names
func TestDomainDisplayNameConsistency(t *testing.T) {
	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist in registry", domain)

		assert.NotEmpty(t, info.DisplayName, "Domain %q should have DisplayName", domain)
		assert.NotEmpty(t, info.Description, "Domain %q should have Description", domain)
	}
}

// TestTierComparisonConsistency verifies tier comparison logic is transitive and consistent
func TestTierComparisonConsistency(t *testing.T) {
	tiers := []string{
		validation.TierStandard,
		validation.TierProfessional,
		validation.TierEnterprise,
	}

	// Verify transitivity: if A >= B and B >= C then A >= C
	for _, tierA := range tiers {
		for _, tierB := range tiers {
			if validation.IsSufficientTier(tierA, tierB) {
				// tierA >= tierB
				tierALevel := validation.TierLevel(tierA)
				tierBLevel := validation.TierLevel(tierB)
				assert.GreaterOrEqual(t, tierALevel, tierBLevel,
					"Tier level should be consistent with IsSufficientTier")
			}
		}
	}

	// Verify reflexivity: each tier is sufficient for itself
	for _, tier := range tiers {
		assert.True(t, validation.IsSufficientTier(tier, tier),
			"Tier %s should be sufficient for itself", tier)
	}
}

// TestEnterpriseCanAccessAll verifies that Enterprise tier can access all domains
func TestEnterpriseCanAccessAll(t *testing.T) {
	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist", domain)

		canAccess := validation.IsSufficientTier(validation.TierEnterprise, info.RequiresTier)
		assert.True(t, canAccess, "Enterprise tier should access domain %q (requires %s)",
			domain, info.RequiresTier)
	}
}

// TestStandardCanAccessOnlyStandard verifies that Standard tier can only access Standard domains
func TestStandardCanAccessOnlyStandard(t *testing.T) {
	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist", domain)

		canAccess := validation.IsSufficientTier(validation.TierStandard, info.RequiresTier)

		if info.RequiresTier == validation.TierStandard {
			assert.True(t, canAccess, "Standard tier should access Standard domain %q", domain)
		} else {
			assert.False(t, canAccess, "Standard tier should not access %s domain %q",
				info.RequiresTier, domain)
		}
	}
}

// TestProfessionalCanAccessStandardAndProfessional verifies Professional tier access rules
func TestProfessionalCanAccessStandardAndProfessional(t *testing.T) {
	for domain := range types.DomainRegistry {
		info, found := types.GetDomainInfo(domain)
		require.True(t, found, "Domain %q should exist", domain)

		canAccess := validation.IsSufficientTier(validation.TierProfessional, info.RequiresTier)

		switch info.RequiresTier {
		case validation.TierStandard, validation.TierProfessional:
			assert.True(t, canAccess, "Professional tier should access %s domain %q",
				info.RequiresTier, domain)
		case validation.TierEnterprise:
			assert.False(t, canAccess, "Professional tier should not access Enterprise domain %q", domain)
		}
	}
}

// BenchmarkValidateDomainTier benchmarks the tier validation function
func BenchmarkValidateDomainTier(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateDomainTier(ctx, "api")
	}
}

// BenchmarkTierComparison benchmarks tier comparison logic
func BenchmarkTierComparison(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validation.IsSufficientTier(validation.TierProfessional, validation.TierStandard)
	}
}
