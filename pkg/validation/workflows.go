package validation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// WorkflowSuggestion represents a suggested workflow involving multiple domains
type WorkflowSuggestion struct {
	Name        string   // e.g., "API Security Workflow"
	Description string   // e.g., "Secure and monitor API endpoints"
	Domains     []string // Domains involved in this workflow
	Category    string   // Primary category (e.g., "Security")
}

// DomainRelationship represents how two domains are related
type DomainRelationship struct {
	Domain1  string // First domain
	Domain2  string // Second domain
	Reason   string // Why they're related (e.g., "Same category", "Complementary use")
	Strength int    // 1-5 scale, 5 being strongest relationship
}

// GetRelatedDomains returns domains that work well with the given domain
func GetRelatedDomains(domain string) []*types.DomainInfo {
	info, found := types.GetDomainInfo(domain)
	if !found {
		return nil
	}

	relatedDomainNames := make(map[string]int) // domain -> strength score

	// Strategy 1: Same category domains (strength 4)
	for _, otherDomain := range GetDomainsByCategory(info.Category) {
		if otherDomain.Name != domain {
			relatedDomainNames[otherDomain.Name] = 4
		}
	}

	// Strategy 2: Complementary use cases (strength 3)
	allUseCases := GetAllUseCases()
	currentDomainUseCases := make(map[string]bool)
	for _, uc := range allUseCases {
		if uc.Domain == domain {
			// Extract key concepts from use cases
			words := strings.Fields(strings.ToLower(uc.Description))
			for _, word := range words {
				currentDomainUseCases[word] = true
			}
		}
	}

	// Find domains with overlapping use case concepts
	for _, uc := range allUseCases {
		if uc.Domain != domain {
			words := strings.Fields(strings.ToLower(uc.Description))
			matchCount := 0
			for _, word := range words {
				if currentDomainUseCases[word] {
					matchCount++
				}
			}
			if matchCount > 0 {
				relatedDomainNames[uc.Domain] = 3
			}
		}
	}

	// Strategy 3: Compatible tier (strength 2)
	for domainName, domainInfo := range types.DomainRegistry {
		if domainName != domain {
			// Domains with same or lower tier requirement are compatible
			if tierCompatible(info.RequiresTier, domainInfo.RequiresTier) {
				if _, exists := relatedDomainNames[domainName]; !exists {
					relatedDomainNames[domainName] = 2
				}
			}
		}
	}

	// Convert to sorted list
	type domainScore struct {
		name  string
		score int
	}
	var scores []domainScore
	for name, score := range relatedDomainNames {
		scores = append(scores, domainScore{name, score})
	}

	// Sort by strength descending
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].score != scores[j].score {
			return scores[i].score > scores[j].score
		}
		return scores[i].name < scores[j].name
	})

	// Limit to top 5 related domains
	maxResults := 5
	if len(scores) > maxResults {
		scores = scores[:maxResults]
	}

	// Convert to DomainInfo
	var related []*types.DomainInfo
	for _, score := range scores {
		if domainInfo, found := types.GetDomainInfo(score.name); found {
			related = append(related, domainInfo)
		}
	}

	return related
}

// GetWorkflowSuggestions returns suggested workflows for a domain
func GetWorkflowSuggestions(domain string) []WorkflowSuggestion {
	info, found := types.GetDomainInfo(domain)
	if !found {
		return nil
	}

	var suggestions []WorkflowSuggestion

	// Based on domain category, suggest relevant workflows
	switch info.Category {
	case "Security":
		suggestions = append(suggestions,
			WorkflowSuggestion{
				Name:        "API Security Workflow",
				Description: "Secure APIs with firewall and threat detection",
				Domains:     []string{"api", "application_firewall", "threat_campaign"},
				Category:    "Security",
			},
			WorkflowSuggestion{
				Name:        "Network Protection Workflow",
				Description: "Protect network infrastructure and applications",
				Domains:     []string{"network_security", "ddos", "infrastructure_protection"},
				Category:    "Security",
			},
		)
	case "Networking":
		suggestions = append(suggestions,
			WorkflowSuggestion{
				Name:        "Load Balancing Workflow",
				Description: "Configure and manage load balancing across regions",
				Domains:     []string{"dns", "virtual", "cdn"},
				Category:    "Networking",
			},
		)
	case "Platform":
		suggestions = append(suggestions,
			WorkflowSuggestion{
				Name:        "Access Management Workflow",
				Description: "Manage users and authentication for platform access",
				Domains:     []string{"authentication", "users", "tenant_and_identity"},
				Category:    "Platform",
			},
		)
	case "Infrastructure":
		suggestions = append(suggestions,
			WorkflowSuggestion{
				Name:        "Kubernetes Management Workflow",
				Description: "Deploy and manage Kubernetes clusters",
				Domains:     []string{"kubernetes", "service_mesh", "observability"},
				Category:    "Infrastructure",
			},
			WorkflowSuggestion{
				Name:        "Cloud Connectivity Workflow",
				Description: "Connect to cloud providers and manage cloud resources",
				Domains:     []string{"cloud_infrastructure", "site", "network"},
				Category:    "Infrastructure",
			},
		)
	case "Operations":
		suggestions = append(suggestions,
			WorkflowSuggestion{
				Name:        "Monitoring and Analytics Workflow",
				Description: "Monitor systems and collect analytics data",
				Domains:     []string{"observability", "statistics", "telemetry_and_insights"},
				Category:    "Operations",
			},
		)
	}

	return suggestions
}

// FormatRelatedDomains formats related domains for display
func FormatRelatedDomains(domains []*types.DomainInfo) string {
	if len(domains) == 0 {
		return ""
	}

	var formatted strings.Builder
	formatted.WriteString("\nRELATED DOMAINS:\n")

	for i, domain := range domains {
		if i < 5 { // Show max 5 related domains
			formatted.WriteString(fmt.Sprintf("  • %s - %s\n", domain.Name, domain.Description))
		}
	}

	return formatted.String()
}

// FormatWorkflowSuggestions formats workflow suggestions for display
func FormatWorkflowSuggestions(workflows []WorkflowSuggestion) string {
	if len(workflows) == 0 {
		return ""
	}

	var formatted strings.Builder
	formatted.WriteString("\nSUGGESTED WORKFLOWS:\n")

	for i, workflow := range workflows {
		if i < 3 { // Show max 3 workflows
			formatted.WriteString(fmt.Sprintf("  • %s\n", workflow.Name))
			formatted.WriteString(fmt.Sprintf("    %s\n", workflow.Description))
			formatted.WriteString(fmt.Sprintf("    Involves: %s\n", strings.Join(workflow.Domains, ", ")))
		}
	}

	return formatted.String()
}

// GetWorkflowsByCategory returns all workflows in a specific category
func GetWorkflowsByCategory(category string) []WorkflowSuggestion {
	// Get all domains in category
	domains := GetDomainsByCategory(category)
	if len(domains) == 0 {
		return nil
	}

	// Collect all workflows for domains in this category
	workflowMap := make(map[string]WorkflowSuggestion)
	for _, domain := range domains {
		suggestions := GetWorkflowSuggestions(domain.Name)
		for _, suggestion := range suggestions {
			if suggestion.Category == category {
				workflowMap[suggestion.Name] = suggestion
			}
		}
	}

	// Convert to slice
	var workflows []WorkflowSuggestion
	for _, workflow := range workflowMap {
		workflows = append(workflows, workflow)
	}

	// Sort by name
	sort.Slice(workflows, func(i, j int) bool {
		return workflows[i].Name < workflows[j].Name
	})

	return workflows
}

// Helper function: check if tiers are compatible
func tierCompatible(tier1, tier2 string) bool {
	tierOrder := map[string]int{
		"Standard":     1,
		"Professional": 2,
		"Enterprise":   3,
	}

	level1 := tierOrder[tier1]
	level2 := tierOrder[tier2]

	if level1 == 0 || level2 == 0 {
		return true // If tier unknown, assume compatible
	}

	// Compatible if tier2 is same or lower than tier1
	return level2 <= level1
}
