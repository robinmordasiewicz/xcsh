package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
	"github.com/robinmordasiewicz/xcsh/pkg/git"
)

// Status bar colors using F5 branding
var (
	// F5 Brand Red for the status bar background
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a1a")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1)

	// Git status icon colors
	cleanStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00c853")) // Green

	dirtyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffc107")) // Yellow

	aheadBehindStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#2196f3")) // Blue

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))
)

// StatusBar represents the status bar state
type StatusBar struct {
	gitInfo git.RepoInfo
	width   int
}

// NewStatusBar creates a new status bar instance
func NewStatusBar() *StatusBar {
	return &StatusBar{
		gitInfo: git.GetRepoInfo(),
	}
}

// Refresh updates the git status information
func (s *StatusBar) Refresh() {
	s.gitInfo = git.GetRepoInfo()
}

// SetWidth sets the terminal width for rendering
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// Render returns the rendered status bar string
func (s *StatusBar) Render() string {
	left := s.renderLeft()
	right := s.renderRight()

	// Calculate available space for padding
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	padding := s.width - leftWidth - rightWidth - 2 // -2 for status bar padding

	if padding < 0 {
		padding = 0
	}

	// Build the status bar with left content, padding, and right content
	content := left
	for i := 0; i < padding; i++ {
		content += " "
	}
	content += right

	return statusBarStyle.Width(s.width).Render(content)
}

// renderLeft renders the left side of the status bar (repo info)
func (s *StatusBar) renderLeft() string {
	// Always show CLI name
	result := branding.CLIName

	// Add git info if in a repo
	if s.gitInfo.InRepo {
		result += " ("

		// Branch name
		result += s.gitInfo.Branch

		// Status icon with color
		icon := s.gitInfo.StatusIcon()
		if icon != "" {
			result += " "
			switch icon {
			case "✓":
				result += cleanStyle.Render(icon)
			case "*":
				result += dirtyStyle.Render(icon)
			case "↑", "↓", "↕":
				result += aheadBehindStyle.Render(icon)
			default:
				result += icon
			}
		}

		result += ")"
	}

	return result
}

// renderRight renders the right side of the status bar (hints)
func (s *StatusBar) renderRight() string {
	return dimStyle.Render("Ctrl+C: quit")
}
