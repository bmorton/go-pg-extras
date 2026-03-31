package pgextras

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	sqlembed "github.com/bmorton/go-pg-extras/sql"
)

var (
	// ErrExtensionNotInstalled is returned when a query requires an extension
	// that is not installed in the database.
	ErrExtensionNotInstalled = errors.New("required PostgreSQL extension is not installed")

	// ErrRequiredParam is returned when a required parameter is missing.
	ErrRequiredParam = errors.New("required parameter is missing")

	// ErrUnsupportedPGVersion is returned when a query requires a newer PG version.
	ErrUnsupportedPGVersion = errors.New("query requires a newer PostgreSQL version")
)

// Config holds connection and behavior settings.
type Config struct {
	// DB is an open *sql.DB handle (required).
	DB *sql.DB

	// Schema defaults to "public". Overridden by PG_EXTRAS_SCHEMA env var.
	Schema string
}

// Client is the primary entrypoint for running pg-extras queries.
type Client struct {
	db     *sql.DB
	schema string
	versionCache
}

// New creates a Client from an existing *sql.DB.
// The caller owns the *sql.DB lifecycle — Client never closes it.
func New(cfg Config) (*Client, error) {
	if cfg.DB == nil {
		return nil, errors.New("pgextras: Config.DB is required")
	}
	schema := cfg.Schema
	if schema == "" {
		schema = os.Getenv("PG_EXTRAS_SCHEMA")
	}
	if schema == "" {
		schema = "public"
	}
	return &Client{
		db:     cfg.DB,
		schema: schema,
	}, nil
}

// DB exposes the underlying *sql.DB for advanced use.
func (c *Client) DB() *sql.DB {
	return c.db
}

// Schema returns the configured schema.
func (c *Client) Schema() string {
	return c.schema
}

// loadSQL reads an embedded SQL file and substitutes parameters.
func (c *Client) loadSQL(name string, params map[string]string) (string, error) {
	data, err := sqlembed.Files.ReadFile(name + ".sql")
	if err != nil {
		return "", fmt.Errorf("pgextras: loading SQL %q: %w", name, err)
	}
	query := string(data)

	// Substitute %{key} placeholders with values.
	if params == nil {
		params = make(map[string]string)
	}
	if _, ok := params["schema"]; !ok {
		params["schema"] = c.schema
	}
	for k, v := range params {
		query = strings.ReplaceAll(query, "%{"+k+"}", v)
	}
	return query, nil
}

// runQuery executes SQL and scans results into a slice of T using db tags.
func runQuery[T any](ctx context.Context, db *sql.DB, query string) ([]T, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []T
	for rows.Next() {
		var item T
		dest := mapColumnsToFields(cols, &item)
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("pgextras: scanning row: %w", err)
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

// mapColumnsToFields maps SQL column names to struct fields via db tags.
func mapColumnsToFields(cols []string, dest any) []any {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	// Build tag -> field index map.
	tagMap := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag != "" {
			tagMap[tag] = i
		}
	}

	scanDest := make([]any, len(cols))
	for i, col := range cols {
		if fieldIdx, ok := tagMap[col]; ok {
			scanDest[i] = v.Field(fieldIdx).Addr().Interface()
		} else {
			// Discard unknown columns.
			scanDest[i] = new(any)
		}
	}
	return scanDest
}
