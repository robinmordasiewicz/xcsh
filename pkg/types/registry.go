package types

import (
	"fmt"
	"strings"
	"sync"
)

// Registry holds all registered resource types
type Registry struct {
	mu        sync.RWMutex
	resources map[string]*ResourceType

	// Domain indexes (lazy-initialized, cached)
	domainResources map[string][]*ResourceType // domain → resources
	primaryDomain   map[string]string          // resourceName → primary domain
}

// Global registry instance
var globalRegistry = &Registry{
	resources: make(map[string]*ResourceType),
}

// Register adds a resource type to the global registry
func Register(rt *ResourceType) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.resources[rt.CLIName] = rt
}

// Get retrieves a resource type from the global registry
func Get(cliName string) (*ResourceType, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	rt, ok := globalRegistry.resources[cliName]
	return rt, ok
}

// All returns all registered resource types
func All() []*ResourceType {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	result := make([]*ResourceType, 0, len(globalRegistry.resources))
	for _, rt := range globalRegistry.resources {
		result = append(result, rt)
	}
	return result
}

// Count returns the number of registered resource types
func Count() int {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	return len(globalRegistry.resources)
}

// buildDomainIndexes creates domain lookup indexes (lazy-initialized, cached)
func (r *Registry) buildDomainIndexes() {
	r.domainResources = make(map[string][]*ResourceType)
	r.primaryDomain = make(map[string]string)

	for _, rt := range r.resources {
		// Index by all domains for cross-domain access
		for _, domain := range rt.Domains {
			r.domainResources[domain] = append(r.domainResources[domain], rt)
		}
		// Track primary domain for each resource
		r.primaryDomain[rt.Name] = rt.PrimaryDomain
	}
}

// GetByDomain returns all resources in a domain (cross-domain enabled)
func GetByDomain(domain string) []*ResourceType {
	globalRegistry.mu.RLock()

	// Check if we need to initialize indexes
	if globalRegistry.domainResources == nil {
		// Release read lock and acquire write lock
		globalRegistry.mu.RUnlock()
		globalRegistry.mu.Lock()

		// Double-check pattern for thread safety
		if globalRegistry.domainResources == nil {
			globalRegistry.buildDomainIndexes()
		}

		globalRegistry.mu.Unlock()
		globalRegistry.mu.RLock()
	}

	// Now we're back in read lock mode
	resources := globalRegistry.domainResources[domain]
	globalRegistry.mu.RUnlock()

	if resources == nil {
		return []*ResourceType{}
	}
	return resources
}

// GetPrimaryDomain returns the primary domain for a resource
func GetPrimaryDomain(resourceName string) string {
	globalRegistry.mu.RLock()

	// Check if we need to initialize indexes
	if globalRegistry.primaryDomain == nil {
		// Release read lock and acquire write lock
		globalRegistry.mu.RUnlock()
		globalRegistry.mu.Lock()

		// Double-check pattern for thread safety
		if globalRegistry.primaryDomain == nil {
			globalRegistry.buildDomainIndexes()
		}

		globalRegistry.mu.Unlock()
		globalRegistry.mu.RLock()
	}

	// Now we're back in read lock mode
	primaryDomain := globalRegistry.primaryDomain[resourceName]
	globalRegistry.mu.RUnlock()

	return primaryDomain
}

// BuildAPIPath constructs the full API path for a resource
func (rt *ResourceType) BuildAPIPath(namespace, name string) string {
	path := rt.APIPath

	if rt.SupportsNamespace && namespace != "" {
		path = strings.Replace(path, "{namespace}", namespace, 1)
	}

	if name != "" {
		path = fmt.Sprintf("%s/%s", path, name)
	}

	return path
}
