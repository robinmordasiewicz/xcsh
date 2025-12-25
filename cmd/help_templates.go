package cmd

import (
	"github.com/spf13/cobra"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
)

// CommandLevel represents the depth/type of a command in the hierarchy
type CommandLevel string

const (
	// LevelRoot is the top-level xcsh command
	LevelRoot CommandLevel = "root"
	// LevelDomain is a domain command (e.g., virtual, load_balancer)
	LevelDomain CommandLevel = "domain"
	// LevelAction is an action command (e.g., list, get, create)
	LevelAction CommandLevel = "action"
	// LevelResource is a resource type command (e.g., http_loadbalancer)
	LevelResource CommandLevel = "resource"
	// LevelCustom is for custom top-level commands (e.g., site, subscription)
	LevelCustom CommandLevel = "custom"
)

// getCommandLevel determines the help level for a command
// Uses annotations for O(1) lookup, falls back to parent chain analysis
func getCommandLevel(cmd *cobra.Command) CommandLevel {
	// Check for explicit annotation override (fastest path)
	if cmd.Annotations != nil {
		if level, ok := cmd.Annotations["help-level"]; ok {
			return CommandLevel(level)
		}
	}

	// Fall back to depth-based detection
	depth := countParents(cmd)

	switch depth {
	case 0:
		return LevelRoot
	case 1:
		// Check if it's a domain command
		if _, isDomain := types.GetDomainInfo(cmd.Name()); isDomain {
			return LevelDomain
		}
		return LevelCustom
	case 2:
		// Check if parent is a domain
		if parent := cmd.Parent(); parent != nil {
			if _, isDomain := types.GetDomainInfo(parent.Name()); isDomain {
				return LevelAction
			}
		}
		return LevelCustom
	default:
		return LevelResource
	}
}

// countParents counts the number of parent commands
func countParents(cmd *cobra.Command) int {
	depth := 0
	parent := cmd.Parent()
	for parent != nil {
		depth++
		parent = parent.Parent()
	}
	return depth
}

// getHelpTemplateForLevel returns the appropriate help template for a command level
func getHelpTemplateForLevel(level CommandLevel) string {
	switch level {
	case LevelRoot:
		return rootHelpTemplate()
	case LevelDomain:
		return domainHelpTemplate()
	case LevelAction:
		return actionHelpTemplate()
	case LevelResource:
		return resourceHelpTemplate()
	case LevelCustom:
		return customHelpTemplate()
	default:
		return defaultHelpTemplate()
	}
}

// rootHelpTemplate returns the full verbose help template for the root command
// Includes environment variables, examples, and configuration sections
func rootHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// domainHelpTemplate returns a focused template for domain commands
// Shows domain info and available actions without global env vars
func domainHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Use "{{.CommandPath}} [action] --help" for action-specific help.
`
}

// actionHelpTemplate returns a focused template for action commands
// Shows flags and available resource types
func actionHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
Use "{{.CommandPath}} [resource-type] --help" for resource-specific help.
`
}

// resourceHelpTemplate returns a focused template for resource type commands
// Shows resource-specific examples and flags (included via UsageString)
func resourceHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`
}

// customHelpTemplate returns a template for custom top-level commands
// Similar to domain but without domain-specific assumptions
func customHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// defaultHelpTemplate returns Cobra's default-like template as fallback
func defaultHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`
}

// applyContextualHelpSystem applies appropriate help templates based on command level
// This replaces the uniform applyHelpTemplateRecursively function
func applyContextualHelpSystem(cmd *cobra.Command) {
	level := getCommandLevel(cmd)
	template := getHelpTemplateForLevel(level)
	cmd.SetHelpTemplate(template)

	for _, subCmd := range cmd.Commands() {
		applyContextualHelpSystem(subCmd)
	}
}
