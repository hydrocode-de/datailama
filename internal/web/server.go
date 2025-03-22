package web

import (
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/web/api"
	"github.com/hydrocode-de/datailama/internal/web/site"
	"github.com/urfave/cli/v2"
)

// Server combines both API and site functionality
type Server struct {
	DB     *db.Manager
	Router *http.ServeMux
}

// SPAHandler handles all requests for the SPA by returning the index.html file
type SPAHandler struct {
	fileSystem http.FileSystem
}

// ServeHTTP implements http.Handler and serves the SPA
func (h SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean path to prevent directory traversal
	urlPath := path.Clean(r.URL.Path)

	// Strip /app/ prefix
	urlPath = strings.TrimPrefix(urlPath, "/app/")
	if urlPath == "" {
		urlPath = "."
	}

	// Try to open the file
	file, err := h.fileSystem.Open(urlPath)
	if err != nil {
		// If file doesn't exist, serve index.html for SPA routing
		indexFile, err := h.fileSystem.Open("index.html")
		if err != nil {
			http.Error(w, "Index file not found", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		// Set proper content type for HTML
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
		return
	}
	defer file.Close()

	// Check if it's a directory
	fi, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If it's a directory, look for index.html
	if fi.IsDir() {
		indexFile, err := h.fileSystem.Open(path.Join(urlPath, "index.html"))
		if err != nil {
			// If no index.html in directory, serve the main index.html for SPA routing
			indexFile, err := h.fileSystem.Open("index.html")
			if err != nil {
				http.Error(w, "Index file not found", http.StatusInternalServerError)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
			return
		}
		defer indexFile.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
		return
	}

	// For non-HTML files, serve the file directly
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), file)
}

// NewServer creates a new server with both API and site routes
func NewServer(dbManager *db.Manager, apiOnly bool, cliCtx *cli.Context) *Server {
	// Create main router
	router := http.NewServeMux()

	server := &Server{
		DB:     dbManager,
		Router: router,
	}

	// Register API routes (using Huma)
	api.RegisterRoutes(router, dbManager, cliCtx)

	// Register frontend routes
	if !apiOnly || true {
		// also handle the exact path "/" to redirect to /app/ for now
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/app/", http.StatusSeeOther)
		})

		// Use our custom SPA handler for all /app/ routes
		frontendFS := site.GetEmbedFrontend()
		router.Handle("/app/", http.StripPrefix("/app", SPAHandler{fileSystem: http.FS(frontendFS)}))
	}

	return server
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
