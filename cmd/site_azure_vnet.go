package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/f5xcctl/pkg/output"
)

var azureVNetFlags struct {
	name            string
	namespace       string
	inputFile       string
	region          string
	resourceGroup   string
	vnetCIDR        string
	machineType     string
	sshKey          string
	cloudCreds      string
	subscriptionID  string
	terraformDir    string
	terraformAction string
	autoApprove     bool
	wait            bool
}

var azureVNetCmd = &cobra.Command{
	Use:   "azure_vnet",
	Short: "Manage Azure VNet site creation through view apis",
	Long: `Manage Azure VNet sites in F5 Distributed Cloud.

Azure VNet sites allow you to deploy F5 XC Customer Edge (CE) nodes in your
Azure Virtual Network, enabling secure connectivity and edge services.`,
	Example: `  # Create an Azure VNet site from a YAML file
  f5xcctl site azure_vnet create -i azure-site.yaml

  # Delete an Azure VNet site
  f5xcctl site azure_vnet delete --name example-site

  # Run Terraform to provision infrastructure
  f5xcctl site azure_vnet run --name example-site --action apply --auto-approve`,
}

var azureVNetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Azure VNet volterra site",
	Long: `Create a new Azure VNet site in F5 Distributed Cloud.

This command registers an Azure VNet site configuration. After creation,
use 'f5xcctl site azure_vnet run --action apply' to provision the infrastructure.

You can provide the site specification via:
- YAML/JSON file using --input-file
- Command line flags for common options`,
	Example: `  # Create from YAML file
  f5xcctl site azure_vnet create -i azure-site.yaml

  # Create with command line flags
  f5xcctl site azure_vnet create --name example-site --region eastus \
    --vnet-cidr 10.0.0.0/16 --cloud-creds example-azure-creds

  # Create with resource group
  f5xcctl site azure_vnet create --name example-site --region eastus \
    --resource-group example-rg --cloud-creds example-azure-creds`,
	RunE: runAzureVNetCreate,
}

var azureVNetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Azure VNet site",
	Long: `Delete an Azure VNet site from F5 Distributed Cloud.

Note: This only removes the site configuration from F5 XC. To fully clean up
Azure resources, first run 'f5xcctl site azure_vnet run --action destroy' before
deleting the site configuration.`,
	Example: `  # Delete a site (after destroying infrastructure)
  f5xcctl site azure_vnet delete --name example-site

  # Delete from a specific namespace
  f5xcctl site azure_vnet delete --name example-site -n system`,
	RunE: runAzureVNetDelete,
}

var azureVNetReplaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace Azure VNet site",
	Long: `Replace an existing Azure VNet site configuration in F5 Distributed Cloud.

This updates the site specification. After replacing, you may need to
run 'f5xcctl site azure_vnet run --action apply' to apply infrastructure changes.`,
	Example: `  # Replace site configuration from file
  f5xcctl site azure_vnet replace -i updated-site.yaml

  # Replace with specific name
  f5xcctl site azure_vnet replace --name example-site -i updated-site.yaml`,
	RunE: runAzureVNetReplace,
}

var azureVNetRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run terraform action, valid actions are plan, apply and destroy",
	Long: `Run Terraform actions to provision or destroy Azure VNet site infrastructure.

This command retrieves Terraform parameters from F5 XC and executes Terraform
to manage the actual Azure resources (VNet, subnets, VMs, etc.).

Available actions:
- plan: Preview changes without applying
- apply: Create or update infrastructure
- destroy: Remove all infrastructure`,
	Example: `  # Preview infrastructure changes
  f5xcctl site azure_vnet run --name example-site --action plan

  # Apply infrastructure (with confirmation prompt)
  f5xcctl site azure_vnet run --name example-site --action apply

  # Apply infrastructure automatically (for CI/CD)
  f5xcctl site azure_vnet run --name example-site --action apply --auto-approve

  # Destroy infrastructure
  f5xcctl site azure_vnet run --name example-site --action destroy --auto-approve`,
	RunE: runAzureVNetTerraform,
}

func init() {
	siteCmd.AddCommand(azureVNetCmd)

	// Create command
	azureVNetCreateCmd.Flags().StringVarP(&azureVNetFlags.inputFile, "input-file", "i", "", "Input file (YAML/JSON) containing site definition")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.name, "name", "", "Site name")
	azureVNetCreateCmd.Flags().StringVarP(&azureVNetFlags.namespace, "namespace", "n", "system", "Namespace")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.region, "region", "", "Azure region")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.resourceGroup, "resource-group", "", "Azure resource group")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.vnetCIDR, "vnet-cidr", "", "VNet CIDR block")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.machineType, "machine-type", "Standard_D3_v2", "Azure VM size")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.sshKey, "ssh-key", "", "SSH public key")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.cloudCreds, "cloud-creds", "", "Cloud credentials name")
	azureVNetCreateCmd.Flags().StringVar(&azureVNetFlags.subscriptionID, "subscription-id", "", "Azure subscription ID")
	azureVNetCmd.AddCommand(azureVNetCreateCmd)

	// Delete command
	azureVNetDeleteCmd.Flags().StringVar(&azureVNetFlags.name, "name", "", "Site name (required)")
	azureVNetDeleteCmd.Flags().StringVarP(&azureVNetFlags.namespace, "namespace", "n", "system", "Namespace")
	_ = azureVNetDeleteCmd.MarkFlagRequired("name")
	azureVNetCmd.AddCommand(azureVNetDeleteCmd)

	// Replace command
	azureVNetReplaceCmd.Flags().StringVarP(&azureVNetFlags.inputFile, "input-file", "i", "", "Input file (YAML/JSON) containing site definition")
	azureVNetReplaceCmd.Flags().StringVar(&azureVNetFlags.name, "name", "", "Site name")
	azureVNetReplaceCmd.Flags().StringVarP(&azureVNetFlags.namespace, "namespace", "n", "system", "Namespace")
	azureVNetCmd.AddCommand(azureVNetReplaceCmd)

	// Run (Terraform) command
	azureVNetRunCmd.Flags().StringVar(&azureVNetFlags.name, "name", "", "Site name (required)")
	azureVNetRunCmd.Flags().StringVarP(&azureVNetFlags.namespace, "namespace", "n", "system", "Namespace")
	azureVNetRunCmd.Flags().StringVar(&azureVNetFlags.terraformAction, "action", "plan", "Terraform action: plan, apply, destroy")
	azureVNetRunCmd.Flags().StringVar(&azureVNetFlags.terraformDir, "terraform-dir", "", "Directory for Terraform files (default: temp dir)")
	azureVNetRunCmd.Flags().BoolVar(&azureVNetFlags.autoApprove, "auto-approve", false, "Auto-approve Terraform apply/destroy")
	azureVNetRunCmd.Flags().BoolVar(&azureVNetFlags.wait, "wait", true, "Wait for operation to complete")
	_ = azureVNetRunCmd.MarkFlagRequired("name")
	azureVNetCmd.AddCommand(azureVNetRunCmd)
}

func runAzureVNetCreate(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	var siteSpec map[string]interface{}

	if azureVNetFlags.inputFile != "" {
		data, err := os.ReadFile(azureVNetFlags.inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		if err := yaml.Unmarshal(data, &siteSpec); err != nil {
			if err := json.Unmarshal(data, &siteSpec); err != nil {
				return fmt.Errorf("failed to parse input file: %w", err)
			}
		}
	} else {
		if azureVNetFlags.name == "" || azureVNetFlags.region == "" {
			return fmt.Errorf("--name and --region are required when not using input file")
		}

		siteSpec = buildAzureVNetSiteSpec()
	}

	// Extract metadata if present, or build from flags
	metadata := map[string]interface{}{
		"name":      azureVNetFlags.name,
		"namespace": azureVNetFlags.namespace,
	}
	if meta, ok := siteSpec["metadata"].(map[string]interface{}); ok {
		metadata = meta
		if azureVNetFlags.name != "" {
			metadata["name"] = azureVNetFlags.name
		}
		if azureVNetFlags.namespace != "" {
			metadata["namespace"] = azureVNetFlags.namespace
		}
	}

	spec := siteSpec
	if s, ok := siteSpec["spec"].(map[string]interface{}); ok {
		spec = s
	}

	requestBody := map[string]interface{}{
		"metadata": metadata,
		"spec":     spec,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	path := fmt.Sprintf("/api/config/namespaces/%s/azure_vnet_sites", azureVNetFlags.namespace)
	resp, err := apiClient.Post(ctx, path, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create Azure VNet site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("Azure VNet site '%s' created successfully", azureVNetFlags.name))
	return output.Print(result, GetOutputFormat())
}

func runAzureVNetDelete(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := fmt.Sprintf("/api/config/namespaces/%s/azure_vnet_sites/%s", azureVNetFlags.namespace, azureVNetFlags.name)
	resp, err := apiClient.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete Azure VNet site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	output.PrintInfo(fmt.Sprintf("Azure VNet site '%s' deleted successfully", azureVNetFlags.name))
	return nil
}

func runAzureVNetReplace(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	if azureVNetFlags.inputFile == "" {
		return fmt.Errorf("--input-file is required for replace operation")
	}

	data, err := os.ReadFile(azureVNetFlags.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	var siteSpec map[string]interface{}
	if err := yaml.Unmarshal(data, &siteSpec); err != nil {
		if err := json.Unmarshal(data, &siteSpec); err != nil {
			return fmt.Errorf("failed to parse input file: %w", err)
		}
	}

	// Extract name from spec or flags
	name := azureVNetFlags.name
	if meta, ok := siteSpec["metadata"].(map[string]interface{}); ok {
		if n, ok := meta["name"].(string); ok && name == "" {
			name = n
		}
	}
	if name == "" {
		return fmt.Errorf("site name is required (via --name or in input file)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	path := fmt.Sprintf("/api/config/namespaces/%s/azure_vnet_sites/%s", azureVNetFlags.namespace, name)
	resp, err := apiClient.Put(ctx, path, siteSpec)
	if err != nil {
		return fmt.Errorf("failed to replace Azure VNet site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("Azure VNet site '%s' replaced successfully", name))
	return output.Print(result, GetOutputFormat())
}

func runAzureVNetTerraform(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Validate action
	action := azureVNetFlags.terraformAction
	switch action {
	case "plan", "apply", "destroy":
		// valid
	default:
		return fmt.Errorf("invalid action: %s (must be plan, apply, or destroy)", action)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// Get Terraform parameters from the site
	output.PrintInfo(fmt.Sprintf("Retrieving Terraform parameters for site '%s'...", azureVNetFlags.name))

	tfParamsPath := fmt.Sprintf("/api/terraform/namespaces/%s/terraform/azure_vnet_site/%s", azureVNetFlags.namespace, azureVNetFlags.name)
	resp, err := apiClient.Get(ctx, tfParamsPath, nil)
	if err != nil {
		return fmt.Errorf("failed to get Terraform parameters: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var tfParams map[string]interface{}
	if err := json.Unmarshal(resp.Body, &tfParams); err != nil {
		return fmt.Errorf("failed to parse Terraform parameters: %w", err)
	}

	// Setup Terraform directory
	tfDir := azureVNetFlags.terraformDir
	if tfDir == "" {
		tfDir = filepath.Join(os.TempDir(), "f5xcctl-terraform", azureVNetFlags.name)
	}

	if err := os.MkdirAll(tfDir, 0755); err != nil {
		return fmt.Errorf("failed to create Terraform directory: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("Terraform directory: %s", tfDir))

	// Write terraform.tfvars.json
	tfVarsPath := filepath.Join(tfDir, "terraform.tfvars.json")
	tfVarsData, err := json.MarshalIndent(tfParams, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Terraform vars: %w", err)
	}
	if err := os.WriteFile(tfVarsPath, tfVarsData, 0644); err != nil {
		return fmt.Errorf("failed to write terraform.tfvars.json: %w", err)
	}

	// Check if Terraform is installed
	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		return fmt.Errorf("terraform not found in PATH - please install Terraform")
	}

	// Run Terraform init if needed
	stateFile := filepath.Join(tfDir, "terraform.tfstate")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		output.PrintInfo("Running terraform init...")
		initCmd := exec.CommandContext(ctx, tfPath, "init")
		initCmd.Dir = tfDir
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		if err := initCmd.Run(); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)
		}
	}

	// Run the Terraform action
	output.PrintInfo(fmt.Sprintf("Running terraform %s...", action))

	tfArgs := []string{action}
	if action == "apply" || action == "destroy" {
		if azureVNetFlags.autoApprove {
			tfArgs = append(tfArgs, "-auto-approve")
		}
	}

	tfCmd := exec.CommandContext(ctx, tfPath, tfArgs...)
	tfCmd.Dir = tfDir
	tfCmd.Stdout = os.Stdout
	tfCmd.Stderr = os.Stderr
	tfCmd.Stdin = os.Stdin

	if err := tfCmd.Run(); err != nil {
		return fmt.Errorf("terraform %s failed: %w", action, err)
	}

	output.PrintInfo(fmt.Sprintf("Terraform %s completed successfully", action))
	return nil
}

func buildAzureVNetSiteSpec() map[string]interface{} {
	spec := map[string]interface{}{
		"azure_region": azureVNetFlags.region,
	}

	if azureVNetFlags.resourceGroup != "" {
		spec["resource_group"] = azureVNetFlags.resourceGroup
	}

	if azureVNetFlags.vnetCIDR != "" {
		spec["vnet"] = map[string]interface{}{
			"new_vnet": map[string]interface{}{
				"primary_ipv4": azureVNetFlags.vnetCIDR,
			},
		}
	}

	if azureVNetFlags.machineType != "" {
		spec["machine_type"] = azureVNetFlags.machineType
	}

	if azureVNetFlags.sshKey != "" {
		spec["ssh_key"] = azureVNetFlags.sshKey
	}

	if azureVNetFlags.subscriptionID != "" {
		spec["azure_subscription_id"] = azureVNetFlags.subscriptionID
	}

	if azureVNetFlags.cloudCreds != "" {
		spec["azure_cred"] = map[string]interface{}{
			"name":      azureVNetFlags.cloudCreds,
			"namespace": "system",
			"tenant":    "ves-io",
		}
	}

	return spec
}
