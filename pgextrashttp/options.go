package pgextrashttp

import (
	"log/slog"

	pgextras "github.com/bmorton/go-pg-extras"
)

// HandlerOptions configures the HTTP dashboard.
type HandlerOptions struct {
	// PathPrefix is prepended to all routes (e.g., "/pg_extras").
	// Default: "/pg_extras"
	PathPrefix string

	// BasicAuthUsername and BasicAuthPassword enable HTTP Basic Auth.
	BasicAuthUsername string
	BasicAuthPassword string

	// PublicDashboard disables authentication entirely.
	PublicDashboard bool

	// EnabledActions lists administrative actions exposed as POST endpoints.
	// Options: "kill_all", "pg_stat_statements_reset", "add_extensions"
	EnabledActions []string

	// Logger receives request logs. If nil, logs are discarded.
	Logger *slog.Logger

	// Client is the pgextras client.
	Client *pgextras.Client
}
