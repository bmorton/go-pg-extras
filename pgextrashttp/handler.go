package pgextrashttp

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"

	pgextras "github.com/bmorton/go-pg-extras"
)

//go:embed templates/*.html
var templateFiles embed.FS

//go:embed static/*
var staticFiles embed.FS

// NewHandler returns an http.Handler that serves the pg-extras dashboard.
func NewHandler(client *pgextras.Client, opts HandlerOptions) http.Handler {
	if opts.PathPrefix == "" {
		opts.PathPrefix = "/pg_extras"
	}
	opts.PathPrefix = strings.TrimRight(opts.PathPrefix, "/")
	opts.Client = client

	if opts.BasicAuthUsername == "" {
		opts.BasicAuthUsername = os.Getenv("PGEXTRAS_AUTH_USER")
	}
	if opts.BasicAuthPassword == "" {
		opts.BasicAuthPassword = os.Getenv("PGEXTRAS_AUTH_PASSWORD")
	}
	if !opts.PublicDashboard && os.Getenv("PGEXTRAS_PUBLIC_DASHBOARD") == "true" {
		opts.PublicDashboard = true
	}
	if opts.Logger == nil {
		opts.Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	funcMap := template.FuncMap{
		"deref": func(s *string) string {
			if s == nil {
				return ""
			}
			return *s
		},
	}
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFS(templateFiles, "templates/*.html"))

	mux := http.NewServeMux()

	prefix := opts.PathPrefix

	mux.HandleFunc(prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		// Only match exact prefix path for index.
		path := strings.TrimPrefix(r.URL.Path, prefix)
		if path != "/" && path != "" {
			http.NotFound(w, r)
			return
		}
		renderIndex(w, tmpl, opts)
	})

	mux.HandleFunc(prefix+"/static/", func(w http.ResponseWriter, r *http.Request) {
		// Strip the prefix to get the file path within the embed.
		filePath := strings.TrimPrefix(r.URL.Path, prefix+"/")
		data, err := staticFiles.ReadFile(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(filePath, ".css") {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(data)
	})

	mux.HandleFunc(prefix+"/query/", func(w http.ResponseWriter, r *http.Request) {
		queryName := strings.TrimPrefix(r.URL.Path, prefix+"/query/")
		handleQuery(w, r, tmpl, client, opts, queryName)
	})

	mux.HandleFunc(prefix+"/diagnose", func(w http.ResponseWriter, r *http.Request) {
		handleDiagnose(w, r, tmpl, client, opts)
	})

	mux.HandleFunc(prefix+"/api/query/", func(w http.ResponseWriter, r *http.Request) {
		queryName := strings.TrimPrefix(r.URL.Path, prefix+"/api/query/")
		handleAPIQuery(w, r, client, queryName)
	})

	mux.HandleFunc(prefix+"/api/diagnose", func(w http.ResponseWriter, r *http.Request) {
		handleAPIDiagnose(w, r, client)
	})

	mux.HandleFunc(prefix+"/action/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		action := strings.TrimPrefix(r.URL.Path, prefix+"/action/")
		handleAction(w, r, client, opts, action)
	})

	// Wrap with middleware.
	var handler http.Handler = mux
	if !opts.PublicDashboard && opts.BasicAuthUsername != "" && opts.BasicAuthPassword != "" {
		handler = basicAuth(handler, opts.BasicAuthUsername, opts.BasicAuthPassword)
	}
	handler = recoveryMiddleware(handler, opts.Logger)
	handler = requestLogger(handler, opts.Logger)

	return handler
}

func renderIndex(w http.ResponseWriter, tmpl *template.Template, opts HandlerOptions) {
	data := map[string]any{
		"Title":     "Dashboard",
		"Prefix":    opts.PathPrefix,
		"Actions":   opts.EnabledActions,
		"ActiveNav": "index",
	}
	renderLayout(w, tmpl, "index", data)
}

func handleQuery(w http.ResponseWriter, r *http.Request, tmpl *template.Template, client *pgextras.Client, opts HandlerOptions, queryName string) {
	if queryName == "table_schemas" {
		handleTableSchemas(w, r, tmpl, client, opts)
		return
	}

	ctx := r.Context()
	headers, rows, err := executeQuery(ctx, client, queryName)

	data := map[string]any{
		"Title":     queryName,
		"Prefix":    opts.PathPrefix,
		"QueryName": queryName,
		"Headers":   headers,
		"Rows":      rows,
		"ActiveNav": queryName,
	}
	if err != nil {
		data["Error"] = err.Error()
	}
	renderLayout(w, tmpl, "query", data)
}

func handleTableSchemas(w http.ResponseWriter, r *http.Request, tmpl *template.Template, client *pgextras.Client, opts HandlerOptions) {
	tables, err := client.TableSchemasGrouped(r.Context())
	data := map[string]any{
		"Title":     "Table Schemas",
		"Prefix":    opts.PathPrefix,
		"Tables":    tables,
		"ActiveNav": "table_schemas",
	}
	if err != nil {
		data["Error"] = err.Error()
	}
	renderLayout(w, tmpl, "table_schemas", data)
}

func handleDiagnose(w http.ResponseWriter, r *http.Request, tmpl *template.Template, client *pgextras.Client, opts HandlerOptions) {
	results, err := client.Diagnose(r.Context())
	data := map[string]any{
		"Title":     "Health Check",
		"Prefix":    opts.PathPrefix,
		"Results":   results,
		"ActiveNav": "diagnose",
	}
	if err != nil {
		data["Error"] = err.Error()
	}
	renderLayout(w, tmpl, "diagnose", data)
}

func handleAPIQuery(w http.ResponseWriter, r *http.Request, client *pgextras.Client, queryName string) {
	ctx := r.Context()
	result, err := executeQueryRaw(ctx, client, queryName)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func handleAPIDiagnose(w http.ResponseWriter, r *http.Request, client *pgextras.Client) {
	results, err := client.Diagnose(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func handleAction(w http.ResponseWriter, r *http.Request, client *pgextras.Client, opts HandlerOptions, action string) {
	// Check if action is enabled.
	enabled := false
	for _, a := range opts.EnabledActions {
		if a == action {
			enabled = true
			break
		}
	}
	if !enabled {
		http.Error(w, "Action not enabled", http.StatusForbidden)
		return
	}

	ctx := r.Context()
	var err error
	switch action {
	case "kill_all":
		err = client.KillAll(ctx)
	case "pg_stat_statements_reset":
		err = client.PgStatStatementsReset(ctx)
	case "add_extensions":
		err = client.AddExtensions(ctx)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, opts.PathPrefix+"/", http.StatusSeeOther)
}

func renderLayout(w http.ResponseWriter, tmpl *template.Template, name string, data map[string]any) {
	var content strings.Builder
	if err := tmpl.ExecuteTemplate(&content, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Content"] = template.HTML(content.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// executeQuery runs a named query and returns headers and string rows for HTML display.
func executeQuery(ctx context.Context, client *pgextras.Client, name string) ([]string, [][]string, error) {
	result, err := executeQueryRaw(ctx, client, name)
	if err != nil {
		return nil, nil, err
	}

	// Use reflection to extract headers and rows.
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Slice || rv.Len() == 0 {
		return nil, nil, nil
	}

	elemType := rv.Type().Elem()
	var headers []string
	for i := 0; i < elemType.NumField(); i++ {
		tag := elemType.Field(i).Tag.Get("db")
		if tag == "" {
			tag = strings.ToLower(elemType.Field(i).Name)
		}
		headers = append(headers, tag)
	}

	var rows [][]string
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		var row []string
		for j := 0; j < elem.NumField(); j++ {
			f := elem.Field(j)
			if f.Kind() == reflect.Ptr {
				if f.IsNil() {
					row = append(row, "")
				} else {
					row = append(row, fmt.Sprintf("%v", f.Elem().Interface()))
				}
			} else {
				row = append(row, fmt.Sprintf("%v", f.Interface()))
			}
		}
		rows = append(rows, row)
	}
	return headers, rows, nil
}

// executeQueryRaw runs a named query and returns the raw result slice.
func executeQueryRaw(ctx context.Context, client *pgextras.Client, name string) (any, error) {
	switch name {
	case "cache_hit":
		return client.CacheHit(ctx)
	case "index_cache_hit":
		return client.IndexCacheHit(ctx)
	case "table_cache_hit":
		return client.TableCacheHit(ctx)
	case "outliers":
		return client.Outliers(ctx)
	case "calls":
		return client.Calls(ctx)
	case "long_running_queries":
		return client.LongRunningQueries(ctx)
	case "index_usage":
		return client.IndexUsage(ctx)
	case "unused_indexes":
		return client.UnusedIndexes(ctx)
	case "duplicate_indexes":
		return client.DuplicateIndexes(ctx)
	case "null_indexes":
		return client.NullIndexes(ctx)
	case "index_size":
		return client.IndexSize(ctx)
	case "total_index_size":
		return client.TotalIndexSize(ctx)
	case "index_scans":
		return client.IndexScans(ctx)
	case "indexes":
		return client.Indexes(ctx)
	case "table_size":
		return client.TableSize(ctx)
	case "total_table_size":
		return client.TotalTableSize(ctx)
	case "table_indexes_size":
		return client.TableIndexesSize(ctx)
	case "records_rank":
		return client.RecordsRank(ctx)
	case "bloat":
		return client.Bloat(ctx)
	case "seq_scans":
		return client.SeqScans(ctx)
	case "table_schemas":
		return client.TableSchemas(ctx)
	case "foreign_keys":
		return client.ForeignKeys(ctx)
	case "table_index_scans":
		return client.TableIndexScans(ctx)
	case "tables":
		return client.Tables(ctx)
	case "vacuum_stats":
		return client.VacuumStats(ctx)
	case "vacuum_progress":
		return client.VacuumProgress(ctx)
	case "analyze_progress":
		return client.AnalyzeProgress(ctx)
	case "vacuum_io_stats":
		return client.VacuumIOStats(ctx)
	case "connections":
		return client.Connections(ctx)
	case "locks":
		return client.Locks(ctx)
	case "all_locks":
		return client.AllLocks(ctx)
	case "blocking":
		return client.Blocking(ctx)
	case "db_settings":
		return client.DBSettings(ctx)
	case "ssl_used":
		return client.SSLUsed(ctx)
	case "extensions":
		return client.Extensions(ctx)
	case "buffercache_stats":
		return client.BuffercacheStats(ctx)
	case "buffercache_usage":
		return client.BuffercacheUsage(ctx)
	case "mandelbrot":
		return client.Mandelbrot(ctx)
	default:
		return nil, fmt.Errorf("unknown query: %s", name)
	}
}

// QueryNames returns all available query names.
func QueryNames() []string {
	return []string{
		"cache_hit",
		"index_cache_hit",
		"table_cache_hit",
		"outliers",
		"calls",
		"long_running_queries",
		"index_usage",
		"unused_indexes",
		"duplicate_indexes",
		"null_indexes",
		"index_size",
		"total_index_size",
		"index_scans",
		"indexes",
		"table_size",
		"total_table_size",
		"table_indexes_size",
		"records_rank",
		"bloat",
		"seq_scans",
		"table_schemas",
		"foreign_keys",
		"table_index_scans",
		"tables",
		"vacuum_stats",
		"vacuum_progress",
		"analyze_progress",
		"vacuum_io_stats",
		"connections",
		"locks",
		"all_locks",
		"blocking",
		"db_settings",
		"ssl_used",
		"extensions",
		"buffercache_stats",
		"buffercache_usage",
		"mandelbrot",
	}
}
