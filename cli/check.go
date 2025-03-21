package cli_interface

import (
	"fmt"

	"github.com/hydrocode-de/datailama/internal/ollama"
	"github.com/urfave/cli/v2"
)

// checkOllamaAction handles the ollama check command
func checkOllamaAction(c *cli.Context) error {
	err := ollama.OllamaConnectionFromContext(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.App.Writer, "Ollama connection successful")
	return nil
}

// GetCheckCommand returns the check command with all its subcommands
func GetCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Check various components of the system",
		Subcommands: []*cli.Command{
			{
				Name:   "ollama",
				Usage:  "Check Ollama connection and status",
				Flags:  GetConnectionFlags(),
				Action: checkOllamaAction,
			},
			// TODO: Add more subcommands here as they are implemented
			// {
			//     Name:   "db",
			//     Usage:  "Check database connection and status",
			//     Flags:  getCheckFlags(),
			//     Action: checkDBAction,
			// },
			// {
			//     Name:   "n8n",
			//     Usage:  "Check n8n connection and status",
			//     Flags:  getCheckFlags(),
			//     Action: checkN8NAction,
			// },
		},
	}
}
