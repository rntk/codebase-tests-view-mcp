package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codebase-view-mcp/internal/files"
	"codebase-view-mcp/internal/mcp"
	"codebase-view-mcp/internal/metadata"
)

func TestSetupRoutes(t *testing.T) {
	t.Run("creates router with all routes", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		if router == nil {
			t.Fatal("SetupRoutes returned nil")
		}

		// Test that the router handles various routes
		testRoutes := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/api/files"},
			{http.MethodGet, "/api/files/test.go"},
			{http.MethodGet, "/api/files/test.go/tests"},
			{http.MethodGet, "/api/files/test.go/suggestions"},
			{http.MethodGet, "/api/files/test.go/sources"},
			{http.MethodGet, "/api/files/test.go/comments"},
			{http.MethodPost, "/api/files/test.go/comments"},
			{http.MethodPut, "/api/files/test.go/comments/123"},
			{http.MethodDelete, "/api/files/test.go/comments/123"},
			{http.MethodPatch, "/api/files/test.go/comments/123/resolved"},
			{http.MethodPost, "/api/files/test.go/export"},
			{http.MethodPost, "/api/mcp"},
			{http.MethodGet, "/"},
		}

		for _, route := range testRoutes {
			t.Run(route.method+" "+route.path, func(t *testing.T) {
				req := httptest.NewRequest(route.method, route.path, nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				// We just check that the router doesn't panic and returns some response
				// The actual response is tested in handler tests
				_ = rr.Code
			})
		}
	})

	t.Run("router handles GET /api/files", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles GET /api/files/{path}", func(t *testing.T) {
		baseDir := t.TempDir()
		fileService := files.NewService(baseDir)
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.txt", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// File doesn't exist, should return 404
		if rr.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("router handles GET /api/files/{path}/tests", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles GET /api/files/{path}/suggestions", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/suggestions", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles GET /api/files/{path}/sources", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/sources", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles GET /api/files/{path}/comments", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/comments", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles POST /api/files/{path}/comments", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Missing body, should return 400
		if rr.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("router handles PUT /api/files/{path}/comments/{commentId}", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodPut, "/api/files/test.go/comments/123", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Missing body, should return 400
		if rr.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("router handles DELETE /api/files/{path}/comments/{commentId}", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodDelete, "/api/files/test.go/comments/123", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Comment doesn't exist, but should return 204 (no content)
		if rr.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}
	})

	t.Run("router handles PATCH /api/files/{path}/comments/{commentId}/resolved", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodPatch, "/api/files/test.go/comments/123/resolved", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Comment doesn't exist, but should return 200
		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles POST /api/files/{path}/export", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/export", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// File doesn't exist, should return 404
		if rr.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("router handles POST /api/mcp", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodPost, "/api/mcp", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Invalid JSON, should return error response
		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("router handles GET /", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Should return fallback HTML since no frontend build
		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})
}

func TestSetupRoutesWithMiddleware(t *testing.T) {
	t.Run("routes work with CORS and Logging middleware", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := NewHandler(fileService, metaStore, mcpHandler)

		router := SetupRoutes(handler)
		chain := Logging(CORS(router))

		req := httptest.NewRequest(http.MethodGet, "/api/files", nil)
		rr := httptest.NewRecorder()

		chain.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Check CORS headers
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS headers not present")
		}
	})
}
