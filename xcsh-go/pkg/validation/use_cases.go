package validation

import (
	"fmt"
	"strings"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// FormatUseCases formats use cases for display in help text or CLI output
func FormatUseCases(useCases []string) string {
	if len(useCases) == 0 {
		return ""
	}

	// Format as a list with bullet points
	var formatted strings.Builder
	formatted.WriteString("\nUSE CASES:\n")

	for _, useCase := range useCases {
		formatted.WriteString(fmt.Sprintf("  • %s\n", useCase))
	}

	return formatted.String()
}

// FormatUseCasesShort formats first N use cases for compact display
func FormatUseCasesShort(useCases []string, maxCount int) string {
	if len(useCases) == 0 {
		return ""
	}

	if maxCount <= 0 {
		maxCount = 1
	}

	if len(useCases) > maxCount {
		useCases = useCases[:maxCount]
	}

	return strings.Join(useCases, ", ")
}

// FormatUseCasesInline formats use cases as a single-line comma-separated list
func FormatUseCasesInline(useCases []string) string {
	if len(useCases) == 0 {
		return ""
	}

	return strings.Join(useCases, ", ")
}

// GetDomainUseCases returns formatted use cases for a domain
func GetDomainUseCases(domain string) (string, error) {
	info, found := types.GetDomainInfo(domain)
	if !found {
		return "", fmt.Errorf("domain not found: %s", domain)
	}

	return FormatUseCases(info.UseCases), nil
}

// GetDomainsWithUseCases returns all domains that have at least one use case
func GetDomainsWithUseCases() []*types.DomainInfo {
	var domains []*types.DomainInfo

	for _, info := range types.DomainRegistry {
		if len(info.UseCases) > 0 {
			domains = append(domains, info)
		}
	}

	return domains
}

// GetDomainsWithoutUseCases returns all domains that have no use cases defined
func GetDomainsWithoutUseCases() []*types.DomainInfo {
	var domains []*types.DomainInfo

	for _, info := range types.DomainRegistry {
		if len(info.UseCases) == 0 {
			domains = append(domains, info)
		}
	}

	return domains
}

// UseCaseStatistics provides statistics about use case coverage
type UseCaseStatistics struct {
	TotalDomains           int
	DomainsWithUseCases    int
	DomainsWithoutUseCases int
	CoveragePercentage     float64
	TotalUseCases          int
	AveragePerDomain       float64
}

// CalculateUseCaseStatistics returns statistics about use case coverage
func CalculateUseCaseStatistics() UseCaseStatistics {
	stats := UseCaseStatistics{
		TotalDomains: len(types.DomainRegistry),
	}

	for _, info := range types.DomainRegistry {
		if len(info.UseCases) > 0 {
			stats.DomainsWithUseCases++
			stats.TotalUseCases += len(info.UseCases)
		} else {
			stats.DomainsWithoutUseCases++
		}
	}

	if stats.TotalDomains > 0 {
		stats.CoveragePercentage = float64(stats.DomainsWithUseCases) / float64(stats.TotalDomains) * 100
	}

	if stats.DomainsWithUseCases > 0 {
		stats.AveragePerDomain = float64(stats.TotalUseCases) / float64(stats.DomainsWithUseCases)
	}

	return stats
}

// FormatUseCaseStatistics formats use case statistics for display
func FormatUseCaseStatistics(stats UseCaseStatistics) string {
	return fmt.Sprintf(`Use Case Coverage Summary:
  Total Domains:              %d
  Domains with Use Cases:     %d
  Domains without Use Cases:  %d
  Coverage:                   %.1f%%
  Total Use Cases:            %d
  Average per Domain:         %.1f`,
		stats.TotalDomains,
		stats.DomainsWithUseCases,
		stats.DomainsWithoutUseCases,
		stats.CoveragePercentage,
		stats.TotalUseCases,
		stats.AveragePerDomain,
	)
}

// UseCase represents a single use case with metadata
type UseCase struct {
	Domain      string
	Description string
	Category    string
}

// GetAllUseCases returns all use cases across all domains with domain info
func GetAllUseCases() []UseCase {
	var useCases []UseCase

	for _, info := range types.DomainRegistry {
		for _, useCase := range info.UseCases {
			useCases = append(useCases, UseCase{
				Domain:      info.Name,
				Description: useCase,
				Category:    info.Category,
			})
		}
	}

	return useCases
}

// SearchUseCases searches for use cases matching a keyword
func SearchUseCases(keyword string) []UseCase {
	if keyword == "" {
		return GetAllUseCases()
	}

	keyword = strings.ToLower(keyword)
	var matching []UseCase

	for _, useCase := range GetAllUseCases() {
		if strings.Contains(strings.ToLower(useCase.Description), keyword) {
			matching = append(matching, useCase)
		}
	}

	return matching
}
