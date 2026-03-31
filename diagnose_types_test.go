package pgextras

import (
	"fmt"
	"testing"
)

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityOK, "OK"},
		{SeverityWarn, "WARN"},
		{SeverityFail, "FAIL"},
		{Severity(99), "Severity(99)"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Severity(%d)", int(tt.sev)), func(t *testing.T) {
			if got := tt.sev.String(); got != tt.want {
				t.Errorf("Severity(%d).String() = %q, want %q", int(tt.sev), got, tt.want)
			}
		})
	}
}
