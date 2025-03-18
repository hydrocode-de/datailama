package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/version"
)

type VersionOutput struct {
	Body struct {
		Version   string `json:"version" example:"1.0.0" doc:"The version of the application"`
		BuildTime string `json:"build_time,omitempty" doc:"The build time of the application"`
		GitCommit string `json:"git_commit,omitempty" doc:"The git commit of the application"`
	}
}

type Server struct {
	db     *db.Manager
	router *http.ServeMux
}

func getVersion(ctx context.Context, input *struct{}) (*VersionOutput, error) {
	resp := &VersionOutput{}
	resp.Body.Version = version.Version
	resp.Body.BuildTime = version.BuildTime
	resp.Body.GitCommit = version.GitCommit
	return resp, nil
}

func DbMiddleware(db *db.Manager) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		newCtx := huma.WithValue(ctx, "db", db)
		next(newCtx)
	}
}

// NewServer creates and configures a new Huma API server
func NewServer(dbManager *db.Manager) *Server {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("DataILama API", version.Version))

	api.UseMiddleware(DbMiddleware(dbManager))

	server := &Server{
		db:     dbManager,
		router: router,
	}

	huma.Register(api, huma.Operation{
		OperationID: "getVersion",
		Method:      http.MethodGet,
		Path:        "/version",
		Summary:     "Get the version",
		Description: "Get the version of DataILama",
	}, getVersion)

	huma.Register(api, huma.Operation{
		OperationID: "searchPaperByTitle",
		Method:      http.MethodGet,
		Path:        "/search",
		Summary:     "Search for Papers by Title",
		Description: "Search by Title. Currently the search is a case insensive excact match to a part of the title.",
	}, searchByTitle)

	return server
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
