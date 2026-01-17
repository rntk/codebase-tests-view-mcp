package api

import "net/http"

// SetupRoutes configures all API routes
func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// File operations
	mux.HandleFunc("GET /api/files", h.ListFiles)
	mux.HandleFunc("GET /api/files/{path...}", h.GetFileOrTests)

	// MCP endpoint
	mux.HandleFunc("POST /api/mcp", h.HandleMCP)

	// Serve static files (will be implemented later with embed)
	mux.HandleFunc("GET /", h.ServeStatic)

	return mux
}
