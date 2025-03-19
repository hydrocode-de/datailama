package cli_interface

import (
	"fmt"
	"strings"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/sql"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func searchAction(c *cli.Context) error {
	// extract the title argument
	title := strings.Join(c.Args().Slice(), " ")

	// get the optional arguments
	author := c.String("author")
	orderBy := c.String("order")
	direction := c.String("direction")
	limit := c.Int("limit")

	dbURL := c.String("database-url")
	dbManager, err := db.New(c.Context, dbURL)
	if err != nil {
		return err
	}
	defer dbManager.Close()

	result, err := dbManager.SearchPaperByTitle(c.Context, sql.SearchPaperByTitleParams{
		Title:     title,
		Author:    author,
		OrderBy:   orderBy,
		Direction: direction,
		Limit:     int32(limit),
	})
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.AppendHeader(table.Row{"ID", "Title", "Author", "Year", "Citations", "Citations / Year"})
	for _, paper := range result {
		t.AppendRow(table.Row{paper.ID, paper.Title[:50] + "...", paper.Author, paper.Published.Time.Year(), paper.Citations, fmt.Sprintf("%.2f", paper.CitationsYear)})
	}

	fmt.Println(t.Render())

	return nil
}
