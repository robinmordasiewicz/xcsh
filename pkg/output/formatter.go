package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format represents the output format type
type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatTable Format = "table"
	FormatTSV   Format = "tsv"
	FormatNone  Format = "none"
)

// Formatter handles output formatting
type Formatter struct {
	format Format
	writer io.Writer
}

// New creates a new formatter
func New(format string) *Formatter {
	f := &Formatter{
		writer: os.Stdout,
	}

	switch strings.ToLower(format) {
	case "json":
		f.format = FormatJSON
	case "yaml", "":
		f.format = FormatYAML
	case "table":
		f.format = FormatTable
	case "tsv":
		f.format = FormatTSV
	case "none":
		f.format = FormatNone
	default:
		f.format = FormatYAML
	}

	return f
}

// SetWriter sets the output writer
func (f *Formatter) SetWriter(w io.Writer) {
	f.writer = w
}

// Format formats and outputs the data
func (f *Formatter) Format(data interface{}) error {
	if f.format == FormatNone {
		return nil
	}

	switch f.format {
	case FormatJSON:
		return f.formatJSON(data)
	case FormatYAML:
		return f.formatYAML(data)
	case FormatTable:
		return f.formatTable(data)
	case FormatTSV:
		return f.formatTSV(data)
	default:
		return f.formatYAML(data)
	}
}

// formatJSON outputs data as pretty-printed JSON
func (f *Formatter) formatJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// formatYAML outputs data as YAML
func (f *Formatter) formatYAML(data interface{}) error {
	encoder := yaml.NewEncoder(f.writer)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

// formatTable outputs data as an ASCII table
func (f *Formatter) formatTable(data interface{}) error {
	rows, headers := extractTableData(data)
	if len(rows) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	separators := make([]string, len(headers))
	for i, h := range headers {
		separators[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	// Print rows
	for _, row := range rows {
		values := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := row[h]; ok {
				values[i] = formatValue(v)
			}
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}

	return w.Flush()
}

// formatTSV outputs data as tab-separated values (no headers)
func (f *Formatter) formatTSV(data interface{}) error {
	rows, headers := extractTableData(data)
	if len(rows) == 0 {
		return nil
	}

	for _, row := range rows {
		values := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := row[h]; ok {
				values[i] = formatValue(v)
			}
		}
		fmt.Fprintln(f.writer, strings.Join(values, "\t"))
	}

	return nil
}

// extractTableData extracts rows and headers from data
func extractTableData(data interface{}) ([]map[string]interface{}, []string) {
	var rows []map[string]interface{}
	headerSet := make(map[string]bool)

	v := reflect.ValueOf(data)

	// Handle pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle slice/array
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			row := toMap(item)
			if row != nil {
				rows = append(rows, row)
				for k := range row {
					headerSet[k] = true
				}
			}
		}
	} else {
		// Single item
		row := toMap(data)
		if row != nil {
			rows = append(rows, row)
			for k := range row {
				headerSet[k] = true
			}
		}
	}

	// Sort headers alphabetically
	var headers []string
	for h := range headerSet {
		headers = append(headers, h)
	}
	sort.Strings(headers)

	// Prioritize common fields
	priorityOrder := []string{"name", "namespace", "status", "created", "modified"}
	headers = prioritizeHeaders(headers, priorityOrder)

	return rows, headers
}

// toMap converts an interface to a map
func toMap(data interface{}) map[string]interface{} {
	// If already a map
	if m, ok := data.(map[string]interface{}); ok {
		return flattenMap(m, "")
	}

	// Try JSON round-trip for structs
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}

	return flattenMap(m, "")
}

// flattenMap flattens nested maps with dot notation
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			// Flatten nested maps
			for fk, fv := range flattenMap(val, key) {
				result[fk] = fv
			}
		case []interface{}:
			// For arrays, just show count or first few items
			if len(val) == 0 {
				result[key] = "[]"
			} else {
				result[key] = fmt.Sprintf("[%d items]", len(val))
			}
		default:
			result[key] = v
		}
	}

	return result
}

// formatValue formats a value for table output
func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%.2f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// prioritizeHeaders moves priority fields to the front
func prioritizeHeaders(headers []string, priority []string) []string {
	headerSet := make(map[string]bool)
	for _, h := range headers {
		headerSet[h] = true
	}

	var result []string
	for _, p := range priority {
		if headerSet[p] {
			result = append(result, p)
			delete(headerSet, p)
		}
	}

	// Add remaining headers
	var remaining []string
	for h := range headerSet {
		remaining = append(remaining, h)
	}
	sort.Strings(remaining)
	result = append(result, remaining...)

	return result
}

// Print is a convenience function to format and print data
func Print(data interface{}, format string) error {
	return New(format).Format(data)
}

// PrintError prints an error message to stderr
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
}

// PrintWarning prints a warning message to stderr
func PrintWarning(msg string) {
	fmt.Fprintf(os.Stderr, "WARNING: %s\n", msg)
}

// PrintInfo prints an info message to stderr
func PrintInfo(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

// PrintSuccess prints a success message with a checkmark
func PrintSuccess(msg string) {
	fmt.Fprintf(os.Stderr, "✓ %s\n", msg)
}

// PrintDebug prints a debug message (only when verbose)
func PrintDebug(msg string, verbose bool) {
	if verbose {
		fmt.Fprintf(os.Stderr, "DEBUG: %s\n", msg)
	}
}

// PrintAPIError formats and prints an API error with helpful context
func PrintAPIError(statusCode int, body []byte, operation string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s failed (HTTP %d)\n", operation, statusCode)

	// Try to parse error body for more details
	var errResp map[string]interface{}
	if err := json.Unmarshal(body, &errResp); err == nil {
		if msg, ok := errResp["message"].(string); ok {
			fmt.Fprintf(os.Stderr, "  Message: %s\n", msg)
		}
		if code, ok := errResp["code"].(string); ok {
			fmt.Fprintf(os.Stderr, "  Code: %s\n", code)
		}
		if details, ok := errResp["details"].(string); ok {
			fmt.Fprintf(os.Stderr, "  Details: %s\n", details)
		}
	} else if len(body) > 0 && len(body) < 500 {
		fmt.Fprintf(os.Stderr, "  Response: %s\n", string(body))
	}

	// Provide helpful hints based on status code
	switch statusCode {
	case 401:
		fmt.Fprintf(os.Stderr, "\nHint: Authentication failed. Check your credentials with 'f5xc configure show'\n")
	case 403:
		fmt.Fprintf(os.Stderr, "\nHint: Permission denied. You may not have access to this resource.\n")
	case 404:
		fmt.Fprintf(os.Stderr, "\nHint: Resource not found. Verify the name and namespace are correct.\n")
	case 409:
		fmt.Fprintf(os.Stderr, "\nHint: Conflict - resource may already exist or be in a conflicting state.\n")
	case 429:
		fmt.Fprintf(os.Stderr, "\nHint: Rate limited. Please wait and try again.\n")
	case 500, 502, 503:
		fmt.Fprintf(os.Stderr, "\nHint: Server error. Please try again later or contact support.\n")
	}
}

// Spinner represents a simple progress indicator
type Spinner struct {
	message string
	done    chan bool
}

// StartSpinner starts a progress indicator
func StartSpinner(message string) *Spinner {
	s := &Spinner{
		message: message,
		done:    make(chan bool),
	}
	go s.run()
	return s
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.done <- true
	close(s.done)
	fmt.Fprintf(os.Stderr, "\r%s... done\n", s.message)
}

// StopWithError stops the spinner with an error indication
func (s *Spinner) StopWithError() {
	s.done <- true
	close(s.done)
	fmt.Fprintf(os.Stderr, "\r%s... failed\n", s.message)
}

func (s *Spinner) run() {
	chars := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-s.done:
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s... %s", s.message, chars[i%len(chars)])
			i++
			// Sleep for a short time (we'll use a simple approach without time import)
		}
	}
}
