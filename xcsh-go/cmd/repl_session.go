package cmd

import (
	"fmt"
	"os"

	"github.com/robinmordasiewicz/xcsh/pkg/client"
)

// REPLSession holds state across the REPL lifetime
type REPLSession struct {
	client       *client.Client
	history      *HistoryManager
	completer    *CobraCompleter
	namespace    string
	lastExitCode int

	// Contextual navigation state
	contextPath *ContextPath      // Current navigation context
	tenant      string            // Extracted tenant name from API URL
	username    string            // Logged-in user's name/email
	validator   *ContextValidator // Domain/action validator
}

// initREPLSession creates a new REPL session with initialized state
func initREPLSession() (*REPLSession, error) {
	session := &REPLSession{
		namespace:   GetValidatedDefaultNamespace(),
		contextPath: &ContextPath{},
	}

	// Extract tenant from server URL if available
	if serverURL != "" {
		session.tenant = client.ExtractTenant(serverURL)
	}

	// Initialize context validator
	session.validator = NewContextValidator()

	// Initialize history manager
	historyPath := getHistoryFilePath()
	history, err := NewHistoryManager(historyPath, 1000)
	if err != nil {
		// Non-fatal: history just won't persist
		fmt.Fprintf(os.Stderr, "Warning: could not initialize history: %v\n", err)
		history = &HistoryManager{
			path:    historyPath,
			maxSize: 1000,
			history: make([]string, 0),
		}
	}
	session.history = history

	// Initialize completer with session reference for context-aware completions
	session.completer = NewCobraCompleter(rootCmd)
	session.completer.SetSession(session)

	// Initialize API client if credentials are available
	// This reuses the same initialization logic from PersistentPreRunE
	if serverURL != "" {
		cfg := &client.Config{
			ServerURL: serverURL,
			Cert:      cert,
			Key:       key,
			CACert:    cacert,
			P12Bundle: p12Bundle,
			Debug:     debug,
			Timeout:   timeout,
		}

		// Handle API token authentication
		if token := os.Getenv("F5XC_API_TOKEN"); token != "" {
			cfg.APIToken = token
		}

		var err error
		apiClient, err = client.New(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: API client initialization failed: %v\n", err)
			fmt.Fprintln(os.Stderr, "Some commands may not work until you configure credentials.")
		} else {
			session.client = apiClient
		}
		// Fetch username for prompt display
		session.username = getUserInfo()
	} else {
		fmt.Fprintln(os.Stderr, "Warning: No API URL configured.")
		fmt.Fprintln(os.Stderr, "Run 'configure' or set F5XC_API_URL environment variable.")
	}

	return session, nil
}

// SetNamespace updates the default namespace for the session
func (s *REPLSession) SetNamespace(ns string) {
	s.namespace = ns
}

// GetNamespace returns the current default namespace
func (s *REPLSession) GetNamespace() string {
	return s.namespace
}

// GetLastExitCode returns the exit code of the last command
func (s *REPLSession) GetLastExitCode() int {
	return s.lastExitCode
}

// SetLastExitCode sets the exit code of the last command
func (s *REPLSession) SetLastExitCode(code int) {
	s.lastExitCode = code
}

// GetContextPath returns the current navigation context
func (s *REPLSession) GetContextPath() *ContextPath {
	return s.contextPath
}

// GetTenant returns the current tenant name
func (s *REPLSession) GetTenant() string {
	return s.tenant
}

// GetUsername returns the logged-in user's name/email
func (s *REPLSession) GetUsername() string {
	return s.username
}

// GetValidator returns the context validator
func (s *REPLSession) GetValidator() *ContextValidator {
	return s.validator
}
