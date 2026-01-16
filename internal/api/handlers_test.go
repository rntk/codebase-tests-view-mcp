package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"codebase-view-mcp/internal/files"
)

func TestHandlerListFiles(t *testing.T) {
	t.Run("defaults to current directory and returns entries", func(t *testing.T) {
		baseDir := t.TempDir()

		subdir := filepath.Join(baseDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		if err := os.WriteFile(filepath.Join(baseDir, "a.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		if err := os.WriteFile(filepath.Join(baseDir, ".secret"), []byte("hidden"), 0644); err != nil {
			t.Fatalf("write hidden file: %v", err)
		}

		h := &Handler{fileService: files.NewService(baseDir)}
		req := httptest.NewRequest(http.MethodGet, "/api/files", nil)
		rr := httptest.NewRecorder()

		h.ListFiles(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.ListFilesResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Path != "." {
			t.Fatalf("path = %q, want %q", response.Path, ".")
		}

		if len(response.Files) != 2 {
			t.Fatalf("files count = %d, want %d", len(response.Files), 2)
		}

		if !response.Files[0].IsDir || response.Files[0].Name != "subdir" {
			t.Fatalf("first entry = %+v, want subdir directory", response.Files[0])
		}

		if response.Files[1].IsDir || response.Files[1].Name != "a.txt" {
			t.Fatalf("second entry = %+v, want a.txt file", response.Files[1])
		}
	})

	t.Run("returns not found for missing path", func(t *testing.T) {
		baseDir := t.TempDir()
		h := &Handler{fileService: files.NewService(baseDir)}

		req := httptest.NewRequest(http.MethodGet, "/api/files?path=missing", nil)
		rr := httptest.NewRecorder()

		h.ListFiles(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})
}
