# go-pg-extras — Complete Package Specification

> A Go port of [rails-pg-extras](https://github.com/pawurb/rails-pg-extras) providing PostgreSQL database performance insights: locks, index usage, buffer cache hit ratios, vacuum stats, and more.

---

## 1. Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| PostgreSQL driver | `database/sql` (stdlib) | Zero third-party driver dependency in the library itself; users register any `database/sql`-compatible PG driver (`lib/pq`, `pgx/v5/stdlib`, etc.) in their own code. The CLI binary imports `lib/pq` as its default driver. |
| Driver registration | Caller's responsibility | The library accepts `*sql.DB`; it never calls `sql.Open()` or imports a driver. This follows the stdlib pattern where the driver `import _ "github.com/lib/pq"` lives in `main`. |
| SQL storage | `embed.FS` | Go 1.16+, zero-cost at runtime, no file I/O |
| HTTP layer | `net/http` handlers | Framework-agnostic; composable with chi, mux, or stdlib |
| Distribution | Single importable module | `go get` only; `cmd/pgextras` binary within the same module |
| Result types | Typed structs per query | Compile-time safety; mirrors the Rust port's approach |
| Context propagation | `context.Context` first param | Standard Go convention for cancellation and deadlines |

---

## 2. Package Structure

```
go-pg-extras/
├── go.mod                       // module github.com/<org>/go-pg-extras
├── go.sum
├── pgextras.go                  // Client, New(), Config (accepts *sql.DB)
├── queries.go                   // Public methods: CacheHit(), Outliers(), etc.
├── queries_params.go            // Param structs: OutliersParams, LongRunningQueriesParams, etc.
├── types.go                     // Result structs: CacheHitResult, OutliersResult, etc.
├── diagnose.go                  // Diagnose() with threshold evaluation
├── diagnose_types.go            // DiagnoseResult, DiagnoseCheck, Severity
├── table_info.go                // TableInfo() aggregation
├── index_info.go                // IndexInfo() aggregation
├── format.go                    // FormatTable (ASCII), FormatJSON, FormatCSV
├── version.go                   // PG version detection, pg_stat_statements routing
├── sql/                         // Embedded .sql files
│   ├── embed.go                 // //go:embed *.sql
│   ├── cache_hit.sql
│   ├── bloat.sql
│   ├── outliers.sql
│   ├── outliers_legacy.sql
│   └── ... (one file per query)
├── pgextrashttp/                // Sub-package for HTTP handlers
│   ├── handler.go               // NewHandler(client, opts) http.Handler
│   ├── middleware.go             // BasicAuth, Recovery, RequestLogger
│   ├── options.go               // HandlerOptions: auth, enabled actions, prefix
│   └── templates/               // Embedded HTML templates
│       ├── embed.go
│       ├── layout.html
│       ├── index.html           // Dashboard with query links
│       ├── query.html           // Query result table
│       └── diagnose.html        // Health check report
└── cmd/
    └── pgextras/                // Standalone CLI binary
        ├── main.go              // Entrypoint, root command, imports _ "github.com/lib/pq"
        ├── serve.go             // `pgextras serve` — standalone HTTP server
        ├── query.go             // `pgextras <query_name>` — run single query to stdout
        └── diagnose.go          // `pgextras diagnose` — health check to stdout
```

---

## 3. Core Library API

### 3.1 Client

```go
package pgextras

import (
    "context"
    "database/sql"
)

// Config holds connection and behavior settings.
type Config struct {
    // DB is an open *sql.DB handle (required).
    // The caller is responsible for opening the connection and registering
    // a PostgreSQL driver (e.g., lib/pq, pgx/v5/stdlib) before calling New.
    DB *sql.DB

    // Schema defaults to "public". Overridden by PG_EXTRAS_SCHEMA env var.
    Schema string

    // MissingFKConstraintsIgnoreList excludes specific columns from missing_fk_constraints.
    // Entries can be "table.column" or just "column".
    MissingFKConstraintsIgnoreList []string

    // MissingFKIndexesIgnoreList excludes specific columns from missing_fk_indexes.
    MissingFKIndexesIgnoreList []string
}

// Client is the primary entrypoint for running pg-extras queries.
type Client struct { /* unexported fields */ }

// New creates a Client from an existing *sql.DB.
// The caller owns the *sql.DB lifecycle — Client never closes it.
func New(cfg Config) (*Client, error)

// DB exposes the underlying *sql.DB for advanced use.
func (c *Client) DB() *sql.DB
```

### 3.2 Query Methods

Every query method follows the same signature pattern:

```go
func (c *Client) QueryName(ctx context.Context, params ...QueryNameParams) ([]QueryNameResult, error)
```

Params are variadic so callers can omit them entirely to use defaults. If provided, only the first element is used.

### 3.3 Version Detection

```go
// pgVersion returns the major PG server version (e.g., 14, 15, 16).
func (c *Client) pgVersion(ctx context.Context) (int, error)

// pgStatStatementsVersion returns the installed pg_stat_statements version string.
// Returns "" if the extension is not installed.
func (c *Client) pgStatStatementsVersion(ctx context.Context) (string, error)
```

These are called lazily and cached for the lifetime of the Client.

---

## 4. Complete Query Catalog

### 4.1 Query Performance

#### `Outliers` / `OutliersLegacy`

Queries from `pg_stat_statements` ordered by total execution time. The library auto-selects the legacy variant when `pg_stat_statements` < 1.8 (PG < 13).

```go
type OutliersParams struct {
    Limit int // default: 10
}

type OutliersResult struct {
    Query        string  `db:"query"`
    ExecTime     string  `db:"exec_time"`       // interval
    PropExecTime string  `db:"prop_exec_time"`  // percentage
    NCalls       string  `db:"ncalls"`           // formatted count
    AvgExecMs    string  `db:"avg_exec_ms"`      // not present in legacy
    SyncIOTime   string  `db:"sync_io_time"`     // interval
}
```

**Requires**: `pg_stat_statements` extension.

#### `Calls` / `CallsLegacy`

Same structure as Outliers but ordered by call count descending.

```go
type CallsParams struct {
    Limit int // default: 10
}

type CallsResult struct {
    Query        string `db:"qry"`
    ExecTime     string `db:"exec_time"`
    PropExecTime string `db:"prop_exec_time"`
    NCalls       string `db:"ncalls"`
    AvgExecMs    string `db:"avg_exec_ms"`
    SyncIOTime   string `db:"sync_io_time"`
}
```

**Requires**: `pg_stat_statements` extension.

#### `LongRunningQueries`

Currently running queries exceeding a duration threshold.

```go
type LongRunningQueriesParams struct {
    Threshold string // default: "500 milliseconds" (PG interval syntax)
}

type LongRunningQueriesResult struct {
    PID      int    `db:"pid"`
    Duration string `db:"duration"`
    Query    string `db:"query"`
}
```

### 4.2 Cache Efficiency

#### `CacheHit`

Overall buffer cache hit ratios for indexes and tables.

```go
type CacheHitResult struct {
    Name  string `db:"name"`   // "index hit rate" or "table hit rate"
    Ratio string `db:"ratio"`  // decimal, e.g. "0.9995"
}
```

#### `IndexCacheHit`

Per-index buffer cache hit breakdown.

```go
type IndexCacheHitResult struct {
    Name       string `db:"name"`
    BufferHits int64  `db:"buffer_hits"`
    BlockReads int64  `db:"block_reads"`
    TotalRead  int64  `db:"total_read"`
    Ratio      string `db:"ratio"`
}
```

#### `TableCacheHit`

Per-table buffer cache hit breakdown. Same struct shape as IndexCacheHitResult.

```go
type TableCacheHitResult struct {
    Name       string `db:"name"`
    BufferHits int64  `db:"buffer_hits"`
    BlockReads int64  `db:"block_reads"`
    TotalRead  int64  `db:"total_read"`
    Ratio      string `db:"ratio"`
}
```

### 4.3 Index Analysis

#### `IndexUsage`

Percentage of scans that used an index, per table.

```go
type IndexUsageResult struct {
    RelName                  string `db:"relname"`
    PercentOfTimesIndexUsed  string `db:"percent_of_times_index_used"`
    RowsInTable              int64  `db:"rows_in_table"`
}
```

#### `UnusedIndexes`

Indexes with fewer than N scans on tables larger than 5 pages.

```go
type UnusedIndexesParams struct {
    MaxScans int // default: 50
}

type UnusedIndexesResult struct {
    Table      string `db:"table"`
    Index      string `db:"index"`
    IndexSize  string `db:"index_size"`
    IndexScans int64  `db:"index_scans"`
}
```

#### `DuplicateIndexes`

Indexes with identical column sets, opclass, expression, and predicate.

```go
type DuplicateIndexesResult struct {
    Size string `db:"size"`
    Idx1 string `db:"idx1"`
    Idx2 string `db:"idx2"`
    Idx3 string `db:"idx3"`
    Idx4 string `db:"idx4"`
}
```

#### `NullIndexes`

Indexes containing significant NULL values with estimated space savings.

```go
type NullIndexesParams struct {
    MinRelationSizeMB int // default: 10
}

type NullIndexesResult struct {
    OID             int64  `db:"oid"`
    Index           string `db:"index"`
    IndexSize       string `db:"index_size"`
    Unique          bool   `db:"unique"`
    IndexedColumn   string `db:"indexed_column"`
    NullFrac        string `db:"null_frac"`
    ExpectedSaving  string `db:"expected_saving"`
}
```

#### `IndexSize`

Size of each index in the database.

```go
type IndexSizeResult struct {
    Name   string `db:"name"`
    Size   string `db:"size"`
    Schema string `db:"schema"`
}
```

#### `TotalIndexSize`

Total size of all indexes combined.

```go
type TotalIndexSizeResult struct {
    Size string `db:"size"`
}
```

#### `MissingFKIndexes`

Foreign key columns lacking a supporting index.

```go
type MissingFKIndexesParams struct {
    TableName  string   // optional: filter to specific table
    IgnoreList []string // merged with Config.MissingFKIndexesIgnoreList
}

type MissingFKIndexesResult struct {
    Table      string `db:"table"`
    ColumnName string `db:"column_name"`
}
```

#### `MissingFKConstraints`

Columns matching `<table_singular>_id` where a related table exists but no FK constraint is defined. Excludes Rails polymorphic columns (those with a matching `_type` sibling).

```go
type MissingFKConstraintsParams struct {
    TableName  string   // optional
    IgnoreList []string // merged with Config.MissingFKConstraintsIgnoreList
}

type MissingFKConstraintsResult struct {
    Table      string `db:"table"`
    ColumnName string `db:"column_name"`
}
```

### 4.4 Table Analysis

#### `TableSize`

Size of each table/materialized view (via `pg_table_size()`).

```go
type TableSizeResult struct {
    Name   string `db:"name"`
    Size   string `db:"size"`
    Schema string `db:"schema"`
}
```

#### `TotalTableSize`

Total size per table including indexes and TOAST (via `pg_total_relation_size()`).

```go
type TotalTableSizeResult struct {
    Name string `db:"name"`
    Size string `db:"size"`
}
```

#### `TableIndexesSize`

Total index size per table (via `pg_indexes_size()`).

```go
type TableIndexesSizeResult struct {
    Table       string `db:"table"`
    IndexesSize string `db:"indexes_size"`
}
```

#### `RecordsRank`

Estimated row counts per table from `n_live_tup`, descending.

```go
type RecordsRankResult struct {
    Name           string `db:"name"`
    EstimatedCount int64  `db:"estimated_count"`
}
```

#### `Bloat`

Table and index bloat estimation via complex CTE against `pg_stats`/`pg_class`.

```go
type BloatResult struct {
    Type       string `db:"type"`        // "table" or "index"
    SchemaName string `db:"schemaname"`
    ObjectName string `db:"object_name"`
    Bloat      string `db:"bloat"`       // numeric ratio
    Waste      string `db:"waste"`       // human-readable size
}
```

#### `SeqScans`

Sequential scan counts per table, descending.

```go
type SeqScansResult struct {
    Name  string `db:"name"`
    Count int64  `db:"count"`
}
```

#### `TableSchema`

Column definitions for a specific table.

```go
type TableSchemaParams struct {
    TableName string // REQUIRED
}

type TableSchemaResult struct {
    ColumnName    string  `db:"column_name"`
    DataType      string  `db:"data_type"`
    IsNullable    string  `db:"is_nullable"`
    ColumnDefault *string `db:"column_default"` // nullable
}
```

#### `TableForeignKeys`

FK constraints for a specific table.

```go
type TableForeignKeysParams struct {
    TableName string // REQUIRED
}

type TableForeignKeysResult struct {
    TableName         string `db:"table_name"`
    ConstraintName    string `db:"constraint_name"`
    ColumnName        string `db:"column_name"`
    ForeignTableName  string `db:"foreign_table_name"`
    ForeignColumnName string `db:"foreign_column_name"`
}
```

### 4.5 Vacuum and Maintenance

#### `VacuumStats`

Dead rows, autovacuum thresholds, and `expect_autovacuum` flag. Includes insert-based thresholds for PG 13+.

```go
type VacuumStatsResult struct {
    Schema                      string  `db:"schema"`
    Table                       string  `db:"table"`
    LastManualVacuum            *string `db:"last_manual_vacuum"`
    ManualVacuumCount           int64   `db:"manual_vacuum_count"`
    LastAutovacuum              *string `db:"last_autovacuum"`
    AutovacuumCount             int64   `db:"autovacuum_count"`
    RowCount                    int64   `db:"rowcount"`
    DeadRowCount                int64   `db:"dead_rowcount"`
    DeadTupAutovacuumThreshold  *int64  `db:"dead_tup_autovacuum_threshold"`
    NInsSinceVacuum             *int64  `db:"n_ins_since_vacuum"`
    InsertAutovacuumThreshold   *int64  `db:"insert_autovacuum_threshold"`
    ExpectAutovacuum            *bool   `db:"expect_autovacuum"`
}
```

#### `VacuumProgress`

Current VACUUM/autovacuum progress from `pg_stat_progress_vacuum`.

```go
type VacuumProgressResult struct {
    Database          string `db:"database"`
    Schema            string `db:"schema"`
    Table             string `db:"table"`
    PID               int    `db:"pid"`
    Phase             string `db:"phase"`
    HeapBlksTotal     int64  `db:"heap_blks_total"`
    HeapBlksScanned   int64  `db:"heap_blks_scanned"`
    HeapBlksVacuumed  int64  `db:"heap_blks_vacuumed"`
    IndexVacuumCount  int64  `db:"index_vacuum_count"`
}
```

#### `AnalyzeProgress`

Current ANALYZE progress from `pg_stat_progress_analyze`.

```go
type AnalyzeProgressResult struct {
    Database            string `db:"database"`
    Schema              string `db:"schema"`
    Table               string `db:"table"`
    PID                 int    `db:"pid"`
    Phase               string `db:"phase"`
    SampleBlksTotal     int64  `db:"sample_blks_total"`
    SampleBlksScanned   int64  `db:"sample_blks_scanned"`
    ExtStatsTotal       int64  `db:"ext_stats_total"`
    ExtStatsComputed    int64  `db:"ext_stats_computed"`
}
```

#### `VacuumIOStats`

Autovacuum I/O statistics from `pg_stat_io`. **Requires PostgreSQL 16+** (returns an informational message on older versions).

```go
type VacuumIOStatsResult struct {
    BackendType string  `db:"backend_type"`
    Object      string  `db:"object"`
    Context     string  `db:"context"`
    Reads       int64   `db:"reads"`
    Writes      int64   `db:"writes"`
    Writebacks  int64   `db:"writebacks"`
    Extends     int64   `db:"extends"`
    Evictions   int64   `db:"evictions"`
    Reuses      int64   `db:"reuses"`
    Fsyncs      int64   `db:"fsyncs"`
    StatsReset  *string `db:"stats_reset"`
}
```

### 4.6 Active Processes and Locks

#### `Connections`

All active connections to the current database.

```go
type ConnectionsResult struct {
    Username        string  `db:"username"`
    PID             int     `db:"pid"`
    ClientAddress   *string `db:"client_address"`
    ApplicationName string  `db:"application_name"`
}
```

#### `Locks`

Queries holding exclusive locks.

```go
type LocksParams struct {
    Limit int // default: 20
}

type LocksResult struct {
    PID           int     `db:"pid"`
    RelName       *string `db:"relname"`
    TransactionID *string `db:"transactionid"`
    Granted       bool    `db:"granted"`
    QuerySnippet  string  `db:"query_snippet"`
    Mode          string  `db:"mode"`
    Age           string  `db:"age"`
    Application   string  `db:"application"`
}
```

#### `AllLocks`

All current locks regardless of mode. Same result type as `LocksResult`.

#### `Blocking`

Statements holding locks that block other statements.

```go
type BlockingResult struct {
    BlockedPID        int    `db:"blocked_pid"`
    BlockingStatement string `db:"blocking_statement"`
    BlockingDuration  string `db:"blocking_duration"`
    BlockingPID       int    `db:"blocking_pid"`
    BlockedStatement  string `db:"blocked_statement"`
    BlockedDuration   string `db:"blocked_duration"`
    BlockedSQLApp     string `db:"blocked_sql_app"`
    BlockingSQLApp    string `db:"blocking_sql_app"`
}
```

### 4.7 System Information

#### `DBSettings`

Selected PostgreSQL configuration values (PGTune-relevant subset).

```go
type DBSettingsResult struct {
    Name    string  `db:"name"`
    Setting string  `db:"setting"`
    Unit    *string `db:"unit"`
}
```

Returns settings for: `max_connections`, `shared_buffers`, `effective_cache_size`, `maintenance_work_mem`, `checkpoint_completion_target`, `wal_buffers`, `default_statistics_target`, `random_page_cost`, `effective_io_concurrency`, `work_mem`, `min_wal_size`, `max_wal_size`.

#### `SSLUsed`

Whether the current connection uses SSL. **Requires** `sslinfo` extension.

```go
type SSLUsedResult struct {
    SSLIsUsed bool `db:"ssl_is_used"`
}
```

#### `Extensions`

Installed PostgreSQL extensions.

```go
type ExtensionsResult struct {
    Name             string  `db:"name"`
    DefaultVersion   string  `db:"default_version"`
    InstalledVersion string  `db:"installed_version"`
    Comment          *string `db:"comment"`
}
```

#### `BuffercacheStats`

Relations buffered in shared memory. **Requires** `pg_buffercache` extension.

```go
type BuffercacheStatsParams struct {
    Limit int // default: 20
}

type BuffercacheStatsResult struct {
    RelName            string `db:"relname"`
    Buffered           string `db:"buffered"`
    BufferPercent      string `db:"buffer_percent"`
    PercentOfRelation  string `db:"percent_of_relation"`
}
```

#### `BuffercacheUsage`

Cached block counts per relation. **Requires** `pg_buffercache` extension.

```go
type BuffercacheUsageParams struct {
    Limit int // default: 20
}

type BuffercacheUsageResult struct {
    RelName string `db:"relname"`
    Buffers int64  `db:"buffers"`
}
```

### 4.8 Administrative Actions

#### `KillAll`

Terminates all connections to the current database except self.

```go
func (c *Client) KillAll(ctx context.Context) error
```

#### `KillPID`

Terminates a specific connection.

```go
func (c *Client) KillPID(ctx context.Context, pid int) error
```

#### `PgStatStatementsReset`

Resets all `pg_stat_statements` statistics.

```go
func (c *Client) PgStatStatementsReset(ctx context.Context) error
```

#### `AddExtensions`

Creates `sslinfo`, `pg_buffercache`, and `pg_stat_statements` if not present.

```go
func (c *Client) AddExtensions(ctx context.Context) error
```

### 4.9 Fun

#### `Mandelbrot`

Renders the Mandelbrot set in ASCII art via recursive SQL CTE.

```go
type MandelbrotResult struct {
    Mandelbrot string `db:"mandelbrot"`
}
```

---

## 5. Aggregated Meta-Queries

These are not single SQL queries but Go-side aggregations calling multiple underlying queries.

### `TableInfo`

```go
type TableInfoParams struct {
    TableName string // optional: filter to single table
}

type TableInfoResult struct {
    TableName      string  `json:"table_name"`
    TableSize      string  `json:"table_size"`
    TableCacheHit  string  `json:"table_cache_hit"`
    IndexCacheHit  string  `json:"index_cache_hit"`
    EstimatedRows  int64   `json:"estimated_rows"`
    SeqScans       int64   `json:"seq_scans"`
    IdxScans       int64   `json:"idx_scans"`
}
```

Internally calls: `TableSize`, `TableCacheHit`, `IndexCacheHit`, `RecordsRank`, `SeqScans`, and an additional query for `pg_stat_user_tables.idx_scan`.

### `IndexInfo`

```go
type IndexInfoParams struct {
    TableName string // optional: filter to single table
}

type IndexInfoResult struct {
    IndexName  string  `json:"index_name"`
    TableName  string  `json:"table_name"`
    Columns    string  `json:"columns"`
    IndexSize  string  `json:"index_size"`
    IndexScans int64   `json:"index_scans"`
    NullFrac   string  `json:"null_frac"`
}
```

---

## 6. Default Parameters

| Query | Parameter | Default |
|---|---|---|
| `Outliers` / `OutliersLegacy` | `Limit` | `10` |
| `Calls` / `CallsLegacy` | `Limit` | `10` |
| `LongRunningQueries` | `Threshold` | `"500 milliseconds"` |
| `Locks` | `Limit` | `20` |
| `UnusedIndexes` | `MaxScans` | `50` |
| `NullIndexes` | `MinRelationSizeMB` | `10` |
| `BuffercacheStats` | `Limit` | `20` |
| `BuffercacheUsage` | `Limit` | `20` |
| All queries | `Schema` (from Config) | `"public"` or `PG_EXTRAS_SCHEMA` env |

---

## 7. `pg_stat_statements` Version Routing

The constant `newPgStatStatementsVersion = "1.8"` controls query selection. When `pg_stat_statements` is installed with version < 1.8 (PG < 13), the `_legacy` SQL variants are used automatically. These use `total_time` instead of `total_exec_time` and lack `mean_exec_time`.

Additionally, PG 17+ renamed `blk_read_time`/`blk_write_time` to `shared_blk_read_time`/`shared_blk_write_time`. The version detection logic must handle this.

---

## 8. Diagnose / Health Check

```go
type Severity int

const (
    SeverityOK   Severity = iota
    SeverityWarn
    SeverityFail
)

type DiagnoseResult struct {
    CheckName string   `json:"check_name"`
    Severity  Severity `json:"severity"`
    OK        bool     `json:"ok"`
    Message   string   `json:"message"`
}

// Diagnose runs all health checks and returns results.
func (c *Client) Diagnose(ctx context.Context) ([]DiagnoseResult, error)
```

### Diagnostic Checks and Thresholds

| # | Check Name | Underlying Query | Pass Condition |
|---|---|---|---|
| 1 | Table Cache Hit Rate | `CacheHit` → "table hit rate" | ratio ≥ 0.985 (98.5%) |
| 2 | Index Cache Hit Rate | `CacheHit` → "index hit rate" | ratio ≥ 0.985 (98.5%) |
| 3 | Unused Indexes | `UnusedIndexes` | No results returned |
| 4 | Null Indexes | `NullIndexes(MinRelationSizeMB: 10)` | No results returned |
| 5 | Bloat | `Bloat` | No tables/indexes with bloat ratio ≥ 10 |
| 6 | Duplicate Indexes | `DuplicateIndexes` | No results returned |
| 7 | Outliers | `Outliers` | No single query > 90% `prop_exec_time` |
| 8 | SSL Connection | `SSLUsed` | `ssl_is_used` = true |
| 9 | Connection Count | `Connections` vs `max_connections` setting | < 90% of max |
| 10 | Long Running Queries | `LongRunningQueries(Threshold: "500 milliseconds")` | No results returned |

---

## 9. Output Formatting

```go
// Format controls how query results are rendered.
type Format int

const (
    FormatTable Format = iota // ASCII bordered table to io.Writer
    FormatJSON                // JSON array
    FormatCSV                 // CSV with header row
)

// FormatResults writes results in the given format.
// T must be a slice of result structs.
func FormatResults[T any](w io.Writer, results []T, format Format, title string) error
```

### ASCII Table Format

Uses box-drawing characters. Title row shows the query description:

```
+----------------+------------------------+
|        Index and table hit rate         |
+----------------+------------------------+
| name           | ratio                  |
+----------------+------------------------+
| index hit rate | 0.99957765013541945832 |
| table hit rate | 1.00                   |
+----------------+------------------------+
```

**Recommended implementation**: `github.com/olekukonenko/tablewriter` or a lightweight custom renderer.

---

## 10. HTTP Dashboard (`pgextrashttp` Sub-Package)

### Handler Setup

```go
package pgextrashttp

// HandlerOptions configures the HTTP dashboard.
type HandlerOptions struct {
    // PathPrefix is prepended to all routes (e.g., "/pg_extras").
    // Default: "/pg_extras"
    PathPrefix string

    // BasicAuthUsername and BasicAuthPassword enable HTTP Basic Auth.
    // If both are empty and PublicDashboard is false, all requests return 401.
    // Can also be set via PGEXTRAS_AUTH_USER / PGEXTRAS_AUTH_PASSWORD env vars.
    BasicAuthUsername string
    BasicAuthPassword string

    // PublicDashboard disables authentication entirely.
    // Can also be set via PGEXTRAS_PUBLIC_DASHBOARD=true env var.
    PublicDashboard bool

    // EnabledActions lists administrative actions exposed as POST endpoints.
    // Options: "kill_all", "pg_stat_statements_reset", "add_extensions"
    // Default: nil (no admin actions exposed)
    EnabledActions []string

    // Logger receives request logs. If nil, logs are discarded.
    Logger *slog.Logger

    // TemplateFuncs allows injecting additional template functions.
    TemplateFuncs template.FuncMap
}

// NewHandler returns an http.Handler that serves the pg-extras dashboard.
// It registers routes under opts.PathPrefix.
func NewHandler(client *pgextras.Client, opts HandlerOptions) http.Handler
```

### Routes

| Method | Path | Description |
|---|---|---|
| GET | `{prefix}/` | Dashboard index — lists all queries as clickable links |
| GET | `{prefix}/query/{name}` | Execute query, render HTML table |
| GET | `{prefix}/diagnose` | Run health check, render color-coded report |
| POST | `{prefix}/action/kill_all` | Kill all connections (if enabled) |
| POST | `{prefix}/action/pg_stat_statements_reset` | Reset stats (if enabled) |
| POST | `{prefix}/action/add_extensions` | Install extensions (if enabled) |
| GET | `{prefix}/api/query/{name}` | JSON API: execute query, return JSON array |
| GET | `{prefix}/api/diagnose` | JSON API: health check results |

### Middleware Stack

Applied in order: Recovery → RequestLogger → BasicAuth (unless PublicDashboard).

### Templates

Embedded via `//go:embed templates/*.html`. Uses Tailwind CSS (CDN link in layout). Diagnose rows use `bg-green-300` for passing and `bg-red-300` for failing checks, matching the Rails version.

---

## 11. Standalone Server Mode (CLI)

The `cmd/pgextras` binary provides both a one-shot query runner **and** a persistent HTTP server, making it trivial to deploy in a Docker container for database monitoring.

### CLI Commands

```
pgextras — PostgreSQL performance insights

USAGE:
    pgextras <command> [flags]

COMMANDS:
    serve       Start the web dashboard as a standalone HTTP server
    diagnose    Run health checks and print results to stdout
    <query>     Run a specific query (e.g., pgextras cache_hit, pgextras bloat)
    list        List all available query names
    version     Print version information

GLOBAL FLAGS:
    --database-url, -d    PostgreSQL connection string (default: $DATABASE_URL)
    --schema              Schema to inspect (default: $PG_EXTRAS_SCHEMA or "public")
    --format, -f          Output format: table, json, csv (default: table)
```

### `pgextras serve`

Starts the HTTP dashboard as a long-running process, suitable for containers and background monitoring.

```
pgextras serve [flags]

FLAGS:
    --addr, -a              Listen address (default: ":8080", env: PGEXTRAS_ADDR)
    --path-prefix           URL path prefix (default: "/", env: PGEXTRAS_PATH_PREFIX)
    --auth-user             Basic auth username (env: PGEXTRAS_AUTH_USER)
    --auth-password         Basic auth password (env: PGEXTRAS_AUTH_PASSWORD)
    --public                Disable authentication (env: PGEXTRAS_PUBLIC_DASHBOARD)
    --enable-actions        Comma-separated admin actions to enable
                            (e.g., "kill_all,pg_stat_statements_reset,add_extensions")
    --tls-cert              Path to TLS certificate file
    --tls-key               Path to TLS private key file
    --read-timeout          HTTP read timeout (default: "5s")
    --write-timeout         HTTP write timeout (default: "30s")
    --idle-timeout          HTTP idle timeout (default: "120s")

EXAMPLES:
    # Start on port 3000, no auth (development)
    pgextras serve -a :3000 --public

    # Production: basic auth, TLS, admin actions disabled
    pgextras serve --auth-user admin --auth-password secret --tls-cert cert.pem --tls-key key.pem

    # Behind a reverse proxy at /monitoring/pg
    pgextras serve --path-prefix /monitoring/pg --public
```

### `pgextras <query_name>`

Runs a single query and prints results to stdout:

```
EXAMPLES:
    pgextras cache_hit
    pgextras outliers --limit 20
    pgextras long_running_queries --threshold "1 second"
    pgextras unused_indexes -f json
    pgextras table_schema --table-name users
```

### `pgextras diagnose`

Runs all diagnostic checks and prints a summary:

```
EXAMPLES:
    pgextras diagnose
    pgextras diagnose -f json
```

Output in table mode uses ANSI colors: green for OK, red for FAIL.

### Implementation Notes

- Use `github.com/spf13/cobra` for command parsing (idiomatic Go CLI framework)
- All flags map to corresponding env vars for container-friendly configuration
- The `serve` command registers SIGINT/SIGTERM handlers for graceful shutdown
- Startup prints the listen address, PG version detected, and schema in use
- The CLI binary is the **only place** that calls `sql.Open()` and imports a driver. It uses `lib/pq` by default. The `main.go` file contains `import _ "github.com/lib/pq"` and opens the connection from `DATABASE_URL`:

```go
// cmd/pgextras/main.go (simplified)
package main

import (
    "database/sql"
    "os"

    _ "github.com/lib/pq"
    "github.com/<org>/go-pg-extras"
)

func openDB() (*sql.DB, error) {
    return sql.Open("postgres", os.Getenv("DATABASE_URL"))
}
```

---

## 12. Docker Support

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /pgextras ./cmd/pgextras

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /pgextras /usr/local/bin/pgextras
EXPOSE 8080
ENTRYPOINT ["pgextras"]
CMD ["serve", "--public", "--addr", ":8080"]
```

### Docker Compose Example

```yaml
services:
  pgextras:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://user:pass@db:5432/myapp?sslmode=disable
      PGEXTRAS_AUTH_USER: admin
      PGEXTRAS_AUTH_PASSWORD: ${PGEXTRAS_PASSWORD}
    depends_on:
      - db
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 64M

  db:
    image: postgres:16
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

### Kubernetes Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pgextras
spec:
  replicas: 1
  selector:
    matchLabels:
      app: pgextras
  template:
    metadata:
      labels:
        app: pgextras
    spec:
      containers:
        - name: pgextras
          image: ghcr.io/<org>/go-pg-extras:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: pg-credentials
                  key: url
            - name: PGEXTRAS_AUTH_USER
              valueFrom:
                secretKeyRef:
                  name: pgextras-auth
                  key: username
            - name: PGEXTRAS_AUTH_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: pgextras-auth
                  key: password
          resources:
            limits:
              memory: 64Mi
              cpu: 100m
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 5
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8080
```

### Health Endpoint

The standalone server exposes `GET /healthz` which returns `200 OK` with body `{"status":"ok"}` if the database connection pool can successfully ping. Returns `503 Service Unavailable` otherwise. This endpoint is always unauthenticated regardless of auth settings.

---

## 13. Environment Variables Summary

| Variable | Used By | Description |
|---|---|---|
| `DATABASE_URL` | CLI only | PostgreSQL connection string (CLI opens the `*sql.DB` on behalf of the user) |
| `PG_EXTRAS_SCHEMA` | Library + CLI | Default schema (default: `"public"`) |
| `PGEXTRAS_ADDR` | CLI `serve` | Listen address (default: `":8080"`) |
| `PGEXTRAS_PATH_PREFIX` | CLI `serve` | URL path prefix (default: `"/"`) |
| `PGEXTRAS_AUTH_USER` | CLI `serve` + HTTP handler | Basic auth username |
| `PGEXTRAS_AUTH_PASSWORD` | CLI `serve` + HTTP handler | Basic auth password |
| `PGEXTRAS_PUBLIC_DASHBOARD` | CLI `serve` + HTTP handler | Set `"true"` to disable auth |

---

## 14. PostgreSQL Version Compatibility

| PG Version | Notes |
|---|---|
| 12–13 | Uses `_legacy` variants for `outliers` and `calls` |
| 13+ | Insert-based autovacuum thresholds in `vacuum_stats` |
| 16+ | `vacuum_io_stats` via `pg_stat_io` |
| 17+ | `shared_blk_read_time`/`shared_blk_write_time` in `pg_stat_statements` |

The library detects the server version at first query and caches it.

---

## 15. Extension Dependencies

| Extension | Queries Enabled | Auto-install via `AddExtensions` |
|---|---|---|
| `pg_stat_statements` | `Outliers`, `Calls`, `PgStatStatementsReset` | Yes (but requires `shared_preload_libraries` in postgresql.conf) |
| `pg_buffercache` | `BuffercacheStats`, `BuffercacheUsage` | Yes |
| `sslinfo` | `SSLUsed` | Yes |

Queries requiring missing extensions return a clear `ErrExtensionNotInstalled` error with the extension name, rather than exposing raw SQL errors.

---

## 16. Error Handling

```go
var (
    // ErrExtensionNotInstalled is returned when a query requires an extension
    // that is not installed in the database.
    ErrExtensionNotInstalled = errors.New("required PostgreSQL extension is not installed")

    // ErrRequiredParam is returned when a required parameter is missing.
    ErrRequiredParam = errors.New("required parameter is missing")

    // ErrUnsupportedPGVersion is returned when a query requires a newer PG version.
    ErrUnsupportedPGVersion = errors.New("query requires a newer PostgreSQL version")
)
```

---

## 17. Testing Strategy

- **Unit tests**: Use `sqlmock` (`github.com/DATA-DOG/go-sqlmock`) against `*sql.DB` for param/SQL construction tests
- **Integration tests**: Use `testcontainers-go` with PostgreSQL 14, 15, 16, 17 images
- **Golden file tests**: Compare ASCII table output against fixtures
- **CI matrix**: Test against PG 14, 15, 16, 17, 18 using GitHub Actions

---

## 18. Go Module Dependencies

```
require (
    github.com/lib/pq              // PostgreSQL driver (used by cmd/pgextras only)
    github.com/spf13/cobra         // CLI framework
    github.com/olekukonenko/tablewriter // ASCII table output (or custom)
)
```

The core `pgextras` package has **zero third-party runtime dependencies** — it uses only `database/sql` from the standard library. The `lib/pq` import lives exclusively in `cmd/pgextras/main.go` as a blank import (`import _ "github.com/lib/pq"`). Users integrating the library into their own applications can use whichever `database/sql`-compatible PostgreSQL driver they prefer (`lib/pq`, `pgx/v5/stdlib`, etc.).

---

## 19. Usage Examples

### Library Usage

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq" // or pgx/v5/stdlib, or any PG driver
    "github.com/<org>/go-pg-extras"
)

func main() {
    ctx := context.Background()
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    client, err := pgextras.New(pgextras.Config{DB: db})
    if err != nil {
        panic(err)
    }

    // Run a query with defaults
    results, _ := client.CacheHit(ctx)
    pgextras.FormatResults(os.Stdout, results, pgextras.FormatTable, "Cache Hit Rates")

    // Run with custom params
    outliers, _ := client.Outliers(ctx, pgextras.OutliersParams{Limit: 5})
    pgextras.FormatResults(os.Stdout, outliers, pgextras.FormatJSON, "")

    // Health check
    diagnoses, _ := client.Diagnose(ctx)
    for _, d := range diagnoses {
        fmt.Printf("[%s] %s: %s\n", d.Severity, d.CheckName, d.Message)
    }
}
```

### HTTP Handler Integration (with existing server)

```go
package main

import (
    "context"
    "database/sql"
    "net/http"

    _ "github.com/lib/pq"
    "github.com/<org>/go-pg-extras"
    "github.com/<org>/go-pg-extras/pgextrashttp"
)

func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    defer db.Close()

    client, _ := pgextras.New(pgextras.Config{DB: db})

    mux := http.NewServeMux()
    mux.Handle("/", yourAppHandler())
    mux.Handle("/pg_extras/", pgextrashttp.NewHandler(client, pgextrashttp.HandlerOptions{
        PathPrefix:        "/pg_extras",
        BasicAuthUsername: "admin",
        BasicAuthPassword: "secret",
        EnabledActions:    []string{"add_extensions"},
    }))
    http.ListenAndServe(":3000", mux)
}
```

### Standalone Docker

```bash
# One-shot query
docker run --rm -e DATABASE_URL=postgres://... ghcr.io/<org>/go-pg-extras cache_hit

# Persistent dashboard
docker run -d -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/db \
  -e PGEXTRAS_AUTH_USER=admin \
  -e PGEXTRAS_AUTH_PASSWORD=secret \
  ghcr.io/<org>/go-pg-extras serve
```
