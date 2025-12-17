package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
)

var (
	validateResourceType string
	validateCount        int
	validateFeature      string
)

var subscriptionValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate resources against subscription capabilities.",
	Long: `Validate planned resources against subscription quotas and feature availability.

Use this command before deploying resources to ensure deployment will succeed.
Checks quota limits for resource types and verifies that required addon services
are subscribed.

Exit codes:
  0 - All validations passed
  1 - Generic error (API failure, invalid arguments)
  2 - Validation failed (quota exceeded or feature unavailable)

AI assistants should run this validation before 'terraform apply' to catch
quota and feature issues early, preventing deployment failures.`,
	Example: `  # Validate if you can create 5 more HTTP load balancers
  f5xcctl subscription validate --resource-type http_loadbalancer --count 5

  # Validate if bot-defense feature is available
  f5xcctl subscription validate --feature bot-defense

  # Validate multiple resources at once
  f5xcctl subscription validate --resource-type origin_pool --count 10

  # Get validation result as JSON for automation
  f5xcctl subscription validate --resource-type http_loadbalancer --count 5 --output-format json

  # Validate in a specific namespace
  f5xcctl subscription validate --resource-type http_loadbalancer --count 5 -n shared`,
	RunE: runSubscriptionValidate,
}

func init() {
	subscriptionCmd.AddCommand(subscriptionValidateCmd)

	subscriptionValidateCmd.Flags().StringVar(&validateResourceType, "resource-type", "", "Resource type to validate quota for (e.g., http_loadbalancer, origin_pool).")
	subscriptionValidateCmd.Flags().IntVar(&validateCount, "count", 1, "Number of resources to create (default: 1).")
	subscriptionValidateCmd.Flags().StringVar(&validateFeature, "feature", "", "Feature/addon to validate availability (e.g., bot-defense, api-security).")

	// Register completions for validate flags
	_ = subscriptionValidateCmd.RegisterFlagCompletionFunc("resource-type", completeResourceType)
	_ = subscriptionValidateCmd.RegisterFlagCompletionFunc("feature", completeFeatureName)
}

func runSubscriptionValidate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	client, err := requireSubscriptionClient()
	if err != nil {
		return err
	}

	// Require at least one validation type
	if validateResourceType == "" && validateFeature == "" {
		return fmt.Errorf("at least one of --resource-type or --feature is required\n\nUsage:\n  f5xcctl subscription validate --resource-type <type> --count <n>\n  f5xcctl subscription validate --feature <name>")
	}

	namespace := GetSubscriptionNamespace()
	result := &subscription.ValidationResult{
		Valid:    true,
		Checks:   []subscription.ValidationCheck{},
		Warnings: []string{},
		Errors:   []string{},
	}

	// Validate resource quota if specified
	if validateResourceType != "" {
		req := subscription.ValidationRequest{
			ResourceType: validateResourceType,
			Count:        validateCount,
			Namespace:    namespace,
		}
		quotaResult, err := client.ValidateResource(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to validate resource quota: %w", err)
		}

		result.Checks = append(result.Checks, quotaResult.Checks...)
		result.Warnings = append(result.Warnings, quotaResult.Warnings...)
		result.Errors = append(result.Errors, quotaResult.Errors...)

		if !quotaResult.Valid {
			result.Valid = false
		}
	}

	// Validate feature availability if specified
	if validateFeature != "" {
		featureResult, err := validateFeatureAvailability(ctx, client, validateFeature)
		if err != nil {
			return fmt.Errorf("failed to validate feature: %w", err)
		}

		result.Checks = append(result.Checks, featureResult.Checks...)
		result.Warnings = append(result.Warnings, featureResult.Warnings...)
		result.Errors = append(result.Errors, featureResult.Errors...)

		if !featureResult.Valid {
			result.Valid = false
		}
	}

	// Output based on format
	format := GetOutputFormatWithDefault("table")
	if err := formatOutputWithTableFallback(result, format, func() error {
		return outputValidationTable(result)
	}); err != nil {
		return err
	}

	// Exit with code 2 if validation failed
	if !result.Valid {
		os.Exit(2)
	}

	return nil
}

func validateFeatureAvailability(ctx context.Context, client *subscription.Client, feature string) (*subscription.ValidationResult, error) {
	result := &subscription.ValidationResult{
		Valid:    true,
		Checks:   []subscription.ValidationCheck{},
		Warnings: []string{},
		Errors:   []string{},
	}

	// Get addon services
	addons, err := client.GetAddonServices(ctx, "system")
	if err != nil {
		return nil, fmt.Errorf("failed to get addon services: %w", err)
	}

	// Find the requested feature
	var foundAddon *subscription.AddonServiceInfo
	featureLower := strings.ToLower(feature)

	for i := range addons {
		nameLower := strings.ToLower(addons[i].Name)
		if nameLower == featureLower || strings.Contains(nameLower, featureLower) {
			foundAddon = &addons[i]
			break
		}
	}

	check := subscription.ValidationCheck{
		Type:    "feature",
		Feature: feature,
	}

	if foundAddon == nil {
		check.Result = "FAIL"
		check.Message = fmt.Sprintf("Feature '%s' not found in available addon services", feature)
		result.Checks = append(result.Checks, check)
		result.Errors = append(result.Errors, check.Message)
		result.Valid = false
		return result, nil
	}

	check.RequiredTier = subscription.TierDescription(foundAddon.Tier)
	check.Status = subscription.StateDescription(foundAddon.State)

	if foundAddon.IsActive() {
		check.Result = "PASS"
		check.Message = fmt.Sprintf("Feature '%s' is actively subscribed", feature)
		result.Checks = append(result.Checks, check)
	} else if foundAddon.IsAvailable() {
		check.Result = "WARN"
		check.Message = fmt.Sprintf("Feature '%s' is available but not yet subscribed", feature)
		result.Checks = append(result.Checks, check)
		result.Warnings = append(result.Warnings, check.Message)
	} else {
		check.Result = "FAIL"

		switch foundAddon.AccessStatus {
		case subscription.AccessUpgradeRequired:
			check.Message = fmt.Sprintf("Feature '%s' requires a plan upgrade to access", feature)
		case subscription.AccessContactSales:
			check.Message = fmt.Sprintf("Feature '%s' requires contacting F5 sales", feature)
		case subscription.AccessDenied:
			check.Message = fmt.Sprintf("Feature '%s' access is denied by policy", feature)
		default:
			check.Message = fmt.Sprintf("Feature '%s' is not available (status: %s)", feature, foundAddon.AccessStatus)
		}

		result.Checks = append(result.Checks, check)
		result.Errors = append(result.Errors, check.Message)
		result.Valid = false
	}

	return result, nil
}

func outputValidationTable(result *subscription.ValidationResult) error {
	// Print header with overall result
	if result.Valid {
		fmt.Println("VALIDATION RESULT: PASS")
	} else {
		fmt.Println("VALIDATION RESULT: FAIL")
	}
	fmt.Println(strings.Repeat("=", 75))
	fmt.Println()

	// Print checks
	if len(result.Checks) > 0 {
		fmt.Println("CHECKS")
		fmt.Println(strings.Repeat("-", 75))

		for _, check := range result.Checks {
			resultIndicator := getResultIndicator(check.Result)

			switch check.Type {
			case "quota":
				fmt.Printf("  %s QUOTA: %s\n", resultIndicator, check.Resource)
				fmt.Printf("      Current: %d | Requested: +%d | Limit: %d | After: %d\n",
					check.Current, check.Requested, check.Limit, check.Current+check.Requested)
				if check.Message != "" {
					fmt.Printf("      %s\n", check.Message)
				}

			case "feature":
				fmt.Printf("  %s FEATURE: %s\n", resultIndicator, check.Feature)
				if check.RequiredTier != "" {
					fmt.Printf("      Tier: %s | Status: %s\n", check.RequiredTier, check.Status)
				}
				if check.Message != "" {
					fmt.Printf("      %s\n", check.Message)
				}
			}
			fmt.Println()
		}
	}

	// Print warnings
	if len(result.Warnings) > 0 {
		fmt.Println("WARNINGS")
		fmt.Println(strings.Repeat("-", 75))
		for _, warning := range result.Warnings {
			fmt.Printf("  (!) %s\n", warning)
		}
		fmt.Println()
	}

	// Print errors
	if len(result.Errors) > 0 {
		fmt.Println("ERRORS")
		fmt.Println(strings.Repeat("-", 75))
		for _, err := range result.Errors {
			fmt.Printf("  (!!!) %s\n", err)
		}
		fmt.Println()
	}

	// Print summary and hints
	fmt.Println("SUMMARY")
	fmt.Println(strings.Repeat("-", 75))
	passCount := 0
	failCount := 0
	warnCount := 0
	for _, check := range result.Checks {
		switch check.Result {
		case "PASS":
			passCount++
		case "FAIL":
			failCount++
		case "WARN":
			warnCount++
		}
	}
	fmt.Printf("  Passed: %d | Warnings: %d | Failed: %d\n", passCount, warnCount, failCount)
	fmt.Println()

	if !result.Valid {
		fmt.Println("HINTS")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Println("  Deployment may fail due to validation errors above.")
		fmt.Println("  Use 'f5xcctl subscription quota' to see current quota usage.")
		fmt.Println("  Use 'f5xcctl subscription addons' to see available addon services.")
		fmt.Println()
	}

	return nil
}

func getResultIndicator(result string) string {
	switch result {
	case "PASS":
		return "[PASS]"
	case "FAIL":
		return "[FAIL]"
	case "WARN":
		return "[WARN]"
	default:
		return "[????]"
	}
}
