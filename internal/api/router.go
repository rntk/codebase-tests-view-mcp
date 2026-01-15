package api

import "net/http"

// SetupRoutes configures all API routes
func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// File operations
	mux.HandleFunc("GET /api/files", h.ListFiles)
	mux.HandleFunc("GET /api/files/{path...}", func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")

		// Check if it's a test request
		if len(path) > 6 && path[len(path)-6:] == "/tests" {
			// Remove "/tests" suffix
			path = path[:len(path)-6]
			r.SetPathValue("path", path)
			h.GetTests(w, r)
		} else {
			h.GetFile(w, r)
		}
	})

	// MCP endpoint
	mux.HandleFunc("POST /api/mcp", h.HandleMCP)

	// Serve static files (will be implemented later with embed)
	mux.HandleFunc("GET /", h.ServeStatic)

	return mux
}
