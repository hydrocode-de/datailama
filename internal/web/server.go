package web

import (
	"net/http"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/web/api"
	"github.com/hydrocode-de/datailama/internal/web/site"
)

// Server combines both API and site functionality
type Server struct {
	DB     *db.Manager
	Router *http.ServeMux
}

// NewServer creates a new server with both API and site routes
func NewServer(dbManager *db.Manager, apiOnly bool) *Server {
	// Create main router
	router := http.NewServeMux()

	server := &Server{
		DB:     dbManager,
		Router: router,
	}

	// Register API routes (using Huma)
	api.RegisterRoutes(router, dbManager)

	// Register frontend routes
	if !apiOnly || true {
		// also handle the exact path "/" to redirect to /app/ for now
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/app/", http.StatusSeeOther)
		})
		//router.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("internal/web/site/frontend"))))
		router.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.FS(site.GetEmbedFrontend()))))
	}

	return server
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
