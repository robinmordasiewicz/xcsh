package cmd

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CobraCompleter adapts Cobra's completion system for go-prompt
type CobraCompleter struct {
	rootCmd *cobra.Command
	session *REPLSession // Session reference for context-aware completions
}

// NewCobraCompleter creates a new Cobra completion adapter
func NewCobraCompleter(root *cobra.Command) *CobraCompleter {
	return &CobraCompleter{rootCmd: root}
}

// SetSession sets the session reference for context-aware completions
func (c *CobraCompleter) SetSession(session *REPLSession) {
	c.session = session
}

// Complete generates suggestions for the current input
func (c *CobraCompleter) Complete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()

	// Handle empty input - show context-aware suggestions
	if strings.TrimSpace(text) == "" {
		return c.getContextualSuggestions()
	}

	// Parse input into command path and current word
	args := parseInputArgs(text)
	if len(args) == 0 {
		return c.getContextualSuggestions()
	}

	// Check if we're completing a flag value or flag name
	currentWord := d.GetWordBeforeCursor()

	// If current word starts with -, complete flags
	if strings.HasPrefix(currentWord, "-") {
		// Build context-aware args for flag completion
		contextArgs := c.prependContextToArgs(args)
		cmd, _, _ := c.rootCmd.Find(contextArgs)
		if cmd != nil {
			return c.getFlagCompletions(cmd, currentWord)
		}
	}

	// Build context-aware args for command finding
	contextArgs := c.prependContextToArgs(args)

	// Find the command being completed
	cmd, remainingArgs, err := c.rootCmd.Find(contextArgs)
	if err != nil || cmd == nil {
		// Could be partial command - filter by prefix from contextual suggestions
		return prompt.FilterHasPrefix(c.getContextualSuggestions(), args[0], true)
	}

	// If we found the root command and have args, filter subcommands
	if cmd == c.rootCmd && len(args) > 0 {
		return prompt.FilterHasPrefix(c.getContextualSuggestions(), args[0], true)
	}

	// Check if the last arg is a partial flag
	if len(remainingArgs) > 0 && strings.HasPrefix(remainingArgs[len(remainingArgs)-1], "-") {
		return c.getFlagCompletions(cmd, remainingArgs[len(remainingArgs)-1])
	}

	// Use Cobra's ValidArgsFunction if available
	if cmd.ValidArgsFunction != nil {
		suggestions := c.callValidArgsFunction(cmd, contextArgs, currentWord)
		if len(suggestions) > 0 {
			return suggestions
		}
	}

	// Fall back to subcommand completion
	suggestions := c.getSubcommandSuggestions(cmd)
	if currentWord != "" {
		return prompt.FilterHasPrefix(suggestions, currentWord, true)
	}
	return suggestions
}

// getContextualSuggestions returns suggestions based on current context
func (c *CobraCompleter) getContextualSuggestions() []prompt.Suggest {
	if c.session == nil {
		return c.getTopLevelSuggestions()
	}

	ctx := c.session.GetContextPath()

	// At root - show domains and built-in commands
	if ctx.IsRoot() {
		return c.getRootContextSuggestions()
	}

	// In domain context - show actions and navigation
	if ctx.IsDomain() {
		return c.getDomainContextSuggestions()
	}

	// In action context - show flags and navigation
	if ctx.IsAction() {
		return c.getActionContextSuggestions()
	}

	return c.getTopLevelSuggestions()
}

// getRootContextSuggestions returns suggestions for root context
func (c *CobraCompleter) getRootContextSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// Add domains from registry
	suggestions = append(suggestions, c.getDomainSuggestions()...)

	// Add navigation and built-in commands
	suggestions = append(suggestions,
		prompt.Suggest{Text: "quit", Description: "Exit the shell"},
		prompt.Suggest{Text: "help", Description: "Show help information"},
		prompt.Suggest{Text: "clear", Description: "Clear the screen"},
		prompt.Suggest{Text: "history", Description: "Show command history"},
		prompt.Suggest{Text: "namespace", Description: "Set default namespace"},
		prompt.Suggest{Text: "ns", Description: "Set default namespace (alias)"},
		prompt.Suggest{Text: "context", Description: "Show current context"},
		prompt.Suggest{Text: "ctx", Description: "Show current context (alias)"},
	)

	return suggestions
}

// getDomainContextSuggestions returns suggestions when in a domain context
func (c *CobraCompleter) getDomainContextSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// Add actions
	suggestions = append(suggestions, c.getActionSuggestions()...)

	// Add navigation commands
	suggestions = append(suggestions,
		prompt.Suggest{Text: "exit", Description: "Go up to root context"},
		prompt.Suggest{Text: "back", Description: "Go up to root context"},
		prompt.Suggest{Text: "..", Description: "Go up to root context"},
		prompt.Suggest{Text: "help", Description: "Show context help"},
	)

	return suggestions
}

// getActionContextSuggestions returns suggestions when in an action context
func (c *CobraCompleter) getActionContextSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	// Get flags for the current command context
	if c.session != nil {
		ctx := c.session.GetContextPath()
		args := []string{ctx.Domain, ctx.Action}
		cmd, _, err := c.rootCmd.Find(args)
		if err == nil && cmd != nil {
			suggestions = append(suggestions, c.getFlagCompletions(cmd, "-")...)
		}
	}

	// Add navigation commands
	suggestions = append(suggestions,
		prompt.Suggest{Text: "exit", Description: "Go up to domain context"},
		prompt.Suggest{Text: "back", Description: "Go up to domain context"},
		prompt.Suggest{Text: "..", Description: "Go up to domain context"},
		prompt.Suggest{Text: "root", Description: "Go to root context"},
		prompt.Suggest{Text: "/", Description: "Go to root context"},
		prompt.Suggest{Text: "help", Description: "Show context help"},
	)

	return suggestions
}

// getDomainSuggestions returns all domain suggestions from the registry
func (c *CobraCompleter) getDomainSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	for domain, meta := range types.DomainRegistry {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        domain,
			Description: meta.Description,
		})
	}

	return suggestions
}

// getActionSuggestions returns action suggestions
func (c *CobraCompleter) getActionSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "list", Description: "List resources"},
		{Text: "get", Description: "Get a specific resource"},
		{Text: "create", Description: "Create a new resource"},
		{Text: "delete", Description: "Delete a resource"},
		{Text: "replace", Description: "Replace a resource"},
		{Text: "apply", Description: "Apply configuration from file"},
		{Text: "status", Description: "Get resource status"},
		{Text: "patch", Description: "Patch a resource"},
		{Text: "add-labels", Description: "Add labels to a resource"},
		{Text: "remove-labels", Description: "Remove labels from a resource"},
	}
}

// prependContextToArgs adds context prefix to args for Cobra command lookup
func (c *CobraCompleter) prependContextToArgs(args []string) []string {
	if c.session == nil || len(args) == 0 {
		return args
	}

	ctx := c.session.GetContextPath()
	if ctx.IsRoot() {
		return args
	}

	// Check if first arg is already a valid root command
	cmd, _, err := c.rootCmd.Find(args)
	if err == nil && cmd != c.rootCmd {
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
	prefixedArgs = append(prefixedArgs, args...)

	return prefixedArgs
}

// getTopLevelSuggestions returns all available top-level commands
func (c *CobraCompleter) getTopLevelSuggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Hidden {
			continue
		}
		suggestions = append(suggestions, prompt.Suggest{
			Text:        cmd.Name(),
			Description: cmd.Short,
		})
	}

	// Add REPL built-in commands
	suggestions = append(suggestions,
		prompt.Suggest{Text: "exit", Description: "Exit the shell"},
		prompt.Suggest{Text: "quit", Description: "Exit the shell"},
		prompt.Suggest{Text: "help", Description: "Show help information"},
		prompt.Suggest{Text: "clear", Description: "Clear the screen"},
		prompt.Suggest{Text: "history", Description: "Show command history"},
		prompt.Suggest{Text: "namespace", Description: "Set default namespace"},
		prompt.Suggest{Text: "ns", Description: "Set default namespace (alias)"},
	)

	return suggestions
}

// getSubcommandSuggestions returns suggestions for subcommands of a command
func (c *CobraCompleter) getSubcommandSuggestions(cmd *cobra.Command) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		suggestions = append(suggestions, prompt.Suggest{
			Text:        sub.Name(),
			Description: sub.Short,
		})
	}

	return suggestions
}

// getFlagCompletions returns flag suggestions for a command
func (c *CobraCompleter) getFlagCompletions(cmd *cobra.Command, prefix string) []prompt.Suggest {
	suggestions := []prompt.Suggest{}

	addFlag := func(flag *pflag.Flag, isGlobal bool) {
		if flag.Hidden {
			return
		}
		desc := flag.Usage
		if isGlobal {
			desc += " (global)"
		}

		// Long flag
		longFlag := "--" + flag.Name
		if strings.HasPrefix(longFlag, prefix) || prefix == "-" || prefix == "--" {
			suggestions = append(suggestions, prompt.Suggest{
				Text:        longFlag,
				Description: desc,
			})
		}

		// Short flag
		if flag.Shorthand != "" {
			shortFlag := "-" + flag.Shorthand
			if strings.HasPrefix(shortFlag, prefix) || prefix == "-" {
				suggestions = append(suggestions, prompt.Suggest{
					Text:        shortFlag,
					Description: desc,
				})
			}
		}
	}

	// Add local flags
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		addFlag(flag, false)
	})

	// Add inherited (persistent) flags
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		addFlag(flag, true)
	})

	return suggestions
}

// callValidArgsFunction invokes Cobra's ValidArgsFunction for dynamic completion
func (c *CobraCompleter) callValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) []prompt.Suggest {
	if cmd.ValidArgsFunction == nil {
		return nil
	}

	// Build the args to pass to ValidArgsFunction
	// We need to pass all args except the one being completed
	var completionArgs []string
	if len(args) > 1 {
		// Find where the subcommand ends and args begin
		_, remaining, _ := c.rootCmd.Find(args)
		if len(remaining) > 1 {
			completionArgs = remaining[:len(remaining)-1]
		}
	}

	// Call Cobra's completion function
	completions, _ := cmd.ValidArgsFunction(cmd, completionArgs, toComplete)

	suggestions := make([]prompt.Suggest, 0, len(completions))
	for _, comp := range completions {
		// Parse Cobra completion format: "value\tdescription"
		parts := strings.SplitN(comp, "\t", 2)
		s := prompt.Suggest{Text: parts[0]}
		if len(parts) > 1 {
			s.Description = parts[1]
		}
		suggestions = append(suggestions, s)
	}

	return suggestions
}
