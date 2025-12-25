package cmd

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/mattn/go-isatty"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
)

// shouldEnterREPL determines if the CLI should enter interactive REPL mode
func shouldEnterREPL() bool {
	// No subcommand provided
	if len(os.Args) > 1 {
		return false
	}

	// Check if stdin is a terminal
	if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		return false
	}

	// Check for explicit non-interactive mode via environment variable
	if os.Getenv("XCSH_NON_INTERACTIVE") != "" {
		return false
	}

	return true
}

// StartREPL initializes and runs the interactive shell
func StartREPL() error {
	// Initialize configuration before starting REPL
	// This ensures viper has loaded config file and environment variables
	initConfig()

	// Display welcome banner
	printWelcomeBanner()

	// Initialize session
	session, err := initREPLSession()
	if err != nil {
		return fmt.Errorf("failed to initialize REPL: %w", err)
	}

	// Create prompt
	p := prompt.New(
		session.executeCommand,
		session.completeInput,
		prompt.OptionTitle(branding.CLIFullName+" - Interactive Shell"),
		prompt.OptionPrefix(buildPrompt(session)),
		prompt.OptionLivePrefix(session.livePrefix),
		prompt.OptionHistory(session.history.GetHistory()),
		prompt.OptionPrefixTextColor(prompt.Cyan),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionMaxSuggestion(10),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlD,
			Fn:  func(*prompt.Buffer) { handleExit(session) },
		}),
	)

	// Run the REPL
	p.Run()

	return nil
}

// printWelcomeBanner displays the welcome message when entering REPL mode
func printWelcomeBanner() {
	fmt.Printf(`
%s - Interactive Mode
Version: %s

Type 'help' for available commands, 'exit' or Ctrl+D to quit.
Tab completion is available for commands and arguments.

`, branding.CLIFullName, Version)
}

// buildPrompt constructs the prompt string based on session state
// Delegates to the colored prompt builder for rich contextual display
func buildPrompt(session *REPLSession) string {
	return buildColoredPrompt(session)
}

// handleExit performs cleanup and exits the REPL
func handleExit(session *REPLSession) {
	// Save history
	if session.history != nil {
		if err := session.history.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save history: %v\n", err)
		}
	}

	fmt.Println("\nGoodbye!")
	os.Exit(0)
}
