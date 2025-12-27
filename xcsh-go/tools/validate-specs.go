// Validate upstream API specifications for quality and consistency
// This script checks for common spec organization issues and reports findings
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

// SpecIndex represents the upstream spec index
type SpecIndex struct {
	Version        string `json:"version"`
	Specifications []struct {
		Domain      string `json:"domain"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PathCount   int    `json:"path_count"`
		SchemaCount int    `json:"schema_count"`
	} `json:"specifications"`
}

// ValidationIssue represents a spec quality issue
type ValidationIssue struct {
	Severity string // "critical", "warning", "info"
	Domain   string
	Issue    string
	Details  string
}

// DomainSummary represents a domain's spec coverage
type DomainSummary struct {
	Domain      string `json:"domain"`
	PathCount   int    `json:"path_count"`
	SchemaCount int    `json:"schema_count"`
}

func main() {
	indexPath := flag.String("index", ".specs/index.json", "Path to spec index")
	report := flag.Bool("report", false, "Generate JSON report")
	verbose := flag.Bool("v", false, "Verbose output")

	flag.Parse()

	log.Println("🔍 Validating upstream API specifications...")

	// Read spec index
	indexData, err := os.ReadFile(*indexPath)
	if err != nil {
		log.Fatalf("Failed to read spec index: %v", err)
	}

	var index SpecIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		log.Fatalf("Failed to parse spec index: %v", err)
	}

	// Run validation checks
	issues := validateSpecs(index, *verbose)

	// Print results
	if len(issues) > 0 {
		log.Printf("\n⚠️  Found %d validation issues:\n", len(issues))

		// Group by severity
		criticals := filterBySeverity(issues, "critical")
		warnings := filterBySeverity(issues, "warning")
		infos := filterBySeverity(issues, "info")

		if len(criticals) > 0 {
			log.Printf("🔴 CRITICAL (%d):", len(criticals))
			for _, issue := range criticals {
				log.Printf("  - %s (%s): %s", issue.Domain, issue.Issue, issue.Details)
			}
		}

		if len(warnings) > 0 {
			log.Printf("🟡 WARNING (%d):", len(warnings))
			for _, issue := range warnings {
				log.Printf("  - %s (%s): %s", issue.Domain, issue.Issue, issue.Details)
			}
		}

		if len(infos) > 0 && *verbose {
			log.Printf("ℹ️  INFO (%d):", len(infos))
			for _, issue := range infos {
				log.Printf("  - %s (%s): %s", issue.Domain, issue.Issue, issue.Details)
			}
		}
	} else {
		log.Println("✅ All validations passed")
	}

	// Print summary
	printSummary(index, issues)

	// Output report if requested
	if *report {
		outputReport(index, issues)
	}

	// Exit with error if critical issues found
	if len(filterBySeverity(issues, "critical")) > 0 {
		os.Exit(1)
	}
}

func validateSpecs(index SpecIndex, verbose bool) []ValidationIssue {
	var issues []ValidationIssue

	// Check version format
	if index.Version == "" {
		issues = append(issues, ValidationIssue{
			Severity: "critical",
			Domain:   "index",
			Issue:    "missing_version",
			Details:  "Spec index missing version field",
		})
	}

	// Validate individual domains
	domainMap := make(map[string]bool)
	for _, spec := range index.Specifications {
		domainMap[spec.Domain] = true

		// Check for empty domains (red flag)
		if spec.PathCount == 0 && spec.SchemaCount == 0 {
			issues = append(issues, ValidationIssue{
				Severity: "info",
				Domain:   spec.Domain,
				Issue:    "empty_domain",
				Details:  "No paths or schemas defined",
			})
		}

		// Check for domains with only paths or only schemas (unusual)
		if (spec.PathCount > 0 && spec.SchemaCount == 0) || (spec.PathCount == 0 && spec.SchemaCount > 0) {
			issues = append(issues, ValidationIssue{
				Severity: "info",
				Domain:   spec.Domain,
				Issue:    "asymmetric_coverage",
				Details:  fmt.Sprintf("Has %d paths but %d schemas", spec.PathCount, spec.SchemaCount),
			})
		}

		// Check naming consistency
		if !isValidDomainName(spec.Domain) {
			issues = append(issues, ValidationIssue{
				Severity: "warning",
				Domain:   spec.Domain,
				Issue:    "invalid_domain_name",
				Details:  "Domain name doesn't follow snake_case convention",
			})
		}

		if verbose {
			log.Printf("  ✓ Validated %s (%d paths, %d schemas)", spec.Domain, spec.PathCount, spec.SchemaCount)
		}
	}

	// Check for duplicates
	if len(domainMap) != len(index.Specifications) {
		issues = append(issues, ValidationIssue{
			Severity: "critical",
			Domain:   "index",
			Issue:    "duplicate_domains",
			Details:  "Duplicate domain names detected",
		})
	}

	return issues
}

func isValidDomainName(domain string) bool {
	// Check if domain follows snake_case pattern (lowercase letters, numbers, underscores)
	for _, char := range domain {
		isLower := char >= 'a' && char <= 'z'
		isDigit := char >= '0' && char <= '9'
		isUnderscore := char == '_'
		if !isLower && !isDigit && !isUnderscore {
			return false
		}
	}
	return true
}

func filterBySeverity(issues []ValidationIssue, severity string) []ValidationIssue {
	var filtered []ValidationIssue
	for _, issue := range issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func printSummary(index SpecIndex, issues []ValidationIssue) {
	totalPaths := 0
	totalSchemas := 0
	nonEmptyDomains := 0

	for _, spec := range index.Specifications {
		totalPaths += spec.PathCount
		totalSchemas += spec.SchemaCount
		if spec.PathCount > 0 || spec.SchemaCount > 0 {
			nonEmptyDomains++
		}
	}

	log.Printf("\nValidation Summary:")
	log.Printf("  Spec version: %s", index.Version)
	log.Printf("  Total domains: %d", len(index.Specifications))
	log.Printf("  Non-empty domains: %d", nonEmptyDomains)
	log.Printf("  Total API paths: %d", totalPaths)
	log.Printf("  Total schemas: %d", totalSchemas)
	log.Printf("  Total issues: %d", len(issues))

	if len(issues) > 0 {
		log.Printf("    - Critical: %d", len(filterBySeverity(issues, "critical")))
		log.Printf("    - Warning: %d", len(filterBySeverity(issues, "warning")))
		log.Printf("    - Info: %d", len(filterBySeverity(issues, "info")))
	}
}

func outputReport(index SpecIndex, issues []ValidationIssue) {
	type Report struct {
		Version string            `json:"version"`
		Summary map[string]int    `json:"summary"`
		Issues  []ValidationIssue `json:"issues"`
		Domains []DomainSummary   `json:"domains"`
	}

	// Build report
	summary := make(map[string]int)
	summary["total_domains"] = len(index.Specifications)
	summary["total_issues"] = len(issues)
	summary["critical_issues"] = len(filterBySeverity(issues, "critical"))
	summary["warning_issues"] = len(filterBySeverity(issues, "warning"))
	summary["info_issues"] = len(filterBySeverity(issues, "info"))

	// Count paths and schemas
	totalPaths := 0
	totalSchemas := 0
	for _, spec := range index.Specifications {
		totalPaths += spec.PathCount
		totalSchemas += spec.SchemaCount
	}
	summary["total_paths"] = totalPaths
	summary["total_schemas"] = totalSchemas

	// Build domain summaries
	var domains []DomainSummary
	for _, spec := range index.Specifications {
		domains = append(domains, DomainSummary{
			Domain:      spec.Domain,
			PathCount:   spec.PathCount,
			SchemaCount: spec.SchemaCount,
		})
	}

	// Sort for reproducibility
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Domain < domains[j].Domain
	})

	report := Report{
		Version: index.Version,
		Summary: summary,
		Issues:  issues,
		Domains: domains,
	}

	// Output JSON
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
}
