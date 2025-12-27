package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
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

// StartREPL initializes and runs the interactive shell using bubbletea
func StartREPL() error {
	// Initialize configuration before starting REPL
	// This ensures viper has loaded config file and environment variables
	initConfig()

	// Initialize session first (this sets up apiClient)
	session, err := initREPLSession()
	if err != nil {
		return fmt.Errorf("failed to initialize REPL: %w", err)
	}

	// Create the TUI model
	model := NewTUIModel(session)

	// Set initial prompt
	model.input.SetPrompt(buildPlainPrompt(session))

	// Initialize history for navigation
	model.history = session.history.GetHistory()

	// Create and run the bubbletea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	// Save history on exit
	if m, ok := finalModel.(*TUIModel); ok {
		if m.session.history != nil {
			if err := m.session.history.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save history: %v\n", err)
			}
		}
	}

	return nil
}

// handleExit performs cleanup and exits the REPL
// Note: This is kept for compatibility with commands that call it directly
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
