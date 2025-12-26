package cmd

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
	"github.com/robinmordasiewicz/xcsh/pkg/client"
)

// renderWelcomeBanner creates the modern CLI welcome banner with F5 logo in a dark red frame
// Vertical layout: logo centered on top, info text below
func renderWelcomeBanner() string {
	var sb strings.Builder
	const frameWidth = 80
	const innerWidth = frameWidth - 2 // Account for left and right borders

	// Add leading newline for visual separation
	sb.WriteString("\n")

	// Build title for top border (title in bold white)
	title := fmt.Sprintf(" %s v%s ", branding.CLIFullName, Version)
	titleLen := len(title)
	dashesAfterTitle := frameWidth - 5 - titleLen

	// Top border: ╭─── Title ───...─╮ with title in bold white
	sb.WriteString(branding.ColorDarkRed + "╭───" + branding.ColorReset)
	sb.WriteString(branding.ColorBoldWhite + title + branding.ColorReset)
	sb.WriteString(branding.ColorDarkRed + strings.Repeat("─", dashesAfterTitle) + "╮" + branding.ColorReset + "\n")

	// Get logo lines
	logoLines := strings.Split(branding.F5Logo, "\n")

	// Render logo lines using embedded spacing from branding.go
	for _, line := range logoLines {
		coloredLine := colorizeLogoLine(line)
		lineWidth := runeWidth(line)

		// Use logo's embedded spacing, only right-pad to fill frame
		rightPad := innerWidth - lineWidth
		if rightPad < 0 {
			rightPad = 0
		}

		sb.WriteString(branding.ColorDarkRed + "│" + branding.ColorReset)
		sb.WriteString(coloredLine)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(branding.ColorDarkRed + "│" + branding.ColorReset + "\n")
	}

	// Add separator line
	sb.WriteString(branding.ColorDarkRed + "├" + strings.Repeat("─", innerWidth) + "┤" + branding.ColorReset + "\n")

	// Info content below logo
	infoLines := []string{
		"Type 'help' for commands, 'exit' or Ctrl+D to quit. Tab completion available.",
		buildConnectionInfo(),
	}

	for _, line := range infoLines {
		lineWidth := runeWidth(line)
		// Center info text
		leftPad := (innerWidth - lineWidth) / 2
		rightPad := innerWidth - lineWidth - leftPad
		if leftPad < 0 {
			leftPad = 0
		}
		if rightPad < 0 {
			rightPad = 0
		}

		sb.WriteString(branding.ColorDarkRed + "│" + branding.ColorReset)
		sb.WriteString(strings.Repeat(" ", leftPad))
		sb.WriteString(branding.ColorBoldWhite + line + branding.ColorReset)
		sb.WriteString(strings.Repeat(" ", rightPad))
		sb.WriteString(branding.ColorDarkRed + "│" + branding.ColorReset + "\n")
	}

	// Bottom border: ╰───...───╯
	bottomBorder := "╰" + strings.Repeat("─", innerWidth) + "╯"
	sb.WriteString(branding.ColorDarkRed + bottomBorder + branding.ColorReset + "\n")

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

// buildConnectionInfo returns tenant and API info string
func buildConnectionInfo() string {
	if serverURL == "" {
		return "Not connected · Set F5XC_API_URL to connect"
	}

	tenant := client.ExtractTenant(serverURL)
	// Extract domain from URL for display
	domain := extractDomain(serverURL)

	return fmt.Sprintf("Tenant: %s · API: %s", tenant, domain)
}

// extractDomain extracts the domain from a URL for compact display
func extractDomain(url string) string {
	// Remove protocol
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	// Remove trailing slashes
	url = strings.TrimSuffix(url, "/")
	return url
}

// runeWidth returns the display width of a string in terminal columns
// Uses default runewidth mode where block elements (▓, ▒, █) are single-width
func runeWidth(s string) int {
	return runewidth.StringWidth(s)
}
