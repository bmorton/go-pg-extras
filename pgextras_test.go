package pgextras

import (
	"database/sql"
	"os"
	"strings"
	"testing"
)

func TestNewNilDB(t *testing.T) {
	_, err := New(Config{DB: nil})
	if err == nil {
		t.Fatal("expected error when DB is nil")
	}
	if !strings.Contains(err.Error(), "Config.DB is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewDefaultSchema(t *testing.T) {
	// Ensure env var doesn't interfere.
	os.Unsetenv("PG_EXTRAS_SCHEMA")

	// We need a non-nil *sql.DB but won't use it for queries.
	db := &sql.DB{}
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Schema() != "public" {
		t.Errorf("Schema() = %q, want %q", c.Schema(), "public")
	}
	if c.DB() != db {
		t.Error("DB() should return the same *sql.DB")
	}
}

func TestNewExplicitSchema(t *testing.T) {
	os.Unsetenv("PG_EXTRAS_SCHEMA")

	db := &sql.DB{}
	c, err := New(Config{DB: db, Schema: "custom"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Schema() != "custom" {
		t.Errorf("Schema() = %q, want %q", c.Schema(), "custom")
	}
}

func TestNewSchemaFromEnv(t *testing.T) {
	t.Setenv("PG_EXTRAS_SCHEMA", "envschema")

	db := &sql.DB{}
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Schema() != "envschema" {
		t.Errorf("Schema() = %q, want %q", c.Schema(), "envschema")
	}
}

func TestNewSchemaExplicitOverridesEnv(t *testing.T) {
	t.Setenv("PG_EXTRAS_SCHEMA", "envschema")

	db := &sql.DB{}
	c, err := New(Config{DB: db, Schema: "explicit"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Schema() != "explicit" {
		t.Errorf("Schema() = %q, want %q", c.Schema(), "explicit")
	}
}

func TestLoadSQL(t *testing.T) {
	db := &sql.DB{}
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// cache_hit.sql should exist in embedded files.
	query, err := c.loadSQL("cache_hit", nil)
	if err != nil {
		t.Fatalf("loadSQL: %v", err)
	}

	// Schema placeholder should have been replaced.
	if strings.Contains(query, "%{schema}") {
		t.Error("query still contains %{schema} placeholder")
	}
}

func TestLoadSQLWithParams(t *testing.T) {
	db := &sql.DB{}
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	query, err := c.loadSQL("unused_indexes", map[string]string{
		"max_scans": "100",
	})
	if err != nil {
		t.Fatalf("loadSQL: %v", err)
	}
	if strings.Contains(query, "%{max_scans}") {
		t.Error("query still contains %{max_scans} placeholder")
	}
}

func TestLoadSQLMissingFile(t *testing.T) {
	db := &sql.DB{}
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = c.loadSQL("nonexistent_query", nil)
	if err == nil {
		t.Fatal("expected error for missing SQL file")
	}
}

func TestMapColumnsToFields(t *testing.T) {
	type sample struct {
		Name  string `db:"name"`
		Value int    `db:"value"`
	}
	s := sample{}
	cols := []string{"name", "value", "unknown_col"}
	dest := mapColumnsToFields(cols, &s)

	if len(dest) != 3 {
		t.Fatalf("expected 3 scan destinations, got %d", len(dest))
	}

	// First two should point to struct fields, third to a discard target.
	// Verify the types are addressable.
	if dest[0] == nil || dest[1] == nil || dest[2] == nil {
		t.Error("scan destinations should not be nil")
	}
}

func TestMapColumnsToFieldsNoMatch(t *testing.T) {
	type sample struct {
		Name string `db:"name"`
	}
	s := sample{}
	cols := []string{"no_match_a", "no_match_b"}
	dest := mapColumnsToFields(cols, &s)

	if len(dest) != 2 {
		t.Fatalf("expected 2 scan destinations, got %d", len(dest))
	}
}
