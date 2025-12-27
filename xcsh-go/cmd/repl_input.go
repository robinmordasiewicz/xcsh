package cmd

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/robinmordasiewicz/xcsh/pkg/branding"
)

// InputBox represents the command input component with horizontal rules
type InputBox struct {
	textInput textinput.Model
	width     int
	prompt    string
}

// Input box styles
var promptStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff")).
	Bold(true)

// NewInputBox creates a new input box with default settings
func NewInputBox() *InputBox {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#E4002B"))

	return &InputBox{
		textInput: ti,
		prompt:    "<xc> ",
	}
}

// SetWidth sets the width of the input box
func (i *InputBox) SetWidth(width int) {
	i.width = width
	i.textInput.Width = width - len(i.prompt) - 2
}

// SetPrompt sets the prompt string (e.g., "<xc.domain.action> ")
func (i *InputBox) SetPrompt(prompt string) {
	i.prompt = prompt
	if i.width > 0 {
		i.textInput.Width = i.width - len(i.prompt) - 2
	}
}

// Focus focuses the input
func (i *InputBox) Focus() tea.Cmd {
	return i.textInput.Focus()
}

// Blur removes focus from the input
func (i *InputBox) Blur() {
	i.textInput.Blur()
}

// Value returns the current input value
func (i *InputBox) Value() string {
	return i.textInput.Value()
}

// SetValue sets the input value
func (i *InputBox) SetValue(value string) {
	i.textInput.SetValue(value)
}

// Reset clears the input
func (i *InputBox) Reset() {
	i.textInput.Reset()
}

// Update handles input events
func (i *InputBox) Update(msg tea.Msg) (*InputBox, tea.Cmd) {
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

// View renders the input box with horizontal rules
func (i *InputBox) View() string {
	var sb strings.Builder

	// Top horizontal rule
	sb.WriteString(i.renderRule())
	sb.WriteString("\n")

	// Prompt and input
	sb.WriteString(promptStyle.Render(i.prompt))
	sb.WriteString(i.textInput.View())
	sb.WriteString("\n")

	// Bottom horizontal rule
	sb.WriteString(i.renderRule())

	return sb.String()
}

// renderRule renders a horizontal rule in F5 red
func (i *InputBox) renderRule() string {
	width := i.width
	if width <= 0 {
		width = 80
	}

	rule := strings.Repeat("─", width)
	return branding.ColorRed + rule + branding.ColorReset
}

// CursorEnd moves cursor to end of input
func (i *InputBox) CursorEnd() {
	i.textInput.CursorEnd()
}

// CursorStart moves cursor to start of input
func (i *InputBox) CursorStart() {
	i.textInput.CursorStart()
}
