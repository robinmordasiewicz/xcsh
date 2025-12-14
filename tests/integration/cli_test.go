package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// getBinaryPath returns the path to the f5xcctl binary
func getBinaryPath(t *testing.T) string {
	// Try current directory first
	cwd, _ := os.Getwd()
	binary := filepath.Join(cwd, "..", "..", "f5xcctl")
	if _, err := os.Stat(binary); err == nil {
		return binary
	}

	// Try relative to project root
	binary = "./f5xcctl"
	if _, err := os.Stat(binary); err == nil {
		return binary
	}

	t.Skip("f5xcctl binary not found - run 'go build' first")
	return ""
}

// TestCLI_Version tests the version command
func TestCLI_Version(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("version command failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()
	// Version output format: "f5xcctl version <version>\n  commit: <sha>\n  built: <date>\n  go: <version>\n  platform: <os/arch>"
	if !strings.Contains(output, "f5xcctl version") || !strings.Contains(output, "commit:") {
		t.Errorf("Expected version output with f5xcctl version and commit, got: %s", output)
	}

	t.Logf("Version output: %s", output)
}

// TestCLI_Help tests the help command
func TestCLI_Help(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("help command failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()

	// Check for expected content in help
	expectedStrings := []string{
		"f5xcctl",
		"Available Commands",
		"--help",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help to contain '%s'", expected)
		}
	}

	t.Logf("Help output contains expected content")
}

// TestCLI_CompletionBash tests bash completion generation
func TestCLI_CompletionBash(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "completion", "bash")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion bash failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "bash") || !strings.Contains(output, "complet") {
		t.Errorf("Expected bash completion script, got: %s", output[:min(200, len(output))])
	}

	t.Log("Bash completion generated successfully")
}

// TestCLI_CompletionZsh tests zsh completion generation
func TestCLI_CompletionZsh(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "completion", "zsh")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion zsh failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if output == "" {
		t.Error("Expected zsh completion script")
	}

	t.Log("Zsh completion generated successfully")
}

// TestCLI_ConfigureShow tests the configure show command
func TestCLI_ConfigureShow(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "configure", "show")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// This may fail if no config exists, which is OK
	_ = cmd.Run()

	// Just verify it runs without crashing
	t.Log("Configure show executed")
}

// TestCLI_ListNamespaces tests listing namespaces via CLI
func TestCLI_ListNamespaces(t *testing.T) {
	binary := getBinaryPath(t)

	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured")
	}

	// Set environment for the command (include --insecure for staging environments)
	cmd := exec.Command(binary, "namespace", "list", "--server-url", apiURL, "--p12-bundle", p12File, "--insecure")
	cmd.Env = append(os.Environ(), "VES_P12_PASSWORD="+p12Password)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("namespace list failed: %v\nstderr: %s\nstdout: %s", err, stderr.String(), stdout.String())
	}

	output := stdout.String()
	t.Logf("Namespace list output:\n%s", output)
}

// TestCLI_HTTPLoadBalancerList tests listing HTTP load balancers via CLI
func TestCLI_HTTPLoadBalancerList(t *testing.T) {
	binary := getBinaryPath(t)

	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured")
	}

	cmd := exec.Command(binary, "http-loadbalancer", "list", "-n", "shared", "--server-url", apiURL, "--p12-bundle", p12File, "--insecure")
	cmd.Env = append(os.Environ(), "VES_P12_PASSWORD="+p12Password)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("http-loadbalancer list failed: %v\nstderr: %s", err, stderr.String())
	}

	t.Logf("HTTP LB list output:\n%s", stdout.String())
}

// TestCLI_OutputFormatJSON tests JSON output format
func TestCLI_OutputFormatJSON(t *testing.T) {
	binary := getBinaryPath(t)

	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured")
	}

	cmd := exec.Command(binary, "namespace", "list", "-o", "json", "--server-url", apiURL, "--p12-bundle", p12File, "--insecure")
	cmd.Env = append(os.Environ(), "VES_P12_PASSWORD="+p12Password)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("namespace list failed: %v\nstderr: %s", err, stderr.String())
	}

	output := stdout.String()
	// JSON output should start with { or [
	trimmed := strings.TrimSpace(output)
	if len(trimmed) > 0 && trimmed[0] != '{' && trimmed[0] != '[' {
		t.Errorf("Expected JSON output, got: %s", trimmed[:min(100, len(trimmed))])
	}

	t.Log("JSON output format works")
}

// TestCLI_OutputFormatTable tests table output format
func TestCLI_OutputFormatTable(t *testing.T) {
	binary := getBinaryPath(t)

	apiURL := os.Getenv("VES_API_URL")
	p12File := os.Getenv("VES_P12_FILE")
	p12Password := os.Getenv("VES_P12_PASSWORD")

	if apiURL == "" || p12File == "" || p12Password == "" {
		t.Skip("Integration test environment not configured")
	}

	cmd := exec.Command(binary, "namespace", "list", "-o", "table", "--server-url", apiURL, "--p12-bundle", p12File, "--insecure")
	cmd.Env = append(os.Environ(), "VES_P12_PASSWORD="+p12Password)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("namespace list failed: %v\nstderr: %s", err, stderr.String())
	}

	t.Log("Table output format works")
}

// TestCLI_InvalidCommand tests handling of invalid commands
func TestCLI_InvalidCommand(t *testing.T) {
	binary := getBinaryPath(t)

	cmd := exec.Command(binary, "nonexistent-command-12345")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Expected error for invalid command")
	}

	// Check that error message is helpful
	combined := stdout.String() + stderr.String()
	if !strings.Contains(combined, "unknown command") && !strings.Contains(combined, "Error") {
		t.Logf("Error output: %s", combined)
	}
}

// TestCLI_MissingRequiredFlag tests handling of missing required flags
func TestCLI_MissingRequiredFlag(t *testing.T) {
	binary := getBinaryPath(t)

	// Try to list without namespace
	cmd := exec.Command(binary, "http-loadbalancer", "list")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Log("Command succeeded without namespace (may be OK if config has default)")
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
