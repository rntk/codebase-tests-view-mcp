package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"codebase-view-mcp/internal/files"
	"codebase-view-mcp/internal/metadata"
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

func TestHandlerGetFile(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		baseDir := t.TempDir()
		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/", nil)
		rr := httptest.NewRecorder()

		h.GetFile(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns not found for missing file", func(t *testing.T) {
		baseDir := t.TempDir()
		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/missing.txt", nil)
		req.SetPathValue("path", "missing.txt")
		rr := httptest.NewRecorder()

		h.GetFile(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("returns file content with metadata", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(baseDir, "hello.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("hello.txt", []metadata.TestReference{
			{
				TestFile: "hello_test.go",
				TestName: "TestHello",
				LineRange: metadata.LineRange{
					Start: 10,
					End:   20,
				},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/hello.txt", nil)
		req.SetPathValue("path", "hello.txt")
		rr := httptest.NewRecorder()

		h.GetFile(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.FileResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.File.Path != "hello.txt" {
			t.Fatalf("path = %q, want %q", response.File.Path, "hello.txt")
		}

		if response.File.Name != "hello.txt" {
			t.Fatalf("name = %q, want %q", response.File.Name, "hello.txt")
		}

		if response.File.Content != "hello" {
			t.Fatalf("content = %q, want %q", response.File.Content, "hello")
		}

		if response.File.MimeType != "text/plain" && response.File.MimeType != "text/plain; charset=utf-8" {
			t.Fatalf("mimeType = %q, want %q or %q", response.File.MimeType, "text/plain", "text/plain; charset=utf-8")
		}

		if response.File.Metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(response.File.Metadata.Tests) != 1 {
			t.Fatalf("tests count = %d, want %d", len(response.File.Metadata.Tests), 1)
		}

		if response.File.Metadata.Tests[0].TestFile != "hello_test.go" {
			t.Fatalf("testFile = %q, want %q", response.File.Metadata.Tests[0].TestFile, "hello_test.go")
		}
	})
}
