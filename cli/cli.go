package cli_interface

import (
	"log"
	"net/http"

	"github.com/hydrocode-de/datailama/internal/api"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/urfave/cli/v2"
)

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
	log.Printf("Starting server on port %s...", port)
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
	}
}
