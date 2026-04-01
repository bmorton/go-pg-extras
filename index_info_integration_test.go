//go:build integration

package pgextras

import (
	"context"
	"testing"
)

func TestIntegrationIndexInfo(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Get info for all indexes.
	results, err := c.IndexInfo(ctx, "")
	if err != nil {
		t.Fatalf("IndexInfo: %v", err)
	}
	if len(results) == 0 {
		t.Error("IndexInfo returned no rows")
	}
}

func TestIntegrationIndexInfoFiltered(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Get info for indexes on orders table.
	results, err := c.IndexInfo(ctx, "orders")
	if err != nil {
		t.Fatalf("IndexInfo(orders): %v", err)
	}
	if len(results) == 0 {
		t.Error("IndexInfo returned no rows for 'orders' table")
	}
	for _, idx := range results {
		if idx.TableName != "orders" {
			t.Errorf("expected table name 'orders', got %q", idx.TableName)
		}
	}
}
