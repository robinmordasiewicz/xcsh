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

var profileShowFlags struct {
	showSensitive bool
}

var loginProfileShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Display profile details.",
	Long: `Display the configuration of a specific profile or the current active profile.

If no profile name is provided, shows the currently active profile.

SECURITY: Sensitive data (API tokens) are masked by default.
Use --show-sensitive to reveal full credentials (use with caution).`,
	Example: `  # Show active profile (tokens masked)
  xcsh login profile show

  # Show specific profile
  xcsh login profile show staging

  # Show as JSON
  xcsh login profile show --output-format json

  # Show with full token (CAUTION: exposes credentials)
  xcsh login profile show --show-sensitive`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLoginProfileShow,
}

func init() {
	loginProfileCmd.AddCommand(loginProfileShowCmd)

	loginProfileShowCmd.Flags().BoolVar(&profileShowFlags.showSensitive, "show-sensitive", false,
		"Show sensitive data (tokens) in plain text (CAUTION: exposes credentials)")
}

func runLoginProfileShow(cmd *cobra.Command, args []string) error {
	manager, err := profile.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize profile manager: %w", err)
	}

	var p *profile.Profile
	var profileName string

	if len(args) > 0 {
		profileName = args[0]
		p, err = manager.Load(profileName)
		if err != nil {
			return fmt.Errorf("failed to load profile %q: %w", profileName, err)
		}
	} else {
		p, err = manager.GetCurrent()
		if err != nil {
			return fmt.Errorf("no active profile: %w", err)
		}
		profileName = p.Name
	}

	isDefault := manager.IsDefault(profileName)
	isCurrent := manager.GetCurrentName() == profileName

	switch outputFormat {
	case "json":
		return outputProfileShowJSON(p, isDefault, isCurrent)
	case "yaml":
		return outputProfileShowYAML(p, isDefault, isCurrent)
	default:
		return outputProfileShowTable(p, isDefault, isCurrent)
	}
}

type profileShowEntry struct {
	Name             string `json:"name" yaml:"name"`
	APIURL           string `json:"api_url" yaml:"api_url"`
	Tenant           string `json:"tenant" yaml:"tenant"`
	AuthMethod       string `json:"auth_method" yaml:"auth_method"`
	APIToken         string `json:"api_token,omitempty" yaml:"api_token,omitempty"`
	P12Bundle        string `json:"p12_bundle,omitempty" yaml:"p12_bundle,omitempty"`
	Cert             string `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key              string `json:"key,omitempty" yaml:"key,omitempty"`
	DefaultNamespace string `json:"default_namespace,omitempty" yaml:"default_namespace,omitempty"`
	IsDefault        bool   `json:"is_default" yaml:"is_default"`
	IsCurrent        bool   `json:"is_current" yaml:"is_current"`
}

// getDisplayToken returns either the masked or full token based on --show-sensitive flag
func getDisplayToken(p *profile.Profile) string {
	if p.APIToken == "" {
		return ""
	}
	if profileShowFlags.showSensitive {
		return p.APIToken
	}
	return p.MaskedToken()
}

func outputProfileShowTable(p *profile.Profile, isDefault, isCurrent bool) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	status := ""
	if isCurrent && isDefault {
		status = "current, default"
	} else if isCurrent {
		status = "current"
	} else if isDefault {
		status = "default"
	}

	fmt.Fprintf(w, "Name:\t%s\n", p.Name)
	fmt.Fprintf(w, "API URL:\t%s\n", p.APIURL)
	fmt.Fprintf(w, "Tenant:\t%s\n", p.TenantName())
	fmt.Fprintf(w, "Auth Method:\t%s\n", p.AuthMethod())

	if p.APIToken != "" {
		fmt.Fprintf(w, "API Token:\t%s\n", getDisplayToken(p))
	}
	if p.P12Bundle != "" {
		fmt.Fprintf(w, "P12 Bundle:\t%s\n", p.P12Bundle)
	}
	if p.Cert != "" {
		fmt.Fprintf(w, "Certificate:\t%s\n", p.Cert)
	}
	if p.Key != "" {
		fmt.Fprintf(w, "Key:\t%s\n", p.Key)
	}
	if p.DefaultNamespace != "" {
		fmt.Fprintf(w, "Default Namespace:\t%s\n", p.DefaultNamespace)
	}
	if status != "" {
		fmt.Fprintf(w, "Status:\t%s\n", status)
	}

	return w.Flush()
}

func outputProfileShowJSON(p *profile.Profile, isDefault, isCurrent bool) error {
	entry := profileShowEntry{
		Name:             p.Name,
		APIURL:           p.APIURL,
		Tenant:           p.TenantName(),
		AuthMethod:       p.AuthMethod(),
		APIToken:         getDisplayToken(p),
		P12Bundle:        p.P12Bundle,
		Cert:             p.Cert,
		Key:              p.Key,
		DefaultNamespace: p.DefaultNamespace,
		IsDefault:        isDefault,
		IsCurrent:        isCurrent,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func outputProfileShowYAML(p *profile.Profile, isDefault, isCurrent bool) error {
	entry := profileShowEntry{
		Name:             p.Name,
		APIURL:           p.APIURL,
		Tenant:           p.TenantName(),
		AuthMethod:       p.AuthMethod(),
		APIToken:         getDisplayToken(p),
		P12Bundle:        p.P12Bundle,
		Cert:             p.Cert,
		Key:              p.Key,
		DefaultNamespace: p.DefaultNamespace,
		IsDefault:        isDefault,
		IsCurrent:        isCurrent,
	}

	data, err := yaml.Marshal(entry)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
