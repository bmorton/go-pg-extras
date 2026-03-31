package pgextras

// OutliersParams configures the Outliers query.
type OutliersParams struct {
	Limit int // default: 10
}

// CallsParams configures the Calls query.
type CallsParams struct {
	Limit int // default: 10
}

// LongRunningQueriesParams configures the LongRunningQueries query.
type LongRunningQueriesParams struct {
	Threshold string // default: "500 milliseconds" (PG interval syntax)
}

// LocksParams configures the Locks query.
type LocksParams struct {
	Limit int // default: 20
}

// BlockingParams configures the Blocking query.
type BlockingParams struct {
	Limit int // default: 20
}

// UnusedIndexesParams configures the UnusedIndexes query.
type UnusedIndexesParams struct {
	MaxScans int // default: 50
}

// NullIndexesParams configures the NullIndexes query.
type NullIndexesParams struct {
	MinRelationSizeMB int // default: 10
}

// BuffercacheStatsParams configures the BuffercacheStats query.
type BuffercacheStatsParams struct {
	Limit int // default: 20
}

// BuffercacheUsageParams configures the BuffercacheUsage query.
type BuffercacheUsageParams struct {
	Limit int // default: 20
}

// TableSchemaParams configures the TableSchema query.
type TableSchemaParams struct {
	TableName string // REQUIRED
}

// TableForeignKeysParams configures the TableForeignKeys query.
type TableForeignKeysParams struct {
	TableName string // REQUIRED
}

// KillPIDParams configures the KillPID action.
type KillPIDParams struct {
	PID int // REQUIRED
}
