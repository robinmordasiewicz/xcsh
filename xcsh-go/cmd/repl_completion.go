package cmd

import (
	"strings"

	"github.com/robinmordasiewicz/xcsh/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Suggestion represents a completion suggestion (go-prompt independent)
type Suggestion struct {
	Text        string
	Description string
}

// Completer provides tab completion for the REPL
type Completer struct {
	rootCmd *cobra.Command
	session *REPLSession
}

// NewCompleter creates a new completer
func NewCompleter(root *cobra.Command) *Completer {
	return &Completer{rootCmd: root}
}

// SetSession sets the session for context-aware completions
func (c *Completer) SetSession(session *REPLSession) {
	c.session = session
}

// Complete returns suggestions for the given input text
func (c *Completer) Complete(text string) []Suggestion {
	text = strings.TrimLeft(text, " ")

	// Handle empty input - show context-aware suggestions
	if text == "" {
		return c.getContextualSuggestions()
	}

	// Parse input into args
	args := parseInputArgs(text)
	if len(args) == 0 {
		return c.getContextualSuggestions()
	}

	// Check for "/" prefix - escape to root context
	isEscapedToRoot := false
	if strings.HasPrefix(args[0], "/") {
		isEscapedToRoot = true
		args[0] = strings.TrimPrefix(args[0], "/")
		if args[0] == "" {
			return c.getRootContextSuggestions()
		}
	}

	// Get current word being typed
	currentWord := ""
	if len(args) > 0 {
		currentWord = args[len(args)-1]
	}

	// If current word starts with -, complete flags
	if strings.HasPrefix(currentWord, "-") {
		var contextArgs []string
		if isEscapedToRoot {
			contextArgs = args
		} else {
			contextArgs = c.prependContextToArgs(args)
		}
		cmd, _, _ := c.rootCmd.Find(contextArgs)
		if cmd != nil {
			return c.getFlagCompletions(cmd, currentWord)
		}
	}

	// Build context-aware args
	var contextArgs []string
	if isEscapedToRoot {
		contextArgs = args
	} else {
		contextArgs = c.prependContextToArgs(args)
	}

	// Find the command being completed
	cmd, remainingArgs, err := c.rootCmd.Find(contextArgs)
	if err != nil || cmd == nil {
		if isEscapedToRoot {
			return filterSuggestions(c.getRootContextSuggestions(), args[0])
		}
		return filterSuggestions(c.getContextualSuggestions(), args[0])
	}

	// If we found the root command and have args, filter subcommands
	if cmd == c.rootCmd && len(args) > 0 {
		if isEscapedToRoot {
			return filterSuggestions(c.getRootContextSuggestions(), args[0])
		}
		return filterSuggestions(c.getContextualSuggestions(), args[0])
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
		return filterSuggestions(suggestions, currentWord)
	}
	return suggestions
}

// getContextualSuggestions returns suggestions based on current context
func (c *Completer) getContextualSuggestions() []Suggestion {
	if c.session == nil {
		return c.getTopLevelSuggestions()
	}

	ctx := c.session.GetContextPath()

	if ctx.IsRoot() {
		return c.getRootContextSuggestions()
	}

	if ctx.IsDomain() {
		return c.getDomainContextSuggestions()
	}

	if ctx.IsAction() {
		return c.getActionContextSuggestions()
	}

	return c.getTopLevelSuggestions()
}

// getRootContextSuggestions returns suggestions for root context
func (c *Completer) getRootContextSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	// Add domains from registry
	suggestions = append(suggestions, c.getDomainSuggestions()...)

	// Add navigation and built-in commands
	suggestions = append(suggestions,
		Suggestion{Text: "quit", Description: "Exit the shell"},
		Suggestion{Text: "help", Description: "Show help information"},
		Suggestion{Text: "clear", Description: "Clear the screen"},
		Suggestion{Text: "history", Description: "Show command history"},
		Suggestion{Text: "namespace", Description: "Set default namespace"},
		Suggestion{Text: "ns", Description: "Set default namespace (alias)"},
		Suggestion{Text: "context", Description: "Show current context"},
		Suggestion{Text: "ctx", Description: "Show current context (alias)"},
	)

	return suggestions
}

// getDomainContextSuggestions returns suggestions when in a domain context
func (c *Completer) getDomainContextSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	suggestions = append(suggestions, c.getActionSuggestions()...)

	suggestions = append(suggestions,
		Suggestion{Text: "exit", Description: "Go up to root context"},
		Suggestion{Text: "back", Description: "Go up to root context"},
		Suggestion{Text: "..", Description: "Go up to root context"},
		Suggestion{Text: "help", Description: "Show context help"},
	)

	return suggestions
}

// getActionContextSuggestions returns suggestions when in an action context
func (c *Completer) getActionContextSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	if c.session != nil {
		ctx := c.session.GetContextPath()
		args := []string{ctx.Domain, ctx.Action}
		cmd, _, err := c.rootCmd.Find(args)
		if err == nil && cmd != nil {
			suggestions = append(suggestions, c.getFlagCompletions(cmd, "-")...)
		}
	}

	suggestions = append(suggestions,
		Suggestion{Text: "exit", Description: "Go up to domain context"},
		Suggestion{Text: "back", Description: "Go up to domain context"},
		Suggestion{Text: "..", Description: "Go up to domain context"},
		Suggestion{Text: "root", Description: "Go to root context"},
		Suggestion{Text: "/", Description: "Go to root context"},
		Suggestion{Text: "help", Description: "Show context help"},
	)

	return suggestions
}

// getDomainSuggestions returns all domain suggestions
func (c *Completer) getDomainSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	for domain, meta := range types.DomainRegistry {
		suggestions = append(suggestions, Suggestion{
			Text:        domain,
			Description: meta.Description,
		})
	}

	return suggestions
}

// getActionSuggestions returns action suggestions
func (c *Completer) getActionSuggestions() []Suggestion {
	return []Suggestion{
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

// prependContextToArgs adds context prefix to args
func (c *Completer) prependContextToArgs(args []string) []string {
	if c.session == nil || len(args) == 0 {
		return args
	}

	ctx := c.session.GetContextPath()
	if ctx.IsRoot() {
		return args
	}

	cmd, _, err := c.rootCmd.Find(args)
	if err == nil && cmd != c.rootCmd {
		return args
	}

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
func (c *Completer) getTopLevelSuggestions() []Suggestion {
	suggestions := []Suggestion{}

	for _, cmd := range c.rootCmd.Commands() {
		if cmd.Hidden {
			continue
		}
		suggestions = append(suggestions, Suggestion{
			Text:        cmd.Name(),
			Description: cmd.Short,
		})
	}

	suggestions = append(suggestions,
		Suggestion{Text: "exit", Description: "Exit the shell"},
		Suggestion{Text: "quit", Description: "Exit the shell"},
		Suggestion{Text: "help", Description: "Show help information"},
		Suggestion{Text: "clear", Description: "Clear the screen"},
		Suggestion{Text: "history", Description: "Show command history"},
		Suggestion{Text: "namespace", Description: "Set default namespace"},
		Suggestion{Text: "ns", Description: "Set default namespace (alias)"},
	)

	return suggestions
}

// getSubcommandSuggestions returns suggestions for subcommands
func (c *Completer) getSubcommandSuggestions(cmd *cobra.Command) []Suggestion {
	suggestions := []Suggestion{}

	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		suggestions = append(suggestions, Suggestion{
			Text:        sub.Name(),
			Description: sub.Short,
		})
	}

	return suggestions
}

// getFlagCompletions returns flag suggestions for a command
func (c *Completer) getFlagCompletions(cmd *cobra.Command, prefix string) []Suggestion {
	suggestions := []Suggestion{}

	addFlag := func(flag *pflag.Flag, isGlobal bool) {
		if flag.Hidden {
			return
		}
		desc := flag.Usage
		if isGlobal {
			desc += " (global)"
		}

		longFlag := "--" + flag.Name
		if strings.HasPrefix(longFlag, prefix) || prefix == "-" || prefix == "--" {
			suggestions = append(suggestions, Suggestion{
				Text:        longFlag,
				Description: desc,
			})
		}

		if flag.Shorthand != "" {
			shortFlag := "-" + flag.Shorthand
			if strings.HasPrefix(shortFlag, prefix) || prefix == "-" {
				suggestions = append(suggestions, Suggestion{
					Text:        shortFlag,
					Description: desc,
				})
			}
		}
	}

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		addFlag(flag, false)
	})

	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		addFlag(flag, true)
	})

	return suggestions
}

// callValidArgsFunction invokes Cobra's ValidArgsFunction
func (c *Completer) callValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) []Suggestion {
	if cmd.ValidArgsFunction == nil {
		return nil
	}

	var completionArgs []string
	if len(args) > 1 {
		_, remaining, _ := c.rootCmd.Find(args)
		if len(remaining) > 1 {
			completionArgs = remaining[:len(remaining)-1]
		}
	}

	completions, _ := cmd.ValidArgsFunction(cmd, completionArgs, toComplete)

	suggestions := make([]Suggestion, 0, len(completions))
	for _, comp := range completions {
		parts := strings.SplitN(comp, "\t", 2)
		s := Suggestion{Text: parts[0]}
		if len(parts) > 1 {
			s.Description = parts[1]
		}
		suggestions = append(suggestions, s)
	}

	return suggestions
}

// filterSuggestions filters suggestions by prefix
func filterSuggestions(suggestions []Suggestion, prefix string) []Suggestion {
	if prefix == "" {
		return suggestions
	}

	prefix = strings.ToLower(prefix)
	filtered := []Suggestion{}
	for _, s := range suggestions {
		if strings.HasPrefix(strings.ToLower(s.Text), prefix) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
