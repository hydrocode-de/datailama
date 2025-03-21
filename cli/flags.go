package cli_interface

import (
	"github.com/urfave/cli/v2"
)

func GetConnectionFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "database-url",
			Value:   "postgresql://postgres:postgres@localhost:5432/vector_db",
			Usage:   "Connection URL for the paper and vector database",
			Aliases: []string{"db"},
			EnvVars: []string{"DATAILAMA_DB_URL"},
		},
		&cli.StringFlag{
			Name:    "ollama-url",
			Value:   "http://localhost:11434",
			Usage:   "URL of the Ollama server",
			Aliases: []string{"o"},
			EnvVars: []string{"DATAILAMA_OLLAMA_URL"},
		},
	}

	return flags
}
