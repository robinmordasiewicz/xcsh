package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
	"github.com/robinmordasiewicz/xcsh/pkg/client"
)

// Padding constants for banner layout
const (
	logoPadLeft   = 2 // Left padding for logo column
	logoPadRight  = 2 // Right padding for logo column
	logoPadTop    = 1 // Empty rows above logo
	logoPadBottom = 1 // Empty rows below logo
	textPadLeft   = 2 // Left padding for text column
	textPadRight  = 2 // Right padding for text column
)

// renderWelcomeBanner creates the modern CLI welcome banner with F5 logo in a dark red frame
// Two-column layout: logo on left, help text on right, connection info in bottom row
func renderWelcomeBanner() string {
	var sb strings.Builder

	// Get logo lines first to calculate widths
	logoLines := strings.Split(branding.F5Logo, "\n")
	logoHeight := len(logoLines)

	// Define text content (3 lines)
	textLines := []string{
		"Type 'help' for documentation.",
		"Run 'namespace' to set namespace.",
		"Press Ctrl+C twice to quit.",
	}

	// Calculate logo content width based on widest logo line
	logoContentWidth := 0
	for _, line := range logoLines {
		if w := runeWidth(line); w > logoContentWidth {
			logoContentWidth = w
		}
	}

	// Calculate text content width based on longest text line
	textContentWidth := 0
	for _, line := range textLines {
		if w := runeWidth(line); w > textContentWidth {
			textContentWidth = w
		}
	}

	// Calculate column widths including padding
	logoColWidth := logoPadLeft + logoContentWidth + logoPadRight
	textColWidth := textPadLeft + textContentWidth + textPadRight

	// Calculate frame dimensions based on content
	// innerWidth = logoColWidth + divider(1) + textColWidth
	innerWidth := logoColWidth + 1 + textColWidth
	frameWidth := innerWidth + 2 // Add 2 for left and right borders

	// Total content rows including padding
	totalRows := logoPadTop + logoHeight + logoPadBottom

	// Add leading newline for visual separation
	sb.WriteString("\n")

	// Build title for top border (title in bold white)
	title := fmt.Sprintf(" %s v%s ", branding.CLIFullName, Version)
	titleLen := len(title)
	dashesAfterTitle := frameWidth - 5 - titleLen

	// Top border: ╭─── Title ───...─╮ with title in bold white
	sb.WriteString(branding.ColorRed + "╭───" + branding.ColorReset)
	sb.WriteString(branding.ColorBoldWhite + title + branding.ColorReset)
	sb.WriteString(branding.ColorRed + strings.Repeat("─", dashesAfterTitle) + "╮" + branding.ColorReset + "\n")

	// Calculate vertical centering for text within total height
	textStartRow := (totalRows - len(textLines)) / 2

	// Two-column rows: logo on left (with padding), text on right (with padding)
	for i := 0; i < totalRows; i++ {
		// Left border (dark red)
		sb.WriteString(branding.ColorRed + "│" + branding.ColorReset)

		// Logo column: leftPad + content/empty + rightPad
		sb.WriteString(strings.Repeat(" ", logoPadLeft))

		// Determine if this row has logo content
		logoRowIdx := i - logoPadTop
		if logoRowIdx >= 0 && logoRowIdx < logoHeight {
			logoLine := logoLines[logoRowIdx]
			coloredLogo := colorizeLogoLine(logoLine)
			logoLineWidth := runeWidth(logoLine)
			logoPad := logoContentWidth - logoLineWidth
			if logoPad < 0 {
				logoPad = 0
			}
			sb.WriteString(coloredLogo)
			sb.WriteString(strings.Repeat(" ", logoPad))
		} else {
			// Empty row (top or bottom padding)
			sb.WriteString(strings.Repeat(" ", logoContentWidth))
		}

		sb.WriteString(strings.Repeat(" ", logoPadRight))

		// Inner divider (red, not bold)
		sb.WriteString(branding.ColorRed + "│" + branding.ColorReset)

		// Text column: leftPad + content/empty + rightPad
		sb.WriteString(strings.Repeat(" ", textPadLeft))

		// Determine if this row has text content
		textLineIdx := i - textStartRow
		if textLineIdx >= 0 && textLineIdx < len(textLines) {
			text := textLines[textLineIdx]
			textLineWidth := runeWidth(text)
			textPad := textContentWidth - textLineWidth
			if textPad < 0 {
				textPad = 0
			}
			sb.WriteString(branding.ColorBoldWhite + text + branding.ColorReset)
			sb.WriteString(strings.Repeat(" ", textPad))
		} else {
			sb.WriteString(strings.Repeat(" ", textContentWidth))
		}

		sb.WriteString(strings.Repeat(" ", textPadRight))

		// Right border (dark red)
		sb.WriteString(branding.ColorRed + "│" + branding.ColorReset + "\n")
	}

	// Full-width separator before connection info
	sb.WriteString(branding.ColorRed + "├" + strings.Repeat("─", innerWidth) + "┤" + branding.ColorReset + "\n")

	// Connection info rows (full width, left-aligned with padding)
	connLines := buildConnectionInfo()
	for _, connLine := range connLines {
		connLineWidth := runeWidth(connLine)
		rightPad := innerWidth - connLineWidth - 2 // 2 for left padding
		if rightPad < 0 {
			rightPad = 0
		}
		sb.WriteString(branding.ColorRed + "│" + branding.ColorReset)
		sb.WriteString("  ") // Left padding
		sb.WriteString(branding.ColorBoldWhite + connLine + branding.ColorReset)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(branding.ColorRed + "│" + branding.ColorReset + "\n")
	}

	// Bottom border: ╰───...───╯
	bottomBorder := "╰" + strings.Repeat("─", innerWidth) + "╯"
	sb.WriteString(branding.ColorRed + bottomBorder + branding.ColorReset + "\n")

	return sb.String()
}

// colorizeLogoLine applies red color to circle background and white color to F5 text
// The logo uses:
// - ▓ for the red circle background
// - █ and ▒ for the white F5 text
// - (, ), |, and _ for the circle outline (rendered in red)
func colorizeLogoLine(line string) string {
	var result strings.Builder
	inRed := false
	inWhite := false

	for _, r := range line {
		switch r {
		case '▓', '▒', '(', ')', '|', '_':
			// Red for circle background and outline
			if !inRed {
				if inWhite {
					result.WriteString(branding.ColorReset)
					inWhite = false
				}
				result.WriteString(branding.ColorRed)
				inRed = true
			}
			if r == '▓' {
				result.WriteRune('█') // Render as solid block
			} else {
				result.WriteRune(r) // Keep ▒, (, ) as-is
			}
		case '█':
			// White for F5 text elements
			if !inWhite {
				if inRed {
					result.WriteString(branding.ColorReset)
					inRed = false
				}
				result.WriteString(branding.ColorBoldWhite)
				inWhite = true
			}
			result.WriteRune(r)
		default:
			// Reset for spaces and other characters
			if inRed || inWhite {
				result.WriteString(branding.ColorReset)
				inRed = false
				inWhite = false
			}
			result.WriteRune(r)
		}
	}

	// Final reset if we ended in a color
	if inRed || inWhite {
		result.WriteString(branding.ColorReset)
	}

	return result.String()
}

// getUserInfo calls the whoami endpoint to get the logged-in user's name
func getUserInfo() string {
	if apiClient == nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := apiClient.Get(ctx, "/api/web/custom/namespaces/system/whoami", nil)
	if err != nil || resp.StatusCode >= 400 {
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return ""
	}

	// Extract user name/email from response
	if name, ok := result["name"].(string); ok {
		return name
	}
	return ""
}

// hasAuthConfig checks if any authentication method is configured
func hasAuthConfig() bool {
	// Check for API token in environment
	if os.Getenv("F5XC_API_TOKEN") != "" {
		return true
	}
	// Check for certificate paths
	if cert != "" || p12Bundle != "" {
		return true
	}
	// Check if apiClient was successfully initialized
	return apiClient != nil
}

// buildConnectionInfo returns connection detail lines for the banner
func buildConnectionInfo() []string {
	// No API URL configured
	if serverURL == "" {
		return []string{
			"Not connected · Set F5XC_API_URL to connect",
			fmt.Sprintf("• Namespace: %s", GetDefaultNamespace()),
		}
	}

	// API URL set but no authentication configured
	if !hasAuthConfig() {
		return []string{
			"To authenticate type: /login",
			fmt.Sprintf("• Namespace: %s", GetDefaultNamespace()),
		}
	}

	tenant := client.ExtractTenant(serverURL)
	username := getUserInfo()

	lines := []string{
		fmt.Sprintf("• Tenant: %s", tenant),
	}

	if username != "" {
		lines = append(lines, fmt.Sprintf("• User: %s", username))
	}

	// Use validated namespace (already validated during session init)
	lines = append(lines, fmt.Sprintf("• Namespace: %s", GetValidatedDefaultNamespace()))

	return lines
}

// runeWidth returns the display width of a string in terminal columns
// Uses default runewidth mode where block elements (▓, ▒, █) are single-width
func runeWidth(s string) int {
	return runewidth.StringWidth(s)
}
