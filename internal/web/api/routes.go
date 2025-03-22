package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/ollama"
	"github.com/hydrocode-de/datailama/internal/version"
	"github.com/urfave/cli/v2"
)

// RegisterRoutes sets up the API routes using Huma
func RegisterRoutes(router *http.ServeMux, dbManager *db.Manager, cliCtx *cli.Context) {
	// Create Huma API on a specific path prefix
	api := humago.New(router, huma.DefaultConfig("DataILama API", version.Version))

	// Register middleware
	api.UseMiddleware(ContextMiddleware(dbManager, cliCtx))

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
func ContextMiddleware(db *db.Manager, cliCtx *cli.Context) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// Add database manager to context
		newCtx := huma.WithValue(ctx, "db", db)

		// Transfer all relevant CLI context values to Huma context
		if cliCtx != nil {
			// Add Ollama URL to context
			ollamaURL := cliCtx.String("ollama-url")
			newCtx = huma.WithValue(newCtx, ollama.OllamaURLKey, ollamaURL)

			// Add any other CLI context values that might be needed
			// e.g., database URL, etc.
			dbURL := cliCtx.String("database-url")
			newCtx = huma.WithValue(newCtx, "database-url", dbURL)
		}

		next(newCtx)
	}
}
