package api

import "net/http"

// SetupRoutes configures all API routes
func SetupRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// File operations
	mux.HandleFunc("GET /api/files", h.ListFiles)
	mux.HandleFunc("GET /api/files/{path...}", h.GetFileOrTests)

	// Comment operations
	mux.HandleFunc("GET /api/files/{path}/comments", h.GetComments)
	mux.HandleFunc("POST /api/files/{path}/comments", h.CreateComment)
	mux.HandleFunc("PUT /api/files/{path}/comments/{commentId}", h.UpdateComment)
	mux.HandleFunc("DELETE /api/files/{path}/comments/{commentId}", h.DeleteComment)
	mux.HandleFunc("PATCH /api/files/{path}/comments/{commentId}/resolved", h.ToggleCommentResolved)

	// Export for AI agents
	mux.HandleFunc("POST /api/files/{path}/export", h.ExportContext)

	// MCP endpoint
	mux.HandleFunc("POST /api/mcp", h.HandleMCP)

	// Serve static files (will be implemented later with embed)
	mux.HandleFunc("GET /", h.ServeStatic)

	return mux
}
