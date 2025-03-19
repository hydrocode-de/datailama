package cli_interface

import (
	"fmt"
	"net/http"

	"github.com/hydrocode-de/datailama/internal/api"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/version"
	"github.com/urfave/cli/v2"
)

func versionAction(c *cli.Context, isShort bool) error {
	if isShort {
		fmt.Fprintln(c.App.Writer, version.Version)
	} else {
		fmt.Fprintln(c.App.Writer, "DataILama CLI")
		fmt.Fprintf(c.App.Writer, "Version: %s\n", version.Version)
		fmt.Fprintf(c.App.Writer, "Build Time: %s\n", version.BuildTime)
		fmt.Fprintf(c.App.Writer, "Git Commit: %s\n", version.GitCommit)
	}
	return nil
}

// serveAction handles the serve command
func serveAction(c *cli.Context) error {
	port := c.String("port")
	dbURL := c.String("database-url")

	dbManager, err := db.New(c.Context, dbURL)
	if err != nil {
		return err
	}
	defer dbManager.Close()

	router := api.NewServer(dbManager)
	fmt.Fprintf(c.App.Writer, "Starting server on port %s...\n", port)
	return http.ListenAndServe(":"+port, router)
}

// GetCommands returns all available CLI commands
func GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "serve",
			Usage: "Start the API server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "port",
					Value:   "8080",
					Usage:   "Port to listen on",
					Aliases: []string{"p"},
					EnvVars: []string{"DATAILAMA_PORT"},
				},
				&cli.StringFlag{
					Name:    "database-url",
					Value:   "postgresql://postgres:postgres@localhost:5432/vector_db",
					Usage:   "Connection URL for the paper and vector database",
					Aliases: []string{"db"},
					EnvVars: []string{"DATAILAMA_DB_URL"},
				},
			},
			Action: serveAction,
		},
		{
			Name:  "stats",
			Usage: "Get statistics about the database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "database-url",
					Value:   "postgresql://postgres:postgres@localhost:5432/vector_db",
					Usage:   "Connection URL for the paper and vector database",
					Aliases: []string{"db"},
					EnvVars: []string{"DATAILAMA_DB_URL"},
				},
			},
			Action: statsAction,
		},
		{
			Name:  "version",
			Usage: "Print extended version information and exit",
			Action: func(c *cli.Context) error {
				return versionAction(c, false)
			},
		},
		{
			Name:  "v",
			Usage: "Print only the version number and exit",
			Action: func(c *cli.Context) error {
				return versionAction(c, true)
			},
		},
	}
}
