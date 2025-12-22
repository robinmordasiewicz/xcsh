// Package branding provides centralized branding information for the CLI.
// This is the single source of truth for CLI names, descriptions, and branding.
// Update this file to rebrand the entire application.
package branding

const (
	// CLIName is the current name of the CLI tool
	CLIName = "xcsh"

	// CLIFullName is the full descriptive name
	CLIFullName = "F5 Distributed Cloud Shell"

	// CLIDescription is the full description used in documentation
	CLIDescription = "Command-line interface for F5 Distributed Cloud services."

	// CLIShortDescription is a brief description
	CLIShortDescription = "F5 Distributed Cloud Shell"

	// ConfigFileName is the name of the configuration file (without path)
	ConfigFileName = ".xcshconfig"

	// EnvPrefix is the prefix for environment variables
	EnvPrefix = "F5XC"

	// Legacy names for backward compatibility documentation
	LegacyCLIName       = "f5xcctl"
	LegacyLegacyCLIName = "vesctl"
)

// RepoOwner returns the GitHub repository owner
func RepoOwner() string {
	return "robinmordasiewicz"
}

// RepoName returns the GitHub repository name
func RepoName() string {
	return CLIName
}

// RepoURL returns the full GitHub repository URL
func RepoURL() string {
	return "https://github.com/" + RepoOwner() + "/" + RepoName()
}

// DocsURL returns the documentation site URL
func DocsURL() string {
	return "https://" + RepoOwner() + ".github.io/" + RepoName() + "/"
}

// ModulePath returns the Go module path
func ModulePath() string {
	return "github.com/" + RepoOwner() + "/" + RepoName()
}
