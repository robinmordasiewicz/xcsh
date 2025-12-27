package validation

import (
	"sort"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// CategoryGrouping organizes domains by their category
type CategoryGrouping struct {
	Category string
	Domains  []*types.DomainInfo
}

// GetAllCategories returns all unique categories sorted alphabetically
func GetAllCategories() []string {
	categories := make(map[string]bool)

	for _, info := range types.DomainRegistry {
		if info.Category != "" {
			categories[info.Category] = true
		}
	}

	result := make([]string, 0, len(categories))
	for cat := range categories {
		result = append(result, cat)
	}
	sort.Strings(result)

	return result
}

// GetDomainsByCategory returns all domains in a specific category, sorted by name
func GetDomainsByCategory(category string) []*types.DomainInfo {
	var domains []*types.DomainInfo

	for _, info := range types.DomainRegistry {
		if info.Category == category {
			domains = append(domains, info)
		}
	}

	// Sort by display name for consistent ordering
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].DisplayName < domains[j].DisplayName
	})

	return domains
}

// GroupDomainsByCategory returns all domains grouped by their category
// Categories are sorted alphabetically, domains within each category are sorted by name
func GroupDomainsByCategory() []CategoryGrouping {
	categories := GetAllCategories()
	result := make([]CategoryGrouping, 0, len(categories))

	for _, cat := range categories {
		grouping := CategoryGrouping{
			Category: cat,
			Domains:  GetDomainsByCategory(cat),
		}
		result = append(result, grouping)
	}

	return result
}

// GetDomainsInCategories returns all domains that belong to any of the specified categories
func GetDomainsInCategories(categoryFilter []string) []*types.DomainInfo {
	if len(categoryFilter) == 0 {
		return nil
	}

	// Create set for fast lookup
	filterSet := make(map[string]bool)
	for _, cat := range categoryFilter {
		filterSet[cat] = true
	}

	var domains []*types.DomainInfo
	for _, info := range types.DomainRegistry {
		if filterSet[info.Category] {
			domains = append(domains, info)
		}
	}

	// Sort by display name for consistent ordering
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].DisplayName < domains[j].DisplayName
	})

	return domains
}

// CategoryCount represents the count of domains in a category
type CategoryCount struct {
	Category string
	Count    int
}

// GetCategoryDistribution returns the count of domains in each category
func GetCategoryDistribution() []CategoryCount {
	categories := GetAllCategories()
	result := make([]CategoryCount, 0, len(categories))

	for _, cat := range categories {
		domains := GetDomainsByCategory(cat)
		result = append(result, CategoryCount{
			Category: cat,
			Count:    len(domains),
		})
	}

	// Sort by count descending, then by category name
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Category < result[j].Category
	})

	return result
}
