package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/sql"
)

type TitleSearchOutput struct {
	Body struct {
		Count int                         `json:"count"`
		Paper []sql.SearchPaperByTitleRow `json:"paper"`
	}
}

func searchByTitle(ctx context.Context, input *struct {
	Title     string `query:"title,omitempty" doc:"The title to search for"`
	Author    string `query:"author,omitempty" doc:"The author to limit the search to"`
	OrderBy   string `query:"order,omitempty" example:"citations_year" doc:"The property to order the results by. Can be citations_year or citations"`
	Direction string `query:"direction,omitempty" example:"desc" doc:"The direction to order the results by. Can be asc or desc"`
}) (*TitleSearchOutput, error) {
	if strings.ToLower(input.OrderBy) != "citations_year" && strings.ToLower(input.OrderBy) != "citations" {
		return nil, fmt.Errorf("invalid order by argument: %v. Has to be citations_year or citations", input.OrderBy)
	}

	if strings.ToLower(input.Direction) != "asc" && strings.ToLower(input.Direction) != "desc" {
		return nil, fmt.Errorf("invalid direction argument: %v. Has to be asc or desc", input.Direction)
	}

	db := ctx.Value("db").(*db.Manager)

	papers, err := db.SearchPaperByTitle(ctx, sql.SearchPaperByTitleParams{
		Title:     input.Title,
		Author:    input.Author,
		OrderBy:   input.OrderBy,
		Direction: input.Direction,
		Limit:     15,
	})
	if err != nil {
		return nil, err
	}

	resp := &TitleSearchOutput{}
	resp.Body.Paper = papers
	resp.Body.Count = len(papers)
	return resp, nil
}
