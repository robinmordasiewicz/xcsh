package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorCyan    = "\033[36m" // tenant
	colorGreen   = "\033[32m" // domain
	colorYellow  = "\033[33m" // action
	colorMagenta = "\033[35m" // namespace
	colorBold    = "\033[1m"  // bold
)

// buildColoredPrompt constructs the prompt with colored segments
// Format: tenant:domain/action@namespace>
func buildColoredPrompt(session *REPLSession) string {
	if !session.IsColorEnabled() {
		return buildPlainPrompt(session)
	}

	var parts []string

	// Tenant segment (cyan)
	tenant := session.GetTenant()
	if tenant != "" && tenant != "unknown" && tenant != "local" {
		parts = append(parts, fmt.Sprintf("%s%s%s", colorCyan, tenant, colorReset))
	}

	// Context path segment (green for domain, yellow for action)
	ctx := session.GetContextPath()
	if ctx.Domain != "" {
		domainPart := fmt.Sprintf("%s%s%s", colorGreen, ctx.Domain, colorReset)
		if ctx.Action != "" {
			domainPart += fmt.Sprintf("/%s%s%s", colorYellow, ctx.Action, colorReset)
		}
		parts = append(parts, domainPart)
	}

	// Namespace segment (magenta)
	ns := session.GetNamespace()
	if ns != "" {
		parts = append(parts, fmt.Sprintf("@%s%s%s", colorMagenta, ns, colorReset))
	}

	// Build final prompt
	if len(parts) == 0 {
		return fmt.Sprintf("%sxcsh%s> ", colorBold, colorReset)
	}

	// Join parts with colons, but namespace uses @ prefix
	prompt := ""
	for i, part := range parts {
		if i == 0 {
			prompt = part
		} else if strings.HasPrefix(part, "@") {
			prompt += part // namespace already has @ prefix
		} else {
			prompt += ":" + part
		}
	}

	return prompt + "> "
}

// buildPlainPrompt constructs a non-colored prompt for non-color terminals
func buildPlainPrompt(session *REPLSession) string {
	var parts []string

	tenant := session.GetTenant()
	if tenant != "" && tenant != "unknown" && tenant != "local" {
		parts = append(parts, tenant)
	}

	ctx := session.GetContextPath()
	if ctx.Domain != "" {
		contextStr := ctx.Domain
		if ctx.Action != "" {
			contextStr += "/" + ctx.Action
		}
		parts = append(parts, contextStr)
	}

	ns := session.GetNamespace()
	if ns != "" {
		parts = append(parts, "@"+ns)
	}

	if len(parts) == 0 {
		return "xcsh> "
	}

	// Join parts with colons, but namespace uses @ prefix
	prompt := ""
	for i, part := range parts {
		if i == 0 {
			prompt = part
		} else if strings.HasPrefix(part, "@") {
			prompt += part
		} else {
			prompt += ":" + part
		}
	}

	return prompt + "> "
}

// detectColorSupport checks if terminal supports colors
func detectColorSupport() bool {
	// Check for NO_COLOR environment variable (standard)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for TERM=dumb
	if os.Getenv("TERM") == "dumb" {
		return false
	}

	// Check if stdout is a terminal
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return false
	}

	return true
}
