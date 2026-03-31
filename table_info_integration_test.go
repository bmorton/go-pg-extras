//go:build integration

package pgextras

import (
	"context"
	"testing"
)

func TestIntegrationTableSchemasGrouped(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableSchemasGrouped(ctx)
	if err != nil {
		t.Fatalf("TableSchemasGrouped: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableSchemasGrouped returned no tables")
	}
	for _, table := range results {
		if table.TableName == "" {
			t.Error("found table with empty name")
		}
		if len(table.Columns) == 0 {
			t.Errorf("table %q has no columns", table.TableName)
		}
	}
}

func TestIntegrationTableInfo(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Get info for all tables.
	results, err := c.TableInfo(ctx, "")
	if err != nil {
		t.Fatalf("TableInfo: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableInfo returned no rows")
	}
}

func TestIntegrationTableInfoFiltered(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Get info for a specific table.
	results, err := c.TableInfo(ctx, "orders")
	if err != nil {
		t.Fatalf("TableInfo(orders): %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for orders, got %d", len(results))
	}
	if len(results) > 0 && results[0].TableName != "orders" {
		t.Errorf("expected table name 'orders', got %q", results[0].TableName)
	}
}
