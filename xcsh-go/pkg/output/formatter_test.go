package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_JSONFormat(t *testing.T) {
	f := New("json")
	if f.format != FormatJSON {
		t.Errorf("Expected JSON format, got %v", f.format)
	}
}

func TestNew_YAMLFormat(t *testing.T) {
	f := New("yaml")
	if f.format != FormatYAML {
		t.Errorf("Expected YAML format, got %v", f.format)
	}
}

func TestNew_TableFormat(t *testing.T) {
	f := New("table")
	if f.format != FormatTable {
		t.Errorf("Expected Table format, got %v", f.format)
	}
}

func TestNew_TSVFormat(t *testing.T) {
	f := New("tsv")
	if f.format != FormatTSV {
		t.Errorf("Expected TSV format, got %v", f.format)
	}
}

func TestNew_NoneFormat(t *testing.T) {
	f := New("none")
	if f.format != FormatNone {
		t.Errorf("Expected None format, got %v", f.format)
	}
}

func TestNew_EmptyDefault(t *testing.T) {
	f := New("")
	if f.format != FormatTable {
		t.Errorf("Expected Table format as default (matching original xcsh), got %v", f.format)
	}
}

func TestNew_InvalidDefault(t *testing.T) {
	f := New("invalid")
	if f.format != FormatTable {
		t.Errorf("Expected Table format for invalid input (matching original xcsh), got %v", f.format)
	}
}

func TestNew_CaseInsensitive(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"JSON", FormatJSON},
		{"Json", FormatJSON},
		{"YAML", FormatYAML},
		{"Yaml", FormatYAML},
		{"TABLE", FormatTable},
		{"Table", FormatTable},
	}

	for _, tc := range cases {
		f := New(tc.input)
		if f.format != tc.expected {
			t.Errorf("For input %q, expected %v, got %v", tc.input, tc.expected, f.format)
		}
	}
}

func TestFormatter_SetWriter(t *testing.T) {
	f := New("json")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := map[string]string{"test": "value"}
	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected output in buffer")
	}
}

func TestFormatter_FormatJSON(t *testing.T) {
	f := New("json")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name 'test', got '%v'", result["name"])
	}
}

func TestFormatter_FormatYAML(t *testing.T) {
	f := New("yaml")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "name: test") {
		t.Errorf("Expected YAML output to contain 'name: test', got: %s", output)
	}
}

func TestFormatter_FormatTable(t *testing.T) {
	f := New("table")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := []map[string]interface{}{
		{"name": "resource1", "namespace": "ns1"},
		{"name": "resource2", "namespace": "ns2"},
	}

	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	// Table format uses box-style headers: NAMESPACE | NAME | LABELS
	if !strings.Contains(output, "NAME") {
		t.Errorf("Expected table header 'NAME', got: %s", output)
	}
	if !strings.Contains(output, "resource1") {
		t.Errorf("Expected row 'resource1', got: %s", output)
	}
}

func TestFormatter_FormatTSV(t *testing.T) {
	f := New("tsv")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := []map[string]interface{}{
		{"name": "resource1", "namespace": "ns1"},
	}

	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\t") {
		t.Errorf("Expected tab-separated output, got: %s", output)
	}
}

func TestFormatter_FormatNone(t *testing.T) {
	f := New("none")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := map[string]string{"test": "value"}
	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if buf.Len() != 0 {
		t.Error("Expected no output for none format")
	}
}

func TestFormatter_FormatEmptySlice(t *testing.T) {
	f := New("table")
	buf := &bytes.Buffer{}
	f.SetWriter(buf)

	data := []map[string]interface{}{}
	if err := f.Format(data); err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	// Empty slice should produce no output
	if buf.Len() != 0 {
		t.Errorf("Expected no output for empty slice, got: %s", buf.String())
	}
}

func TestPrint(t *testing.T) {
	// This is a convenience function that writes to stdout
	// Just verify it doesn't panic
	data := map[string]string{"test": "value"}
	err := Print(data, "none") // Use none to suppress actual output
	if err != nil {
		t.Errorf("Print failed: %v", err)
	}
}

func TestFormatValue(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{nil, ""},
		{"test", "test"},
		{true, "true"},
		{false, "false"},
		{float64(123), "123"},
		{float64(123.45), "123.45"},
		{[]string{"a", "b"}, "[a b]"},
	}

	for _, tc := range cases {
		result := formatValue(tc.input)
		if result != tc.expected {
			t.Errorf("For input %v, expected %q, got %q", tc.input, tc.expected, result)
		}
	}
}

func TestExtractTableData_SingleItem(t *testing.T) {
	data := map[string]interface{}{
		"name":      "test",
		"namespace": "default",
	}

	rows, headers := extractTableData(data)

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	if len(headers) < 2 {
		t.Errorf("Expected at least 2 headers, got %d", len(headers))
	}
}

func TestExtractTableData_Slice(t *testing.T) {
	data := []map[string]interface{}{
		{"name": "resource1", "namespace": "ns1"},
		{"name": "resource2", "namespace": "ns2"},
	}

	rows, headers := extractTableData(data)

	if len(rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(rows))
	}

	// Check headers contain expected fields
	headerSet := make(map[string]bool)
	for _, h := range headers {
		headerSet[h] = true
	}

	if !headerSet["name"] {
		t.Error("Expected 'name' header")
	}
	if !headerSet["namespace"] {
		t.Error("Expected 'namespace' header")
	}
}

func TestExtractTableData_NestedMap(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"metadata": map[string]interface{}{
			"namespace": "default",
			"labels": map[string]interface{}{
				"app": "test-app",
			},
		},
	}

	rows, headers := extractTableData(data)

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	// Nested maps should be flattened
	headerSet := make(map[string]bool)
	for _, h := range headers {
		headerSet[h] = true
	}

	if !headerSet["name"] {
		t.Error("Expected 'name' header")
	}
	if !headerSet["metadata.namespace"] {
		t.Error("Expected 'metadata.namespace' header (flattened)")
	}
}

func TestPrioritizeHeaders(t *testing.T) {
	headers := []string{"created", "status", "modified", "name", "namespace", "id"}
	priority := []string{"name", "namespace", "status", "created", "modified"}

	result := prioritizeHeaders(headers, priority)

	// name should be first if present
	if result[0] != "name" {
		t.Errorf("Expected 'name' first, got '%s'", result[0])
	}

	// namespace should be second
	if result[1] != "namespace" {
		t.Errorf("Expected 'namespace' second, got '%s'", result[1])
	}

	// All original headers should be present
	if len(result) != len(headers) {
		t.Errorf("Expected %d headers, got %d", len(headers), len(result))
	}
}

func TestFlattenMap_Empty(t *testing.T) {
	input := map[string]interface{}{}
	result := flattenMap(input, "")

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestFlattenMap_Array(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{"a", "b", "c"},
	}

	result := flattenMap(input, "")

	// Arrays should be shown as count
	if result["items"] != "[3 items]" {
		t.Errorf("Expected '[3 items]', got '%v'", result["items"])
	}
}

func TestFlattenMap_EmptyArray(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{},
	}

	result := flattenMap(input, "")

	if result["items"] != "[]" {
		t.Errorf("Expected '[]', got '%v'", result["items"])
	}
}
