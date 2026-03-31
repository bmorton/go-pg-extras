package pgextras

// CacheHitResult holds index and table hit rate data.
type CacheHitResult struct {
	Name  string `db:"name" json:"name"`
	Ratio string `db:"ratio" json:"ratio"`
}

// IndexCacheHitResult holds per-index buffer cache hit data.
type IndexCacheHitResult struct {
	Name       string `db:"name" json:"name"`
	BufferHits int64  `db:"buffer_hits" json:"buffer_hits"`
	BlockReads int64  `db:"block_reads" json:"block_reads"`
	TotalRead  int64  `db:"total_read" json:"total_read"`
	Ratio      string `db:"ratio" json:"ratio"`
}

// TableCacheHitResult holds per-table buffer cache hit data.
type TableCacheHitResult struct {
	Name       string `db:"name" json:"name"`
	BufferHits int64  `db:"buffer_hits" json:"buffer_hits"`
	BlockReads int64  `db:"block_reads" json:"block_reads"`
	TotalRead  int64  `db:"total_read" json:"total_read"`
	Ratio      string `db:"ratio" json:"ratio"`
}

// OutliersResult holds queries ordered by total execution time.
type OutliersResult struct {
	Query        string `db:"query" json:"query"`
	ExecTime     string `db:"total_exec_time" json:"exec_time"`
	PropExecTime string `db:"prop_exec_time" json:"prop_exec_time"`
	NCalls       string `db:"ncalls" json:"ncalls"`
	AvgExecMs    string `db:"avg_exec_ms" json:"avg_exec_ms"`
	SyncIOTime   string `db:"sync_io_time" json:"sync_io_time"`
}

// CallsResult holds queries ordered by call count.
type CallsResult struct {
	Query        string `db:"qry" json:"query"`
	ExecTime     string `db:"exec_time" json:"exec_time"`
	PropExecTime string `db:"prop_exec_time" json:"prop_exec_time"`
	NCalls       string `db:"ncalls" json:"ncalls"`
	AvgExecMs    string `db:"avg_exec_ms" json:"avg_exec_ms"`
	SyncIOTime   string `db:"sync_io_time" json:"sync_io_time"`
}

// LongRunningQueriesResult holds currently running long queries.
type LongRunningQueriesResult struct {
	PID      int    `db:"pid" json:"pid"`
	Duration string `db:"duration" json:"duration"`
	Query    string `db:"query" json:"query"`
}

// IndexUsageResult holds index usage data per table.
type IndexUsageResult struct {
	RelName                 string `db:"relname" json:"relname"`
	PercentOfTimesIndexUsed string `db:"percent_of_times_index_used" json:"percent_of_times_index_used"`
	RowsInTable             int64  `db:"rows_in_table" json:"rows_in_table"`
}

// UnusedIndexesResult holds unused index data.
type UnusedIndexesResult struct {
	Table      string `db:"table" json:"table"`
	Index      string `db:"index" json:"index"`
	IndexSize  string `db:"index_size" json:"index_size"`
	IndexScans int64  `db:"index_scans" json:"index_scans"`
}

// DuplicateIndexesResult holds duplicate index data.
type DuplicateIndexesResult struct {
	Size string  `db:"size" json:"size"`
	Idx1 string  `db:"idx1" json:"idx1"`
	Idx2 string  `db:"idx2" json:"idx2"`
	Idx3 *string `db:"idx3" json:"idx3"`
	Idx4 *string `db:"idx4" json:"idx4"`
}

// NullIndexesResult holds index null fraction data.
type NullIndexesResult struct {
	OID            int64  `db:"oid" json:"oid"`
	Index          string `db:"index" json:"index"`
	IndexSize      string `db:"index_size" json:"index_size"`
	Unique         bool   `db:"unique" json:"unique"`
	IndexedColumn  string `db:"indexed_column" json:"indexed_column"`
	Table          string `db:"table" json:"table"`
	NullFrac       string `db:"null_frac" json:"null_frac"`
	ExpectedSaving string `db:"expected_saving" json:"expected_saving"`
	Schema         string `db:"schema" json:"schema"`
}

// IndexSizeResult holds index size data.
type IndexSizeResult struct {
	Name   string `db:"name" json:"name"`
	Size   string `db:"size" json:"size"`
	Schema string `db:"schema" json:"schema"`
}

// TotalIndexSizeResult holds total index size.
type TotalIndexSizeResult struct {
	Size string `db:"size" json:"size"`
}

// IndexScansResult holds index scan data.
type IndexScansResult struct {
	SchemaName string `db:"schemaname" json:"schemaname"`
	Table      string `db:"table" json:"table"`
	Index      string `db:"index" json:"index"`
	IndexSize  string `db:"index_size" json:"index_size"`
	IndexScans int64  `db:"index_scans" json:"index_scans"`
}

// IndexesResult holds index metadata.
type IndexesResult struct {
	SchemaName string `db:"schemaname" json:"schemaname"`
	IndexName  string `db:"indexname" json:"indexname"`
	TableName  string `db:"tablename" json:"tablename"`
	Columns    string `db:"columns" json:"columns"`
}

// TableSizeResult holds table size data.
type TableSizeResult struct {
	Name   string `db:"name" json:"name"`
	Size   string `db:"size" json:"size"`
	Schema string `db:"schema" json:"schema"`
}

// TotalTableSizeResult holds total table size including indexes.
type TotalTableSizeResult struct {
	Name string `db:"name" json:"name"`
	Size string `db:"size" json:"size"`
}

// TableIndexesSizeResult holds total index size per table.
type TableIndexesSizeResult struct {
	Table     string `db:"table" json:"table"`
	IndexSize string `db:"index_size" json:"index_size"`
}

// RecordsRankResult holds estimated row counts per table.
type RecordsRankResult struct {
	Name           string `db:"name" json:"name"`
	EstimatedCount int64  `db:"estimated_count" json:"estimated_count"`
}

// BloatResult holds bloat estimation per table/index.
type BloatResult struct {
	Type       string `db:"type" json:"type"`
	SchemaName string `db:"schemaname" json:"schemaname"`
	ObjectName string `db:"object_name" json:"object_name"`
	Bloat      string `db:"bloat" json:"bloat"`
	Waste      string `db:"waste" json:"waste"`
}

// SeqScansResult holds sequential scan counts per table.
type SeqScansResult struct {
	Name  string `db:"name" json:"name"`
	Count int64  `db:"count" json:"count"`
}

// TableSchemaResult holds column definitions for a table.
type TableSchemaResult struct {
	ColumnName    string  `db:"column_name" json:"column_name"`
	DataType      string  `db:"data_type" json:"data_type"`
	IsNullable    string  `db:"is_nullable" json:"is_nullable"`
	ColumnDefault *string `db:"column_default" json:"column_default"`
}

// TableSchemasResult holds column definitions for all tables.
type TableSchemasResult struct {
	TableName     string  `db:"table_name" json:"table_name"`
	ColumnName    string  `db:"column_name" json:"column_name"`
	DataType      string  `db:"data_type" json:"data_type"`
	IsNullable    string  `db:"is_nullable" json:"is_nullable"`
	ColumnDefault *string `db:"column_default" json:"column_default"`
}

// TableForeignKeysResult holds foreign key data for a table.
type TableForeignKeysResult struct {
	TableName         string `db:"table_name" json:"table_name"`
	ConstraintName    string `db:"constraint_name" json:"constraint_name"`
	ColumnName        string `db:"column_name" json:"column_name"`
	ForeignTableName  string `db:"foreign_table_name" json:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name" json:"foreign_column_name"`
}

// ForeignKeysResult holds foreign key data for all tables.
type ForeignKeysResult struct {
	TableName         string `db:"table_name" json:"table_name"`
	ConstraintName    string `db:"constraint_name" json:"constraint_name"`
	ColumnName        string `db:"column_name" json:"column_name"`
	ForeignTableName  string `db:"foreign_table_name" json:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name" json:"foreign_column_name"`
}

// TableIndexScansResult holds index scan counts per table.
type TableIndexScansResult struct {
	Name  string `db:"name" json:"name"`
	Count int64  `db:"count" json:"count"`
}

// TablesResult holds table listing data.
type TablesResult struct {
	TableName  string `db:"tablename" json:"tablename"`
	SchemaName string `db:"schemaname" json:"schemaname"`
}

// VacuumStatsResult holds vacuum statistics.
type VacuumStatsResult struct {
	Schema                     string  `db:"schema" json:"schema"`
	Table                      string  `db:"table" json:"table"`
	LastManualVacuum           *string `db:"last_manual_vacuum" json:"last_manual_vacuum"`
	ManualVacuumCount          string  `db:"manual_vacuum_count" json:"manual_vacuum_count"`
	LastAutovacuum             *string `db:"last_autovacuum" json:"last_autovacuum"`
	AutovacuumCount            string  `db:"autovacuum_count" json:"autovacuum_count"`
	RowCount                   string  `db:"rowcount" json:"rowcount"`
	DeadRowCount               string  `db:"dead_rowcount" json:"dead_rowcount"`
	DeadTupAutovacuumThreshold *string `db:"dead_tup_autovacuum_threshold" json:"dead_tup_autovacuum_threshold"`
	NInsSinceVacuum            *string `db:"n_ins_since_vacuum" json:"n_ins_since_vacuum"`
	InsertAutovacuumThreshold  *string `db:"insert_autovacuum_threshold" json:"insert_autovacuum_threshold"`
	ExpectAutovacuum           *string `db:"expect_autovacuum" json:"expect_autovacuum"`
}

// VacuumProgressResult holds current VACUUM progress data.
type VacuumProgressResult struct {
	Database         string `db:"database" json:"database"`
	Schema           string `db:"schema" json:"schema"`
	Table            string `db:"table" json:"table"`
	PID              int    `db:"pid" json:"pid"`
	Phase            string `db:"phase" json:"phase"`
	HeapBlksTotal    int64  `db:"heap_blks_total" json:"heap_blks_total"`
	HeapBlksScanned  int64  `db:"heap_blks_scanned" json:"heap_blks_scanned"`
	HeapBlksVacuumed int64  `db:"heap_blks_vacuumed" json:"heap_blks_vacuumed"`
	IndexVacuumCount int64  `db:"index_vacuum_count" json:"index_vacuum_count"`
}

// AnalyzeProgressResult holds current ANALYZE progress data.
type AnalyzeProgressResult struct {
	Database          string `db:"database" json:"database"`
	Schema            string `db:"schema" json:"schema"`
	Table             string `db:"table" json:"table"`
	PID               int    `db:"pid" json:"pid"`
	Phase             string `db:"phase" json:"phase"`
	SampleBlksTotal   int64  `db:"sample_blks_total" json:"sample_blks_total"`
	SampleBlksScanned int64  `db:"sample_blks_scanned" json:"sample_blks_scanned"`
	ExtStatsTotal     int64  `db:"ext_stats_total" json:"ext_stats_total"`
	ExtStatsComputed  int64  `db:"ext_stats_computed" json:"ext_stats_computed"`
}

// VacuumIOStatsResult holds autovacuum I/O statistics (PG 16+).
type VacuumIOStatsResult struct {
	BackendType   string  `db:"backend_type" json:"backend_type"`
	Object        string  `db:"object" json:"object"`
	Context       string  `db:"context" json:"context"`
	Reads         int64   `db:"reads" json:"reads"`
	ReadTime      string  `db:"read_time" json:"read_time"`
	Writes        int64   `db:"writes" json:"writes"`
	WriteTime     string  `db:"write_time" json:"write_time"`
	Writebacks    int64   `db:"writebacks" json:"writebacks"`
	WritebackTime string  `db:"writeback_time" json:"writeback_time"`
	Extends       int64   `db:"extends" json:"extends"`
	ExtendTime    string  `db:"extend_time" json:"extend_time"`
	Fsyncs        int64   `db:"fsyncs" json:"fsyncs"`
	FsyncTime     string  `db:"fsync_time" json:"fsync_time"`
	Reuses        int64   `db:"reuses" json:"reuses"`
	Evictions     int64   `db:"evictions" json:"evictions"`
	StatsReset    *string `db:"stats_reset" json:"stats_reset"`
}

// VacuumIOStatsLegacyResult holds the message for PG < 16.
type VacuumIOStatsLegacyResult struct {
	FeatureNotAvailable string `db:"feature not available" json:"feature_not_available"`
}

// ConnectionsResult holds active connection data.
type ConnectionsResult struct {
	Username        string  `db:"username" json:"username"`
	PID             int     `db:"pid" json:"pid"`
	ClientAddress   *string `db:"client_address" json:"client_address"`
	ApplicationName string  `db:"application_name" json:"application_name"`
}

// LocksResult holds lock data.
type LocksResult struct {
	PID           int     `db:"pid" json:"pid"`
	RelName       *string `db:"relname" json:"relname"`
	TransactionID *string `db:"transactionid" json:"transactionid"`
	LockType      string  `db:"locktype" json:"locktype"`
	Database      *string `db:"database" json:"database"`
	Granted       bool    `db:"granted" json:"granted"`
	Mode          string  `db:"mode" json:"mode"`
	QuerySnippet  string  `db:"query_snippet" json:"query_snippet"`
	Age           string  `db:"age" json:"age"`
	Application   string  `db:"application" json:"application"`
}

// BlockingResult holds blocking query data.
type BlockingResult struct {
	BlockedPID        int    `db:"blocked_pid" json:"blocked_pid"`
	BlockingStatement string `db:"blocking_statement" json:"blocking_statement"`
	BlockingDuration  string `db:"blocking_duration" json:"blocking_duration"`
	BlockingPID       int    `db:"blocking_pid" json:"blocking_pid"`
	BlockedStatement  string `db:"blocked_statement" json:"blocked_statement"`
	BlockedDuration   string `db:"blocked_duration" json:"blocked_duration"`
	BlockedSQLApp     string `db:"blocked_sql_app" json:"blocked_sql_app"`
	BlockingSQLApp    string `db:"blocking_sql_app" json:"blocking_sql_app"`
}

// DBSettingsResult holds PostgreSQL configuration settings.
type DBSettingsResult struct {
	Name      string  `db:"name" json:"name"`
	Setting   string  `db:"setting" json:"setting"`
	Unit      *string `db:"unit" json:"unit"`
	ShortDesc string  `db:"short_desc" json:"short_desc"`
}

// SSLUsedResult holds SSL status.
type SSLUsedResult struct {
	SSLIsUsed bool `db:"ssl_is_used" json:"ssl_is_used"`
}

// ExtensionsResult holds installed extension data.
type ExtensionsResult struct {
	Name             string  `db:"name" json:"name"`
	DefaultVersion   string  `db:"default_version" json:"default_version"`
	InstalledVersion *string `db:"installed_version" json:"installed_version"`
	Comment          *string `db:"comment" json:"comment"`
}

// BuffercacheStatsResult holds buffer cache statistics.
type BuffercacheStatsResult struct {
	RelName           string `db:"relname" json:"relname"`
	Buffered          string `db:"buffered" json:"buffered"`
	BufferPercent     string `db:"buffer_percent" json:"buffer_percent"`
	PercentOfRelation string `db:"percent_of_relation" json:"percent_of_relation"`
}

// BuffercacheUsageResult holds buffer cache usage data.
type BuffercacheUsageResult struct {
	RelName string `db:"relname" json:"relname"`
	Buffers int64  `db:"buffers" json:"buffers"`
}

// MandelbrotResult holds the Mandelbrot set ASCII art.
type MandelbrotResult struct {
	Mandelbrot string `db:"array_to_string" json:"mandelbrot"`
}
