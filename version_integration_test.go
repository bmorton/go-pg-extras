//go:build integration

package pgextras

import (
	"context"
	"testing"
)

func TestIntegrationPgVersion(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	version, err := c.pgVersion(ctx)
	if err != nil {
		t.Fatalf("pgVersion: %v", err)
	}
	if version < 12 || version > 30 {
		t.Errorf("pgVersion returned unexpected value: %d", version)
	}
	t.Logf("PostgreSQL major version: %d", version)
}

func TestIntegrationPgStatStatementsVersion(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	version, err := c.pgStatStatementsVersion(ctx)
	if err != nil {
		t.Fatalf("pgStatStatementsVersion: %v", err)
	}
	// May be empty if not installed.
	t.Logf("pg_stat_statements version: %q", version)
}

func TestIntegrationSelectSQLVariantDefault(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// A query with no special variants should return the base name.
	variant, err := c.selectSQLVariant(ctx, "cache_hit")
	if err != nil {
		t.Fatalf("selectSQLVariant(cache_hit): %v", err)
	}
	if variant != "cache_hit" {
		t.Errorf("expected %q, got %q", "cache_hit", variant)
	}
}

func TestIntegrationSelectSQLVariantOutliers(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	variant, err := c.selectSQLVariant(ctx, "outliers")
	if err != nil {
		// May fail if pg_stat_statements is not installed.
		t.Logf("selectSQLVariant(outliers): %v", err)
		return
	}
	// Should be one of: "outliers", "outliers_legacy", or "outliers_17"
	valid := map[string]bool{"outliers": true, "outliers_legacy": true, "outliers_17": true}
	if !valid[variant] {
		t.Errorf("unexpected variant: %q", variant)
	}
	t.Logf("outliers variant: %s", variant)
}

func TestIntegrationSelectSQLVariantVacuumProgress(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	variant, err := c.selectSQLVariant(ctx, "vacuum_progress")
	if err != nil {
		t.Fatalf("selectSQLVariant(vacuum_progress): %v", err)
	}
	valid := map[string]bool{"vacuum_progress": true, "vacuum_progress_17": true}
	if !valid[variant] {
		t.Errorf("unexpected variant: %q", variant)
	}
	t.Logf("vacuum_progress variant: %s", variant)
}
