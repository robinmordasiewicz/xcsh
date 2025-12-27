package cmd

import (
	"strings"
)

// buildPlainPrompt constructs the prompt string
// Format: <xc.domain.action>
func buildPlainPrompt(session *REPLSession) string {
	parts := []string{"xc"}

	ctx := session.GetContextPath()
	if ctx.Domain != "" {
		parts = append(parts, ctx.Domain)
		if ctx.Action != "" {
			parts = append(parts, ctx.Action)
		}
	}

	return "<" + strings.Join(parts, ".") + "> "
}
