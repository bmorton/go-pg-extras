package pgextras

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Diagnose runs all health checks and returns results.
func (c *Client) Diagnose(ctx context.Context) ([]DiagnoseResult, error) {
	var results []DiagnoseResult

	checks := []func(ctx context.Context) DiagnoseResult{
		c.diagnoseTableCacheHitRate,
		c.diagnoseIndexCacheHitRate,
		c.diagnoseUnusedIndexes,
		c.diagnoseNullIndexes,
		c.diagnoseBloat,
		c.diagnoseDuplicateIndexes,
		c.diagnoseOutliers,
		c.diagnoseSSL,
		c.diagnoseConnectionCount,
		c.diagnoseLongRunningQueries,
	}

	for _, check := range checks {
		results = append(results, check(ctx))
	}
	return results, nil
}

func (c *Client) diagnoseTableCacheHitRate(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Table Cache Hit Rate"}
	hits, err := c.CacheHit(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	for _, h := range hits {
		if h.Name == "table hit rate" {
			ratio, err := strconv.ParseFloat(strings.TrimSpace(h.Ratio), 64)
			if err != nil {
				result.Severity = SeverityWarn
				result.Message = fmt.Sprintf("Could not parse ratio: %s", h.Ratio)
				return result
			}
			if ratio >= 0.985 {
				result.OK = true
				result.Severity = SeverityOK
				result.Message = fmt.Sprintf("Table cache hit rate: %.4f", ratio)
			} else {
				result.Severity = SeverityFail
				result.Message = fmt.Sprintf("Table cache hit rate is %.4f (threshold: 0.985)", ratio)
			}
			return result
		}
	}
	result.Severity = SeverityWarn
	result.Message = "No table hit rate data found"
	return result
}

func (c *Client) diagnoseIndexCacheHitRate(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Index Cache Hit Rate"}
	hits, err := c.CacheHit(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	for _, h := range hits {
		if h.Name == "index hit rate" {
			ratio, err := strconv.ParseFloat(strings.TrimSpace(h.Ratio), 64)
			if err != nil {
				result.Severity = SeverityWarn
				result.Message = fmt.Sprintf("Could not parse ratio: %s", h.Ratio)
				return result
			}
			if ratio >= 0.985 {
				result.OK = true
				result.Severity = SeverityOK
				result.Message = fmt.Sprintf("Index cache hit rate: %.4f", ratio)
			} else {
				result.Severity = SeverityFail
				result.Message = fmt.Sprintf("Index cache hit rate is %.4f (threshold: 0.985)", ratio)
			}
			return result
		}
	}
	result.Severity = SeverityWarn
	result.Message = "No index hit rate data found"
	return result
}

func (c *Client) diagnoseUnusedIndexes(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Unused Indexes"}
	indexes, err := c.UnusedIndexes(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	if len(indexes) == 0 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "No unused indexes found"
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Found %d unused indexes", len(indexes))
	}
	return result
}

func (c *Client) diagnoseNullIndexes(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Null Indexes"}
	indexes, err := c.NullIndexes(ctx, NullIndexesParams{MinRelationSizeMB: 10})
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	if len(indexes) == 0 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "No indexes with high null fraction found"
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Found %d indexes with high null fraction", len(indexes))
	}
	return result
}

func (c *Client) diagnoseBloat(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Bloat"}
	bloat, err := c.Bloat(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	var highBloat int
	for _, b := range bloat {
		ratio, err := strconv.ParseFloat(strings.TrimSpace(b.Bloat), 64)
		if err != nil {
			continue
		}
		if ratio >= 10 {
			highBloat++
		}
	}
	if highBloat == 0 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "No significant bloat detected"
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Found %d tables/indexes with bloat ratio >= 10", highBloat)
	}
	return result
}

func (c *Client) diagnoseDuplicateIndexes(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Duplicate Indexes"}
	dups, err := c.DuplicateIndexes(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	if len(dups) == 0 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "No duplicate indexes found"
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Found %d sets of duplicate indexes", len(dups))
	}
	return result
}

func (c *Client) diagnoseOutliers(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Outliers"}
	outliers, err := c.Outliers(ctx)
	if err != nil {
		// pg_stat_statements might not be installed - just warn.
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Could not check (pg_stat_statements may not be installed): %v", err)
		result.OK = true
		return result
	}
	for _, o := range outliers {
		pct := strings.TrimSpace(strings.TrimSuffix(o.PropExecTime, "%"))
		prop, err := strconv.ParseFloat(pct, 64)
		if err != nil {
			continue
		}
		if prop > 90 {
			result.Severity = SeverityFail
			result.Message = fmt.Sprintf("Query consuming %.1f%% of total execution time", prop)
			return result
		}
	}
	result.OK = true
	result.Severity = SeverityOK
	result.Message = "No single query dominates execution time"
	return result
}

func (c *Client) diagnoseSSL(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "SSL Connection"}
	ssl, err := c.SSLUsed(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Could not check SSL (sslinfo extension may not be installed): %v", err)
		return result
	}
	if len(ssl) > 0 && ssl[0].SSLIsUsed {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "SSL is in use"
	} else {
		result.Severity = SeverityFail
		result.Message = "SSL is not in use"
	}
	return result
}

func (c *Client) diagnoseConnectionCount(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Connection Count"}
	conns, err := c.Connections(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}

	settings, err := c.DBSettings(ctx)
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error getting settings: %v", err)
		return result
	}

	var maxConn int
	for _, s := range settings {
		if s.Name == "max_connections" {
			maxConn, _ = strconv.Atoi(s.Setting)
			break
		}
	}

	if maxConn == 0 {
		result.Severity = SeverityWarn
		result.Message = "Could not determine max_connections"
		return result
	}

	usage := float64(len(conns)) / float64(maxConn)
	if usage < 0.9 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = fmt.Sprintf("Using %d of %d connections (%.0f%%)", len(conns), maxConn, usage*100)
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Using %d of %d connections (%.0f%%) — exceeds 90%% threshold", len(conns), maxConn, usage*100)
	}
	return result
}

func (c *Client) diagnoseLongRunningQueries(ctx context.Context) DiagnoseResult {
	result := DiagnoseResult{CheckName: "Long Running Queries"}
	queries, err := c.LongRunningQueries(ctx, LongRunningQueriesParams{Threshold: "500 milliseconds"})
	if err != nil {
		result.Severity = SeverityWarn
		result.Message = fmt.Sprintf("Error: %v", err)
		return result
	}
	if len(queries) == 0 {
		result.OK = true
		result.Severity = SeverityOK
		result.Message = "No long running queries"
	} else {
		result.Severity = SeverityFail
		result.Message = fmt.Sprintf("Found %d long running queries (> 500ms)", len(queries))
	}
	return result
}
