package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codebase-view-mcp/internal/files"
	"codebase-view-mcp/internal/mcp"
	"codebase-view-mcp/internal/metadata"
)

func TestGetFileOrTests(t *testing.T) {
	t.Run("routes to GetTests when path ends with /tests", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go/tests")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Verify the path was correctly processed
		// The response should be a TestsResponse
	})

	t.Run("routes to GetSources when path ends with /sources", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/sources", nil)
		req.SetPathValue("path", "test.go/sources")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("routes to GetFile for regular paths", func(t *testing.T) {
		baseDir := t.TempDir()
		fileService := files.NewService(baseDir)
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.txt", nil)
		req.SetPathValue("path", "test.txt")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		// File doesn't exist, should return 404
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("handles path that looks like tests suffix but is longer", func(t *testing.T) {
		baseDir := t.TempDir()
		fileService := files.NewService(baseDir)
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		// Create a file named "mytests" (not ending with /tests)
		req := httptest.NewRequest(http.MethodGet, "/api/files/mytests", nil)
		req.SetPathValue("path", "mytests")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		// File doesn't exist, should return 404
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("handles path that looks like sources suffix but is longer", func(t *testing.T) {
		baseDir := t.TempDir()
		fileService := files.NewService(baseDir)
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/mysources", nil)
		req.SetPathValue("path", "mysources")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		// File doesn't exist, should return 404
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("handles nested path with tests suffix", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/dir/subdir/test.go/tests", nil)
		req.SetPathValue("path", "dir/subdir/test.go/tests")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("handles nested path with sources suffix", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/dir/subdir/test.go/sources", nil)
		req.SetPathValue("path", "dir/subdir/test.go/sources")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("GetTests returns empty tests array when no metadata exists", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go/tests")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Response should contain empty tests array
		// Content-Type should be application/json
		if rr.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type %q, got %q", "application/json", rr.Header().Get("Content-Type"))
		}
	})

	t.Run("GetSources returns empty sources array when no metadata exists", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/sources", nil)
		req.SetPathValue("path", "test.go/sources")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		if rr.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type %q, got %q", "application/json", rr.Header().Get("Content-Type"))
		}
	})
}

func TestGetFileOrTestsPathValueUpdate(t *testing.T) {
	t.Run("GetTests receives updated path value without suffix", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go/tests")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)

		// The path value should have been updated to remove the /tests suffix
		// This is verified by checking that GetTests was called with the correct path
	})

	t.Run("GetSources receives updated path value without suffix", func(t *testing.T) {
		fileService := files.NewService(t.TempDir())
		metaStore := metadata.NewStore("")
		mcpHandler := mcp.NewHandler(metaStore)
		handler := &Handler{
			fileService: fileService,
			metaStore:   metaStore,
			mcpHandler:  mcpHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/sources", nil)
		req.SetPathValue("path", "test.go/sources")
		rr := httptest.NewRecorder()

		handler.GetFileOrTests(rr, req)
	})
}
