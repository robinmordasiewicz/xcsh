package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	yamlv2 "gopkg.in/yaml.v2"
)

// Format represents the output format type
type Format string

const (
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatTable Format = "table"
	FormatText  Format = "text" // Alias for table (human-readable)
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
	case "yaml":
		f.format = FormatYAML
	case "table", "text", "":
		f.format = FormatTable // Default to table for compatibility with original f5xcctl
	case "tsv":
		f.format = FormatTSV
	case "none":
		f.format = FormatNone
	default:
		f.format = FormatTable
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

// formatYAML outputs data as YAML (matching original f5xcctl format)
func (f *Formatter) formatYAML(data interface{}) error {
	// Use yaml.v2 Marshal for compatibility with original f5xcctl output format
	// yaml.v2 uses 2-space indent and doesn't indent array items under parent keys
	out, err := yamlv2.Marshal(data)
	if err != nil {
		return err
	}
	_, err = f.writer.Write(out)
	if err != nil {
		return err
	}
	// Add trailing newline to match original f5xcctl format
	_, err = f.writer.Write([]byte("\n"))
	return err
}

// formatTable outputs data as an ASCII box table matching original f5xcctl format
func (f *Formatter) formatTable(data interface{}) error {
	// Extract items from list response
	items := extractItems(data)
	if len(items) == 0 {
		return nil
	}

	// Define columns for namespace list (matching original f5xcctl)
	// Original format: NAMESPACE | NAME | LABELS
	headers := []string{"NAMESPACE", "NAME", "LABELS"}

	// Fixed column widths matching original f5xcctl output
	// Original f5xcctl uses fixed widths: NAMESPACE=9, NAME=27, LABELS=30
	fixedWidths := []int{9, 27, 30}

	// Calculate actual column widths (max of fixed width and header/content width)
	widths := make([]int, len(headers))
	for i, h := range headers {
		if len(h) > fixedWidths[i] {
			widths[i] = len(h)
		} else {
			widths[i] = fixedWidths[i]
		}
	}

	// Build rows - each item may produce multiple display rows if content wraps
	var displayRows [][][]string // each item can have multiple lines
	for _, item := range items {
		row := make([]string, len(headers))

		// NAMESPACE column - namespace field (empty for namespaces themselves)
		ns := getStringField(item, "namespace")
		if ns == "" {
			row[0] = "<None>"
		} else {
			row[0] = ns
		}

		// NAME column
		row[1] = getStringField(item, "name")
		if row[1] == "" {
			row[1] = "<None>"
		}

		// LABELS column
		labels := getLabelsString(item)
		if labels == "" {
			row[2] = "<None>"
		} else {
			row[2] = labels
		}

		// Wrap row into multiple lines if needed
		wrappedRows := wrapRow(row, widths)
		displayRows = append(displayRows, wrappedRows)
	}

	// Print ASCII box table
	printBoxLine(f.writer, widths)
	printBoxRowCentered(f.writer, headers, widths) // Headers centered
	printBoxLine(f.writer, widths)
	for _, wrappedRows := range displayRows {
		for _, row := range wrappedRows {
			printBoxRowLeft(f.writer, row, widths) // Data left-aligned
		}
		printBoxLine(f.writer, widths)
	}

	return nil
}

// wrapRow wraps a row into multiple lines based on column widths
func wrapRow(row []string, widths []int) [][]string {
	var result [][]string

	// Split each cell into lines based on max width
	cellLines := make([][]string, len(row))
	maxLines := 1

	for i, cell := range row {
		cellLines[i] = wrapText(cell, widths[i])
		if len(cellLines[i]) > maxLines {
			maxLines = len(cellLines[i])
		}
	}

	// Create rows for each line
	for lineNum := 0; lineNum < maxLines; lineNum++ {
		newRow := make([]string, len(row))
		for colNum := 0; colNum < len(row); colNum++ {
			if lineNum < len(cellLines[colNum]) {
				newRow[colNum] = cellLines[colNum][lineNum]
			} else {
				newRow[colNum] = "" // Empty for continuation lines
			}
		}
		result = append(result, newRow)
	}

	return result
}

// wrapText wraps text to fit within maxWidth characters
func wrapText(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}

	var lines []string
	remaining := text

	for len(remaining) > 0 {
		if len(remaining) <= maxWidth {
			lines = append(lines, remaining)
			break
		}

		// Find a good break point (preferring space)
		breakPoint := maxWidth
		for i := maxWidth - 1; i > 0; i-- {
			if remaining[i] == ' ' {
				breakPoint = i
				break
			}
		}

		lines = append(lines, remaining[:breakPoint])
		remaining = remaining[breakPoint:]
		// Trim leading space from next line
		for len(remaining) > 0 && remaining[0] == ' ' {
			remaining = remaining[1:]
		}
	}

	return lines
}

// extractItems extracts the items array from a list response
func extractItems(data interface{}) []map[string]interface{} {
	// Handle map with "items" key
	if m, ok := data.(map[string]interface{}); ok {
		if items, ok := m["items"]; ok {
			if itemSlice, ok := items.([]interface{}); ok {
				var result []map[string]interface{}
				for _, item := range itemSlice {
					if itemMap, ok := item.(map[string]interface{}); ok {
						result = append(result, itemMap)
					}
				}
				return result
			}
		}
		// Single item, wrap it
		return []map[string]interface{}{m}
	}

	// Handle slice directly
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		var result []map[string]interface{}
		for i := 0; i < v.Len(); i++ {
			if itemMap, ok := v.Index(i).Interface().(map[string]interface{}); ok {
				result = append(result, itemMap)
			}
		}
		return result
	}

	return nil
}

// getStringField gets a string field from a map
func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getLabelsString formats labels as map[key:value key:value] (matching original f5xcctl format)
func getLabelsString(m map[string]interface{}) string {
	labels, ok := m["labels"]
	if !ok {
		return ""
	}

	labelMap, ok := labels.(map[string]interface{})
	if !ok || len(labelMap) == 0 {
		return ""
	}

	var parts []string
	for k, v := range labelMap {
		parts = append(parts, fmt.Sprintf("%s:%v", k, v))
	}
	sort.Strings(parts)
	return "map[" + strings.Join(parts, " ") + "]"
}

// printBoxLine prints a horizontal line: +-----+-----+
func printBoxLine(w io.Writer, widths []int) {
	_, _ = fmt.Fprint(w, "+")
	for _, width := range widths {
		_, _ = fmt.Fprint(w, strings.Repeat("-", width+2))
		_, _ = fmt.Fprint(w, "+")
	}
	_, _ = fmt.Fprintln(w)
}

// printBoxRowCentered prints a row with centered content (for headers)
func printBoxRowCentered(w io.Writer, cells []string, widths []int) {
	_, _ = fmt.Fprint(w, "|")
	for i, cell := range cells {
		// Center the cell content
		padding := widths[i] - len(cell)
		leftPad := padding / 2
		rightPad := padding - leftPad
		_, _ = fmt.Fprintf(w, " %s%s%s |", strings.Repeat(" ", leftPad), cell, strings.Repeat(" ", rightPad))
	}
	_, _ = fmt.Fprintln(w)
}

// printBoxRowLeft prints a row with left-aligned content (for data)
func printBoxRowLeft(w io.Writer, cells []string, widths []int) {
	_, _ = fmt.Fprint(w, "|")
	for i, cell := range cells {
		// Left-align the cell content
		padding := widths[i] - len(cell)
		_, _ = fmt.Fprintf(w, " %s%s |", cell, strings.Repeat(" ", padding))
	}
	_, _ = fmt.Fprintln(w)
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
		_, _ = fmt.Fprintln(f.writer, strings.Join(values, "\t"))
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
		fmt.Fprintf(os.Stderr, "\nHint: Authentication failed. Check your credentials with 'f5xcctl configure show'\n")
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
