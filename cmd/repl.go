package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/mattn/go-isatty"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
)

// printWelcomeBanner displays the modern welcome banner when entering REPL mode
func printWelcomeBanner() {
	fmt.Print(renderWelcomeBanner())
}

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

	// Initialize session first (this sets up apiClient)
	session, err := initREPLSession()
	if err != nil {
		return fmt.Errorf("failed to initialize REPL: %w", err)
	}

	// Display welcome banner after session init so we have apiClient for user info
	printWelcomeBanner()

	// Create prompt
	p := prompt.New(
		session.executeCommand,
		session.completeInput,
		prompt.OptionTitle(branding.CLIFullName+" - Interactive Shell"),
		prompt.OptionPrefix(buildPrompt(session)),
		prompt.OptionLivePrefix(session.livePrefix),
		prompt.OptionHistory(session.history.GetHistory()),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionMaxSuggestion(10),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn:  func(*prompt.Buffer) { handleCtrlC(session) },
		}),
	)

	// Run the REPL
	p.Run()

	return nil
}

// buildPrompt constructs the prompt string based on session state
func buildPrompt(session *REPLSession) string {
	return buildPlainPrompt(session)
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

// handleCtrlC handles Ctrl+C with double-press to exit
func handleCtrlC(session *REPLSession) {
	now := time.Now()
	// 500ms window for double-press detection
	if now.Sub(session.lastCtrlCTime) < 500*time.Millisecond {
		handleExit(session)
	}
	session.lastCtrlCTime = now
	// Show hint message
	fmt.Print("\nPress Ctrl+C again to exit, or continue typing\n")
}
