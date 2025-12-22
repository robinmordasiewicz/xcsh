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

	"github.com/robinmordasiewicz/xcsh/pkg/output"
)

var awsVPCFlags struct {
	name            string
	namespace       string
	inputFile       string
	region          string
	azs             []string
	vpcCIDR         string
	instanceType    string
	sshKey          string
	cloudCreds      string
	terraformDir    string
	terraformAction string
	autoApprove     bool
	wait            bool
}

var awsVPCCmd = &cobra.Command{
	Use:   "aws_vpc",
	Short: "Manage AWS VPC site creation through view apis",
	Long: `Manage AWS VPC sites in F5 Distributed Cloud.

AWS VPC sites allow you to deploy F5 XC Customer Edge (CE) nodes in your
AWS Virtual Private Cloud, enabling secure connectivity and edge services.`,
	Example: `  # Create an AWS VPC site from a YAML file
  f5xcctl site aws_vpc create -i aws-site.yaml

  # Delete an AWS VPC site
  f5xcctl site aws_vpc delete --name example-site

  # Run Terraform to provision infrastructure
  f5xcctl site aws_vpc run --name example-site --action apply --auto-approve`,
}

var awsVPCCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create AWS VPC volterra site",
	Long: `Create a new AWS VPC site in F5 Distributed Cloud.

This command registers an AWS VPC site configuration. After creation,
use 'f5xcctl site aws_vpc run --action apply' to provision the infrastructure.

You can provide the site specification via:
- YAML/JSON file using --input-file
- Command line flags for common options`,
	Example: `  # Create from YAML file
  f5xcctl site aws_vpc create -i aws-site.yaml

  # Create with command line flags
  f5xcctl site aws_vpc create --name example-site --region us-west-2 \
    --vpc-cidr 10.0.0.0/16 --cloud-creds example-aws-creds

  # Create with availability zones
  f5xcctl site aws_vpc create --name example-site --region us-west-2 \
    --azs us-west-2a,us-west-2b --cloud-creds example-aws-creds`,
	RunE: runAWSVPCCreate,
}

var awsVPCDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete AWS VPC volterra site",
	Long: `Delete an AWS VPC site from F5 Distributed Cloud.

Note: This only removes the site configuration from F5 XC. To fully clean up
AWS resources, first run 'f5xcctl site aws_vpc run --action destroy' before
deleting the site configuration.`,
	Example: `  # Delete a site (after destroying infrastructure)
  f5xcctl site aws_vpc delete --name example-site

  # Delete from a specific namespace
  f5xcctl site aws_vpc delete --name example-site -n system`,
	RunE: runAWSVPCDelete,
}

var awsVPCReplaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace AWS VPC volterra site",
	Long: `Replace an existing AWS VPC site configuration in F5 Distributed Cloud.

This updates the site specification. After replacing, you may need to
run 'f5xcctl site aws_vpc run --action apply' to apply infrastructure changes.`,
	Example: `  # Replace site configuration from file
  f5xcctl site aws_vpc replace -i updated-site.yaml

  # Replace with specific name
  f5xcctl site aws_vpc replace --name example-site -i updated-site.yaml`,
	RunE: runAWSVPCReplace,
}

var awsVPCRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run terraform action, valid actions are plan, apply and destroy",
	Long: `Run Terraform actions to provision or destroy AWS VPC site infrastructure.

This command retrieves Terraform parameters from F5 XC and executes Terraform
to manage the actual AWS resources (VPC, subnets, EC2 instances, etc.).

Available actions:
- plan: Preview changes without applying
- apply: Create or update infrastructure
- destroy: Remove all infrastructure`,
	Example: `  # Preview infrastructure changes
  f5xcctl site aws_vpc run --name example-site --action plan

  # Apply infrastructure (with confirmation prompt)
  f5xcctl site aws_vpc run --name example-site --action apply

  # Apply infrastructure automatically (for CI/CD)
  f5xcctl site aws_vpc run --name example-site --action apply --auto-approve

  # Destroy infrastructure
  f5xcctl site aws_vpc run --name example-site --action destroy --auto-approve`,
	RunE: runAWSVPCTerraform,
}

func init() {
	siteCmd.AddCommand(awsVPCCmd)

	// Create command
	awsVPCCreateCmd.Flags().StringVarP(&awsVPCFlags.inputFile, "input-file", "i", "", "Input file (YAML/JSON) containing site definition")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.name, "name", "", "Site name")
	awsVPCCreateCmd.Flags().StringVarP(&awsVPCFlags.namespace, "namespace", "n", "system", "Namespace")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.region, "region", "", "AWS region")
	awsVPCCreateCmd.Flags().StringSliceVar(&awsVPCFlags.azs, "azs", nil, "Availability zones (comma-separated)")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.vpcCIDR, "vpc-cidr", "", "VPC CIDR block")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.instanceType, "instance-type", "t3.xlarge", "EC2 instance type")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.sshKey, "ssh-key", "", "SSH key name")
	awsVPCCreateCmd.Flags().StringVar(&awsVPCFlags.cloudCreds, "cloud-creds", "", "Cloud credentials name")
	awsVPCCmd.AddCommand(awsVPCCreateCmd)

	// Delete command
	awsVPCDeleteCmd.Flags().StringVar(&awsVPCFlags.name, "name", "", "Site name (required)")
	awsVPCDeleteCmd.Flags().StringVarP(&awsVPCFlags.namespace, "namespace", "n", "system", "Namespace")
	_ = awsVPCDeleteCmd.MarkFlagRequired("name")
	awsVPCCmd.AddCommand(awsVPCDeleteCmd)

	// Replace command
	awsVPCReplaceCmd.Flags().StringVarP(&awsVPCFlags.inputFile, "input-file", "i", "", "Input file (YAML/JSON) containing site definition")
	awsVPCReplaceCmd.Flags().StringVar(&awsVPCFlags.name, "name", "", "Site name")
	awsVPCReplaceCmd.Flags().StringVarP(&awsVPCFlags.namespace, "namespace", "n", "system", "Namespace")
	awsVPCCmd.AddCommand(awsVPCReplaceCmd)

	// Run (Terraform) command
	awsVPCRunCmd.Flags().StringVar(&awsVPCFlags.name, "name", "", "Site name (required)")
	awsVPCRunCmd.Flags().StringVarP(&awsVPCFlags.namespace, "namespace", "n", "system", "Namespace")
	awsVPCRunCmd.Flags().StringVar(&awsVPCFlags.terraformAction, "action", "plan", "Terraform action: plan, apply, destroy")
	awsVPCRunCmd.Flags().StringVar(&awsVPCFlags.terraformDir, "terraform-dir", "", "Directory for Terraform files (default: temp dir)")
	awsVPCRunCmd.Flags().BoolVar(&awsVPCFlags.autoApprove, "auto-approve", false, "Auto-approve Terraform apply/destroy")
	awsVPCRunCmd.Flags().BoolVar(&awsVPCFlags.wait, "wait", true, "Wait for operation to complete")
	_ = awsVPCRunCmd.MarkFlagRequired("name")
	awsVPCCmd.AddCommand(awsVPCRunCmd)
}

func runAWSVPCCreate(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	var siteSpec map[string]interface{}

	if awsVPCFlags.inputFile != "" {
		data, err := os.ReadFile(awsVPCFlags.inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		if err := yaml.Unmarshal(data, &siteSpec); err != nil {
			if err := json.Unmarshal(data, &siteSpec); err != nil {
				return fmt.Errorf("failed to parse input file: %w", err)
			}
		}
	} else {
		if awsVPCFlags.name == "" || awsVPCFlags.region == "" {
			return fmt.Errorf("--name and --region are required when not using input file")
		}

		siteSpec = buildAWSVPCSiteSpec()
	}

	// Extract metadata if present, or build from flags
	metadata := map[string]interface{}{
		"name":      awsVPCFlags.name,
		"namespace": awsVPCFlags.namespace,
	}
	if meta, ok := siteSpec["metadata"].(map[string]interface{}); ok {
		metadata = meta
		if awsVPCFlags.name != "" {
			metadata["name"] = awsVPCFlags.name
		}
		if awsVPCFlags.namespace != "" {
			metadata["namespace"] = awsVPCFlags.namespace
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

	path := fmt.Sprintf("/api/config/namespaces/%s/aws_vpc_sites", awsVPCFlags.namespace)
	resp, err := apiClient.Post(ctx, path, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create AWS VPC site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("AWS VPC site '%s' created successfully", awsVPCFlags.name))
	return output.Print(result, GetOutputFormat())
}

func runAWSVPCDelete(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	path := fmt.Sprintf("/api/config/namespaces/%s/aws_vpc_sites/%s", awsVPCFlags.namespace, awsVPCFlags.name)
	resp, err := apiClient.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete AWS VPC site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	output.PrintInfo(fmt.Sprintf("AWS VPC site '%s' deleted successfully", awsVPCFlags.name))
	return nil
}

func runAWSVPCReplace(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	if awsVPCFlags.inputFile == "" {
		return fmt.Errorf("--input-file is required for replace operation")
	}

	data, err := os.ReadFile(awsVPCFlags.inputFile)
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
	name := awsVPCFlags.name
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

	path := fmt.Sprintf("/api/config/namespaces/%s/aws_vpc_sites/%s", awsVPCFlags.namespace, name)
	resp, err := apiClient.Put(ctx, path, siteSpec)
	if err != nil {
		return fmt.Errorf("failed to replace AWS VPC site: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	var result interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	output.PrintInfo(fmt.Sprintf("AWS VPC site '%s' replaced successfully", name))
	return output.Print(result, GetOutputFormat())
}

func runAWSVPCTerraform(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Validate action
	action := awsVPCFlags.terraformAction
	switch action {
	case "plan", "apply", "destroy":
		// valid
	default:
		return fmt.Errorf("invalid action: %s (must be plan, apply, or destroy)", action)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// Get Terraform parameters from the site
	output.PrintInfo(fmt.Sprintf("Retrieving Terraform parameters for site '%s'...", awsVPCFlags.name))

	tfParamsPath := fmt.Sprintf("/api/terraform/namespaces/%s/terraform/aws_vpc_site/%s", awsVPCFlags.namespace, awsVPCFlags.name)
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
	tfDir := awsVPCFlags.terraformDir
	if tfDir == "" {
		tfDir = filepath.Join(os.TempDir(), "f5xcctl-terraform", awsVPCFlags.name)
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
		if awsVPCFlags.autoApprove {
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

func buildAWSVPCSiteSpec() map[string]interface{} {
	spec := map[string]interface{}{
		"aws_region": awsVPCFlags.region,
	}

	if awsVPCFlags.vpcCIDR != "" {
		spec["vpc"] = map[string]interface{}{
			"new_vpc": map[string]interface{}{
				"primary_ipv4": awsVPCFlags.vpcCIDR,
			},
		}
	}

	if len(awsVPCFlags.azs) > 0 {
		azNodes := make([]map[string]interface{}, len(awsVPCFlags.azs))
		for i, az := range awsVPCFlags.azs {
			azNodes[i] = map[string]interface{}{
				"aws_az_name": az,
			}
		}
		spec["ingress_egress_gw"] = map[string]interface{}{
			"aws_certified_hw": "aws-byol-voltmesh",
			"az_nodes":         azNodes,
		}
	}

	if awsVPCFlags.instanceType != "" {
		spec["instance_type"] = awsVPCFlags.instanceType
	}

	if awsVPCFlags.sshKey != "" {
		spec["ssh_key"] = awsVPCFlags.sshKey
	}

	if awsVPCFlags.cloudCreds != "" {
		spec["aws_cred"] = map[string]interface{}{
			"name":      awsVPCFlags.cloudCreds,
			"namespace": "system",
			"tenant":    "ves-io",
		}
	}

	return spec
}
