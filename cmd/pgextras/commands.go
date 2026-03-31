package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	pgextras "github.com/bmorton/go-pg-extras"
	"github.com/bmorton/go-pg-extras/pgextrashttp"
)

func diagnoseCommand() *cli.Command {
	return &cli.Command{
		Name:  "diagnose",
		Usage: "Run health checks and print results to stdout",
		Action: func(c *cli.Context) error {
			client, db, err := openClient(c)
			if err != nil {
				return err
			}
			defer db.Close()

			ctx := context.Background()
			results, err := client.Diagnose(ctx)
			if err != nil {
				return err
			}

			format := parseFormat(c.String("format"))
			return pgextras.FormatResults(os.Stdout, results, format, "Health Check")
		},
	}
}

func listCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all available query names",
		Action: func(c *cli.Context) error {
			names := pgextrashttp.QueryNames()
			for _, name := range names {
				fmt.Println(name)
			}
			return nil
		},
	}
}

func queryCommand() *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "Run a specific query (e.g., pgextras query cache_hit)",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "limit",
				Usage: "Limit number of results",
			},
			&cli.StringFlag{
				Name:  "threshold",
				Usage: "Duration threshold (e.g., '1 second')",
			},
			&cli.StringFlag{
				Name:  "table-name",
				Usage: "Table name (for table_schema, table_foreign_keys)",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return fmt.Errorf("query name required. Run 'pgextras list' to see all available queries")
			}
			return runQueryCLI(c, c.Args().First())
		},
	}
}
