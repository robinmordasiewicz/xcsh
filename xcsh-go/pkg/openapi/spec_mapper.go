package openapi

import (
	"sort"
	"strings"
)

// SpecMapper maps resource names to their corresponding OpenAPI spec files
type SpecMapper struct {
	specs       map[string]*Spec  // filename -> spec
	resourceMap map[string]string // resourceName -> filename
}

// NewSpecMapper creates a new spec mapper from loaded specs
func NewSpecMapper(specs map[string]*Spec) *SpecMapper {
	mapper := &SpecMapper{
		specs:       specs,
		resourceMap: make(map[string]string),
	}
	mapper.buildResourceMap()
	return mapper
}

// buildResourceMap builds the mapping from resource names to spec files
// Uses schema-based discovery: scans all schemas for CreateRequest patterns
func (m *SpecMapper) buildResourceMap() {
	// Sort filenames for deterministic iteration (maps have random iteration order in Go)
	var sortedFilenames []string
	for filename := range m.specs {
		sortedFilenames = append(sortedFilenames, filename)
	}
	sort.Strings(sortedFilenames)

	for _, filename := range sortedFilenames {
		spec := m.specs[filename]
		// Extract resource names from schema names (domain-organized specs)
		// Each spec file can contain multiple resources
		resources := spec.FindAllResourceSchemas()
		for resourceName, schema := range resources {
			if schema != nil {
				// Only store if not already mapped (prefer first/exact match)
				lowerName := strings.ToLower(resourceName)
				if _, exists := m.resourceMap[lowerName]; !exists {
					m.resourceMap[lowerName] = filename
				}
			}
		}
	}
}

// FindSpec finds the OpenAPI spec for a given resource name
func (m *SpecMapper) FindSpec(resourceName string) *Spec {
	// Try exact match first (most common case)
	filename, ok := m.resourceMap[strings.ToLower(resourceName)]
	if ok {
		return m.specs[filename]
	}

	// Try fuzzy match with underscore variations
	// Example: "http_loadbalancer" matches "httploadbalancer"
	lowerResource := strings.ToLower(resourceName)
	normalized := strings.ReplaceAll(lowerResource, "_", "")

	// Sort mapped names for deterministic fuzzy matching (maps have random iteration order in Go)
	var sortedMappedNames []string
	for mappedName := range m.resourceMap {
		sortedMappedNames = append(sortedMappedNames, mappedName)
	}
	sort.Strings(sortedMappedNames)

	for _, mappedName := range sortedMappedNames {
		if strings.ReplaceAll(mappedName, "_", "") == normalized {
			return m.specs[m.resourceMap[mappedName]]
		}
	}

	return nil
}

// FindSpecFile returns the filename of the spec for a given resource
func (m *SpecMapper) FindSpecFile(resourceName string) string {
	// Try exact match first
	filename, ok := m.resourceMap[strings.ToLower(resourceName)]
	if ok {
		return filename
	}

	// Try fuzzy match with underscore variations
	lowerResource := strings.ToLower(resourceName)
	normalized := strings.ReplaceAll(lowerResource, "_", "")

	for mappedName, filename := range m.resourceMap {
		if strings.ReplaceAll(mappedName, "_", "") == normalized {
			return filename
		}
	}

	return ""
}

// GetMappedResources returns all resource names that have been mapped
func (m *SpecMapper) GetMappedResources() []string {
	resources := make([]string, 0, len(m.resourceMap))
	seen := make(map[string]bool)
	for name := range m.resourceMap {
		if !seen[name] {
			resources = append(resources, name)
			seen[name] = true
		}
	}
	return resources
}

// GetSpecCount returns the number of loaded specs
func (m *SpecMapper) GetSpecCount() int {
	return len(m.specs)
}

// GenerateExampleForResource generates a JSON example for the given resource
func (m *SpecMapper) GenerateExampleForResource(resourceName string) (string, error) {
	spec := m.FindSpec(resourceName)
	if spec == nil {
		return "", nil
	}

	generator := NewExampleGenerator(spec)
	return generator.GenerateCreateRequestExample(resourceName)
}

// ResourceSpecInfo contains information about a resource's spec
type ResourceSpecInfo struct {
	ResourceName    string
	SpecFile        string
	HasCreateSchema bool
	HasSpecSchema   bool
}

// GetResourceInfo returns information about a resource's OpenAPI spec
func (m *SpecMapper) GetResourceInfo(resourceName string) *ResourceSpecInfo {
	info := &ResourceSpecInfo{
		ResourceName: resourceName,
	}

	spec := m.FindSpec(resourceName)
	if spec == nil {
		return info
	}

	info.SpecFile = m.FindSpecFile(resourceName)
	info.HasCreateSchema = spec.FindCreateRequestSchema(resourceName) != nil
	info.HasSpecSchema = spec.FindCreateSpecTypeSchema(resourceName) != nil

	return info
}
