package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/version"
)

// RegisterRoutes sets up the API routes using Huma
func RegisterRoutes(router *http.ServeMux, dbManager *db.Manager) {
	// Create Huma API on a specific path prefix
	api := humago.New(router, huma.DefaultConfig("DataILama API", version.Version))

	// Register middleware
	api.UseMiddleware(DbMiddleware(dbManager))

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
		Path:        "/api/paper/search",
		Summary:     "Search for Papers by Title",
		Description: "Search by Title. Currently the search is a case insensive excact match to a part of the title.",
	}, searchByTitle)
}

// DbMiddleware provides database access to API handlers
func DbMiddleware(db *db.Manager) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		newCtx := huma.WithValue(ctx, "db", db)
		next(newCtx)
	}
}
