# go-pg-extras

[![Go Reference](https://pkg.go.dev/badge/github.com/bmorton/go-pg-extras.svg)](https://pkg.go.dev/github.com/bmorton/go-pg-extras)
[![CI](https://github.com/bmorton/go-pg-extras/actions/workflows/main.yml/badge.svg)](https://github.com/bmorton/go-pg-extras/actions/workflows/main.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

PostgreSQL performance diagnostics for Go — as a library, CLI, or web dashboard.

A Go port of [rails-pg-extras](https://github.com/pawurb/rails-pg-extras), providing 50+ useful PostgreSQL diagnostic queries to help you identify performance issues, missing indexes, cache inefficiencies, lock contention, and more.

<p align="center">
  <img src="docs/screenshots/web-health-check.png" alt="Web dashboard — Health Check" width="800">
</p>

## Features

- **50+ diagnostic queries** — cache hit ratios, index usage, bloat, vacuum stats, locks, table sizes, and more
- **Automated health checks** — the `diagnose` command evaluates 10 key metrics with OK/FAIL status
- **Three usage modes** — Go library, standalone CLI binary, or HTTP web dashboard
- **Version-aware SQL** — automatically selects the right query variant for PostgreSQL 13, 17, and legacy versions
- **Multiple output formats** — JSON, CSV, and ASCII table
- **Secure web dashboard** — basic auth, TLS, and configurable path prefix
- **No ORM dependency** — uses stdlib `database/sql`; bring your own `*sql.DB`

## Acknowledgments

The SQL queries used in this project originate from several open-source projects and community resources. All credit for the original queries and concepts belongs to their respective authors:

- **[heroku-pg-extras](https://github.com/heroku/heroku-pg-extras)** — the original collection of PostgreSQL diagnostic queries
- **[rails-pg-extras](https://github.com/pawurb/rails-pg-extras)** and **[ruby-pg-extras](https://github.com/pawurb/ruby-pg-extras)** — the Ruby/Rails implementations that this project is ported from
- **[PostgreSQL Unused Index Size](https://hakibenita.com/postgresql-unused-index-size)** by Haki Benita
- **[Useful SQLs to Check Contents of PostgreSQL shared_buffers](https://sites.google.com/site/itmyshare/database-tips-and-examples/postgres/useful-sqls-to-check-contents-of-postgresql-shared_buffer)**
- **[Index Maintenance](https://wiki.postgresql.org/wiki/Index_Maintenance)** — PostgreSQL Wiki

This project is an independent Go port. It is **not** affiliated with, endorsed by, or officially connected to Heroku, the rails-pg-extras project, or any of the other sources listed above.

This project is also built with the use of GitHub Copilot tooling under human oversight and review.

## Installation

### Go library

```sh
go get github.com/bmorton/go-pg-extras
```

### CLI binary

```sh
go install github.com/bmorton/go-pg-extras/cmd/pgextras@latest
```

### Docker

```sh
docker run --rm ghcr.io/bmorton/go-pg-extras \
  --database-url "postgres://user:pass@host:5432/dbname" cache_hit
```

## Quick Start — Library

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/bmorton/go-pg-extras"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	client, err := pgextras.New(pgextras.Config{DB: db})
	if err != nil {
		log.Fatal(err)
	}

	results, err := client.CacheHit(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range results {
		fmt.Printf("%s: %s\n", r.Name, r.Ratio)
	}
}
```

## Quick Start — CLI

```sh
# Check buffer cache hit ratios
pgextras -d "postgres://user:pass@localhost:5432/mydb" cache_hit

# Show top 10 slowest queries
pgextras -d "postgres://user:pass@localhost:5432/mydb" outliers --limit 10

# Run automated health checks
pgextras -d "postgres://user:pass@localhost:5432/mydb" diagnose

# List all available queries
pgextras -d "postgres://user:pass@localhost:5432/mydb" list
```

You can also set `DATABASE_URL` as an environment variable instead of passing `-d` each time.

Output format can be changed with `--format` (`table`, `json`, or `csv`).

## Quick Start — Web Dashboard

Start the built-in web dashboard as a standalone HTTP server:

```sh
pgextras serve -d "postgres://user:pass@localhost:5432/mydb" --public --addr :8080
```

Then visit `http://localhost:8080/pg_extras` to browse all queries from the sidebar, run health checks, and inspect your database.

<p align="center">
  <img src="docs/screenshots/web-table-overview.png" alt="Web dashboard — Table Sizes" width="800">
</p>

<details>
<summary>More screenshots</summary>

| | |
|---|---|
| [![Unused Indexes](docs/screenshots/web-unused-indexes.png)](docs/screenshots/web-unused-indexes.png) | [![Extensions](docs/screenshots/web-extensions.png)](docs/screenshots/web-extensions.png) |
| Unused Indexes | Extensions |

</details>

## Embedding in Your Application

Mount the dashboard into an existing Go HTTP server:

```go
import (
	"net/http"

	"github.com/bmorton/go-pg-extras/pgextrashttp"
)

handler := pgextrashttp.NewHandler(client, pgextrashttp.HandlerOptions{
	PathPrefix:        "/pg_extras",
	BasicAuthUsername: "admin",
	BasicAuthPassword: "secret",
	EnabledActions:    []string{"kill_all", "pg_stat_statements_reset"},
})

mux := http.NewServeMux()
mux.Handle("/pg_extras/", handler)
// ... add your other routes
http.ListenAndServe(":8080", mux)
```

Set `PublicDashboard: true` to disable authentication entirely (useful for internal/development environments).

## Available Queries

### Query Performance

| Query | Description |
|---|---|
| `outliers` | Queries with longest execution time (requires `pg_stat_statements`) |
| `calls` | Most frequently executed queries |
| `long_running_queries` | Currently running queries exceeding a duration threshold |

### Cache Efficiency

| Query | Description |
|---|---|
| `cache_hit` | Overall buffer cache hit ratios for indexes and tables |
| `index_cache_hit` | Per-index buffer cache hit breakdown |
| `table_cache_hit` | Per-table buffer cache hit breakdown |
| `buffercache_stats` | Shared buffer cache statistics |
| `buffercache_usage` | Shared buffer cache usage by relation |

### Index Analysis

| Query | Description |
|---|---|
| `index_usage` | Index hit rate per table |
| `index_size` | Size of each index |
| `index_scans` | Number of scans per index |
| `total_index_size` | Total index size per table |
| `indexes` | All indexes with size and scan count |
| `unused_indexes` | Indexes with low scan counts on large tables |
| `duplicate_indexes` | Indexes with redundant column sets |
| `null_indexes` | Indexes with high NULL ratios |

### Table Analysis

| Query | Description |
|---|---|
| `table_size` | Table size excluding indexes |
| `total_table_size` | Table size including indexes and TOAST |
| `table_indexes_size` | Combined index size per table |
| `table_overview` | Table sizes with estimated row counts |
| `table_index_scans` | Index scan counts per table |
| `records_rank` | Estimated row counts (descending) |
| `scan_activity` | Sequential vs. index scans per table |
| `seq_scans` | Sequential scan counts |
| `bloat` | Table and index bloat estimation |
| `table_schema` | Column definitions for a specific table |
| `table_schemas` | Column definitions for all tables |
| `tables` | All tables with size information |

### Foreign Keys

| Query | Description |
|---|---|
| `foreign_keys` | All foreign key constraints |
| `table_foreign_keys` | Foreign key constraints for a specific table |

### Vacuum & Maintenance

| Query | Description |
|---|---|
| `vacuum_stats` | Dead rows, autovacuum thresholds, and health |
| `vacuum_progress` | Progress of running VACUUM operations |
| `vacuum_io_stats` | Vacuum I/O statistics (PG 13+) |
| `analyze_progress` | Progress of running ANALYZE operations |

### Locks & Connections

| Query | Description |
|---|---|
| `locks` | Current locks (waiting and held) |
| `all_locks` | Extended lock information |
| `blocking` | Queries blocking other queries |
| `connections` | Active connection counts |
| `kill_pid` | Terminate a specific backend by PID |
| `kill_all` | Terminate all active backends |

### Extensions & Settings

| Query | Description |
|---|---|
| `extensions` | Installed PostgreSQL extensions |
| `db_settings` | Server configuration parameters |
| `ssl_used` | Whether the current connection uses SSL |

### Other

| Query | Description |
|---|---|
| `mandelbrot` | ASCII Mandelbrot set (just for fun) |

## Health Check (Diagnose)

The `diagnose` command runs automated health checks against your database and reports each with an OK, WARNING, or FAIL status:

| Check | Threshold |
|---|---|
| Table cache hit rate | ≥ 98.5% |
| Index cache hit rate | ≥ 98.5% |
| Unused indexes | None on large tables |
| Null indexes | No indexes with high NULL fraction |
| Bloat | No significant bloat |
| Duplicate indexes | No redundant indexes |
| Outliers | No single query dominating execution time |
| SSL connection | SSL in use |
| Connection count | Within limits |
| Long running queries | None exceeding threshold |

Run via CLI (`pgextras diagnose`), library (`client.Diagnose(ctx)`), or web dashboard (`/pg_extras/diagnose`).

## Configuration Reference

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string | — |
| `PG_EXTRAS_SCHEMA` | Schema to inspect | `public` |
| `PGEXTRAS_AUTH_USER` | HTTP dashboard Basic Auth username | — |
| `PGEXTRAS_AUTH_PASSWORD` | HTTP dashboard Basic Auth password | — |
| `PGEXTRAS_PUBLIC_DASHBOARD` | Set to `"true"` to disable dashboard auth | — |

### CLI Flags

| Flag | Alias | Description |
|---|---|---|
| `--database-url` | `-d` | PostgreSQL connection string (overrides `DATABASE_URL`) |
| `--schema` | | Schema to inspect (overrides `PG_EXTRAS_SCHEMA`) |
| `--format` | `-f` | Output format: `table`, `json`, or `csv` (default: `table`) |
| `--limit` | | Limit number of results (for `outliers`, `calls`, etc.) |
| `--threshold` | | Duration threshold, e.g. `"500ms"`, `"1s"`, `"2m"` (for `long_running_queries`) |
| `--table-name` | | Table name (for `table_schema`, `table_foreign_keys`) |

### Serve Subcommand

| Flag | Description | Default |
|---|---|---|
| `--addr`, `-a` | Listen address | `:8080` |
| `--path-prefix` | URL path prefix | `/` |
| `--public` | Disable authentication | `false` |
| `--auth-user` | Basic Auth username | `PGEXTRAS_AUTH_USER` env |
| `--auth-password` | Basic Auth password | `PGEXTRAS_AUTH_PASSWORD` env |
| `--enable-action` | Admin action to enable (repeatable) | — |
| `--tls-cert` | Path to TLS certificate file | — |
| `--tls-key` | Path to TLS private key file | — |

## Troubleshooting

### `pg_stat_statements` not installed

Queries like `outliers` and `calls` require the `pg_stat_statements` extension. If you get an error about a missing relation:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

You may also need to add it to `shared_preload_libraries` in `postgresql.conf` and restart PostgreSQL.

### Permission errors

Some queries access system catalogs that require superuser or `pg_monitor` role membership. Grant the necessary role:

```sql
GRANT pg_monitor TO your_user;
```

### Connection issues

Make sure your `DATABASE_URL` is correct and the database is reachable. The connection string format is:

```
postgres://user:password@host:port/dbname?sslmode=disable
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Make your changes and add tests
4. Run `make check` to verify formatting, linting, and tests
5. Commit and push your branch
6. Open a Pull Request

## License

MIT
