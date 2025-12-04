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
