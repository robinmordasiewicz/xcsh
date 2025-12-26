package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/robinmordasiewicz/xcsh/pkg/errors"
)

// replBuiltins maps built-in REPL commands to their handlers
var replBuiltins = map[string]func(*REPLSession, []string) error{
	"quit":      cmdQuit, // Always exits REPL
	"help":      cmdHelp,
	"clear":     cmdClear,
	"history":   cmdHistory,
	"namespace": cmdNamespace,
	"ns":        cmdNamespace,
	"context":   cmdContext,
	"ctx":       cmdContext,
}

// executeCommand processes a command line within the REPL
func (s *REPLSession) executeCommand(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Add to history
	s.history.Add(input)

	// Parse input into arguments
	args := parseInputArgs(input)
	if len(args) == 0 {
		return
	}

	// Check for navigation commands first (domain entry, exit, back, etc.)
	if s.handleNavigationCommand(args) {
		s.lastExitCode = 0
		return
	}

	// Check for REPL built-in commands
	if handler, ok := replBuiltins[args[0]]; ok {
		if err := handler(s, args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			s.lastExitCode = 1
		} else {
			s.lastExitCode = 0
		}
		return
	}

	// Execute context-aware command via Cobra
	s.lastExitCode = s.executeContextAwareCommand(args)
}

// handleNavigationCommand handles context navigation commands
// Returns true if the command was handled as navigation, false otherwise
func (s *REPLSession) handleNavigationCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	cmd := args[0]
	ctx := s.GetContextPath()
	validator := s.GetValidator()

	switch cmd {
	case "exit":
		// At root, exit the shell; otherwise navigate up
		if ctx.IsRoot() {
			handleExit(s)
			return true // Never reached
		}
		ctx.NavigateUp()
		return true

	case "back", "..":
		// Navigate up one level
		if ctx.NavigateUp() {
			return true
		}
		fmt.Println("Already at root context")
		return true

	case "root", "/":
		// Return to root context
		ctx.Reset()
		return true

	default:
		// Check if this is a domain entry (at root or navigating to new domain)
		if ctx.IsRoot() && validator.IsValidDomain(cmd) {
			// Resolve alias to canonical domain name
			if canonical, ok := validator.ResolveDomain(cmd); ok {
				ctx.SetDomain(canonical)
			} else {
				ctx.SetDomain(cmd)
			}
			return true
		}

		// Check if this is an action entry (in domain context)
		if ctx.IsDomain() && validator.IsValidAction(cmd) {
			ctx.SetAction(cmd)
			return true
		}

		// Not a navigation command
		return false
	}
}

// executeContextAwareCommand runs a command with context prepended
func (s *REPLSession) executeContextAwareCommand(args []string) int {
	// Prepend context if applicable
	args = s.prependContext(args)

	// Inject default namespace if not specified
	args = s.injectDefaultNamespace(args)

	// Reset the command state for fresh execution
	rootCmd.SetArgs(args)

	// Capture stdout/stderr for proper REPL display
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	// Execute the command
	_, err := rootCmd.ExecuteC()
	if err != nil {
		return errors.GetExitCode(err)
	}

	return 0
}

// prependContext adds domain/action prefix to args based on current context
func (s *REPLSession) prependContext(args []string) []string {
	if len(args) == 0 {
		return args
	}

	// Check for "/" prefix - escape to root context (bypass context prepending)
	if strings.HasPrefix(args[0], "/") {
		// Strip the "/" prefix and return as root-level command
		args[0] = strings.TrimPrefix(args[0], "/")
		if args[0] == "" && len(args) > 1 {
			// Handle case like "/ command" (space after /)
			return args[1:]
		}
		return args
	}

	ctx := s.GetContextPath()

	// If at root, no context to prepend
	if ctx.IsRoot() {
		return args
	}

	// Check if first arg is already a valid root command (full path)
	cmd, _, err := rootCmd.Find(args)
	if err == nil && cmd != rootCmd {
		// It's already a valid command path, don't prepend
		return args
	}

	// Build context-prefixed args
	var prefixedArgs []string

	if ctx.Domain != "" {
		prefixedArgs = append(prefixedArgs, ctx.Domain)
	}

	if ctx.Action != "" {
		prefixedArgs = append(prefixedArgs, ctx.Action)
	}

	// Append original args
	prefixedArgs = append(prefixedArgs, args...)

	return prefixedArgs
}

// injectDefaultNamespace adds -n flag if not present and command supports it
func (s *REPLSession) injectDefaultNamespace(args []string) []string {
	if s.namespace == "" {
		return args
	}

	// Check if namespace flag is already present
	for _, arg := range args {
		if arg == "-n" || arg == "--namespace" {
			return args // Already specified
		}
		if strings.HasPrefix(arg, "-n=") || strings.HasPrefix(arg, "--namespace=") {
			return args // Already specified
		}
	}

	// Find the command to check if it supports namespace
	cmd, _, err := rootCmd.Find(args)
	if err != nil {
		return args
	}

	// Check if command has namespace flag
	if cmd.Flags().Lookup("namespace") != nil {
		// Append namespace at the end
		return append(args, "-n", s.namespace)
	}

	return args
}

// parseInputArgs splits input into arguments, respecting quotes
func parseInputArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

// cmdQuit handles the quit command (always exits immediately)
func cmdQuit(s *REPLSession, args []string) error {
	handleExit(s)
	return nil // Never reached
}

// cmdContext shows current context information
func cmdContext(s *REPLSession, args []string) error {
	ctx := s.GetContextPath()
	tenant := s.GetTenant()
	ns := s.GetNamespace()

	fmt.Println("Current Context:")
	fmt.Printf("  Tenant:    %s\n", valueOrDefault(tenant, "(not set)"))
	fmt.Printf("  Domain:    %s\n", valueOrDefault(ctx.Domain, "(root)"))
	fmt.Printf("  Action:    %s\n", valueOrDefault(ctx.Action, "(none)"))
	fmt.Printf("  Namespace: %s\n", valueOrDefault(ns, "(not set)"))
	fmt.Printf("  Path:      %s\n", valueOrDefault(ctx.String(), "/"))

	return nil
}

// valueOrDefault returns the value if non-empty, otherwise the default
func valueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// cmdHelp shows REPL help
func cmdHelp(s *REPLSession, args []string) error {
	if len(args) == 0 {
		ctx := s.GetContextPath()

		// Show context-specific help
		if !ctx.IsRoot() {
			showContextHelp(s)
			return nil
		}

		// Root help
		fmt.Print(`
xcsh Interactive Shell

Context Navigation (Cisco IOS style):
  <domain>          Enter a domain context (e.g., http_loadbalancer)
  <action>          Enter action context when in domain (e.g., list, create)
  exit              Go up one level (or exit shell at root)
  back, ..          Go up one level
  root, /           Return to root context
  quit              Exit shell immediately (bypass context)

Built-in Commands:
  help [command]    Show help for a command
  clear             Clear the screen
  history           Show command history
  namespace <ns>    Set default namespace (alias: ns)
  context           Show current context info (alias: ctx)

Keyboard Shortcuts:
  Tab               Auto-complete commands and arguments
  Ctrl+D            Exit the shell
  Up/Down           Navigate command history

Prompt Format: tenant:domain/action@namespace>

For command help, use: <command> --help
`)
		return nil
	}

	// Delegate to Cobra help for specific commands
	rootCmd.SetArgs(append(args, "--help"))
	return rootCmd.Execute()
}

// showContextHelp shows help relevant to current context
func showContextHelp(s *REPLSession) {
	ctx := s.GetContextPath()

	fmt.Printf("\nContext: %s\n", ctx.String())
	fmt.Println()

	if ctx.IsDomain() {
		// In domain context - show available actions
		fmt.Printf("Available actions in '%s':\n", ctx.Domain)
		fmt.Println("  list              List resources")
		fmt.Println("  get <name>        Get a specific resource")
		fmt.Println("  create            Create a new resource")
		fmt.Println("  delete <name>     Delete a resource")
		fmt.Println("  replace           Replace a resource")
		fmt.Println("  apply             Apply configuration from file")
		fmt.Println("  status <name>     Get resource status")
		fmt.Println()
		fmt.Println("Navigation:")
		fmt.Println("  exit, back, ..    Return to root")
		fmt.Println("  <action>          Enter action context")
		fmt.Println()
		fmt.Printf("Example: list -n production\n")
	} else if ctx.IsAction() {
		// In action context - show action-specific help
		fmt.Printf("Action: %s/%s\n", ctx.Domain, ctx.Action)
		fmt.Println()
		fmt.Println("Commands execute with this context prepended.")
		fmt.Println("Use flags and arguments directly.")
		fmt.Println()
		fmt.Println("Navigation:")
		fmt.Println("  exit, back, ..    Return to domain context")
		fmt.Println("  root, /           Return to root")
		fmt.Println()
		fmt.Printf("Example: -n production --output-format json\n")
	}
}

// cmdClear clears the terminal (matches bash/zsh clear behavior)
func cmdClear(s *REPLSession, args []string) error {
	// \033[H    - Move cursor to home position (top-left)
	// \033[2J   - Clear entire screen
	// \033[3J   - Clear scrollback buffer
	fmt.Print("\033[H\033[2J\033[3J")
	return nil
}

// cmdHistory shows command history
func cmdHistory(s *REPLSession, args []string) error {
	for i, cmd := range s.history.GetHistory() {
		fmt.Printf("%4d  %s\n", i+1, cmd)
	}
	return nil
}

// cmdNamespace sets the default namespace and persists to config
func cmdNamespace(s *REPLSession, args []string) error {
	if len(args) == 0 {
		ns := s.GetNamespace()
		source := getNamespaceSource()
		if ns == "" {
			fmt.Println("No default namespace set")
		} else {
			fmt.Printf("Default namespace: %s\n", ns)
			fmt.Printf("Source: %s\n", source)
		}
		return nil
	}

	newNamespace := args[0]

	// Validate namespace exists before setting
	if !NamespaceExists(newNamespace) {
		return fmt.Errorf("namespace '%s' does not exist", newNamespace)
	}

	// Persist to config file
	if err := setDefaultNamespace(newNamespace); err != nil {
		return err
	}

	// Update session state
	s.SetNamespace(newNamespace)

	// Clear cache so validation picks up the change
	ClearNamespaceCache()

	fmt.Printf("Default namespace set to: %s\n", newNamespace)
	return nil
}
