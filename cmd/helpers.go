package cmd

import (
	"fmt"

	"github.com/robinmordasiewicz/f5xcctl/pkg/cloudstatus"
	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
	"github.com/robinmordasiewicz/f5xcctl/pkg/subscription"
)

// requireSubscriptionClient returns the subscription client or an error if not initialized.
// This centralizes the nil-check pattern used across subscription commands.
func requireSubscriptionClient() (*subscription.Client, error) {
	client := GetSubscriptionClient()
	if client == nil {
		return nil, fmt.Errorf("subscription client not initialized - check authentication")
	}
	return client, nil
}

// requireCloudStatusClient returns the cloudstatus client or an error if not initialized.
// This centralizes the nil-check pattern used across cloudstatus commands.
func requireCloudStatusClient() (*cloudstatus.Client, error) {
	client := GetCloudStatusClient()
	if client == nil {
		return nil, fmt.Errorf("cloudstatus client not initialized")
	}
	return client, nil
}

// formatOutputWithTableFallback outputs data in JSON or YAML format,
// or calls the provided table formatter function for table/text output.
func formatOutputWithTableFallback(data interface{}, format string, tableFn func() error) error {
	switch format {
	case "json", "yaml":
		return output.Print(data, format)
	default:
		return tableFn()
	}
}
