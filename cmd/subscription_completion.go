package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
)

// Completion client with caching for subscription completions
var subscriptionCompletionClient *subscription.Client

// getSubscriptionCompletionClient returns a subscription client optimized for tab completion.
// Uses a short timeout to avoid blocking the shell.
func getSubscriptionCompletionClient() *subscription.Client {
	if subscriptionCompletionClient == nil {
		subscriptionCompletionClient = GetSubscriptionClient()
	}
	return subscriptionCompletionClient
}

// Static completion functions

// completeAddonFilter provides completion for the --filter flag in subscription addons.
func completeAddonFilter(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"active\tShow only active/subscribed addons",
		"available\tShow addons available for activation",
		"all\tShow all addon services including denied",
	}, cobra.ShellCompDirectiveNoFileComp
}

// Dynamic completion functions

// completeAddonName provides completion for the --addon flag in subscription activate.
// Fetches available addons from the API with caching.
func completeAddonName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getSubscriptionCompletionClient()
	if client == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Use a short timeout context for completion
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	namespace := GetSubscriptionNamespace()
	addons, err := client.GetAddonServices(ctx, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, addon := range addons {
		// Only suggest addons that can be activated
		if addon.CanActivate() {
			description := addon.DisplayName
			if description == "" {
				description = addon.Name
			}
			// Add activation type hint
			if addon.IsSelfActivation() {
				description += " (instant)"
			} else if addon.IsManagedActivation() {
				description += " (managed)"
			}
			completions = append(completions, addon.Name+"\t"+description)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completePendingAddonName provides completion for addons with pending activation.
// Used for activation-status --addon flag.
func completePendingAddonName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getSubscriptionCompletionClient()
	if client == nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Use a short timeout context for completion
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	namespace := GetSubscriptionNamespace()
	addons, err := client.GetAddonServices(ctx, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, addon := range addons {
		// Suggest both pending and active addons for status check
		if addon.IsPending() || addon.State == subscription.StateSubscribed {
			description := addon.DisplayName
			if description == "" {
				description = addon.Name
			}
			if addon.IsPending() {
				description += " (pending)"
			} else {
				description += " (active)"
			}
			completions = append(completions, addon.Name+"\t"+description)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeResourceType provides completion for --resource-type in subscription validate.
func completeResourceType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"http_loadbalancer\tHTTP Load Balancer",
		"tcp_loadbalancer\tTCP Load Balancer",
		"origin_pool\tOrigin Pool",
		"healthcheck\tHealth Check",
		"app_firewall\tApplication Firewall (WAF)",
		"service_policy\tService Policy",
		"rate_limiter\tRate Limiter",
		"api_definition\tAPI Definition",
		"api_group\tAPI Group",
		"forward_proxy_policy\tForward Proxy Policy",
		"network_policy\tNetwork Policy",
		"network_firewall\tNetwork Firewall",
		"dns_zone\tDNS Zone",
		"dns_lb_pool\tDNS Load Balancer Pool",
	}, cobra.ShellCompDirectiveNoFileComp
}

// completeFeatureName provides completion for --feature in subscription validate.
func completeFeatureName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client := getSubscriptionCompletionClient()
	if client == nil {
		// Fall back to common features if no client
		return []string{
			"bot-defense\tBot Defense",
			"api-security\tAPI Security",
			"client-side-defense\tClient-Side Defense",
			"malicious-user-detection\tMalicious User Detection",
			"synthetic-monitoring\tSynthetic Monitoring",
			"web-app-scanning\tWeb Application Scanning",
		}, cobra.ShellCompDirectiveNoFileComp
	}

	// Use a short timeout context for completion
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	namespace := GetSubscriptionNamespace()
	addons, err := client.GetAddonServices(ctx, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	completions := []string{}
	for _, addon := range addons {
		description := addon.DisplayName
		if description == "" {
			description = addon.Name
		}
		// Add status hint
		switch addon.State {
		case subscription.StateSubscribed:
			description += " [active]"
		case subscription.StatePending:
			description += " [pending]"
		default:
			if addon.CanActivate() {
				description += " [available]"
			}
		}
		completions = append(completions, addon.Name+"\t"+description)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
