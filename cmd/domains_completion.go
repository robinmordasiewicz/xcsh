package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/naming"
	"github.com/robinmordasiewicz/f5xcctl/pkg/types"
)

// Domain completion client with caching
var domainCompletionClient *DomainCompletionClient

// DomainCompletionClient handles completion API calls with caching
type DomainCompletionClient struct {
	cache   *CompletionCache
	timeout time.Duration
}

// CompletionCache stores completion results with TTL
type CompletionCache struct {
	namespaces   []string
	namespaceTTL time.Time
	resources    map[string][]string // domain+resourceType -> names
	resourceTTL  map[string]time.Time
}

// getDomainCompletionClient returns a client optimized for tab completion
func getDomainCompletionClient() *DomainCompletionClient {
	if domainCompletionClient == nil {
		domainCompletionClient = &DomainCompletionClient{
			cache: &CompletionCache{
				resources:   make(map[string][]string),
				resourceTTL: make(map[string]time.Time),
			},
			timeout: 3 * time.Second,
		}
	}
	return domainCompletionClient
}

// Static completion functions (no API calls)

// completeDomainResourceType provides completion for resource types within a domain
func completeDomainResourceType(domain string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		resources := types.GetByDomain(domain)
		completions := make([]string, 0, len(resources))

		for _, rt := range resources {
			displayName := naming.ToHumanReadable(rt.Name)
			description := rt.Description
			if description == "" {
				description = displayName
			}
			// Add tier annotation if applicable
			if tier := getTierAnnotation(rt.Name); tier != "" {
				description += " " + tier
			}
			completions = append(completions, rt.Name+"\t"+description)
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

// completeOutputFormat provides completion for --output-format flag
func completeOutputFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"json\tJSON format",
		"yaml\tYAML format",
		"table\tTable format (default)",
	}, cobra.ShellCompDirectiveNoFileComp
}

// Dynamic completion functions (require API calls)

// completeNamespace provides completion for --namespace flag
func completeNamespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getDomainCompletionClient()

	// Check cache first (5-minute TTL)
	if time.Now().Before(client.cache.namespaceTTL) && len(client.cache.namespaces) > 0 {
		completions := make([]string, len(client.cache.namespaces))
		for i, ns := range client.cache.namespaces {
			completions[i] = ns + "\tNamespace"
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	// Fetch from API with timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	namespaces, err := fetchNamespaces(ctx)
	if err != nil {
		// Fall back to common namespaces on error
		return []string{
			"default\tDefault namespace",
			"system\tSystem namespace",
		}, cobra.ShellCompDirectiveNoFileComp
	}

	// Update cache
	client.cache.namespaces = namespaces
	client.cache.namespaceTTL = time.Now().Add(5 * time.Minute)

	completions := make([]string, len(namespaces))
	for i, ns := range namespaces {
		completions[i] = ns + "\tNamespace"
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeResourceName provides completion for resource names (for get/delete/status)
func completeResourceName(domain, resourceType string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Resource name is a positional arg, not a flag
		// Only complete if we haven't provided the name yet
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		client := getDomainCompletionClient()
		cacheKey := domain + ":" + resourceType

		// Check cache first (5-minute TTL)
		if ttl, ok := client.cache.resourceTTL[cacheKey]; ok && time.Now().Before(ttl) {
			if names, ok := client.cache.resources[cacheKey]; ok {
				completions := make([]string, len(names))
				for i, name := range names {
					completions[i] = name + "\tResource name"
				}
				return completions, cobra.ShellCompDirectiveNoFileComp
			}
		}

		// Get namespace from flag
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}

		// Fetch from API with timeout
		ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
		defer cancel()

		names, err := fetchResourceNames(ctx, resourceType, namespace)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// Update cache
		client.cache.resources[cacheKey] = names
		client.cache.resourceTTL[cacheKey] = time.Now().Add(5 * time.Minute)

		completions := make([]string, len(names))
		for i, name := range names {
			completions[i] = name + "\tResource name"
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

// completeLabelKey provides completion for label keys
func completeLabelKey(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Common F5 XC label keys
	return []string{
		"environment\tEnvironment label",
		"application\tApplication label",
		"owner\tOwner label",
		"cost-center\tCost center label",
		"tier\tTier label",
		"version\tVersion label",
	}, cobra.ShellCompDirectiveNoFileComp
}

// Helper functions

// fetchNamespaces retrieves namespace list from API
func fetchNamespaces(ctx context.Context) ([]string, error) {
	// Use existing API client to fetch namespaces
	apiClient := GetClient()
	if apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	result, err := apiClient.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, len(result.Items))
	for i, ns := range result.Items {
		namespaces[i] = ns.Name
	}
	return namespaces, nil
}

// fetchResourceNames retrieves resource names from API
func fetchResourceNames(ctx context.Context, resourceType, namespace string) ([]string, error) {
	// Use existing API client to list resources
	apiClient := GetClient()
	if apiClient == nil {
		return nil, fmt.Errorf("API client not initialized")
	}

	rt, ok := types.Get(resourceType)
	if !ok || rt == nil {
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	result, err := apiClient.ListResources(ctx, rt, namespace)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(result.Items))
	for i, item := range result.Items {
		names[i] = item.Metadata.Name
	}
	return names, nil
}
