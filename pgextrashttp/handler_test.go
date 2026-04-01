package pgextrashttp

import "testing"

func TestQueryNames(t *testing.T) {
	names := QueryNames()
	if len(names) == 0 {
		t.Fatal("expected non-empty query names list")
	}

	// Verify a few known query names are present.
	expected := map[string]bool{
		"cache_hit":         false,
		"index_usage":       false,
		"bloat":             false,
		"locks":             false,
		"connections":       false,
		"vacuum_stats":      false,
		"db_settings":       false,
		"extensions":        false,
		"tables":            false,
		"mandelbrot":        false,
		"outliers":          false,
		"seq_scans":         false,
		"table_size":        false,
		"index_size":        false,
		"duplicate_indexes": false,
	}

	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("expected query name %q not found in QueryNames()", name)
		}
	}
}
