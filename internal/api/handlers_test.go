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
	"codebase-view-mcp/internal/mcp"
	"codebase-view-mcp/internal/metadata"
)

func setupTestHandler(t *testing.T) (*Handler, string, func()) {
	tmpDir := t.TempDir()
	
	// Create test files
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	
	testDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	metadataFile := filepath.Join(tmpDir, "metadata.json")
	fileService := files.NewService(tmpDir)
	metaStore := metadata.NewStore(metadataFile)
	mcpHandler := mcp.NewHandler(metaStore, fileService)
	handler := NewHandler(fileService, metaStore, mcpHandler)
	
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	
	return handler, tmpDir, cleanup
}

func TestListFiles_Success(t *testing.T) {
	handler, tmpDir, cleanup := setupTestHandler(t)
	defer cleanup()
	
	req := httptest.NewRequest("GET", "/api/files?path=.", nil)
	w := httptest.NewRecorder()
	
	handler.ListFiles(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.ListFilesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Files) < 2 {
		t.Errorf("expected at least 2 files, got %d", len(resp.Files))
	}
	
	foundFile := false
	foundDir := false
	for _, f := range resp.Files {
		if f.Name == "test.go" && !f.IsDir {
			foundFile = true
		}
		if f.Name == "subdir" && f.IsDir {
			foundDir = true
		}
	}
	
	if !foundFile {
		t.Error("test.go not found in listing")
	}
	if !foundDir {
		t.Error("subdir not found in listing")
	}
	
	_ = tmpDir
}

func TestListFiles_NotFound(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	req := httptest.NewRequest("GET", "/api/files?path=nonexistent", nil)
	w := httptest.NewRecorder()
	
	handler.ListFiles(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
	
	var errResp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if errResp.Code != ErrFileNotFound {
		t.Errorf("error code = %s, want %s", errResp.Code, ErrFileNotFound)
	}
}

func TestGetFile_Success(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path...}", handler.GetFile)
	
	req := httptest.NewRequest("GET", "/api/files/test.go", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.FileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.File.Name != "test.go" {
		t.Errorf("name = %s, want test.go", resp.File.Name)
	}
	
	if !strings.Contains(resp.File.Content, "package main") {
		t.Error("expected file content to contain 'package main'")
	}
}

func TestGetFile_NotFound(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path...}", handler.GetFile)
	
	req := httptest.NewRequest("GET", "/api/files/nonexistent.go", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetFile_WithMetadata(t *testing.T) {
	handler, tmpDir, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Add test metadata
	testRef := files.TestReference{
		FunctionName: "main",
		TestFile:     "test_test.go",
		TestName:     "TestMain",
		Comment:      "tests main function",
		LineRange:    files.LineRange{Start: 1, End: 5},
		CoveredLines: files.LineRange{Start: 3, End: 3},
	}
	
	handler.metaStore.AddTestMetadata("test.go", []files.TestReference{testRef})
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path...}", handler.GetFile)
	
	req := httptest.NewRequest("GET", "/api/files/test.go", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.FileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.File.Metadata == nil {
		t.Fatal("expected metadata, got nil")
	}
	
	if len(resp.File.Metadata.Tests) != 1 {
		t.Errorf("expected 1 test, got %d", len(resp.File.Metadata.Tests))
	}
	
	if resp.File.CoverageDepth == nil {
		t.Fatal("expected coverage depth, got nil")
	}
	
	if len(resp.File.CoverageDepth[3]) != 1 {
		t.Errorf("expected 1 test covering line 3, got %d", len(resp.File.CoverageDepth[3]))
	}
	
	_ = tmpDir
}

func TestGetTests_Success(t *testing.T) {
	handler, tmpDir, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Create test file
	testFile := filepath.Join(tmpDir, "test_test.go")
	testContent := `package main
import "testing"
func TestMain(t *testing.T) {
	// input data
	x := 1
	// expected output
	if x != 1 {
		t.Error("fail")
	}
}`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Add test metadata
	testRef := files.TestReference{
		FunctionName: "main",
		TestFile:     "test_test.go",
		TestName:     "TestMain",
		Comment:      "tests main function",
		LineRange:    files.LineRange{Start: 3, End: 9},
		CoveredLines: files.LineRange{Start: 1, End: 3},
		InputLines:   files.LineRange{Start: 4, End: 5},
		OutputLines:  files.LineRange{Start: 6, End: 8},
	}
	
	handler.metaStore.AddTestMetadata("test.go", []files.TestReference{testRef})
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/tests", handler.GetTests)
	
	req := httptest.NewRequest("GET", "/api/files/test.go/tests", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.TestsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.SourceFile != "test.go" {
		t.Errorf("source file = %s, want test.go", resp.SourceFile)
	}
	
	if len(resp.Tests) != 1 {
		t.Fatalf("expected 1 test, got %d", len(resp.Tests))
	}
	
	test := resp.Tests[0]
	if test.TestName != "TestMain" {
		t.Errorf("test name = %s, want TestMain", test.TestName)
	}
	
	if test.InputData == "" {
		t.Error("expected input data to be extracted")
	}
	
	if test.ExpectedOutput == "" {
		t.Error("expected output data to be extracted")
	}
	
	if !strings.Contains(test.Content, "TestMain") {
		t.Error("expected test content to be included")
	}
}

func TestGetTests_NoTests(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/tests", handler.GetTests)
	
	req := httptest.NewRequest("GET", "/api/files/test.go/tests", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.TestsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Tests) != 0 {
		t.Errorf("expected 0 tests, got %d", len(resp.Tests))
	}
}

func TestGetTests_FilterByFunction(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	handler.metaStore.AddTestMetadata("test.go", []files.TestReference{
		{
			FunctionName: "main",
			TestFile:     "test_test.go",
			TestName:     "TestMain",
			LineRange:    files.LineRange{Start: 1, End: 5},
			CoveredLines: files.LineRange{Start: 1, End: 3},
		},
		{
			FunctionName: "helper",
			TestFile:     "test_test.go",
			TestName:     "TestHelper",
			LineRange:    files.LineRange{Start: 7, End: 10},
			CoveredLines: files.LineRange{Start: 5, End: 7},
		},
	})
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/tests", handler.GetTests)
	
	req := httptest.NewRequest("GET", "/api/files/test.go/tests?functionName=main", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.TestsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Tests) != 1 {
		t.Fatalf("expected 1 test, got %d", len(resp.Tests))
	}
	
	if resp.Tests[0].FunctionName != "main" {
		t.Errorf("function name = %s, want main", resp.Tests[0].FunctionName)
	}
}
func TestGetSources_Success(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Add test metadata
	testRef := files.TestReference{
		FunctionName: "main",
		TestFile:     "test_test.go",
		TestName:     "TestMain",
		LineRange:    files.LineRange{Start: 1, End: 5},
		CoveredLines: files.LineRange{Start: 1, End: 3},
	}
	
	handler.metaStore.AddTestMetadata("test.go", []files.TestReference{testRef})
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/sources", handler.GetSources)
	
	req := httptest.NewRequest("GET", "/api/files/test_test.go/sources", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.TestFileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.TestFile != "test_test.go" {
		t.Errorf("test file = %s, want test_test.go", resp.TestFile)
	}
	
	if len(resp.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(resp.Sources))
	}
	
	source := resp.Sources[0]
	if source.SourceFile != "test.go" {
		t.Errorf("source file = %s, want test.go", source.SourceFile)
	}
	
	if source.FunctionName != "main" {
		t.Errorf("function name = %s, want main", source.FunctionName)
	}
}

func TestGetSources_NoSources(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/sources", handler.GetSources)
	
	req := httptest.NewRequest("GET", "/api/files/test_test.go/sources", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.TestFileResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Sources) != 0 {
		t.Errorf("expected 0 sources, got %d", len(resp.Sources))
	}
}

func TestGetSuggestions_Success(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Add test suggestion
	suggestion := files.TestSuggestion{
		SourceFile:    "test.go",
		FunctionName:  "main",
		TargetLines:   files.LineRange{Start: 1, End: 3},
		Reason:        "needs test coverage",
		SuggestedName: "TestMainBasic",
		TestSkeleton:  "func TestMainBasic(t *testing.T) {}",
		Priority:      "high",
	}
	
	handler.metaStore.AddSuggestions("test.go", []files.TestSuggestion{suggestion})
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/suggestions", handler.GetSuggestions)
	
	req := httptest.NewRequest("GET", "/api/files/test.go/suggestions", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.SuggestionsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.SourceFile != "test.go" {
		t.Errorf("source file = %s, want test.go", resp.SourceFile)
	}
	
	if len(resp.Suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(resp.Suggestions))
	}
	
	sugg := resp.Suggestions[0]
	if sugg.SuggestedName != "TestMainBasic" {
		t.Errorf("suggested name = %s, want TestMainBasic", sugg.SuggestedName)
	}
	
	if sugg.Priority != "high" {
		t.Errorf("priority = %s, want high", sugg.Priority)
	}
}

func TestGetSuggestions_NoSuggestions(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/suggestions", handler.GetSuggestions)
	
	req := httptest.NewRequest("GET", "/api/files/test.go/suggestions", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.SuggestionsResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(resp.Suggestions))
	}
}
func TestCommentCRUD(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/{path}/comments", handler.GetComments)
	mux.HandleFunc("POST /api/files/{path}/comments", handler.CreateComment)
	mux.HandleFunc("PUT /api/files/{path}/comments/{commentId}", handler.UpdateComment)
	mux.HandleFunc("DELETE /api/files/{path}/comments/{commentId}", handler.DeleteComment)
	mux.HandleFunc("PATCH /api/files/{path}/comments/{commentId}/resolved", handler.ToggleCommentResolved)
	
	// Test GET comments (empty)
	req := httptest.NewRequest("GET", "/api/files/test.go/comments", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("GET comments status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var getResp files.CommentsResponse
	if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(getResp.Comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(getResp.Comments))
	}
	
	// Test POST comment
	createBody := `{"line": 2, "content": "test comment"}`
	req = httptest.NewRequest("POST", "/api/files/test.go/comments", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("POST comment status = %d, want %d", w.Code, http.StatusCreated)
	}
	
	var createResp files.CommentResponse
	if err := json.NewDecoder(w.Body).Decode(&createResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	commentID := createResp.Comment.ID
	if commentID == "" {
		t.Error("expected comment ID to be set")
	}
	
	if createResp.Comment.Line != 2 {
		t.Errorf("comment line = %d, want 2", createResp.Comment.Line)
	}
	
	if createResp.Comment.Content != "test comment" {
		t.Errorf("comment content = %s, want 'test comment'", createResp.Comment.Content)
	}
	
	// Test GET comments (with comment)
	req = httptest.NewRequest("GET", "/api/files/test.go/comments", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("GET comments status = %d, want %d", w.Code, http.StatusOK)
	}
	
	if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(getResp.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(getResp.Comments))
	}
	
	// Test PUT comment
	updateBody := `{"content": "updated comment"}`
	req = httptest.NewRequest("PUT", "/api/files/test.go/comments/"+commentID, strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("PUT comment status = %d, want %d", w.Code, http.StatusOK)
	}
	
	// Test PATCH resolved
	req = httptest.NewRequest("PATCH", "/api/files/test.go/comments/"+commentID+"/resolved", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("PATCH resolved status = %d, want %d", w.Code, http.StatusOK)
	}
	
	// Test DELETE comment
	req = httptest.NewRequest("DELETE", "/api/files/test.go/comments/"+commentID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("DELETE comment status = %d, want %d", w.Code, http.StatusNoContent)
	}
	
	// Verify comment is deleted
	req = httptest.NewRequest("GET", "/api/files/test.go/comments", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	
	if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(getResp.Comments) != 0 {
		t.Errorf("expected 0 comments after delete, got %d", len(getResp.Comments))
	}
}

func TestCreateComment_ValidationErrors(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/files/{path}/comments", handler.CreateComment)
	
	tests := []struct {
		name     string
		body     string
		wantCode int
		wantErr  ErrorCode
	}{
		{
			name:     "invalid JSON",
			body:     "invalid json",
			wantCode: http.StatusBadRequest,
			wantErr:  ErrInvalidRequest,
		},
		{
			name:     "line too small",
			body:     `{"line": 0, "content": "test"}`,
			wantCode: http.StatusBadRequest,
			wantErr:  ErrValidation,
		},
		{
			name:     "empty content",
			body:     `{"line": 1, "content": ""}`,
			wantCode: http.StatusBadRequest,
			wantErr:  ErrValidation,
		},
		{
			name:     "whitespace content",
			body:     `{"line": 1, "content": "   "}`,
			wantCode: http.StatusBadRequest,
			wantErr:  ErrValidation,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/files/test.go/comments", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}

			var errResp ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if errResp.Code != tt.wantErr {
				t.Errorf("error code = %s, want %s", errResp.Code, tt.wantErr)
			}
		})
	}
}
func TestExportContext_Success(t *testing.T) {
	handler, tmpDir, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Create a file with more content
	testFile := filepath.Join(tmpDir, "example.go")
	content := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}

func helper() {
	// helper function
}`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Add comment
	comment := files.Comment{
		Line:    6,
		Content: "This line needs review",
	}
	handler.metaStore.AddComment("example.go", comment)
	
	// Add test
	testRef := files.TestReference{
		FunctionName: "main",
		TestFile:     "example_test.go",
		TestName:     "TestMain",
		Comment:      "tests main function",
		LineRange:    files.LineRange{Start: 1, End: 5},
		CoveredLines: files.LineRange{Start: 5, End: 7},
	}
	handler.metaStore.AddTestMetadata("example.go", []files.TestReference{testRef})
	
	// Add suggestion
	suggestion := files.TestSuggestion{
		SourceFile:    "example.go",
		FunctionName:  "helper",
		TargetLines:   files.LineRange{Start: 9, End: 11},
		Reason:        "needs test coverage",
		SuggestedName: "TestHelper",
		TestSkeleton:  "func TestHelper(t *testing.T) {}",
		Priority:      "medium",
	}
	handler.metaStore.AddSuggestions("example.go", []files.TestSuggestion{suggestion})
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/files/{path}/export", handler.ExportContext)

	reqBody := `{"contextLines": 2, "includeTests": true, "includeSuggestions": true}`
	req := httptest.NewRequest("POST", "/api/files/example.go/export", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.ExportContextResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.SourceFile != "example.go" {
		t.Errorf("source file = %s, want example.go", resp.SourceFile)
	}
	
	if len(resp.CodeContext) != 1 {
		t.Errorf("expected 1 code context block, got %d", len(resp.CodeContext))
	}
	
	if len(resp.Tests) != 1 {
		t.Errorf("expected 1 test, got %d", len(resp.Tests))
	}
	
	if len(resp.Suggestions) != 1 {
		t.Errorf("expected 1 suggestion, got %d", len(resp.Suggestions))
	}
	
	if resp.Formatted == "" {
		t.Error("expected formatted export to be non-empty")
	}
	
	// Check that formatted export contains expected sections
	if !strings.Contains(resp.Formatted, "Code Review Export") {
		t.Error("formatted export should contain header")
	}
	
	if !strings.Contains(resp.Formatted, "This line needs review") {
		t.Error("formatted export should contain comment")
	}
	
	_ = tmpDir
}

func TestExportContext_DefaultValues(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/files/{path}/export", handler.ExportContext)
	
	// Test with no body (should use defaults)
	req := httptest.NewRequest("POST", "/api/files/test.go/export", nil)
	w := httptest.NewRecorder()
	
	mux.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.ExportContextResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.SourceFile != "test.go" {
		t.Errorf("source file = %s, want test.go", resp.SourceFile)
	}
}

func TestGetOverview_Success(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Add test metadata for multiple files
	handler.metaStore.AddTestMetadata("file1.go", []files.TestReference{
		{
			FunctionName: "func1",
			TestFile:     "file1_test.go",
			TestName:     "TestFunc1",
			LineRange:    files.LineRange{Start: 1, End: 5},
			CoveredLines: files.LineRange{Start: 1, End: 3},
		},
		{
			FunctionName: "func2",
			TestFile:     "file1_test.go",
			TestName:     "TestFunc2",
			LineRange:    files.LineRange{Start: 7, End: 10},
			CoveredLines: files.LineRange{Start: 5, End: 7},
		},
	})
	
	handler.metaStore.AddTestMetadata("file2.go", []files.TestReference{
		{
			FunctionName: "func3",
			TestFile:     "file2_test.go",
			TestName:     "TestFunc3",
			LineRange:    files.LineRange{Start: 1, End: 5},
			CoveredLines: files.LineRange{Start: 1, End: 3},
		},
	})
	
	req := httptest.NewRequest("GET", "/api/overview", nil)
	w := httptest.NewRecorder()
	
	handler.GetOverview(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.OverviewResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.TotalTests != 3 {
		t.Errorf("total tests = %d, want 3", resp.TotalTests)
	}
	
	if resp.TotalFunctions != 3 {
		t.Errorf("total functions = %d, want 3", resp.TotalFunctions)
	}
	
	if resp.TotalSourceFiles != 2 {
		t.Errorf("total source files = %d, want 2", resp.TotalSourceFiles)
	}
	
	if resp.TotalTestFiles != 2 {
		t.Errorf("total test files = %d, want 2", resp.TotalTestFiles)
	}
	
	if len(resp.Functions) != 3 {
		t.Errorf("functions count = %d, want 3", len(resp.Functions))
	}
	
	if len(resp.TestsBySourceFile) != 2 {
		t.Errorf("tests by source file count = %d, want 2", len(resp.TestsBySourceFile))
	}
}
func TestExtractLines(t *testing.T) {
	lines := []string{
		"line 1",
		"line 2", 
		"line 3",
		"line 4",
		"line 5",
	}
	
	tests := []struct {
		name     string
		start    int
		end      int
		expected string
	}{
		{
			name:     "valid range",
			start:    2,
			end:      4,
			expected: "line 2\nline 3\nline 4",
		},
		{
			name:     "single line",
			start:    3,
			end:      3,
			expected: "line 3",
		},
		{
			name:     "invalid start",
			start:    0,
			end:      2,
			expected: "",
		},
		{
			name:     "invalid end",
			start:    2,
			end:      10,
			expected: "",
		},
		{
			name:     "start > end",
			start:    4,
			end:      2,
			expected: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLines(lines, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("extractLines(%d, %d) = %q, want %q", tt.start, tt.end, result, tt.expected)
			}
		})
	}
}

func TestBuildFormattedExport(t *testing.T) {
	data := files.ExportContextResponse{
		SourceFile: "test.go",
		CodeContext: []files.CodeContextBlock{
			{
				LineRange: files.LineRange{Start: 5, End: 7},
				Code:      "5: func main() {\n6:   fmt.Println(\"hello\")\n7: }",
				Comments: []files.Comment{
					{
						Line:    6,
						Content: "This needs review",
					},
				},
			},
		},
		Tests: []files.TestDetail{
			{
				FunctionName: "main",
				TestFile:     "test_test.go",
				TestName:     "TestMain",
				Comment:      "tests main function",
				LineRange:    files.LineRange{Start: 1, End: 5},
			},
		},
		Suggestions: []files.TestSuggestion{
			{
				SuggestedName: "TestHelper",
				Priority:      "high",
				Reason:        "needs coverage",
				TargetLines:   files.LineRange{Start: 10, End: 15},
				TestSkeleton:  "func TestHelper(t *testing.T) {}",
			},
		},
	}
	
	result := buildFormattedExport(data)
	
	if !strings.Contains(result, "Code Review Export") {
		t.Error("should contain header")
	}
	
	if !strings.Contains(result, "test.go") {
		t.Error("should contain file name")
	}
	
	if !strings.Contains(result, "This needs review") {
		t.Error("should contain comment")
	}
	
	if !strings.Contains(result, "TestMain") {
		t.Error("should contain test name")
	}
	
	if !strings.Contains(result, "TestHelper") {
		t.Error("should contain suggestion name")
	}
	
	if !strings.Contains(result, "Priority: high") {
		t.Error("should contain priority")
	}
}

func TestGetMetadataIssues_Success(t *testing.T) {
	handler, tmpDir, cleanup := setupTestHandler(t)
	defer cleanup()
	
	// Add metadata with invalid paths
	handler.metaStore.AddTestMetadata("nonexistent.go", []files.TestReference{
		{
			FunctionName: "test",
			TestFile:     "nonexistent_test.go",
			TestName:     "TestSomething",
			LineRange:    files.LineRange{Start: 1, End: 5},
			CoveredLines: files.LineRange{Start: 1, End: 3},
		},
	})
	
	req := httptest.NewRequest("GET", "/api/metadata/issues", nil)
	w := httptest.NewRecorder()
	
	handler.GetMetadataIssues(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.MetadataIssuesResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if len(resp.Issues) == 0 {
		t.Error("expected at least one issue")
	}
	
	issue := resp.Issues[0]
	if issue.SourceFile != "nonexistent.go" {
		t.Errorf("source file = %s, want nonexistent.go", issue.SourceFile)
	}
	
	if issue.SourceValid {
		t.Error("expected source to be invalid")
	}
	
	_ = tmpDir
}

func TestHandleMCP_MethodNotAllowed(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	req := httptest.NewRequest("GET", "/api/mcp", nil)
	w := httptest.NewRecorder()
	
	handler.HandleMCP(w, req)
	
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
	
	var errResp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if errResp.Code != ErrMethodNotAllowed {
		t.Errorf("error code = %s, want %s", errResp.Code, ErrMethodNotAllowed)
	}
}

func TestSearch_Success(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()
	
	req := httptest.NewRequest("GET", "/api/search?q=test", nil)
	w := httptest.NewRecorder()
	
	handler.Search(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var resp files.SearchResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	
	if resp.Query != "test" {
		t.Errorf("query = %s, want test", resp.Query)
	}
	
	// Should find test.go file
	found := false
	for _, result := range resp.Results {
		if strings.Contains(result.Title, "test.go") {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("expected to find test.go in search results")
	}
}

// Test error response helper
func TestHandlerJSONErrors(t *testing.T) {
	handler, _, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		setupRoute     func(*http.ServeMux)
		wantStatus     int
		wantErrorCode  ErrorCode
		wantErrorMsg   string
		checkDetails   bool
		wantDetailsKey string
	}{
		{
			name:   "ListFiles - file not found",
			method: "GET",
			path:   "/api/files?path=nonexistent-dir-12345",
			setupRoute: func(mux *http.ServeMux) {
				mux.HandleFunc("GET /api/files", handler.ListFiles)
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: ErrFileNotFound,
		},
		{
			name:   "GetFile - file not found",
			method: "GET",
			path:   "/api/files/nonexistent-file-12345.txt",
			setupRoute: func(mux *http.ServeMux) {
				mux.HandleFunc("GET /api/files/{path...}", handler.GetFile)
			},
			wantStatus:    http.StatusNotFound,
			wantErrorCode: ErrFileNotFound,
		},
		{
			name:   "CreateComment - invalid body",
			method: "POST",
			path:   "/api/files/test.go/comments",
			body:   "invalid json",
			setupRoute: func(mux *http.ServeMux) {
				mux.HandleFunc("POST /api/files/{path}/comments", handler.CreateComment)
			},
			wantStatus:    http.StatusBadRequest,
			wantErrorCode: ErrInvalidRequest,
			wantErrorMsg:  "invalid request body",
		},
		{
			name:   "CreateComment - validation error with details",
			method: "POST",
			path:   "/api/files/test.go/comments",
			body:   `{"line": 0, "content": "test"}`,
			setupRoute: func(mux *http.ServeMux) {
				mux.HandleFunc("POST /api/files/{path}/comments", handler.CreateComment)
			},
			wantStatus:     http.StatusBadRequest,
			wantErrorCode:  ErrValidation,
			wantErrorMsg:   "line must be >= 1",
			checkDetails:   true,
			wantDetailsKey: "field",
		},
		{
			name:   "HandleMCP - method not allowed",
			method: "GET",
			path:   "/api/mcp",
			setupRoute: func(mux *http.ServeMux) {
				mux.HandleFunc("/api/mcp", handler.HandleMCP)
			},
			wantStatus:    http.StatusMethodNotAllowed,
			wantErrorCode: ErrMethodNotAllowed,
			wantErrorMsg:  "method not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			tt.setupRoute(mux)

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status code = %d, want %d", w.Code, tt.wantStatus)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %s, want application/json", contentType)
			}

			var errResp ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
				t.Fatalf("failed to decode error response: %v, body: %s", err, w.Body.String())
			}

			if errResp.Code != tt.wantErrorCode {
				t.Errorf("error code = %s, want %s", errResp.Code, tt.wantErrorCode)
			}

			if tt.wantErrorMsg != "" && errResp.Error != tt.wantErrorMsg {
				t.Errorf("error message = %s, want %s", errResp.Error, tt.wantErrorMsg)
			}

			if tt.checkDetails {
				if errResp.Details == nil {
					t.Error("expected details but got nil")
				} else if _, ok := errResp.Details[tt.wantDetailsKey]; !ok {
					t.Errorf("expected details key %s but not found", tt.wantDetailsKey)
				}
			}
		})
	}
}
