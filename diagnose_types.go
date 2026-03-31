package pgextras

import "fmt"

// Severity represents the health check severity level.
type Severity int

const (
	SeverityOK Severity = iota
	SeverityWarn
	SeverityFail
)

// String returns the string representation of Severity.
func (s Severity) String() string {
	switch s {
	case SeverityOK:
		return "OK"
	case SeverityWarn:
		return "WARN"
	case SeverityFail:
		return "FAIL"
	default:
		return fmt.Sprintf("Severity(%d)", int(s))
	}
}

// DiagnoseResult holds the result of a single health check.
type DiagnoseResult struct {
	CheckName string   `json:"check_name"`
	Severity  Severity `json:"severity"`
	OK        bool     `json:"ok"`
	Message   string   `json:"message"`
}
