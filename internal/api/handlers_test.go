package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestHandlerGetFileOrTests(t *testing.T) {
	t.Run("routes tests requests to GetTests", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/hello.txt/tests", nil)
		req.SetPathValue("path", "hello.txt/tests")
		rr := httptest.NewRecorder()

		h.GetFileOrTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.SourceFile != "hello.txt" {
			t.Fatalf("sourceFile = %q, want %q", response.SourceFile, "hello.txt")
		}
	})

	t.Run("routes file requests to GetFile", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(baseDir, "hello.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/hello.txt", nil)
		req.SetPathValue("path", "hello.txt")
		rr := httptest.NewRecorder()

		h.GetFileOrTests(rr, req)

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
	})
}

func TestHandlerGetTests(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files//tests", nil)
		req.SetPathValue("path", "")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns empty tests array when no metadata exists", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.SourceFile != "test.go" {
			t.Fatalf("sourceFile = %q, want %q", response.SourceFile, "test.go")
		}

		if len(response.Tests) != 0 {
			t.Fatalf("expected 0 tests, got %d", len(response.Tests))
		}
	})

	t.Run("returns tests with metadata", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "test_test.go")
		testContent := `package main

func TestSomething(t *testing.T) {
	input := "hello"
	result := process(input)
	if result != "world" {
		t.Errorf("expected world, got %s", result)
	}
}`
		if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("test.go", []metadata.TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				FunctionName: "Process",
				Comment:      "Tests the process function",
				LineRange:    metadata.LineRange{Start: 3, End: 9},
				CoveredLines: metadata.LineRange{Start: 1, End: 10},
				InputLines:   metadata.LineRange{Start: 4, End: 4},
				OutputLines:  metadata.LineRange{Start: 5, End: 7},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(response.Tests))
		}

		test := response.Tests[0]
		if test.TestName != "TestSomething" {
			t.Errorf("expected testName %q, got %q", "TestSomething", test.TestName)
		}

		if test.FunctionName != "Process" {
			t.Errorf("expected functionName %q, got %q", "Process", test.FunctionName)
		}

		if test.Comment != "Tests the process function" {
			t.Errorf("expected comment %q, got %q", "Tests the process function", test.Comment)
		}

		if test.InputData == "" {
			t.Error("expected inputData to be populated")
		}

		if test.ExpectedOutput == "" {
			t.Error("expected expectedOutput to be populated")
		}
	})

	t.Run("handles missing test file gracefully", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("test.go", []metadata.TestReference{
			{
				TestFile:     "missing_test.go",
				TestName:     "TestSomething",
				FunctionName: "Process",
				Comment:      "Tests the process function",
				LineRange:    metadata.LineRange{Start: 3, End: 9},
				CoveredLines: metadata.LineRange{Start: 1, End: 10},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(t.TempDir()),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(response.Tests))
		}

		// Content should be empty since test file doesn't exist
		if response.Tests[0].Content != "" {
			t.Error("expected empty content for missing test file")
		}
	})

	t.Run("filters tests by functionName query param", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("test.go", []metadata.TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestFoo",
				FunctionName: "Foo",
				Comment:      "Tests Foo",
				LineRange:    metadata.LineRange{Start: 1, End: 5},
				CoveredLines: metadata.LineRange{Start: 1, End: 3},
			},
			{
				TestFile:     "test_test.go",
				TestName:     "TestBar",
				FunctionName: "Bar",
				Comment:      "Tests Bar",
				LineRange:    metadata.LineRange{Start: 6, End: 10},
				CoveredLines: metadata.LineRange{Start: 4, End: 6},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(t.TempDir()),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests?functionName=Foo", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Tests) != 1 {
			t.Fatalf("expected 1 test after filter, got %d", len(response.Tests))
		}
		if response.Tests[0].FunctionName != "Foo" {
			t.Errorf("expected functionName %q, got %q", "Foo", response.Tests[0].FunctionName)
		}
	})

	t.Run("returns empty tests when functionName filter matches nothing", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("test.go", []metadata.TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestFoo",
				FunctionName: "Foo",
				Comment:      "Tests Foo",
				LineRange:    metadata.LineRange{Start: 1, End: 5},
				CoveredLines: metadata.LineRange{Start: 1, End: 3},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(t.TempDir()),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/tests?functionName=NonExistent", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetTests(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Tests) != 0 {
			t.Fatalf("expected 0 tests for non-matching filter, got %d", len(response.Tests))
		}
	})
}

func TestHandlerGetSuggestions(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files//suggestions", nil)
		req.SetPathValue("path", "")
		rr := httptest.NewRecorder()

		h.GetSuggestions(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns empty suggestions when no metadata exists", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/suggestions", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetSuggestions(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.SuggestionsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.SourceFile != "test.go" {
			t.Fatalf("sourceFile = %q, want %q", response.SourceFile, "test.go")
		}

		if len(response.Suggestions) != 0 {
			t.Fatalf("expected 0 suggestions, got %d", len(response.Suggestions))
		}
	})

	t.Run("returns suggestions with metadata", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.AddSuggestions("test.go", []metadata.TestSuggestion{
			{
				FunctionName:  "Process",
				TargetLines:   metadata.LineRange{Start: 10, End: 20},
				Reason:        "Missing test for error case",
				SuggestedName: "TestProcessError",
				TestSkeleton:  "func TestProcessError(t *testing.T) {}",
				Priority:      "high",
			},
		}); err != nil {
			t.Fatalf("add suggestions: %v", err)
		}

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/suggestions", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetSuggestions(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.SuggestionsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Suggestions) != 1 {
			t.Fatalf("expected 1 suggestion, got %d", len(response.Suggestions))
		}

		suggestion := response.Suggestions[0]
		if suggestion.SuggestedName != "TestProcessError" {
			t.Errorf("expected suggestedName %q, got %q", "TestProcessError", suggestion.SuggestedName)
		}

		if suggestion.Priority != "high" {
			t.Errorf("expected priority %q, got %q", "high", suggestion.Priority)
		}
	})
}

func TestHandlerGetSources(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files//sources", nil)
		req.SetPathValue("path", "")
		rr := httptest.NewRecorder()

		h.GetSources(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns empty sources when test file has no references", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test_test.go/sources", nil)
		req.SetPathValue("path", "test_test.go")
		rr := httptest.NewRecorder()

		h.GetSources(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestFileResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.TestFile != "test_test.go" {
			t.Fatalf("testFile = %q, want %q", response.TestFile, "test_test.go")
		}

		if len(response.Sources) != 0 {
			t.Fatalf("expected 0 sources, got %d", len(response.Sources))
		}
	})

	t.Run("returns source references for test file", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("source.go", []metadata.TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				FunctionName: "Process",
				Comment:      "Tests the process function",
				LineRange:    metadata.LineRange{Start: 3, End: 9},
				CoveredLines: metadata.LineRange{Start: 1, End: 10},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test_test.go/sources", nil)
		req.SetPathValue("path", "test_test.go")
		rr := httptest.NewRecorder()

		h.GetSources(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.TestFileResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Sources) != 1 {
			t.Fatalf("expected 1 source, got %d", len(response.Sources))
		}

		source := response.Sources[0]
		if source.SourceFile != "source.go" {
			t.Errorf("expected sourceFile %q, got %q", "source.go", source.SourceFile)
		}

		if source.FunctionName != "Process" {
			t.Errorf("expected functionName %q, got %q", "Process", source.FunctionName)
		}
	})
}

func TestHandlerGetComments(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files//comments", nil)
		req.SetPathValue("path", "")
		rr := httptest.NewRecorder()

		h.GetComments(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns empty comments when no metadata exists", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/comments", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetComments(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.CommentsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.SourceFile != "test.go" {
			t.Fatalf("sourceFile = %q, want %q", response.SourceFile, "test.go")
		}

		if len(response.Comments) != 0 {
			t.Fatalf("expected 0 comments, got %d", len(response.Comments))
		}
	})

	t.Run("returns comments with metadata", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		if _, err := metaStore.AddComment("test.go", files.Comment{
			Line:    10,
			Content: "This needs improvement",
		}); err != nil {
			t.Fatalf("add comment: %v", err)
		}

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/files/test.go/comments", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.GetComments(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.CommentsResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(response.Comments))
		}

		comment := response.Comments[0]
		if comment.Line != 10 {
			t.Errorf("expected line %d, got %d", 10, comment.Line)
		}

		if comment.Content != "This needs improvement" {
			t.Errorf("expected content %q, got %q", "This needs improvement", comment.Content)
		}
	})
}

func TestHandlerCreateComment(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files//comments", nil)
		req.SetPathValue("path", "")
		req.Body = nil
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request for invalid JSON", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", nil)
		req.SetPathValue("path", "test.go")
		req.Body = http.MaxBytesReader(nil, req.Body, 1)
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when line is less than 1", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		body := `{"line": 0, "content": "test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when content is empty", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		body := `{"line": 10, "content": ""}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when content is whitespace only", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		body := `{"line": 10, "content": "   "}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("creates comment successfully", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		h := &Handler{
			metaStore: metaStore,
		}

		body := `{"line": 10, "content": "This needs improvement"}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
		}

		var response files.CommentResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Comment.Line != 10 {
			t.Errorf("expected line %d, got %d", 10, response.Comment.Line)
		}

		if response.Comment.Content != "This needs improvement" {
			t.Errorf("expected content %q, got %q", "This needs improvement", response.Comment.Content)
		}

		if response.Comment.ID == "" {
			t.Error("expected comment ID to be generated")
		}
	})

	t.Run("trims whitespace from content", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		h := &Handler{
			metaStore: metaStore,
		}

		body := `{"line": 10, "content": "  trimmed content  "}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/comments", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.CreateComment(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
		}

		var response files.CommentResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.Comment.Content != "trimmed content" {
			t.Errorf("expected trimmed content %q, got %q", "trimmed content", response.Comment.Content)
		}
	})
}

func TestHandlerUpdateComment(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPut, "/api/files//comments/123", nil)
		req.SetPathValue("path", "")
		req.SetPathValue("commentId", "123")
		rr := httptest.NewRecorder()

		h.UpdateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when commentId missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPut, "/api/files/test.go/comments/", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "")
		rr := httptest.NewRecorder()

		h.UpdateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request for invalid JSON", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPut, "/api/files/test.go/comments/123", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "123")
		rr := httptest.NewRecorder()

		h.UpdateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when content is empty", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		body := `{"content": ""}`
		req := httptest.NewRequest(http.MethodPut, "/api/files/test.go/comments/123", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "123")
		rr := httptest.NewRecorder()

		h.UpdateComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("updates comment successfully", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		created, _ := metaStore.AddComment("test.go", files.Comment{
			Line:    10,
			Content: "Original content",
		})

		h := &Handler{
			metaStore: metaStore,
		}

		body := `{"content": "Updated content"}`
		req := httptest.NewRequest(http.MethodPut, "/api/files/test.go/comments/"+created.ID, strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", created.ID)
		rr := httptest.NewRecorder()

		h.UpdateComment(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Verify the comment was updated
		comments := metaStore.GetComments("test.go")
		if len(comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(comments))
		}

		if comments[0].Content != "Updated content" {
			t.Errorf("expected content %q, got %q", "Updated content", comments[0].Content)
		}
	})
}

func TestHandlerDeleteComment(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/files//comments/123", nil)
		req.SetPathValue("path", "")
		req.SetPathValue("commentId", "123")
		rr := httptest.NewRecorder()

		h.DeleteComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when commentId missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/files/test.go/comments/", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "")
		rr := httptest.NewRecorder()

		h.DeleteComment(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("deletes comment successfully", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		created, _ := metaStore.AddComment("test.go", files.Comment{
			Line:    10,
			Content: "Content to delete",
		})

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/files/test.go/comments/"+created.ID, nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", created.ID)
		rr := httptest.NewRecorder()

		h.DeleteComment(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}

		// Verify the comment was deleted
		comments := metaStore.GetComments("test.go")
		if len(comments) != 0 {
			t.Fatalf("expected 0 comments after delete, got %d", len(comments))
		}
	})

	t.Run("handles deleting non-existent comment gracefully", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/files/test.go/comments/nonexistent", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "nonexistent")
		rr := httptest.NewRecorder()

		h.DeleteComment(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}
	})
}

func TestHandlerToggleCommentResolved(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPatch, "/api/files//comments/123/resolved", nil)
		req.SetPathValue("path", "")
		req.SetPathValue("commentId", "123")
		rr := httptest.NewRecorder()

		h.ToggleCommentResolved(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns bad request when commentId missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPatch, "/api/files/test.go/comments//resolved", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "")
		rr := httptest.NewRecorder()

		h.ToggleCommentResolved(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("toggles comment from unresolved to resolved", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		created, _ := metaStore.AddComment("test.go", files.Comment{
			Line:     10,
			Content:  "To be resolved",
			Resolved: false,
		})

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodPatch, "/api/files/test.go/comments/"+created.ID+"/resolved", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", created.ID)
		rr := httptest.NewRecorder()

		h.ToggleCommentResolved(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Verify the comment was toggled
		comments := metaStore.GetComments("test.go")
		if len(comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(comments))
		}

		if !comments[0].Resolved {
			t.Error("expected comment to be resolved")
		}
	})

	t.Run("toggles comment from resolved to unresolved", func(t *testing.T) {
		metaStore := metadata.NewStore("")
		created, _ := metaStore.AddComment("test.go", files.Comment{
			Line:     10,
			Content:  "Already resolved",
			Resolved: true,
		})

		h := &Handler{
			metaStore: metaStore,
		}

		req := httptest.NewRequest(http.MethodPatch, "/api/files/test.go/comments/"+created.ID+"/resolved", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", created.ID)
		rr := httptest.NewRecorder()

		h.ToggleCommentResolved(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		// Verify the comment was toggled
		comments := metaStore.GetComments("test.go")
		if len(comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(comments))
		}

		if comments[0].Resolved {
			t.Error("expected comment to be unresolved")
		}
	})

	t.Run("handles non-existent comment gracefully", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPatch, "/api/files/test.go/comments/nonexistent/resolved", nil)
		req.SetPathValue("path", "test.go")
		req.SetPathValue("commentId", "nonexistent")
		rr := httptest.NewRecorder()

		h.ToggleCommentResolved(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})
}

func TestHandlerExportContext(t *testing.T) {
	t.Run("returns bad request when path missing", func(t *testing.T) {
		h := &Handler{
			metaStore: metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files//export", nil)
		req.SetPathValue("path", "")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns not found for missing file", func(t *testing.T) {
		h := &Handler{
			fileService: files.NewService(t.TempDir()),
			metaStore:   metadata.NewStore(""),
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files/missing.go/export", nil)
		req.SetPathValue("path", "missing.go")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("exports context with default options", func(t *testing.T) {
		baseDir := t.TempDir()
		fileContent := `package main

func Process(input string) string {
	return "result"
}`
		if err := os.WriteFile(filepath.Join(baseDir, "test.go"), []byte(fileContent), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		metaStore.AddComment("test.go", files.Comment{
			Line:    3,
			Content: "This function needs tests",
		})

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/export", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.ExportContextResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if response.SourceFile != "test.go" {
			t.Errorf("expected sourceFile %q, got %q", "test.go", response.SourceFile)
		}

		if len(response.CodeContext) == 0 {
			t.Error("expected code context blocks")
		}

		if response.Formatted == "" {
			t.Error("expected formatted output")
		}
	})

	t.Run("exports context with custom options", func(t *testing.T) {
		baseDir := t.TempDir()
		fileContent := `package main

func Process(input string) string {
	return "result"
}`
		if err := os.WriteFile(filepath.Join(baseDir, "test.go"), []byte(fileContent), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		metaStore.AddComment("test.go", files.Comment{
			Line:    3,
			Content: "This function needs tests",
		})

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		body := `{"includeTests": false, "includeSuggestions": false, "contextLines": 2}`
		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/export", strings.NewReader(body))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.ExportContextResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		// Tests and suggestions should be empty
		if response.Tests != nil && len(response.Tests) > 0 {
			t.Error("expected no tests when includeTests is false")
		}

		if response.Suggestions != nil && len(response.Suggestions) > 0 {
			t.Error("expected no suggestions when includeSuggestions is false")
		}
	})

	t.Run("handles invalid JSON body gracefully", func(t *testing.T) {
		baseDir := t.TempDir()
		fileContent := `package main`
		if err := os.WriteFile(filepath.Join(baseDir, "test.go"), []byte(fileContent), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/export", strings.NewReader("invalid"))
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		// Should use defaults and succeed
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}
	})

	t.Run("skips resolved comments in export", func(t *testing.T) {
		baseDir := t.TempDir()
		fileContent := `package main

func Process(input string) string {
	return "result"
}`
		if err := os.WriteFile(filepath.Join(baseDir, "test.go"), []byte(fileContent), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		metaStore.AddComment("test.go", files.Comment{
			Line:     3,
			Content:  "Resolved comment",
			Resolved: true,
		})
		metaStore.AddComment("test.go", files.Comment{
			Line:     4,
			Content:  "Unresolved comment",
			Resolved: false,
		})

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodPost, "/api/files/test.go/export", nil)
		req.SetPathValue("path", "test.go")
		rr := httptest.NewRecorder()

		h.ExportContext(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.ExportContextResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		// Should only have 1 block (unresolved comment)
		if len(response.CodeContext) != 1 {
			t.Errorf("expected 1 code context block, got %d", len(response.CodeContext))
		}
	})
}

func TestBuildFormattedExport(t *testing.T) {
	t.Run("formats export with comments", func(t *testing.T) {
		data := files.ExportContextResponse{
			SourceFile: "test.go",
			CodeContext: []files.CodeContextBlock{
				{
					LineRange: files.LineRange{Start: 1, End: 5},
					Code:      "1: package main\n2: \n3: func main() {}\n",
					Comments: []files.Comment{
						{Line: 3, Content: "Add tests here"},
					},
				},
			},
		}

		formatted := buildFormattedExport(data)

		if !strings.Contains(formatted, "# Code Review Export") {
			t.Error("formatted output missing header")
		}

		if !strings.Contains(formatted, "**File:** `test.go`") {
			t.Error("formatted output missing file path")
		}

		if !strings.Contains(formatted, "## Comments and Code Context") {
			t.Error("formatted output missing comments section")
		}

		if !strings.Contains(formatted, "Add tests here") {
			t.Error("formatted output missing comment content")
		}
	})

	t.Run("formats export with tests", func(t *testing.T) {
		data := files.ExportContextResponse{
			SourceFile: "test.go",
			Tests: []files.TestDetail{
				{
					TestName:     "TestMain",
					FunctionName: "Main",
					TestFile:     "test_test.go",
					Comment:      "Tests the main function",
					LineRange:    files.LineRange{Start: 1, End: 10},
					CoveredLines: files.LineRange{Start: 1, End: 5},
				},
			},
		}

		formatted := buildFormattedExport(data)

		if !strings.Contains(formatted, "## Related Tests") {
			t.Error("formatted output missing tests section")
		}

		if !strings.Contains(formatted, "TestMain") {
			t.Error("formatted output missing test name")
		}
	})

	t.Run("formats export with suggestions", func(t *testing.T) {
		data := files.ExportContextResponse{
			SourceFile: "test.go",
			Suggestions: []files.TestSuggestion{
				{
					SuggestedName: "TestEdgeCase",
					Priority:      "high",
					Reason:        "Missing edge case test",
					TargetLines:   files.LineRange{Start: 10, End: 20},
					TestSkeleton:  "func TestEdgeCase(t *testing.T) {}",
				},
			},
		}

		formatted := buildFormattedExport(data)

		if !strings.Contains(formatted, "## Test Suggestions") {
			t.Error("formatted output missing suggestions section")
		}

		if !strings.Contains(formatted, "TestEdgeCase") {
			t.Error("formatted output missing suggestion name")
		}

		if !strings.Contains(formatted, "Priority: high") {
			t.Error("formatted output missing priority")
		}
	})

	t.Run("formats export with all sections", func(t *testing.T) {
		data := files.ExportContextResponse{
			SourceFile: "test.go",
			CodeContext: []files.CodeContextBlock{
				{
					LineRange: files.LineRange{Start: 1, End: 5},
					Code:      "code here",
					Comments:  []files.Comment{{Line: 3, Content: "comment"}},
				},
			},
			Tests: []files.TestDetail{
				{TestName: "TestSomething", FunctionName: "Something", TestFile: "test_test.go", LineRange: files.LineRange{Start: 1, End: 10}, CoveredLines: files.LineRange{Start: 1, End: 5}},
			},
			Suggestions: []files.TestSuggestion{
				{SuggestedName: "TestMissing", Priority: "medium", Reason: "missing", TargetLines: files.LineRange{Start: 1, End: 5}, TestSkeleton: "func TestMissing(t *testing.T) {}"},
			},
		}

		formatted := buildFormattedExport(data)

		// Check all sections are present
		sections := []string{
			"# Code Review Export",
			"## Comments and Code Context",
			"## Related Tests",
			"## Test Suggestions",
		}

		for _, section := range sections {
			if !strings.Contains(formatted, section) {
				t.Errorf("formatted output missing section: %s", section)
			}
		}
	})
}

func TestExtractLines(t *testing.T) {
	t.Run("extracts lines correctly", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3", "line4", "line5"}

		result := extractLines(lines, 2, 4)

		expected := "line2\nline3\nline4"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("handles single line", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3"}

		result := extractLines(lines, 2, 2)

		if result != "line2" {
			t.Errorf("expected %q, got %q", "line2", result)
		}
	})

	t.Run("returns empty for invalid start", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3"}

		result := extractLines(lines, 0, 2)

		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("returns empty for start greater than end", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3"}

		result := extractLines(lines, 3, 1)

		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("returns empty for start beyond lines", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3"}

		result := extractLines(lines, 10, 15)

		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("returns empty for end beyond lines", func(t *testing.T) {
		lines := []string{"line1", "line2", "line3"}

		result := extractLines(lines, 1, 10)

		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("handles empty lines slice", func(t *testing.T) {
		lines := []string{}

		result := extractLines(lines, 1, 1)

		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})
}

func TestMetadataIssueHandlers(t *testing.T) {
	t.Run("lists invalid metadata issues", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(baseDir, "valid.go"), []byte("package main"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("/app/invalid.go", []metadata.TestReference{
			{
				TestFile:     "missing_test.go",
				TestName:     "TestBroken",
				FunctionName: "Broken",
				Comment:      "broken metadata",
				LineRange:    metadata.LineRange{Start: 1, End: 2},
				CoveredLines: metadata.LineRange{Start: 1, End: 1},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		req := httptest.NewRequest(http.MethodGet, "/api/metadata/issues", nil)
		rr := httptest.NewRecorder()

		h.GetMetadataIssues(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
		}

		var response files.MetadataIssuesResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(response.Issues) != 1 {
			t.Fatalf("expected 1 issue, got %d", len(response.Issues))
		}
		if response.Issues[0].SourceValid {
			t.Fatal("expected invalid source entry")
		}
		if len(response.Issues[0].InvalidTestIssues) != 1 {
			t.Fatalf("expected 1 invalid test issue, got %d", len(response.Issues[0].InvalidTestIssues))
		}
	})

	t.Run("updates source path to canonical relative value", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(baseDir, "pkg"), 0755); err != nil {
			t.Fatalf("mkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "pkg", "source.go"), []byte("package pkg"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("/app/pkg/source.go", []metadata.TestReference{
			{
				TestFile:     "pkg/source_test.go",
				TestName:     "TestSource",
				FunctionName: "Source",
				Comment:      "renamed source",
				LineRange:    metadata.LineRange{Start: 1, End: 2},
				CoveredLines: metadata.LineRange{Start: 1, End: 1},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		body := `{"oldPath":"/app/pkg/source.go","newPath":"pkg/source.go"}`
		req := httptest.NewRequest(http.MethodPut, "/api/metadata/source-path", strings.NewReader(body))
		rr := httptest.NewRecorder()

		h.UpdateSourcePath(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
		}
		if metaStore.GetTestMetadata("/app/pkg/source.go") != nil {
			t.Fatal("expected old source path to be removed")
		}
		if metaStore.GetTestMetadata("pkg/source.go") == nil {
			t.Fatal("expected canonical source path to exist")
		}
	})

	t.Run("updates and deletes invalid test paths", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(baseDir, "pkg"), 0755); err != nil {
			t.Fatalf("mkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "pkg", "source.go"), []byte("package pkg"), 0644); err != nil {
			t.Fatalf("write source file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "pkg", "source_test.go"), []byte("package pkg"), 0644); err != nil {
			t.Fatalf("write test file: %v", err)
		}

		metaStore := metadata.NewStore("")
		if err := metaStore.SetTestMetadata("pkg/source.go", []metadata.TestReference{
			{
				TestFile:     "/app/pkg/source_test.go",
				TestName:     "TestSource",
				FunctionName: "Source",
				Comment:      "broken test path",
				LineRange:    metadata.LineRange{Start: 1, End: 2},
				CoveredLines: metadata.LineRange{Start: 1, End: 1},
			},
			{
				TestFile:     "missing_test.go",
				TestName:     "TestMissing",
				FunctionName: "Source",
				Comment:      "missing test path",
				LineRange:    metadata.LineRange{Start: 3, End: 4},
				CoveredLines: metadata.LineRange{Start: 1, End: 1},
			},
		}); err != nil {
			t.Fatalf("set metadata: %v", err)
		}

		h := &Handler{
			fileService: files.NewService(baseDir),
			metaStore:   metaStore,
		}

		updateBody := `{"sourceFile":"pkg/source.go","testFile":"/app/pkg/source_test.go","testName":"TestSource","newTestFile":"pkg/source_test.go"}`
		updateReq := httptest.NewRequest(http.MethodPut, "/api/metadata/test-path", strings.NewReader(updateBody))
		updateRR := httptest.NewRecorder()

		h.UpdateTestPath(updateRR, updateReq)

		if updateRR.Code != http.StatusNoContent {
			t.Fatalf("update status = %d, want %d", updateRR.Code, http.StatusNoContent)
		}

		deleteBody := `{"sourceFile":"pkg/source.go","testFile":"missing_test.go","testName":"TestMissing"}`
		deleteReq := httptest.NewRequest(http.MethodDelete, "/api/metadata/test-path", strings.NewReader(deleteBody))
		deleteRR := httptest.NewRecorder()

		h.DeleteTestPath(deleteRR, deleteReq)

		if deleteRR.Code != http.StatusNoContent {
			t.Fatalf("delete status = %d, want %d", deleteRR.Code, http.StatusNoContent)
		}

		meta := metaStore.GetTestMetadata("pkg/source.go")
		if meta == nil || len(meta.Tests) != 1 {
			t.Fatalf("expected 1 remaining test, got %+v", meta)
		}
		if meta.Tests[0].TestFile != "pkg/source_test.go" {
			t.Fatalf("expected updated test path %q, got %q", "pkg/source_test.go", meta.Tests[0].TestFile)
		}
	})
}
