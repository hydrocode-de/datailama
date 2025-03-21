package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/ollama"
	"github.com/hydrocode-de/datailama/internal/sql"
	"github.com/hydrocode-de/datailama/internal/version"
)

// VersionOutput defines the response structure for the version endpoint
type VersionOutput struct {
	Body struct {
		Version   string `json:"version" example:"1.0.0" doc:"The version of the application"`
		BuildTime string `json:"build_time,omitempty" doc:"The build time of the application"`
		GitCommit string `json:"git_commit,omitempty" doc:"The git commit of the application"`
	}
}

// PaperSearchOutput defines the response structure for the paper search endpoint
type PaperSearchOutput struct {
	Body struct {
		Count int                         `json:"count"`
		Paper []sql.SearchPaperByTitleRow `json:"paper"`
	}
}

type PaperSearchBodyOutput struct {
	Body struct {
		Count          int                      `json:"count"`
		Paper          []sql.SearchPaperBodyRow `json:"paper"`
		EmbedDuration  time.Duration            `json:"embed_duration"`
		SearchDuration time.Duration            `json:"search_duration"`
	}
}

// getVersion handles the version endpoint
func getVersion(ctx context.Context, input *struct{}) (*VersionOutput, error) {
	resp := &VersionOutput{}
	resp.Body.Version = version.Version
	resp.Body.BuildTime = version.BuildTime
	resp.Body.GitCommit = version.GitCommit
	return resp, nil
}

// searchByTitle handles the paper search endpoint
func searchByTitle(ctx context.Context, input *struct {
	Title     string `query:"title,omitempty" doc:"The title to search for"`
	Author    string `query:"author,omitempty" doc:"The author to limit the search to"`
	OrderBy   string `query:"order,omitempty" example:"citations_year" doc:"The property to order the results by. Can be citations_year or citations"`
	Direction string `query:"direction,omitempty" example:"desc" doc:"The direction to order the results by. Can be asc or desc"`
}) (*PaperSearchOutput, error) {
	if strings.ToLower(input.OrderBy) != "citations_year" && strings.ToLower(input.OrderBy) != "citations" && input.OrderBy != "" {
		return nil, fmt.Errorf("invalid order by argument: %v. Has to be citations_year or citations", input.OrderBy)
	}

	if strings.ToLower(input.Direction) != "asc" && strings.ToLower(input.Direction) != "desc" && input.Direction != "" {
		return nil, fmt.Errorf("invalid direction argument: %v. Has to be asc or desc", input.Direction)
	}

	// Set defaults
	orderBy := input.OrderBy
	if orderBy == "" {
		orderBy = "citations_year"
	}

	direction := input.Direction
	if direction == "" {
		direction = "desc"
	}

	db := ctx.Value("db").(*db.Manager)

	papers, err := db.SearchPaperByTitle(ctx, sql.SearchPaperByTitleParams{
		Title:     input.Title,
		Author:    input.Author,
		OrderBy:   orderBy,
		Direction: direction,
		Limit:     15,
	})
	if err != nil {
		return nil, err
	}

	resp := &PaperSearchOutput{}
	resp.Body.Paper = papers
	resp.Body.Count = len(papers)
	return resp, nil
}

// searchByBody handles the paper search endpoint
func searchPaperBody(ctx context.Context, input *struct {
	Prompt string `query:"prompt" doc:"The prompt to search the "`
	Limit  int    `query:"limit" doc:"Limit the results (deaults to 15)"`
}) (*PaperSearchBodyOutput, error) {
	if input.Prompt == "" {
		return nil, fmt.Errorf("empty prompts are not allowed in body searches")
	}

	limit := input.Limit
	if limit == 0 {
		limit = 15
	}

	db := ctx.Value("db").(*db.Manager)
	start := time.Now()
	embedding, err := ollama.EmbedText(ctx, input.Prompt)
	if err != nil {
		return nil, err
	}
	embed_duration := time.Since(start)

	start = time.Now()
	papers, err := db.SearchPaperBody(ctx, sql.SearchPaperBodyParams{
		Embedding: embedding,
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}
	search_duration := time.Since(start)

	resp := &PaperSearchBodyOutput{}
	resp.Body.Count = len(papers)
	resp.Body.Paper = papers
	resp.Body.EmbedDuration = embed_duration
	resp.Body.SearchDuration = search_duration
	return resp, nil
}
