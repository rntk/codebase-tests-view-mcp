package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetStaticFS(t *testing.T) {
	t.Run("returns filesystem without error", func(t *testing.T) {
		fs, err := GetStaticFS()
		if err != nil {
			t.Fatalf("GetStaticFS returned error: %v", err)
		}

		if fs == nil {
			t.Fatal("GetStaticFS returned nil filesystem")
		}
	})
}

func TestServeStaticFiles(t *testing.T) {
	t.Run("returns fallback handler when no frontend build exists", func(t *testing.T) {
		handler := ServeStaticFiles()
		if handler == nil {
			t.Fatal("ServeStaticFiles returned nil")
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		body := rr.Body.String()
		if !strings.Contains(body, "Codebase Test Viewer") {
			t.Error("response does not contain expected fallback HTML")
		}

		if !strings.Contains(body, "Frontend build not found") {
			t.Error("response does not contain frontend build not found message")
		}

		if !strings.Contains(body, "npm install") {
			t.Error("response does not contain npm install instructions")
		}

		if !strings.Contains(body, "/api/*") {
			t.Error("response does not contain API path information")
		}
	})

	t.Run("fallback handler returns HTML content type", func(t *testing.T) {
		handler := ServeStaticFiles()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		contentType := rr.Header().Get("Content-Type")
		if !strings.HasPrefix(contentType, "text/html") {
			t.Errorf("expected Content-Type to start with %q, got %q", "text/html", contentType)
		}
	})

	t.Run("fallback handler works for any path", func(t *testing.T) {
		handler := ServeStaticFiles()
		paths := []string{
			"/",
			"/index.html",
			"/app",
			"/files",
			"/some/deep/path",
		}

		for _, path := range paths {
			t.Run(path, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rr := httptest.NewRecorder()

				handler.ServeHTTP(rr, req)

				if rr.Code != http.StatusOK {
					t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
				}
			})
		}
	})
}

func TestFallbackHTML(t *testing.T) {
	t.Run("contains expected title", func(t *testing.T) {
		if !strings.Contains(fallbackHTML, "<title>Codebase Test Viewer</title>") {
			t.Error("fallback HTML does not contain expected title")
		}
	})

	t.Run("contains expected heading", func(t *testing.T) {
		if !strings.Contains(fallbackHTML, "<h1>Codebase Test Viewer</h1>") {
			t.Error("fallback HTML did not contain expected heading")
		}
	})

	t.Run("contains build instructions", func(t *testing.T) {
		if !strings.Contains(fallbackHTML, "cd frontend && npm install && npm run build") {
			t.Error("fallback HTML did not contain build instructions")
		}
	})

	t.Run("contains API information", func(t *testing.T) {
		if !strings.Contains(fallbackHTML, "API is running at") {
			t.Error("fallback HTML did not contain API information")
		}
	})

	t.Run("is valid HTML structure", func(t *testing.T) {
		if !strings.Contains(fallbackHTML, "<!DOCTYPE html>") {
			t.Error("fallback HTML missing DOCTYPE")
		}

		if !strings.Contains(fallbackHTML, "<html>") {
			t.Error("fallback HTML missing html tag")
		}

		if !strings.Contains(fallbackHTML, "<head>") {
			t.Error("fallback HTML missing head tag")
		}

		if !strings.Contains(fallbackHTML, "<body>") {
			t.Error("fallback HTML missing body tag")
		}
	})
}

func TestStaticHandlerWithDifferentMethods(t *testing.T) {
	t.Run("handles GET requests", func(t *testing.T) {
		handler := ServeStaticFiles()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("handles HEAD requests", func(t *testing.T) {
		handler := ServeStaticFiles()
		req := httptest.NewRequest(http.MethodHead, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// HEAD requests should return headers but no body
		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("handles POST requests with fallback", func(t *testing.T) {
		handler := ServeStaticFiles()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// POST to static handler should still return fallback
		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})
}

func TestGetStaticFSSubDirectory(t *testing.T) {
	t.Run("filesystem is properly configured", func(t *testing.T) {
		fs, err := GetStaticFS()
		if err != nil {
			t.Fatalf("GetStaticFS returned error: %v", err)
		}

		// The embedded FS should exist but index.html won't be found since no build exists
		// We just verify the filesystem is not nil
		if fs == nil {
			t.Fatal("expected non-nil filesystem")
		}
	})
}

func TestStaticHandlerIntegration(t *testing.T) {
	t.Run("works in full middleware chain", func(t *testing.T) {
		fileService := NewHandler(
			nil, // fileService - not used for static
			nil, // metaStore - not used for static
			nil, // mcpHandler - not used for static
		)

		// Create a handler that uses ServeStatic
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fileService.ServeStatic(w, r)
		})

		chain := Logging(CORS(handler))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		chain.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Check CORS headers are present
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS headers not present")
		}
	})
}
