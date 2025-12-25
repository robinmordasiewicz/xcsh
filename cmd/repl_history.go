package cmd

import (
	"bufio"
	"os"
	"path/filepath"
)

// HistoryManager handles command history persistence
type HistoryManager struct {
	path    string
	maxSize int
	history []string
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(path string, maxSize int) (*HistoryManager, error) {
	h := &HistoryManager{
		path:    path,
		maxSize: maxSize,
		history: make([]string, 0),
	}

	// Load existing history
	if err := h.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return h, nil
}

// getHistoryFilePath returns the default history file path
func getHistoryFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".xcsh_history"
	}
	return filepath.Join(home, ".xcsh_history")
}

// Load reads history from the history file
func (h *HistoryManager) Load() error {
	file, err := os.Open(h.path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		h.history = append(h.history, scanner.Text())
	}

	// Trim to max size
	if len(h.history) > h.maxSize {
		h.history = h.history[len(h.history)-h.maxSize:]
	}

	return scanner.Err()
}

// Save writes history to the history file
func (h *HistoryManager) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(h.path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(h.path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	for _, cmd := range h.history {
		if _, err := file.WriteString(cmd + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Add adds a command to history
func (h *HistoryManager) Add(cmd string) {
	// Don't add empty commands or duplicates of the last command
	if cmd == "" {
		return
	}
	if len(h.history) > 0 && h.history[len(h.history)-1] == cmd {
		return
	}

	h.history = append(h.history, cmd)

	// Trim if necessary
	if len(h.history) > h.maxSize {
		h.history = h.history[1:]
	}
}

// GetHistory returns all history entries
func (h *HistoryManager) GetHistory() []string {
	return h.history
}
