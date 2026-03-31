package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/bmorton/go-pg-extras/pgextrashttp"
)

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the web dashboard as a standalone HTTP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Value:   ":8080",
				EnvVars: []string{"PGEXTRAS_ADDR"},
				Usage:   "Listen address",
			},
			&cli.StringFlag{
				Name:    "path-prefix",
				Value:   "/",
				EnvVars: []string{"PGEXTRAS_PATH_PREFIX"},
				Usage:   "URL path prefix",
			},
			&cli.StringFlag{
				Name:    "auth-user",
				EnvVars: []string{"PGEXTRAS_AUTH_USER"},
				Usage:   "Basic auth username",
			},
			&cli.StringFlag{
				Name:    "auth-password",
				EnvVars: []string{"PGEXTRAS_AUTH_PASSWORD"},
				Usage:   "Basic auth password",
			},
			&cli.BoolFlag{
				Name:    "public",
				EnvVars: []string{"PGEXTRAS_PUBLIC_DASHBOARD"},
				Usage:   "Disable authentication",
			},
			&cli.StringFlag{
				Name:  "enable-actions",
				Usage: "Comma-separated admin actions to enable (kill_all,pg_stat_statements_reset,add_extensions)",
			},
			&cli.StringFlag{
				Name:  "tls-cert",
				Usage: "Path to TLS certificate file",
			},
			&cli.StringFlag{
				Name:  "tls-key",
				Usage: "Path to TLS private key file",
			},
			&cli.DurationFlag{
				Name:  "read-timeout",
				Value: 5 * time.Second,
				Usage: "HTTP read timeout",
			},
			&cli.DurationFlag{
				Name:  "write-timeout",
				Value: 30 * time.Second,
				Usage: "HTTP write timeout",
			},
			&cli.DurationFlag{
				Name:  "idle-timeout",
				Value: 120 * time.Second,
				Usage: "HTTP idle timeout",
			},
		},
		Action: func(c *cli.Context) error {
			client, db, err := openClient(c)
			if err != nil {
				return err
			}
			defer db.Close()

			logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

			prefix := c.String("path-prefix")
			if prefix == "/" {
				prefix = ""
			}

			var enabledActions []string
			if actions := c.String("enable-actions"); actions != "" {
				enabledActions = strings.Split(actions, ",")
			}

			handler := pgextrashttp.NewHandler(client, pgextrashttp.HandlerOptions{
				PathPrefix:        prefix,
				BasicAuthUsername: c.String("auth-user"),
				BasicAuthPassword: c.String("auth-password"),
				PublicDashboard:   c.Bool("public"),
				EnabledActions:    enabledActions,
				Logger:            logger,
			})

			mux := http.NewServeMux()
			if prefix != "" {
				mux.Handle(prefix+"/", handler)
			} else {
				mux.Handle("/", handler)
			}

			// Health endpoint — always unauthenticated.
			mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
				if err := db.PingContext(r.Context()); err != nil {
					w.WriteHeader(http.StatusServiceUnavailable)
					json.NewEncoder(w).Encode(map[string]string{"status": "error", "error": err.Error()})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			})

			addr := c.String("addr")
			server := &http.Server{
				Addr:         addr,
				Handler:      mux,
				ReadTimeout:  c.Duration("read-timeout"),
				WriteTimeout: c.Duration("write-timeout"),
				IdleTimeout:  c.Duration("idle-timeout"),
			}

			logger.Info("starting server",
				"addr", addr,
				"prefix", prefix,
				"schema", client.Schema(),
			)

			// Graceful shutdown.
			errC := make(chan error, 1)
			go func() {
				tlsCert := c.String("tls-cert")
				tlsKey := c.String("tls-key")
				if tlsCert != "" && tlsKey != "" {
					errC <- server.ListenAndServeTLS(tlsCert, tlsKey)
				} else {
					errC <- server.ListenAndServe()
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case err := <-errC:
				return fmt.Errorf("server error: %w", err)
			case sig := <-quit:
				logger.Info("shutting down", "signal", sig)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				return server.Shutdown(ctx)
			}
		},
	}
}
