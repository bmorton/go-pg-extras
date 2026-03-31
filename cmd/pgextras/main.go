package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"

	pgextras "github.com/bmorton/go-pg-extras"
)

func main() {
	app := &cli.App{
		Name:  "pgextras",
		Usage: "PostgreSQL performance insights",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "database-url",
				Aliases: []string{"d"},
				EnvVars: []string{"DATABASE_URL"},
				Usage:   "PostgreSQL connection string",
			},
			&cli.StringFlag{
				Name:    "schema",
				EnvVars: []string{"PG_EXTRAS_SCHEMA"},
				Value:   "public",
				Usage:   "Schema to inspect",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "table",
				Usage:   "Output format: table, json, csv",
			},
		},
		Commands: []*cli.Command{
			serveCommand(),
			diagnoseCommand(),
			listCommand(),
			queryCommand(),
		},
		// Default action: treat first arg as query name.
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return cli.ShowAppHelp(c)
			}
			return runQueryCLI(c, c.Args().First())
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func openClient(c *cli.Context) (*pgextras.Client, *sql.DB, error) {
	dbURL := c.String("database-url")
	if dbURL == "" {
		return nil, nil, fmt.Errorf("database URL is required (set DATABASE_URL or use --database-url)")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("opening database: %w", err)
	}
	client, err := pgextras.New(pgextras.Config{
		DB:     db,
		Schema: c.String("schema"),
	})
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	return client, db, nil
}

func parseFormat(s string) pgextras.Format {
	switch strings.ToLower(s) {
	case "json":
		return pgextras.FormatJSON
	case "csv":
		return pgextras.FormatCSV
	default:
		return pgextras.FormatTable
	}
}

func runQueryCLI(c *cli.Context, name string) error {
	client, db, err := openClient(c)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()
	format := parseFormat(c.String("format"))

	return executeAndFormat(ctx, client, name, format, c)
}

func executeAndFormat(ctx context.Context, client *pgextras.Client, name string, format pgextras.Format, c *cli.Context) error {
	switch name {
	case "cache_hit":
		r, err := client.CacheHit(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Index and table hit rate")
	case "index_cache_hit":
		r, err := client.IndexCacheHit(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Index cache hit rate")
	case "table_cache_hit":
		r, err := client.TableCacheHit(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table cache hit rate")
	case "outliers":
		limit := c.Int("limit")
		if limit == 0 {
			limit = 10
		}
		r, err := client.Outliers(ctx, pgextras.OutliersParams{Limit: limit})
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Queries with longest execution time")
	case "calls":
		limit := c.Int("limit")
		if limit == 0 {
			limit = 10
		}
		r, err := client.Calls(ctx, pgextras.CallsParams{Limit: limit})
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Queries with highest frequency of execution")
	case "long_running_queries":
		threshold := c.String("threshold")
		if threshold == "" {
			threshold = "500 milliseconds"
		}
		r, err := client.LongRunningQueries(ctx, pgextras.LongRunningQueriesParams{Threshold: threshold})
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Long running queries")
	case "index_usage":
		r, err := client.IndexUsage(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Index usage")
	case "unused_indexes":
		r, err := client.UnusedIndexes(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Unused indexes")
	case "duplicate_indexes":
		r, err := client.DuplicateIndexes(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Duplicate indexes")
	case "null_indexes":
		r, err := client.NullIndexes(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Null indexes")
	case "index_size":
		r, err := client.IndexSize(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Index size")
	case "total_index_size":
		r, err := client.TotalIndexSize(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Total index size")
	case "index_scans":
		r, err := client.IndexScans(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Index scans")
	case "indexes":
		r, err := client.Indexes(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Indexes")
	case "table_size":
		r, err := client.TableSize(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table size")
	case "total_table_size":
		r, err := client.TotalTableSize(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Total table size")
	case "table_indexes_size":
		r, err := client.TableIndexesSize(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table indexes size")
	case "records_rank":
		r, err := client.RecordsRank(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Records rank")
	case "bloat":
		r, err := client.Bloat(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table and index bloat")
	case "seq_scans":
		r, err := client.SeqScans(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Sequential scans")
	case "table_schema":
		tableName := c.String("table-name")
		if tableName == "" {
			return fmt.Errorf("--table-name is required for table_schema")
		}
		r, err := client.TableSchema(ctx, pgextras.TableSchemaParams{TableName: tableName})
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table schema: "+tableName)
	case "table_schemas":
		r, err := client.TableSchemas(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "All table schemas")
	case "table_foreign_keys":
		tableName := c.String("table-name")
		if tableName == "" {
			return fmt.Errorf("--table-name is required for table_foreign_keys")
		}
		r, err := client.TableForeignKeys(ctx, pgextras.TableForeignKeysParams{TableName: tableName})
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Foreign keys: "+tableName)
	case "foreign_keys":
		r, err := client.ForeignKeys(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Foreign keys")
	case "table_index_scans":
		r, err := client.TableIndexScans(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Table index scans")
	case "tables":
		r, err := client.Tables(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Tables")
	case "vacuum_stats":
		r, err := client.VacuumStats(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Vacuum stats")
	case "vacuum_progress":
		r, err := client.VacuumProgress(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Vacuum progress")
	case "analyze_progress":
		r, err := client.AnalyzeProgress(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Analyze progress")
	case "vacuum_io_stats":
		r, err := client.VacuumIOStats(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Vacuum I/O stats")
	case "connections":
		r, err := client.Connections(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Connections")
	case "locks":
		r, err := client.Locks(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Locks")
	case "all_locks":
		r, err := client.AllLocks(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "All locks")
	case "blocking":
		r, err := client.Blocking(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Blocking queries")
	case "db_settings":
		r, err := client.DBSettings(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Database settings")
	case "ssl_used":
		r, err := client.SSLUsed(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "SSL usage")
	case "extensions":
		r, err := client.Extensions(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Extensions")
	case "buffercache_stats":
		r, err := client.BuffercacheStats(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Buffercache stats")
	case "buffercache_usage":
		r, err := client.BuffercacheUsage(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Buffercache usage")
	case "mandelbrot":
		r, err := client.Mandelbrot(ctx)
		if err != nil {
			return err
		}
		return pgextras.FormatResults(os.Stdout, r, format, "Mandelbrot set")
	case "kill_all":
		return client.KillAll(ctx)
	case "pg_stat_statements_reset":
		return client.PgStatStatementsReset(ctx)
	case "add_extensions":
		return client.AddExtensions(ctx)
	default:
		return fmt.Errorf("unknown query: %s. Run 'pgextras list' to see all available queries", name)
	}
}
