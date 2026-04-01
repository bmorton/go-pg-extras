package pgextras

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// testRow is a simple struct used for formatting tests.
type testRow struct {
	Name  string `db:"name" json:"name"`
	Value int    `db:"value" json:"value"`
}

// testRowWithPtr has a pointer field to test nil handling.
type testRowWithPtr struct {
	Name  string  `db:"name" json:"name"`
	Alias *string `db:"alias" json:"alias"`
}

func TestFormatResultsJSON(t *testing.T) {
	rows := []testRow{
		{Name: "alpha", Value: 1},
		{Name: "beta", Value: 2},
	}

	var buf bytes.Buffer
	if err := FormatResults(&buf, rows, FormatJSON, ""); err != nil {
		t.Fatalf("FormatResults JSON: %v", err)
	}

	var decoded []testRow
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON decode: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(decoded))
	}
	if decoded[0].Name != "alpha" || decoded[0].Value != 1 {
		t.Errorf("row 0: got %+v", decoded[0])
	}
	if decoded[1].Name != "beta" || decoded[1].Value != 2 {
		t.Errorf("row 1: got %+v", decoded[1])
	}
}

func TestFormatResultsCSV(t *testing.T) {
	rows := []testRow{
		{Name: "alpha", Value: 10},
		{Name: "beta", Value: 20},
	}

	var buf bytes.Buffer
	if err := FormatResults(&buf, rows, FormatCSV, ""); err != nil {
		t.Fatalf("FormatResults CSV: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 CSV lines (header + 2 rows), got %d: %v", len(lines), lines)
	}
	if lines[0] != "name,value" {
		t.Errorf("CSV header: got %q, want %q", lines[0], "name,value")
	}
	if lines[1] != "alpha,10" {
		t.Errorf("CSV row 1: got %q, want %q", lines[1], "alpha,10")
	}
}

func TestFormatResultsTable(t *testing.T) {
	rows := []testRow{
		{Name: "alpha", Value: 1},
	}

	var buf bytes.Buffer
	if err := FormatResults(&buf, rows, FormatTable, "test title"); err != nil {
		t.Fatalf("FormatResults Table: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("table output missing 'alpha': %s", out)
	}
	if !strings.Contains(out, "name") {
		t.Errorf("table output missing header 'name': %s", out)
	}
	if !strings.Contains(out, "test title") {
		t.Errorf("table output missing title: %s", out)
	}
}

func TestFormatResultsUnsupportedFormat(t *testing.T) {
	rows := []testRow{{Name: "a", Value: 1}}
	var buf bytes.Buffer
	err := FormatResults(&buf, rows, Format(99), "")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("error should mention unsupported format: %v", err)
	}
}

func TestFormatResultsEmptySlice(t *testing.T) {
	var rows []testRow
	var buf bytes.Buffer
	if err := FormatResults(&buf, rows, FormatJSON, ""); err != nil {
		t.Fatalf("FormatResults JSON empty: %v", err)
	}
}

func TestFieldToStringNilPointer(t *testing.T) {
	alias := "bob"
	rows := []testRowWithPtr{
		{Name: "alice", Alias: nil},
		{Name: "carol", Alias: &alias},
	}

	var buf bytes.Buffer
	if err := FormatResults(&buf, rows, FormatCSV, ""); err != nil {
		t.Fatalf("FormatResults CSV with ptr: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// header + 2 data rows
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	// nil pointer should produce empty string
	if !strings.Contains(lines[1], "alice,") {
		t.Errorf("row with nil pointer: got %q", lines[1])
	}
	if !strings.Contains(lines[2], "carol,bob") {
		t.Errorf("row with non-nil pointer: got %q", lines[2])
	}
}

func TestExtractHeadersAndRows(t *testing.T) {
	rows := []testRow{{Name: "x", Value: 42}}
	headers, data := extractHeadersAndRows(rows)

	if len(headers) != 2 || headers[0] != "name" || headers[1] != "value" {
		t.Errorf("headers: got %v", headers)
	}
	if len(data) != 1 || data[0][0] != "x" || data[0][1] != "42" {
		t.Errorf("rows: got %v", data)
	}
}
