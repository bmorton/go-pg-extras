package pgextras

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// CacheHit returns overall buffer cache hit ratios for indexes and tables.
func (c *Client) CacheHit(ctx context.Context) ([]CacheHitResult, error) {
	query, err := c.loadSQL("cache_hit", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[CacheHitResult](ctx, c.db, query)
}

// IndexCacheHit returns per-index buffer cache hit breakdown.
func (c *Client) IndexCacheHit(ctx context.Context) ([]IndexCacheHitResult, error) {
	query, err := c.loadSQL("index_cache_hit", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[IndexCacheHitResult](ctx, c.db, query)
}

// TableCacheHit returns per-table buffer cache hit breakdown.
func (c *Client) TableCacheHit(ctx context.Context) ([]TableCacheHitResult, error) {
	query, err := c.loadSQL("table_cache_hit", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableCacheHitResult](ctx, c.db, query)
}

// Outliers returns queries with longest execution time from pg_stat_statements.
func (c *Client) Outliers(ctx context.Context, params ...OutliersParams) ([]OutliersResult, error) {
	p := OutliersParams{Limit: 10}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	variant, err := c.selectSQLVariant(ctx, "outliers")
	if err != nil {
		return nil, err
	}
	query, err := c.loadSQL(variant, map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[OutliersResult](ctx, c.db, query)
}

// Calls returns queries with highest frequency of execution.
func (c *Client) Calls(ctx context.Context, params ...CallsParams) ([]CallsResult, error) {
	p := CallsParams{Limit: 10}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	variant, err := c.selectSQLVariant(ctx, "calls")
	if err != nil {
		return nil, err
	}
	query, err := c.loadSQL(variant, map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[CallsResult](ctx, c.db, query)
}

// LongRunningQueries returns currently running queries exceeding a duration threshold.
func (c *Client) LongRunningQueries(ctx context.Context, params ...LongRunningQueriesParams) ([]LongRunningQueriesResult, error) {
	p := LongRunningQueriesParams{Threshold: 500 * time.Millisecond}
	if len(params) > 0 && params[0].Threshold > 0 {
		p.Threshold = params[0].Threshold
	}
	query, err := c.loadSQL("long_running_queries", map[string]string{
		"threshold": formatPGInterval(p.Threshold),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[LongRunningQueriesResult](ctx, c.db, query)
}

// IndexUsage returns percentage of scans that used an index, per table.
func (c *Client) IndexUsage(ctx context.Context) ([]IndexUsageResult, error) {
	query, err := c.loadSQL("index_usage", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[IndexUsageResult](ctx, c.db, query)
}

// UnusedIndexes returns indexes with fewer than N scans on tables larger than 5 pages.
func (c *Client) UnusedIndexes(ctx context.Context, params ...UnusedIndexesParams) ([]UnusedIndexesResult, error) {
	p := UnusedIndexesParams{MaxScans: 50}
	if len(params) > 0 && params[0].MaxScans > 0 {
		p.MaxScans = params[0].MaxScans
	}
	query, err := c.loadSQL("unused_indexes", map[string]string{
		"max_scans": strconv.Itoa(p.MaxScans),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[UnusedIndexesResult](ctx, c.db, query)
}

// DuplicateIndexes returns indexes with identical column sets.
func (c *Client) DuplicateIndexes(ctx context.Context) ([]DuplicateIndexesResult, error) {
	query, err := c.loadSQL("duplicate_indexes", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[DuplicateIndexesResult](ctx, c.db, query)
}

// NullIndexes returns indexes containing significant NULL values.
func (c *Client) NullIndexes(ctx context.Context, params ...NullIndexesParams) ([]NullIndexesResult, error) {
	p := NullIndexesParams{MinRelationSizeMB: 10}
	if len(params) > 0 && params[0].MinRelationSizeMB > 0 {
		p.MinRelationSizeMB = params[0].MinRelationSizeMB
	}
	query, err := c.loadSQL("null_indexes", map[string]string{
		"min_relation_size_mb": strconv.Itoa(p.MinRelationSizeMB),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[NullIndexesResult](ctx, c.db, query)
}

// IndexSize returns the size of each index.
func (c *Client) IndexSize(ctx context.Context) ([]IndexSizeResult, error) {
	query, err := c.loadSQL("index_size", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[IndexSizeResult](ctx, c.db, query)
}

// TotalIndexSize returns the total size of all indexes combined.
func (c *Client) TotalIndexSize(ctx context.Context) ([]TotalIndexSizeResult, error) {
	query, err := c.loadSQL("total_index_size", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TotalIndexSizeResult](ctx, c.db, query)
}

// IndexScans returns the number of scans performed on indexes.
func (c *Client) IndexScans(ctx context.Context) ([]IndexScansResult, error) {
	query, err := c.loadSQL("index_scans", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[IndexScansResult](ctx, c.db, query)
}

// Indexes returns all indexes with their corresponding tables and columns.
func (c *Client) Indexes(ctx context.Context) ([]IndexesResult, error) {
	query, err := c.loadSQL("indexes", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[IndexesResult](ctx, c.db, query)
}

// TableSize returns the size of each table excluding indexes.
func (c *Client) TableSize(ctx context.Context) ([]TableSizeResult, error) {
	query, err := c.loadSQL("table_size", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableSizeResult](ctx, c.db, query)
}

// TotalTableSize returns total size per table including indexes and TOAST.
func (c *Client) TotalTableSize(ctx context.Context) ([]TotalTableSizeResult, error) {
	query, err := c.loadSQL("total_table_size", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TotalTableSizeResult](ctx, c.db, query)
}

// TableIndexesSize returns total index size per table.
func (c *Client) TableIndexesSize(ctx context.Context) ([]TableIndexesSizeResult, error) {
	query, err := c.loadSQL("table_indexes_size", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableIndexesSizeResult](ctx, c.db, query)
}

// RecordsRank returns estimated row counts per table from n_live_tup, descending.
func (c *Client) RecordsRank(ctx context.Context) ([]RecordsRankResult, error) {
	query, err := c.loadSQL("records_rank", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[RecordsRankResult](ctx, c.db, query)
}

// TableOverview returns combined size and row count data for all tables.
func (c *Client) TableOverview(ctx context.Context) ([]TableOverviewResult, error) {
	query, err := c.loadSQL("table_overview", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableOverviewResult](ctx, c.db, query)
}

// ScanActivity returns combined index and sequential scan counts per table.
func (c *Client) ScanActivity(ctx context.Context) ([]ScanActivityResult, error) {
	query, err := c.loadSQL("scan_activity", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[ScanActivityResult](ctx, c.db, query)
}

// Bloat returns table and index bloat estimation.
func (c *Client) Bloat(ctx context.Context) ([]BloatResult, error) {
	query, err := c.loadSQL("bloat", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[BloatResult](ctx, c.db, query)
}

// SeqScans returns sequential scan counts per table, descending.
func (c *Client) SeqScans(ctx context.Context) ([]SeqScansResult, error) {
	query, err := c.loadSQL("seq_scans", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[SeqScansResult](ctx, c.db, query)
}

// TableSchema returns column definitions for a specific table.
func (c *Client) TableSchema(ctx context.Context, params TableSchemaParams) ([]TableSchemaResult, error) {
	if params.TableName == "" {
		return nil, fmt.Errorf("%w: TableName", ErrRequiredParam)
	}
	if err := ValidateIdentifier(params.TableName); err != nil {
		return nil, err
	}
	query, err := c.loadSQL("table_schema", map[string]string{
		"table_name": params.TableName,
	})
	if err != nil {
		return nil, err
	}
	return runQuery[TableSchemaResult](ctx, c.db, query)
}

// TableSchemas returns column definitions for all tables.
func (c *Client) TableSchemas(ctx context.Context) ([]TableSchemasResult, error) {
	query, err := c.loadSQL("table_schemas", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableSchemasResult](ctx, c.db, query)
}

// TableForeignKeys returns FK constraints for a specific table.
func (c *Client) TableForeignKeys(ctx context.Context, params TableForeignKeysParams) ([]TableForeignKeysResult, error) {
	if params.TableName == "" {
		return nil, fmt.Errorf("%w: TableName", ErrRequiredParam)
	}
	if err := ValidateIdentifier(params.TableName); err != nil {
		return nil, err
	}
	query, err := c.loadSQL("table_foreign_keys", map[string]string{
		"table_name": params.TableName,
	})
	if err != nil {
		return nil, err
	}
	return runQuery[TableForeignKeysResult](ctx, c.db, query)
}

// ForeignKeys returns FK constraints for all tables.
func (c *Client) ForeignKeys(ctx context.Context) ([]ForeignKeysResult, error) {
	query, err := c.loadSQL("foreign_keys", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[ForeignKeysResult](ctx, c.db, query)
}

// TableIndexScans returns index scan counts per table.
func (c *Client) TableIndexScans(ctx context.Context) ([]TableIndexScansResult, error) {
	query, err := c.loadSQL("table_index_scans", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TableIndexScansResult](ctx, c.db, query)
}

// Tables returns all tables in the schema.
func (c *Client) Tables(ctx context.Context) ([]TablesResult, error) {
	query, err := c.loadSQL("tables", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[TablesResult](ctx, c.db, query)
}

// VacuumStats returns dead rows, autovacuum thresholds, and expect_autovacuum flag.
func (c *Client) VacuumStats(ctx context.Context) ([]VacuumStatsResult, error) {
	query, err := c.loadSQL("vacuum_stats", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[VacuumStatsResult](ctx, c.db, query)
}

// VacuumProgress returns current VACUUM/autovacuum progress.
func (c *Client) VacuumProgress(ctx context.Context) ([]VacuumProgressResult, error) {
	variant, err := c.selectSQLVariant(ctx, "vacuum_progress")
	if err != nil {
		return nil, err
	}
	query, err := c.loadSQL(variant, nil)
	if err != nil {
		return nil, err
	}
	return runQuery[VacuumProgressResult](ctx, c.db, query)
}

// AnalyzeProgress returns current ANALYZE progress.
func (c *Client) AnalyzeProgress(ctx context.Context) ([]AnalyzeProgressResult, error) {
	query, err := c.loadSQL("analyze_progress", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[AnalyzeProgressResult](ctx, c.db, query)
}

// VacuumIOStats returns autovacuum I/O statistics (PG 16+).
func (c *Client) VacuumIOStats(ctx context.Context) ([]VacuumIOStatsResult, error) {
	variant, err := c.selectSQLVariant(ctx, "vacuum_io_stats")
	if err != nil {
		return nil, err
	}
	// Legacy variant returns a different result type, but we handle it as a message.
	if variant == "vacuum_io_stats_legacy" {
		return nil, fmt.Errorf("%w: vacuum_io_stats requires PostgreSQL 16+", ErrUnsupportedPGVersion)
	}
	query, err := c.loadSQL(variant, nil)
	if err != nil {
		return nil, err
	}
	return runQuery[VacuumIOStatsResult](ctx, c.db, query)
}

// Connections returns all active connections to the current database.
func (c *Client) Connections(ctx context.Context) ([]ConnectionsResult, error) {
	query, err := c.loadSQL("connections", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[ConnectionsResult](ctx, c.db, query)
}

// Locks returns queries holding exclusive locks.
func (c *Client) Locks(ctx context.Context, params ...LocksParams) ([]LocksResult, error) {
	p := LocksParams{Limit: 20}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	query, err := c.loadSQL("locks", map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[LocksResult](ctx, c.db, query)
}

// AllLocks returns all current locks regardless of mode.
func (c *Client) AllLocks(ctx context.Context) ([]LocksResult, error) {
	query, err := c.loadSQL("all_locks", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[LocksResult](ctx, c.db, query)
}

// Blocking returns statements holding locks that block other statements.
func (c *Client) Blocking(ctx context.Context, params ...BlockingParams) ([]BlockingResult, error) {
	p := BlockingParams{Limit: 20}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	query, err := c.loadSQL("blocking", map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[BlockingResult](ctx, c.db, query)
}

// DBSettings returns selected PostgreSQL configuration values.
func (c *Client) DBSettings(ctx context.Context) ([]DBSettingsResult, error) {
	query, err := c.loadSQL("db_settings", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[DBSettingsResult](ctx, c.db, query)
}

// SSLUsed returns whether the current connection uses SSL.
func (c *Client) SSLUsed(ctx context.Context) ([]SSLUsedResult, error) {
	query, err := c.loadSQL("ssl_used", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[SSLUsedResult](ctx, c.db, query)
}

// Extensions returns installed PostgreSQL extensions.
func (c *Client) Extensions(ctx context.Context) ([]ExtensionsResult, error) {
	query, err := c.loadSQL("extensions", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[ExtensionsResult](ctx, c.db, query)
}

// BuffercacheStats returns relations buffered in shared memory.
func (c *Client) BuffercacheStats(ctx context.Context, params ...BuffercacheStatsParams) ([]BuffercacheStatsResult, error) {
	p := BuffercacheStatsParams{Limit: 20}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	query, err := c.loadSQL("buffercache_stats", map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[BuffercacheStatsResult](ctx, c.db, query)
}

// BuffercacheUsage returns cached block counts per relation.
func (c *Client) BuffercacheUsage(ctx context.Context, params ...BuffercacheUsageParams) ([]BuffercacheUsageResult, error) {
	p := BuffercacheUsageParams{Limit: 20}
	if len(params) > 0 && params[0].Limit > 0 {
		p.Limit = params[0].Limit
	}
	query, err := c.loadSQL("buffercache_usage", map[string]string{
		"limit": strconv.Itoa(p.Limit),
	})
	if err != nil {
		return nil, err
	}
	return runQuery[BuffercacheUsageResult](ctx, c.db, query)
}

// Mandelbrot renders the Mandelbrot set in ASCII art via SQL CTE.
func (c *Client) Mandelbrot(ctx context.Context) ([]MandelbrotResult, error) {
	query, err := c.loadSQL("mandelbrot", nil)
	if err != nil {
		return nil, err
	}
	return runQuery[MandelbrotResult](ctx, c.db, query)
}

// KillAll terminates all connections to the current database except self.
func (c *Client) KillAll(ctx context.Context) error {
	query, err := c.loadSQL("kill_all", nil)
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, query)
	return err
}

// KillPID terminates a specific connection.
func (c *Client) KillPID(ctx context.Context, pid int) error {
	query, err := c.loadSQL("kill_pid", map[string]string{
		"pid": strconv.Itoa(pid),
	})
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, query)
	return err
}

// PgStatStatementsReset resets all pg_stat_statements statistics.
func (c *Client) PgStatStatementsReset(ctx context.Context) error {
	query, err := c.loadSQL("pg_stat_statements_reset", nil)
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, query)
	return err
}

// AddExtensions creates sslinfo, pg_buffercache, and pg_stat_statements if not present.
func (c *Client) AddExtensions(ctx context.Context) error {
	query, err := c.loadSQL("add_extensions", nil)
	if err != nil {
		return err
	}
	_, err = c.db.ExecContext(ctx, query)
	return err
}
