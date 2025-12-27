package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/robinmordasiewicz/xcsh/pkg/profile"
)

var loginProfileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all authentication profiles.",
	Long: `Display all configured authentication profiles.

Shows profile names, associated tenants, authentication methods, and status.
Use --output-format json for programmatic parsing.`,
	Example: `  # List profiles in table format
  xcsh login profile list

  # List profiles as JSON
  xcsh login profile list --output-format json

  # List profiles as YAML
  xcsh login profile list --output-format yaml`,
	RunE: runLoginProfileList,
}

func init() {
	loginProfileCmd.AddCommand(loginProfileListCmd)
}

func runLoginProfileList(cmd *cobra.Command, args []string) error {
	manager, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize profile manager: %w", err)
	}

	profiles, err := manager.ListProfiles()
	if err != nil {
		return fmt.Errorf("failed to list profiles: %w", err)
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles configured.")
		fmt.Println("\nCreate a profile with:")
		fmt.Println("  xcsh login profile create --name <name> --api-url <url> --api-token <token>")
		return nil
	}

	currentName := manager.GetCurrentName()
	defaultName := manager.GetDefault()

	switch outputFormat {
	case "json":
		return outputProfileListJSON(profiles, currentName, defaultName)
	case "yaml":
		return outputProfileListYAML(profiles, currentName, defaultName)
	default:
		return outputProfileListTable(profiles, currentName, defaultName)
	}
}

type profileListEntry struct {
	Name       string `json:"name" yaml:"name"`
	Tenant     string `json:"tenant" yaml:"tenant"`
	APIURL     string `json:"api_url" yaml:"api_url"`
	AuthMethod string `json:"auth_method" yaml:"auth_method"`
	IsDefault  bool   `json:"is_default" yaml:"is_default"`
	IsCurrent  bool   `json:"is_current" yaml:"is_current"`
}

func outputProfileListTable(profiles []*profile.Profile, currentName, defaultName string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTENANT\tAUTH METHOD\tSTATUS")
	fmt.Fprintln(w, "----\t------\t-----------\t------")

	for _, p := range profiles {
		status := ""
		if p.Name == currentName && p.Name == defaultName {
			status = "current, default"
		} else if p.Name == currentName {
			status = "current"
		} else if p.Name == defaultName {
			status = "default"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			p.Name,
			p.TenantName(),
			p.AuthMethod(),
			status,
		)
	}

	return w.Flush()
}

func outputProfileListJSON(profiles []*profile.Profile, currentName, defaultName string) error {
	entries := make([]profileListEntry, len(profiles))
	for i, p := range profiles {
		entries[i] = profileListEntry{
			Name:       p.Name,
			Tenant:     p.TenantName(),
			APIURL:     p.APIURL,
			AuthMethod: p.AuthMethod(),
			IsDefault:  p.Name == defaultName,
			IsCurrent:  p.Name == currentName,
		}
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputProfileListYAML(profiles []*profile.Profile, currentName, defaultName string) error {
	entries := make([]profileListEntry, len(profiles))
	for i, p := range profiles {
		entries[i] = profileListEntry{
			Name:       p.Name,
			Tenant:     p.TenantName(),
			APIURL:     p.APIURL,
			AuthMethod: p.AuthMethod(),
			IsDefault:  p.Name == defaultName,
			IsCurrent:  p.Name == currentName,
		}
	}

	data, err := yaml.Marshal(entries)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
