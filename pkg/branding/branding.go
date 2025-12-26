// Package branding provides centralized branding information for the CLI.
// This is the single source of truth for CLI names, descriptions, and branding.
// Update this file to rebrand the entire application.
package branding

const (
	// CLIName is the current name of the CLI tool
	CLIName = "xcsh"

	// F5Logo is a compact circular logo with F5 text
	// - ▓ characters represent the red circle background
	// - █ and ▒ characters represent the white F5 text
	// - (, ), |, and _ are the circle outline (rendered in red)
	// - Spacing pre-calculated for 78-char frame width (17 spaces offset)
	F5Logo = `                   ________
              (▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒)
         (▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒)
      (▒▒▓▓▓▓██████████▓▓▓▓█████████████)
    (▒▓▓▓▓██████▒▒▒▒▒███▓▓██████████████▒)
   (▒▓▓▓▓██████▒▓▓▓▓▓▒▒▒▓██▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒)
  (▒▓▓▓▓▓██████▓▓▓▓▓▓▓▓▓██▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒)
 (▒▓▓███████████████▓▓▓▓█████████████▓▓▓▓▓▓▒)
(▒▓▓▓▒▒▒███████▒▒▒▒▒▓▓▓████████████████▓▓▓▓▓▒)
|▒▓▓▓▓▓▓▒██████▓▓▓▓▓▓▓████████████████████▓▓▒|
|▒▓▓▓▓▓▓▓██████▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒██████████▓▒|
(▒▓▓▓▓▓▓▓██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒████████▒▒)
 (▒▓▓▓▓▓▓██████▓▓▓▓▓▓▓███▓▓▓▓▓▓▓▓▓▓▒▒▒████▒▒)
  (▒▓▓▓▓▓██████▓▓▓▓▓▓█████▓▓▓▓▓▓▓▓▓▓▓▓███▒▒)
   (▒▒██████████▓▓▓▓▓▒██████▓▓▓▓▓▓▓▓███▒▒▒)
    (▒▒▒▒▒██████████▓▓▒▒█████████████▒▒▓▒)
      (▒▓▓▒▒▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▓▒)
         (▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒)
              (▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒)`

	// ANSI color codes for terminal output
	ColorRed       = "\033[38;2;228;0;43m" // F5 Brand Red (#E4002B)
	ColorBoldWhite = "\033[1;97m"          // Bold bright white
	ColorReset     = "\033[0m"             // Reset to default

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
