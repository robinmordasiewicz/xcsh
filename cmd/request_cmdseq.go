package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/client"
	"github.com/robinmordasiewicz/xcsh/pkg/output"
	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

var cmdSeqFlags struct {
	inputFile string
	oldFile   string
	operation string
}

var commandSequenceCmd = &cobra.Command{
	Use:     "command-sequence",
	Aliases: []string{"cmdseq"},
	Short:   "Execute a command sequence",
	Long: `Execute a sequence of commands from a YAML file.

A command sequence file contains a list of resource definitions that will
be processed in order. This is useful for deploying multiple related
resources in a single operation.

Supported operations:
  - create: Create all resources in the sequence
  - delete: Delete all resources in the sequence
  - replace: Replace existing resources (requires --old-file)`,
	Example: `  # Create resources from a sequence file
  f5xcctl request command-sequence -i objects.yaml --operation create

  # Delete resources from a sequence file
  f5xcctl request command-sequence -i objects.yaml --operation delete

  # Replace resources (requires old file for comparison)
  f5xcctl request command-sequence -i newobjects.yaml --old-file oldobjects.yaml --operation replace`,
	RunE: runCommandSequence,
}

// CommandSequence represents a sequence of commands to execute
type CommandSequence struct {
	Items []CommandSequenceItem `yaml:"items" json:"items"`
}

// CommandSequenceItem represents a single item in the command sequence
type CommandSequenceItem struct {
	Kind      string                 `yaml:"kind" json:"kind"`
	Metadata  types.ResourceMetadata `yaml:"metadata" json:"metadata"`
	Spec      map[string]interface{} `yaml:"spec" json:"spec"`
	Operation string                 `yaml:"operation,omitempty" json:"operation,omitempty"`
}

func init() {
	requestCmd.AddCommand(commandSequenceCmd)

	commandSequenceCmd.Flags().StringVarP(&cmdSeqFlags.inputFile, "input-file", "i", "", "File containing command sequence data (required)")
	commandSequenceCmd.Flags().StringVar(&cmdSeqFlags.oldFile, "old-file", "", "File containing old command sequence data (for replace operation)")
	commandSequenceCmd.Flags().StringVar(&cmdSeqFlags.operation, "operation", "", "Operation to perform: create, delete, replace")
	_ = commandSequenceCmd.MarkFlagRequired("input-file")
}

func runCommandSequence(cmd *cobra.Command, args []string) error {
	apiClient := GetClient()
	if apiClient == nil {
		return fmt.Errorf("client not initialized - check configuration")
	}

	// Load the command sequence
	sequence, err := loadCommandSequence(cmdSeqFlags.inputFile)
	if err != nil {
		return fmt.Errorf("failed to load command sequence: %w", err)
	}

	// Load old sequence if provided (for replace operation)
	var oldSequence *CommandSequence
	if cmdSeqFlags.oldFile != "" {
		oldSequence, err = loadCommandSequence(cmdSeqFlags.oldFile)
		if err != nil {
			return fmt.Errorf("failed to load old command sequence: %w", err)
		}
	}

	// Determine the operation
	operation := strings.ToLower(cmdSeqFlags.operation)
	if operation == "" {
		operation = "create"
	}

	// Validate operation
	switch operation {
	case "create", "delete", "replace":
		// valid
	default:
		return fmt.Errorf("invalid operation: %s (must be create, delete, or replace)", operation)
	}

	// Process each item in the sequence
	results := make([]map[string]interface{}, 0)
	errors := make([]string, 0)

	for i, item := range sequence.Items {
		itemOp := operation
		if item.Operation != "" {
			itemOp = strings.ToLower(item.Operation)
		}

		output.PrintInfo(fmt.Sprintf("[%d/%d] %s %s/%s", i+1, len(sequence.Items), itemOp, item.Kind, item.Metadata.Name))

		result, err := executeSequenceItem(apiClient, &item, itemOp, oldSequence)
		if err != nil {
			errMsg := fmt.Sprintf("Item %d (%s/%s): %v", i+1, item.Kind, item.Metadata.Name, err)
			errors = append(errors, errMsg)
			output.PrintError(fmt.Errorf("%s", errMsg))
			continue
		}

		if result != nil {
			results = append(results, map[string]interface{}{
				"item":      fmt.Sprintf("%s/%s", item.Kind, item.Metadata.Name),
				"operation": itemOp,
				"status":    "success",
				"result":    result,
			})
		}
	}

	// Summary
	fmt.Println()
	output.PrintInfo(fmt.Sprintf("Command sequence completed: %d succeeded, %d failed", len(results), len(errors)))

	if len(errors) > 0 {
		return fmt.Errorf("%d items failed", len(errors))
	}

	return nil
}

func loadCommandSequence(path string) (*CommandSequence, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as a sequence
	var sequence CommandSequence
	if err := yaml.Unmarshal(data, &sequence); err != nil {
		// Try parsing as a list of items directly
		var items []CommandSequenceItem
		if err := yaml.Unmarshal(data, &items); err != nil {
			// Try parsing as a single item
			var item CommandSequenceItem
			if err := yaml.Unmarshal(data, &item); err != nil {
				return nil, fmt.Errorf("failed to parse file: %w", err)
			}
			sequence.Items = []CommandSequenceItem{item}
		} else {
			sequence.Items = items
		}
	}

	// If items is empty but we have top-level fields, treat as single item
	if len(sequence.Items) == 0 {
		var item CommandSequenceItem
		if err := yaml.Unmarshal(data, &item); err == nil && item.Kind != "" {
			sequence.Items = []CommandSequenceItem{item}
		}
	}

	return &sequence, nil
}

func executeSequenceItem(apiClient *client.Client, item *CommandSequenceItem, operation string, oldSequence *CommandSequence) (interface{}, error) {
	// Get the resource type
	rt, ok := types.Get(item.Kind)
	if !ok {
		// Try converting kind to CLI name format
		cliName := strings.ToLower(strings.ReplaceAll(item.Kind, "_", "-"))
		rt, ok = types.Get(cliName)
		if !ok {
			return nil, fmt.Errorf("unknown resource type: %s", item.Kind)
		}
	}

	// Build the API path
	path := rt.BuildAPIPath(item.Metadata.Namespace, "")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Build the request body
	requestBody := map[string]interface{}{
		"metadata": item.Metadata,
		"spec":     item.Spec,
	}

	var resp *client.Response
	var err error

	switch operation {
	case "create":
		resp, err = apiClient.Post(ctx, path, requestBody)
	case "delete":
		deletePath := rt.BuildAPIPath(item.Metadata.Namespace, item.Metadata.Name)
		resp, err = apiClient.Delete(ctx, deletePath)
	case "replace":
		replacePath := rt.BuildAPIPath(item.Metadata.Namespace, item.Metadata.Name)
		resp, err = apiClient.Put(ctx, replacePath, requestBody)
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(resp.Body))
	}

	// Parse response
	var result interface{}
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &result); err != nil {
			return string(resp.Body), nil
		}
	}

	return result, nil
}
