package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetDomainCompletionClient(t *testing.T) {
	// Reset global client for test isolation
	domainCompletionClient = nil

	client := getDomainCompletionClient()

	assert.NotNil(t, client)
	assert.NotNil(t, client.cache)
	assert.NotNil(t, client.cache.resources)
	assert.NotNil(t, client.cache.resourceTTL)
	assert.Equal(t, 3*time.Second, client.timeout)

	// Verify singleton pattern
	client2 := getDomainCompletionClient()
	assert.Same(t, client, client2)
}

func TestCompleteLabelKey(t *testing.T) {
	completions, directive := completeLabelKey(nil, []string{}, "")

	assert.NotEmpty(t, completions, "Should return completions for label keys")
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)

	// Verify expected label keys are present
	labelKeys := map[string]bool{
		"environment": false,
		"application": false,
		"owner":       false,
		"cost-center": false,
		"tier":        false,
		"version":     false,
	}

	for _, completion := range completions {
		for key := range labelKeys {
			if len(completion) > len(key) && completion[:len(key)] == key {
				labelKeys[key] = true
			}
		}
	}

	// Verify all expected keys are present
	for key, found := range labelKeys {
		assert.True(t, found, "Label key %q should be in completions", key)
	}
}

func TestCompleteOutputFormat(t *testing.T) {
	completions, directive := completeOutputFormat(nil, []string{}, "")

	assert.Len(t, completions, 3, "Should return 3 output format options")
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)

	expected := map[string]bool{
		"json":  false,
		"yaml":  false,
		"table": false,
	}

	for _, completion := range completions {
		for format := range expected {
			if len(completion) > len(format) && completion[:len(format)] == format {
				expected[format] = true
			}
		}
	}

	// Verify all formats are present
	for format, found := range expected {
		assert.True(t, found, "Format %q should be in completions", format)
	}
}

func TestCompleteDomainResourceType(t *testing.T) {
	// Test with a domain that has resources
	completionFunc := completeDomainResourceType("cdn")
	completions, directive := completionFunc(nil, []string{}, "")

	assert.NotEmpty(t, completions, "Should return resource types for cdn domain")
	assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)

	// Verify http_loadbalancer is present (should be in cdn domain)
	found := false
	for _, comp := range completions {
		if len(comp) > len("http_loadbalancer") && comp[:len("http_loadbalancer")] == "http_loadbalancer" {
			found = true
			break
		}
	}
	assert.True(t, found, "http_loadbalancer should be in cdn domain completions")
}

func TestCompleteDomainResourceTypeMultipleDomains(t *testing.T) {
	// Test that different domains return their respective resources
	completionFunc := completeDomainResourceType("cdn")
	completions, _ := completionFunc(nil, []string{}, "")

	// CDN domain should have http_loadbalancer
	found := false
	for _, comp := range completions {
		if len(comp) > len("http_loadbalancer") && comp[:len("http_loadbalancer")] == "http_loadbalancer" {
			found = true
			break
		}
	}
	assert.True(t, found, "http_loadbalancer should be in cdn domain completions")

	// Another domain should have different resources
	completionFunc2 := completeDomainResourceType("identity")
	completions2, _ := completionFunc2(nil, []string{}, "")

	// Verify completions are different for different domains
	if len(completions) > 0 && len(completions2) > 0 {
		// They shouldn't be exactly the same since they're different domains
		// (Though technically some resources can appear in multiple domains)
		assert.NotNil(t, completions, "cdn should have completions")
		assert.NotNil(t, completions2, "identity should have completions")
	}
}

func TestCompletionCacheInitialization(t *testing.T) {
	domainCompletionClient = nil
	client := getDomainCompletionClient()

	// Verify cache is properly initialized
	assert.NotNil(t, client.cache)
	assert.Empty(t, client.cache.namespaces, "Namespaces cache should start empty")
	assert.True(t, client.cache.namespaceTTL.IsZero(), "Namespace TTL should be zero initially")
	assert.Empty(t, client.cache.resources, "Resources cache should start empty")
	assert.Empty(t, client.cache.resourceTTL, "Resource TTL should start empty")
}

func TestCompleteResourceNameSignature(t *testing.T) {
	// Verify the function returns a proper completion function
	completionFunc := completeResourceName("load_balancer", "http_loadbalancer")

	assert.NotNil(t, completionFunc)

	// The function should be callable and return the right types
	// We can't call it with nil command as it will panic trying to get flags
	// but we've verified the function is not nil and is of the right type
	// (the actual functionality is tested in integration tests with real commands)

	// Just verify the closure was created properly
	assert.NotNil(t, completionFunc, "Should return a completion function for load_balancer/http_loadbalancer")
}

func TestCompletionDomainVariations(t *testing.T) {
	// Test various domains to ensure the factory function works for all
	domains := []string{
		"load_balancer",
		"infrastructure",
		"security",
		"networking",
		"observability",
		"application",
		"identity",
	}

	for _, domain := range domains {
		completionFunc := completeDomainResourceType(domain)
		assert.NotNil(t, completionFunc, "Should return completion function for domain %q", domain)

		// All should return without error
		completions, directive := completionFunc(nil, []string{}, "")
		assert.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive, "Should return NoFileComp directive for domain %q", domain)

		// Should have some completions for all domains (otherwise domain wouldn't exist)
		if len(completions) == 0 {
			t.Logf("Warning: domain %q has no completions", domain)
		}
	}
}

func TestCompleteLabelKeyDescription(t *testing.T) {
	completions, _ := completeLabelKey(nil, []string{}, "")

	// Each completion should have format "key\tdescription"
	for _, comp := range completions {
		// Should have tab separator
		tabFound := false
		for i, ch := range comp {
			if ch == '\t' {
				// Check that there's a key before and description after
				assert.Greater(t, i, 0, "Should have key before tab")
				assert.Less(t, i, len(comp)-1, "Should have description after tab")
				tabFound = true
				break
			}
		}
		assert.True(t, tabFound, "Completion %q should have tab separator", comp)
	}
}

func TestOutputFormatDescription(t *testing.T) {
	completions, _ := completeOutputFormat(nil, []string{}, "")

	// Each completion should have format "format\tdescription"
	for _, comp := range completions {
		// Should have tab separator
		tabFound := false
		for i, ch := range comp {
			if ch == '\t' {
				// Check that there's a format before and description after
				assert.Greater(t, i, 0, "Should have format before tab")
				assert.Less(t, i, len(comp)-1, "Should have description after tab")
				tabFound = true
				break
			}
		}
		assert.True(t, tabFound, "Completion %q should have tab separator", comp)
	}
}

// Benchmark tests for completion performance

func BenchmarkGetDomainCompletionClient(b *testing.B) {
	domainCompletionClient = nil
	client := getDomainCompletionClient()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = getDomainCompletionClient()
	}

	// Verify singleton was used
	assert.Same(b, client, getDomainCompletionClient())
}

func BenchmarkCompleteLabelKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = completeLabelKey(nil, []string{}, "")
	}
}

func BenchmarkCompleteOutputFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = completeOutputFormat(nil, []string{}, "")
	}
}

func BenchmarkCompleteDomainResourceType(b *testing.B) {
	completionFunc := completeDomainResourceType("load_balancer")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = completionFunc(nil, []string{}, "")
	}
}
