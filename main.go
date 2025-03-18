package main

import (
	"log"
	"os"

	cli_interface "github.com/hydrocode-de/datailama/cli"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	// Load the dotenv file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	app := &cli.App{
		Name:     "datailama",
		Usage:    "DataILama API server",
		Commands: cli_interface.GetCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
