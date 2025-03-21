package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/version"
	"github.com/urfave/cli/v2"
)

// RegisterRoutes sets up the API routes using Huma
func RegisterRoutes(router *http.ServeMux, dbManager *db.Manager, cliCtx *cli.Context) {
	// Create Huma API on a specific path prefix
	api := humago.New(router, huma.DefaultConfig("DataILama API", version.Version))

	// Register middleware
	api.UseMiddleware(DbMiddleware(dbManager, cliCtx))

	// Register API endpoints
	huma.Register(api, huma.Operation{
		OperationID: "getVersion",
		Method:      http.MethodGet,
		Path:        "/api/version",
		Summary:     "Get the version",
		Description: "Get the version of DataILama",
	}, getVersion)

	huma.Register(api, huma.Operation{
		OperationID: "searchPaperByTitle",
		Method:      http.MethodGet,
		Path:        "/api/paper/search/title",
		Summary:     "Search for Papers by Title",
		Description: "Search by Title. Currently the search is a case insensive excact match to a part of the title.",
	}, searchByTitle)

	huma.Register(api, huma.Operation{
		OperationID: "searchPaperBody",
		Method:      http.MethodGet,
		Path:        "/api/paper/search/body",
		Summary:     "Search for Papers by Body",
		Description: "Search by Body. Currently the search is a case insensive excact match to a part of the body.",
	}, searchPaperBody)
}

// DbMiddleware provides database access to API handlers
func DbMiddleware(db *db.Manager, cliCtx *cli.Context) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		newCtx := huma.WithValue(ctx, "db", db)

		// Get Ollama URL from CLI context
		ollamaURL := cliCtx.String("ollama-url")

		// Add Ollama URL to context
		newCtx = huma.WithValue(newCtx, "ollama-url", ollamaURL)

		next(newCtx)
	}
}
