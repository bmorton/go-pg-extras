package pgextras

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// pgVersion returns the major PG server version (e.g., 14, 15, 16).
func (c *Client) pgVersion(ctx context.Context) (int, error) {
	c.versionOnce.Do(func() {
		var versionStr string
		c.versionErr = c.db.QueryRowContext(ctx, "SHOW server_version").Scan(&versionStr)
		if c.versionErr != nil {
			return
		}
		parts := strings.Split(versionStr, ".")
		c.versionVal, c.versionErr = strconv.Atoi(parts[0])
	})
	return c.versionVal, c.versionErr
}

// pgStatStatementsVersion returns the installed pg_stat_statements version string.
// Returns "" if the extension is not installed.
func (c *Client) pgStatStatementsVersion(ctx context.Context) (string, error) {
	c.pssVersionOnce.Do(func() {
		var version sql.NullString
		err := c.db.QueryRowContext(ctx,
			"SELECT installed_version FROM pg_available_extensions WHERE name = 'pg_stat_statements'",
		).Scan(&version)
		if err != nil {
			if err == sql.ErrNoRows {
				c.pssVersionVal = ""
				return
			}
			c.pssVersionErr = err
			return
		}
		if version.Valid {
			c.pssVersionVal = version.String
		}
	})
	return c.pssVersionVal, c.pssVersionErr
}

const newPgStatStatementsVersion = "1.8"

// selectSQLVariant picks the appropriate SQL file based on PG/extension version.
// For queries using pg_stat_statements, it checks the extension version and PG major version.
func (c *Client) selectSQLVariant(ctx context.Context, base string) (string, error) {
	switch base {
	case "outliers", "calls":
		pssVersion, err := c.pgStatStatementsVersion(ctx)
		if err != nil {
			return "", fmt.Errorf("detecting pg_stat_statements version: %w", err)
		}
		if pssVersion == "" {
			return "", fmt.Errorf("%w: pg_stat_statements", ErrExtensionNotInstalled)
		}
		if compareVersions(pssVersion, newPgStatStatementsVersion) < 0 {
			return base + "_legacy", nil
		}
		pgVer, err := c.pgVersion(ctx)
		if err != nil {
			return "", err
		}
		if pgVer >= 17 {
			return base + "_17", nil
		}
		return base, nil

	case "vacuum_progress":
		pgVer, err := c.pgVersion(ctx)
		if err != nil {
			return "", err
		}
		if pgVer >= 17 {
			return "vacuum_progress_17", nil
		}
		return base, nil

	case "vacuum_io_stats":
		pgVer, err := c.pgVersion(ctx)
		if err != nil {
			return "", err
		}
		if pgVer < 16 {
			return "vacuum_io_stats_legacy", nil
		}
		return base, nil

	default:
		return base, nil
	}
}

// compareVersions compares two semver-like version strings.
// Returns -1, 0, or 1.
func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	for i := 0; i < maxLen; i++ {
		var av, bv int
		if i < len(aParts) {
			av, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bv, _ = strconv.Atoi(bParts[i])
		}
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}
	return 0
}

// versionCache fields are stored in Client.
type versionCache struct {
	versionOnce sync.Once
	versionVal  int
	versionErr  error

	pssVersionOnce sync.Once
	pssVersionVal  string
	pssVersionErr  error
}
