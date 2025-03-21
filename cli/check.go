package cli_interface

import (
	"fmt"
	"time"

	"github.com/hydrocode-de/datailama/internal/ollama"
	"github.com/urfave/cli/v2"
)

// checkOllamaAction handles the ollama check command
func checkOllamaAction(c *cli.Context) error {
	found, err := ollama.CheckOllamaConnection(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.App.Writer, "Ollama connection successful: %s\n", found)
	return nil
}

func checkEmbeddingAction(c *cli.Context) error {
	raw := c.Bool("raw-response")
	prompt := c.String("prompt")
	start := time.Now()
	embedding, err := ollama.EmbedText(c, prompt)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	if raw {
		fmt.Fprintf(c.App.Writer, "%v\n", embedding)
	} else {
		fmt.Fprintf(c.App.Writer, "Embedding length: %d (took %v ms)\n", len(embedding.Slice()), duration.Milliseconds())
	}
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
			{
				Name:  "embedding",
				Usage: "Check embedding model connection and status",
				Flags: append(
					GetConnectionFlags(),
					&cli.BoolFlag{
						Name:    "raw-response",
						Usage:   "Show raw response from embedding model instead of a text message",
						Aliases: []string{"raw"},
					},
					&cli.StringFlag{
						Name:  "prompt",
						Usage: "Prompt to embed",
						Value: "Soil moisture is highly redundant in time",
					},
				),
				Action: checkEmbeddingAction,
			},
		},
	}
}
