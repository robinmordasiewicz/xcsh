package cmd

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TUIModel is the main bubbletea model for the REPL
type TUIModel struct {
	session            *REPLSession
	input              *InputBox
	statusBar          *StatusBar
	completer          *Completer   // Tab completion provider
	output             []string     // Output lines from command execution
	width              int          // Terminal width
	height             int          // Terminal height
	ready              bool         // Whether the terminal size is known
	quitting           bool         // Whether we're exiting
	history            []string     // Command history for up/down navigation
	historyIndex       int          // Current position in history (-1 = not navigating)
	lastCtrlC          time.Time    // For double-press detection
	suggestions        []Suggestion // Current completion suggestions
	suggestionIdx      int          // Selected suggestion index
	showingSuggestions bool         // Whether suggestions are visible
}

// Output area style
var outputStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ffffff"))

// NewTUIModel creates a new TUI model with the given session
func NewTUIModel(session *REPLSession) *TUIModel {
	// Initialize with welcome banner as first output
	banner := renderWelcomeBanner()
	initialOutput := strings.Split(strings.TrimSuffix(banner, "\n"), "\n")

	// Initialize completer
	completer := NewCompleter(rootCmd)
	completer.SetSession(session)

	return &TUIModel{
		session:      session,
		input:        NewInputBox(),
		statusBar:    NewStatusBar(),
		completer:    completer,
		output:       initialOutput,
		historyIndex: -1,
	}
}

// Init initializes the model
func (m *TUIModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.input.Focus(),
	)
}

// Update handles messages and updates the model
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(msg.Width)
		m.statusBar.SetWidth(msg.Width)
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			now := time.Now()
			if now.Sub(m.lastCtrlC) < 500*time.Millisecond {
				// Double Ctrl+C - exit
				m.quitting = true
				return m, tea.Quit
			}
			m.lastCtrlC = now
			m.addOutput("\nPress Ctrl+C again to exit, or continue typing")
			return m, nil

		case tea.KeyCtrlD:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			input := strings.TrimSpace(m.input.Value())

			// Special case: "/" alone navigates to root context
			if input == "/" {
				m.showingSuggestions = false
				m.suggestions = nil
				m.input.Reset()

				// Navigate to root
				ctx := m.session.GetContextPath()
				if !ctx.IsRoot() {
					ctx.Reset()
					m.addOutput("Navigated to root context")
				}
				m.input.SetPrompt(buildPlainPrompt(m.session))
				return m, nil
			}

			// If showing suggestions, apply the selected one
			if m.showingSuggestions && len(m.suggestions) > 0 {
				m.applyCompletion(m.suggestions[m.suggestionIdx])
				m.suggestions = nil
				m.showingSuggestions = false
				return m, nil
			}

			input = m.input.Value()
			m.input.Reset()

			if input != "" {
				// Add to display
				prompt := buildPlainPrompt(m.session)
				m.addOutput(prompt + input)

				// Execute command and capture output
				m.executeAndCapture(input)

				// Refresh git status after command
				m.statusBar.Refresh()

				// Reset history navigation
				m.historyIndex = -1
				m.history = m.session.history.GetHistory()
			}

			// Update prompt based on new context
			m.input.SetPrompt(buildPlainPrompt(m.session))
			return m, nil

		case tea.KeyUp:
			// If showing suggestions, navigate up in list
			if m.showingSuggestions && len(m.suggestions) > 0 {
				m.suggestionIdx--
				if m.suggestionIdx < 0 {
					m.suggestionIdx = len(m.suggestions) - 1
				}
				return m, nil
			}
			// Otherwise navigate history up
			if len(m.history) > 0 {
				if m.historyIndex < len(m.history)-1 {
					m.historyIndex++
					idx := len(m.history) - 1 - m.historyIndex
					m.input.SetValue(m.history[idx])
					m.input.CursorEnd()
				}
			}
			return m, nil

		case tea.KeyDown:
			// If showing suggestions, navigate down in list
			if m.showingSuggestions && len(m.suggestions) > 0 {
				m.suggestionIdx = (m.suggestionIdx + 1) % len(m.suggestions)
				return m, nil
			}
			// Otherwise navigate history down
			if m.historyIndex > 0 {
				m.historyIndex--
				idx := len(m.history) - 1 - m.historyIndex
				m.input.SetValue(m.history[idx])
				m.input.CursorEnd()
			} else if m.historyIndex == 0 {
				m.historyIndex = -1
				m.input.SetValue("")
			}
			return m, nil

		case tea.KeyTab:
			// Tab: cycle through suggestions or trigger completion
			if m.showingSuggestions && len(m.suggestions) > 0 {
				// Cycle through suggestions
				m.suggestionIdx = (m.suggestionIdx + 1) % len(m.suggestions)
			} else {
				// Get new suggestions
				m.triggerCompletion()
			}
			return m, nil

		case tea.KeyShiftTab:
			// Reverse cycle through suggestions
			if m.showingSuggestions && len(m.suggestions) > 0 {
				m.suggestionIdx--
				if m.suggestionIdx < 0 {
					m.suggestionIdx = len(m.suggestions) - 1
				}
			}
			return m, nil

		case tea.KeyEscape:
			// Cancel suggestions
			m.showingSuggestions = false
			m.suggestions = nil
			return m, nil

		case tea.KeyRunes:
			// "/" key is handled after input update to allow both typing and completion trigger

		case tea.KeyRight:
			// Apply selected suggestion on right arrow when showing suggestions
			if m.showingSuggestions && len(m.suggestions) > 0 {
				m.applyCompletion(m.suggestions[m.suggestionIdx])
				m.suggestions = nil
				m.showingSuggestions = false
				return m, nil
			}
		}
	}

	// Pass other messages to input
	var cmd tea.Cmd
	oldValue := m.input.Value()
	m.input, cmd = m.input.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Live filtering: if suggestions are showing and input changed, re-filter
	newValue := m.input.Value()
	if m.showingSuggestions {
		if newValue != oldValue {
			// Re-filter suggestions based on new input
			m.suggestions = m.completer.Complete(newValue)
			if len(m.suggestions) == 0 {
				m.showingSuggestions = false
			} else {
				m.suggestionIdx = 0
			}
		}
	}

	// Check if "/" was just typed at start of input or after space - trigger completion
	if newValue != oldValue && strings.HasSuffix(newValue, "/") {
		beforeSlash := strings.TrimSuffix(newValue, "/")
		if beforeSlash == "" || strings.HasSuffix(beforeSlash, " ") {
			if !m.showingSuggestions {
				m.triggerCompletion()
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the full UI
func (m *TUIModel) View() string {
	if m.quitting {
		return "\nGoodbye!\n"
	}

	if !m.ready {
		return "Initializing..."
	}

	var sb strings.Builder

	// Calculate available height for output
	// Layout: banner + output + input (3 lines) + statusbar (1 line)
	inputHeight := 3 // top rule + input + bottom rule
	statusHeight := 1
	bannerHeight := 0 // Banner is shown only once, not in the continuous view

	outputHeight := m.height - inputHeight - statusHeight - bannerHeight - 1 // -1 for safety

	// Calculate suggestion height if showing
	suggestionHeight := 0
	if m.showingSuggestions && len(m.suggestions) > 0 {
		suggestionHeight = min(len(m.suggestions), 10) + 1 // +1 for possible scroll indicator
	}

	// Adjust output height for suggestions
	adjustedOutputHeight := outputHeight - suggestionHeight

	// Render output area (scrollable, shows most recent)
	outputView := m.renderOutput(adjustedOutputHeight)
	sb.WriteString(outputView)

	// Render suggestions if showing
	if m.showingSuggestions {
		sb.WriteString(m.renderSuggestions())
	}

	// Render input box with horizontal rules
	sb.WriteString(m.input.View())
	sb.WriteString("\n")

	// Render status bar
	sb.WriteString(m.statusBar.Render())

	return sb.String()
}

// renderOutput renders the output area with proper height
func (m *TUIModel) renderOutput(height int) string {
	if height <= 0 {
		return ""
	}

	var lines []string

	// Get the most recent lines that fit
	startIdx := 0
	if len(m.output) > height {
		startIdx = len(m.output) - height
	}

	lines = m.output[startIdx:]

	// Pad with empty lines if needed
	for len(lines) < height {
		lines = append([]string{""}, lines...)
	}

	var sb strings.Builder
	for _, line := range lines {
		sb.WriteString(outputStyle.Render(line))
		sb.WriteString("\n")
	}

	return sb.String()
}

// addOutput adds a line to the output buffer
func (m *TUIModel) addOutput(line string) {
	// Split by newlines and add each
	m.output = append(m.output, strings.Split(line, "\n")...)

	// Keep output buffer from growing too large
	maxLines := 1000
	if len(m.output) > maxLines {
		m.output = m.output[len(m.output)-maxLines:]
	}
}

// executeAndCapture executes a command and captures its output
func (m *TUIModel) executeAndCapture(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Special handling for clear command
	if input == "clear" {
		m.output = make([]string, 0)
		return
	}

	// Special handling for quit command
	if input == "quit" || input == "exit" {
		ctx := m.session.GetContextPath()
		if input == "quit" || ctx.IsRoot() {
			m.quitting = true
			return
		}
	}

	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	// Execute command
	m.session.executeCommand(input)

	// Restore and read output
	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var outBuf, errBuf bytes.Buffer
	_, _ = io.Copy(&outBuf, rOut)
	_, _ = io.Copy(&errBuf, rErr)
	_ = rOut.Close()
	_ = rErr.Close()

	// Add captured output
	if outBuf.Len() > 0 {
		m.addOutput(strings.TrimSuffix(outBuf.String(), "\n"))
	}
	if errBuf.Len() > 0 {
		m.addOutput(strings.TrimSuffix(errBuf.String(), "\n"))
	}
}

// triggerCompletion gets suggestions and shows them or auto-completes
func (m *TUIModel) triggerCompletion() {
	m.suggestions = m.completer.Complete(m.input.Value())
	if len(m.suggestions) == 1 {
		// Single match - auto-complete
		m.applyCompletion(m.suggestions[0])
		m.suggestions = nil
		m.showingSuggestions = false
	} else if len(m.suggestions) > 0 {
		// Multiple matches - show suggestions
		m.showingSuggestions = true
		m.suggestionIdx = 0
	}
}

// applyCompletion applies a suggestion to the current input
func (m *TUIModel) applyCompletion(s Suggestion) {
	currentInput := m.input.Value()

	// Find the last word being completed
	words := strings.Fields(currentInput)
	if len(words) == 0 {
		// Empty input - just set the suggestion
		m.input.SetValue(s.Text + " ")
	} else {
		// Replace the last word with the suggestion
		lastWord := words[len(words)-1]
		if strings.HasPrefix(s.Text, lastWord) || strings.HasPrefix(strings.ToLower(s.Text), strings.ToLower(lastWord)) {
			// The suggestion matches the partial word - replace it
			prefix := strings.TrimSuffix(currentInput, lastWord)
			m.input.SetValue(prefix + s.Text + " ")
		} else {
			// Just append the suggestion
			m.input.SetValue(currentInput + s.Text + " ")
		}
	}
	m.input.CursorEnd()
}

// renderSuggestions renders the completion suggestions popup
func (m *TUIModel) renderSuggestions() string {
	if !m.showingSuggestions || len(m.suggestions) == 0 {
		return ""
	}

	var sb strings.Builder

	// Style for suggestions
	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#E4002B"))

	// Show up to 10 suggestions
	maxShow := 10
	start := 0
	if m.suggestionIdx >= maxShow {
		start = m.suggestionIdx - maxShow + 1
	}
	end := start + maxShow
	if end > len(m.suggestions) {
		end = len(m.suggestions)
	}

	for i := start; i < end; i++ {
		s := m.suggestions[i]
		line := s.Text
		if s.Description != "" {
			line += " - " + s.Description
		}

		// Truncate if too long
		if len(line) > m.width-4 {
			line = line[:m.width-7] + "..."
		}

		if i == m.suggestionIdx {
			sb.WriteString("  " + selectedStyle.Render(line))
		} else {
			sb.WriteString("  " + normalStyle.Render(line))
		}
		sb.WriteString("\n")
	}

	// Show scroll indicator if needed
	if len(m.suggestions) > maxShow {
		indicator := "  (" + strconv.Itoa(len(m.suggestions)) + " total, Tab/Shift+Tab to cycle)"
		sb.WriteString(normalStyle.Render(indicator))
		sb.WriteString("\n")
	}

	return sb.String()
}
