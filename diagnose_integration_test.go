//go:build integration

package pgextras

import (
	"context"
	"testing"
)

func TestIntegrationDiagnose(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	results, err := c.Diagnose(ctx)
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if len(results) != 10 {
		t.Errorf("expected 10 diagnostic checks, got %d", len(results))
	}
	for _, r := range results {
		t.Logf("  %s: severity=%s ok=%v message=%s", r.CheckName, r.Severity, r.OK, r.Message)
	}
}
