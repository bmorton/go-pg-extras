//go:build integration

package pgextras

import (
	"context"
	"testing"
)

func TestIntegrationCacheHit(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.CacheHit(ctx)
	if err != nil {
		t.Fatalf("CacheHit: %v", err)
	}
	if len(results) == 0 {
		t.Log("CacheHit returned no rows (may be expected on fresh database)")
	}
}

func TestIntegrationIndexCacheHit(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.IndexCacheHit(ctx)
	if err != nil {
		t.Fatalf("IndexCacheHit: %v", err)
	}
	_ = results
}

func TestIntegrationTableCacheHit(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableCacheHit(ctx)
	if err != nil {
		t.Fatalf("TableCacheHit: %v", err)
	}
	_ = results
}

func TestIntegrationIndexUsage(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.IndexUsage(ctx)
	if err != nil {
		t.Fatalf("IndexUsage: %v", err)
	}
	_ = results
}

func TestIntegrationUnusedIndexes(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.UnusedIndexes(ctx)
	if err != nil {
		t.Fatalf("UnusedIndexes: %v", err)
	}
	_ = results
}

func TestIntegrationDuplicateIndexes(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.DuplicateIndexes(ctx)
	if err != nil {
		t.Fatalf("DuplicateIndexes: %v", err)
	}
	_ = results
}

func TestIntegrationBloat(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Bloat(ctx)
	if err != nil {
		t.Fatalf("Bloat: %v", err)
	}
	_ = results
}

func TestIntegrationTableSize(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableSize(ctx)
	if err != nil {
		t.Fatalf("TableSize: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableSize returned no rows, expected at least some tables")
	}
}

func TestIntegrationTotalTableSize(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TotalTableSize(ctx)
	if err != nil {
		t.Fatalf("TotalTableSize: %v", err)
	}
	_ = results
}

func TestIntegrationIndexSize(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.IndexSize(ctx)
	if err != nil {
		t.Fatalf("IndexSize: %v", err)
	}
	_ = results
}

func TestIntegrationTotalIndexSize(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TotalIndexSize(ctx)
	if err != nil {
		t.Fatalf("TotalIndexSize: %v", err)
	}
	_ = results
}

func TestIntegrationIndexScans(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.IndexScans(ctx)
	if err != nil {
		t.Fatalf("IndexScans: %v", err)
	}
	_ = results
}

func TestIntegrationIndexes(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Indexes(ctx)
	if err != nil {
		t.Fatalf("Indexes: %v", err)
	}
	if len(results) == 0 {
		t.Error("Indexes returned no rows, expected at least some indexes")
	}
}

func TestIntegrationRecordsRank(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.RecordsRank(ctx)
	if err != nil {
		t.Fatalf("RecordsRank: %v", err)
	}
	_ = results
}

func TestIntegrationSeqScans(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.SeqScans(ctx)
	if err != nil {
		t.Fatalf("SeqScans: %v", err)
	}
	_ = results
}

func TestIntegrationTableOverview(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableOverview(ctx)
	if err != nil {
		t.Fatalf("TableOverview: %v", err)
	}
	_ = results
}

func TestIntegrationScanActivity(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.ScanActivity(ctx)
	if err != nil {
		t.Fatalf("ScanActivity: %v", err)
	}
	_ = results
}

func TestIntegrationConnections(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Connections(ctx)
	if err != nil {
		t.Fatalf("Connections: %v", err)
	}
	if len(results) == 0 {
		t.Error("Connections returned no rows, expected at least the test connection")
	}
}

func TestIntegrationLocks(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Locks(ctx)
	if err != nil {
		t.Fatalf("Locks: %v", err)
	}
	_ = results
}

func TestIntegrationAllLocks(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.AllLocks(ctx)
	if err != nil {
		t.Fatalf("AllLocks: %v", err)
	}
	_ = results
}

func TestIntegrationBlocking(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Blocking(ctx)
	if err != nil {
		t.Fatalf("Blocking: %v", err)
	}
	_ = results
}

func TestIntegrationDBSettings(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.DBSettings(ctx)
	if err != nil {
		t.Fatalf("DBSettings: %v", err)
	}
	if len(results) == 0 {
		t.Error("DBSettings returned no rows")
	}
}

func TestIntegrationExtensions(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Extensions(ctx)
	if err != nil {
		t.Fatalf("Extensions: %v", err)
	}
	if len(results) == 0 {
		t.Error("Extensions returned no rows")
	}
}

func TestIntegrationVacuumStats(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.VacuumStats(ctx)
	if err != nil {
		t.Fatalf("VacuumStats: %v", err)
	}
	_ = results
}

func TestIntegrationVacuumProgress(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.VacuumProgress(ctx)
	if err != nil {
		t.Fatalf("VacuumProgress: %v", err)
	}
	_ = results
}

func TestIntegrationAnalyzeProgress(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.AnalyzeProgress(ctx)
	if err != nil {
		t.Fatalf("AnalyzeProgress: %v", err)
	}
	_ = results
}

func TestIntegrationLongRunningQueries(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.LongRunningQueries(ctx)
	if err != nil {
		t.Fatalf("LongRunningQueries: %v", err)
	}
	_ = results
}

func TestIntegrationTables(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Tables(ctx)
	if err != nil {
		t.Fatalf("Tables: %v", err)
	}
	if len(results) == 0 {
		t.Error("Tables returned no rows")
	}
}

func TestIntegrationTableSchemas(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableSchemas(ctx)
	if err != nil {
		t.Fatalf("TableSchemas: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableSchemas returned no rows")
	}
}

func TestIntegrationForeignKeys(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.ForeignKeys(ctx)
	if err != nil {
		t.Fatalf("ForeignKeys: %v", err)
	}
	if len(results) == 0 {
		t.Error("ForeignKeys returned no rows")
	}
}

func TestIntegrationTableSchema(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableSchema(ctx, TableSchemaParams{TableName: "orders"})
	if err != nil {
		t.Fatalf("TableSchema: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableSchema returned no rows for 'orders' table")
	}
}

func TestIntegrationTableSchemaRequiresParam(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	_, err := c.TableSchema(ctx, TableSchemaParams{})
	if err == nil {
		t.Fatal("expected error when TableName is empty")
	}
}

func TestIntegrationTableForeignKeys(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableForeignKeys(ctx, TableForeignKeysParams{TableName: "orders"})
	if err != nil {
		t.Fatalf("TableForeignKeys: %v", err)
	}
	if len(results) == 0 {
		t.Error("TableForeignKeys returned no rows for 'orders' table")
	}
}

func TestIntegrationTableIndexesSize(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableIndexesSize(ctx)
	if err != nil {
		t.Fatalf("TableIndexesSize: %v", err)
	}
	_ = results
}

func TestIntegrationTableIndexScans(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.TableIndexScans(ctx)
	if err != nil {
		t.Fatalf("TableIndexScans: %v", err)
	}
	_ = results
}

func TestIntegrationNullIndexes(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Use MinRelationSizeMB: 0 to catch any index
	results, err := c.NullIndexes(ctx, NullIndexesParams{MinRelationSizeMB: 0})
	if err != nil {
		t.Fatalf("NullIndexes: %v", err)
	}
	_ = results
}

func TestIntegrationMandelbrot(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Mandelbrot(ctx)
	if err != nil {
		t.Fatalf("Mandelbrot: %v", err)
	}
	if len(results) == 0 {
		t.Error("Mandelbrot returned no rows")
	}
}

func TestIntegrationOutliers(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Outliers requires pg_stat_statements; skip if unavailable.
	results, err := c.Outliers(ctx)
	if err != nil {
		t.Logf("Outliers: %v (may require pg_stat_statements)", err)
		return
	}
	_ = results
}

func TestIntegrationCalls(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Calls requires pg_stat_statements; skip if unavailable.
	results, err := c.Calls(ctx)
	if err != nil {
		t.Logf("Calls: %v (may require pg_stat_statements)", err)
		return
	}
	_ = results
}

func TestIntegrationSSLUsed(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.SSLUsed(ctx)
	if err != nil {
		t.Logf("SSLUsed: %v (may require sslinfo extension)", err)
		return
	}
	_ = results
}
