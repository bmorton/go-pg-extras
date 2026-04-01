//go:build integration

package pgextras

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testDB     *sql.DB
	testDBOnce sync.Once
	testDBErr  error
)

// getTestDB returns a shared *sql.DB for integration tests.
// It reads the DATABASE_URL environment variable.
func getTestDB(t *testing.T) *sql.DB {
	t.Helper()
	testDBOnce.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			testDBErr = fmt.Errorf("DATABASE_URL is not set")
			return
		}
		testDB, testDBErr = sql.Open("postgres", dsn)
		if testDBErr != nil {
			return
		}
		testDBErr = testDB.Ping()
	})
	if testDBErr != nil {
		t.Fatalf("failed to connect to test database: %v", testDBErr)
	}
	return testDB
}

// newTestClient creates a Client connected to the test database.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	db := getTestDB(t)
	c, err := New(Config{DB: db})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}
