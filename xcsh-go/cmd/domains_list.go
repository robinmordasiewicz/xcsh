package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/robinmordasiewicz/xcsh/pkg/validation"
)

// domainsListCmd lists all available domains organized by category
var domainsListCmd = &cobra.Command{
	Use:   "domains",
	Short: "List all available domains organized by category",
	Long: `List all available domains in F5 Distributed Cloud organized by category.

Each domain is displayed with:
- Category classification
- Subscription tier requirement (if applicable)
- Preview status (if applicable)
- Brief description

Use this command to discover available domains and understand their organization.

EXAMPLES:
  # List all domains by category
  xcsh domains

  # List domains by subscription tier
  xcsh domains by-tier

  # List domains in specific category/categories
  xcsh domains --category Security
  xcsh domains --category Security,Platform,Networking

  # List domains by tier
  xcsh domains by-tier`,
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryFilter, _ := cmd.Flags().GetString("category")
		return listDomainsInteractiveWithFilter(cmd, categoryFilter)
	},
}

func initDomainsListCmd() {
	domainsListCmd.Flags().String("category", "", "Filter by category (comma-separated for multiple)")
}

// domainsListByTierCmd lists domains filtered by subscription tier
var domainsListByTierCmd = &cobra.Command{
	Use:   "by-tier",
	Short: "List domains filtered by subscription tier",
	Long: `List all available domains organized by subscription tier requirement.

Subscription tiers:
- Standard: Available with base subscription (Default)
- Advanced: Requires Advanced tier subscription

This helps you understand which features are available with your subscription level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listDomainsByTier(cmd)
	},
}

// listDomainsInteractiveWithFilter displays domains with optional category filtering
func listDomainsInteractiveWithFilter(cmd *cobra.Command, categoryFilter string) error {
	var groupings []validation.CategoryGrouping

	// If no filter, show all categories
	if categoryFilter == "" {
		groupings = validation.GroupDomainsByCategory()
	} else {
		// Parse category filter (comma-separated)
		categories := strings.Split(categoryFilter, ",")
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}

		// Get domains in specified categories
		filteredDomains := validation.GetDomainsInCategories(categories)

		// Group them by category for display
		groupingMap := make(map[string][]*types.DomainInfo)
		for _, domain := range filteredDomains {
			groupingMap[domain.Category] = append(groupingMap[domain.Category], domain)
		}

		// Convert to CategoryGrouping format
		for _, cat := range categories {
			if domains, exists := groupingMap[cat]; exists {
				groupings = append(groupings, validation.CategoryGrouping{
					Category: cat,
					Domains:  domains,
				})
			}
		}
	}

	return displayGroupedDomains(cmd, groupings)
}

// displayGroupedDomains renders grouped domains to output
func displayGroupedDomains(cmd *cobra.Command, groupings []validation.CategoryGrouping) error {

	// Print header
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", strings.Repeat("═", 80))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Available Domains by Category\n")
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", strings.Repeat("═", 80))

	// Print each category grouping
	for _, grouping := range groupings {
		// Category header
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "📁 %s (%d domain%s)\n", grouping.Category, len(grouping.Domains), pluralize(len(grouping.Domains)))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", strings.Repeat("─", 80))

		// List domains in category
		for _, domain := range grouping.Domains {
			tier := ""
			if domain.RequiresTier != "" && domain.RequiresTier != "Standard" {
				tier = fmt.Sprintf(" [Requires %s]", domain.RequiresTier)
			}

			preview := ""
			if domain.IsPreview {
				preview = " [PREVIEW]"
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %-30s%s%s\n", domain.Name, tier, preview)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    %s\n", domain.Description)

			if domain.Complexity != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Complexity: %s\n", domain.Complexity)
			}

			if len(domain.UseCases) > 0 {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Use cases: %s\n", strings.Join(domain.UseCases[:minInt(2, len(domain.UseCases))], ", "))
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n")
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n")
	}

	// Print summary
	distribution := validation.GetCategoryDistribution()
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", strings.Repeat("═", 80))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Summary\n")
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", strings.Repeat("─", 80))

	totalDomains := 0
	for _, cd := range distribution {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %-20s %2d domain%s\n", cd.Category, cd.Count, pluralize(cd.Count))
		totalDomains += cd.Count
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", strings.Repeat("─", 80))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  Total:             %2d domains\n", totalDomains)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", strings.Repeat("═", 80))

	return nil
}

// listDomainsByTier displays domains organized by subscription tier requirement
func listDomainsByTier(cmd *cobra.Command, args ...interface{}) error {
	tiers := []string{"Standard", "Advanced"}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", strings.Repeat("═", 80))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Available Domains by Subscription Tier\n")
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", strings.Repeat("═", 80))

	for _, tier := range tiers {
		// Filter domains by tier
		domains := make([]*types.DomainInfo, 0)
		for _, info := range types.DomainRegistry {
			if info.RequiresTier == tier {
				domains = append(domains, info)
			}
		}

		// Sort domains by display name for consistent ordering
		sort.Slice(domains, func(i, j int) bool {
			return domains[i].DisplayName < domains[j].DisplayName
		})

		// Tier header
		tierSymbol := getTierSymbol(tier)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s %s Tier (%d domain%s)\n", tierSymbol, tier, len(domains), pluralize(len(domains)))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", strings.Repeat("─", 80))

		if len(domains) == 0 {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  (No domains available)\n\n")
			continue
		}

		// List domains in tier
		for _, domain := range domains {
			category := ""
			if domain.Category != "" {
				category = fmt.Sprintf(" [%s]", domain.Category)
			}

			preview := ""
			if domain.IsPreview {
				preview = " [PREVIEW]"
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  %-30s%s%s\n", domain.Name, category, preview)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    %s\n\n", domain.Description)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n")
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n", strings.Repeat("═", 80))

	return nil
}

// getTierSymbol returns a symbol representing the tier level
func getTierSymbol(tier string) string {
	switch tier {
	case "Standard":
		return "🟢"
	case "Advanced":
		return "🟡"
	default:
		return "⭕"
	}
}

// pluralize returns appropriate singular/plural suffix
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	initDomainsListCmd()
	rootCmd.AddCommand(domainsListCmd)
	domainsListCmd.AddCommand(domainsListByTierCmd)
}
