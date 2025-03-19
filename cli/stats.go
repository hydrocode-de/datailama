package cli_interface

import (
	"fmt"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func statsAction(c *cli.Context) error {
	dbURL := c.String("database-url")
	dbManager, err := db.New(c.Context, dbURL)
	if err != nil {
		return err
	}
	defer dbManager.Close()

	stats, err := dbManager.GetPaperStatistics(c.Context)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredDark)
	t.AppendHeader(table.Row{"Title", "Year", "Paper"})
	for _, stat := range stats {
		t.AppendRow(table.Row{stat.Title, stat.Year, stat.Count})
	}
	fmt.Println(t.Render())

	return nil
}
