package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"codebase-view-mcp/internal/files"
	"codebase-view-mcp/internal/mcp"
	"codebase-view-mcp/internal/metadata"
)

// Handler handles HTTP requests
type Handler struct {
	fileService *files.Service
	metaStore   *metadata.Store
	mcpHandler  *mcp.Handler
}

// NewHandler creates a new HTTP handler
func NewHandler(fileService *files.Service, metaStore *metadata.Store, mcpHandler *mcp.Handler) *Handler {
	return &Handler{
		fileService: fileService,
		metaStore:   metaStore,
		mcpHandler:  mcpHandler,
	}
}

// ListFiles handles GET /api/files
func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	response, err := h.fileService.ListFiles(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetFile handles GET /api/files/{path}
func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read file content
	fileContent, err := h.fileService.ReadFile(canonicalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Attach metadata if available
	metadata := h.metaStore.GetTestMetadata(canonicalPath)
	if metadata != nil {
		fileContent.Metadata = metadata

		// Calculate coverage depth: map of line number -> list of test names covering it
		if len(metadata.Tests) > 0 {
			coverageDepth := make(map[int][]string)
			for _, test := range metadata.Tests {
				for line := test.CoveredLines.Start; line <= test.CoveredLines.End; line++ {
					coverageDepth[line] = append(coverageDepth[line], test.TestName)
				}
			}
			fileContent.CoverageDepth = coverageDepth
		}
	}

	response := files.FileResponse{
		File: *fileContent,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetTests handles GET /api/files/{path}/tests
func (h *Handler) GetTests(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get metadata for the file
	fileMeta := h.metaStore.GetTestMetadata(canonicalPath)
	if fileMeta == nil || len(fileMeta.Tests) == 0 {
		// No tests found
		response := files.TestsResponse{
			SourceFile: canonicalPath,
			Tests:      []files.TestDetail{},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	// Build detailed test information
	var testDetails []files.TestDetail
	for _, testRef := range fileMeta.Tests {
		detail := files.TestDetail{
			FunctionName: testRef.FunctionName,
			TestFile:     testRef.TestFile,
			TestName:     testRef.TestName,
			Comment:      testRef.Comment,
			LineRange:    testRef.LineRange,
			CoveredLines: testRef.CoveredLines,
			InputLines:   testRef.InputLines,
			OutputLines:  testRef.OutputLines,
		}

		// Read test file content
		testContent, err := h.fileService.ReadFile(testRef.TestFile)
		if err == nil {
			detail.Content = testContent.Content

			// Extract input and output data if line ranges are provided
			lines := strings.Split(testContent.Content, "\n")

			if testRef.InputLines.Start > 0 && testRef.InputLines.End > 0 {
				detail.InputData = extractLines(lines, testRef.InputLines.Start, testRef.InputLines.End)
			}

			if testRef.OutputLines.Start > 0 && testRef.OutputLines.End > 0 {
				detail.ExpectedOutput = extractLines(lines, testRef.OutputLines.Start, testRef.OutputLines.End)
			}
		}

		testDetails = append(testDetails, detail)
	}

	response := files.TestsResponse{
		SourceFile: canonicalPath,
		Tests:      testDetails,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetSources handles GET /api/files/{path}/sources (reverse lookup: test file → source references)
func (h *Handler) GetSources(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sources := h.metaStore.GetTestFileMetadata(canonicalPath)
	if sources == nil {
		sources = []files.SourceReference{}
	}

	response := files.TestFileResponse{
		TestFile: canonicalPath,
		Sources:  sources,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetSuggestions handles GET /api/files/{path}/suggestions
func (h *Handler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get suggestions for the file
	suggestions := h.metaStore.GetSuggestions(canonicalPath)
	if suggestions == nil {
		suggestions = []files.TestSuggestion{}
	}

	response := files.SuggestionsResponse{
		SourceFile:  canonicalPath,
		Suggestions: suggestions,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// HandleMCP handles POST /api/mcp
func (h *Handler) HandleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Delegate to MCP handler
	h.mcpHandler.Handle(w, r)
}

// ServeStatic serves the static frontend files
func (h *Handler) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Serve embedded static files
	staticHandler := ServeStaticFiles()
	staticHandler.ServeHTTP(w, r)
}

// ==================== COMMENT HANDLERS ====================

// GetComments handles GET /api/files/{path}/comments
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments := h.metaStore.GetComments(canonicalPath)
	if comments == nil {
		comments = []files.Comment{}
	}

	response := files.CommentsResponse{
		SourceFile: canonicalPath,
		Comments:   comments,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateComment handles POST /api/files/{path}/comments
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req files.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Line < 1 {
		http.Error(w, "line must be >= 1", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	comment := files.Comment{
		Line:         req.Line,
		Content:      strings.TrimSpace(req.Content),
		ContextLines: req.ContextLines,
	}

	created, err := h.metaStore.AddComment(canonicalPath, comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := files.CommentResponse{Comment: created}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateComment handles PUT /api/files/{path}/comments/{commentId}
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	commentID := r.PathValue("commentId")
	if path == "" || commentID == "" {
		http.Error(w, "path and commentId are required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	if err := h.metaStore.UpdateComment(canonicalPath, commentID, strings.TrimSpace(req.Content)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteComment handles DELETE /api/files/{path}/comments/{commentId}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	commentID := r.PathValue("commentId")
	if path == "" || commentID == "" {
		http.Error(w, "path and commentId are required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.metaStore.DeleteComment(canonicalPath, commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleCommentResolved handles PATCH /api/files/{path}/comments/{commentId}/resolved
func (h *Handler) ToggleCommentResolved(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	commentID := r.PathValue("commentId")
	if path == "" || commentID == "" {
		http.Error(w, "path and commentId are required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.metaStore.ToggleCommentResolved(canonicalPath, commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ExportContext handles POST /api/files/{path}/export
// Returns formatted code context with comments for AI agents
func (h *Handler) ExportContext(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req files.ExportContextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use defaults if no body
		req.ContextLines = 5
		req.IncludeTests = true
		req.IncludeSuggestions = true
	}

	if req.ContextLines < 1 {
		req.ContextLines = 5
	}

	// Get file content
	fileContent, err := h.fileService.ReadFile(canonicalPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get comments
	comments := h.metaStore.GetComments(canonicalPath)
	if comments == nil {
		comments = []files.Comment{}
	}

	// Build code context blocks
	lines := strings.Split(fileContent.Content, "\n")
	var codeBlocks []files.CodeContextBlock

	for _, comment := range comments {
		if comment.Resolved {
			continue // Skip resolved comments
		}

		// Calculate context range
		start := comment.Line - req.ContextLines
		if start < 1 {
			start = 1
		}
		end := comment.Line + req.ContextLines
		if end > len(lines) {
			end = len(lines)
		}

		// Extract code with line numbers
		var codeBuilder strings.Builder
		for i := start - 1; i < end; i++ {
			codeBuilder.WriteString(fmt.Sprintf("%d: %s\n", i+1, lines[i]))
		}

		codeBlocks = append(codeBlocks, files.CodeContextBlock{
			LineRange: files.LineRange{Start: start, End: end},
			Code:      codeBuilder.String(),
			Comments:  []files.Comment{comment},
		})
	}

	// Build response
	response := files.ExportContextResponse{
		SourceFile:  canonicalPath,
		CodeContext: codeBlocks,
	}

	// Include tests if requested
	if req.IncludeTests {
		testsMeta := h.metaStore.GetTestMetadata(canonicalPath)
		if testsMeta != nil {
			for _, testRef := range testsMeta.Tests {
				detail := files.TestDetail{
					FunctionName: testRef.FunctionName,
					TestFile:     testRef.TestFile,
					TestName:     testRef.TestName,
					Comment:      testRef.Comment,
					LineRange:    testRef.LineRange,
					CoveredLines: testRef.CoveredLines,
				}

				testContent, err := h.fileService.ReadFile(testRef.TestFile)
				if err == nil {
					detail.Content = testContent.Content
					lines := strings.Split(testContent.Content, "\n")
					if testRef.InputLines.Start > 0 && testRef.InputLines.End > 0 {
						detail.InputData = extractLines(lines, testRef.InputLines.Start, testRef.InputLines.End)
					}
					if testRef.OutputLines.Start > 0 && testRef.OutputLines.End > 0 {
						detail.ExpectedOutput = extractLines(lines, testRef.OutputLines.Start, testRef.OutputLines.End)
					}
				}

				response.Tests = append(response.Tests, detail)
			}
		}
	}

	// Include suggestions if requested
	if req.IncludeSuggestions {
		response.Suggestions = h.metaStore.GetSuggestions(canonicalPath)
	}

	// Build formatted string for easy copying
	response.Formatted = buildFormattedExport(response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetMetadataIssues handles GET /api/metadata/issues
func (h *Handler) GetMetadataIssues(w http.ResponseWriter, r *http.Request) {
	allMetadata := h.metaStore.GetAllMetadata()

	sourceFiles := make([]string, 0, len(allMetadata))
	for sourceFile := range allMetadata {
		sourceFiles = append(sourceFiles, sourceFile)
	}
	sort.Strings(sourceFiles)

	issues := make([]files.MetadataIssue, 0)
	for _, sourceFile := range sourceFiles {
		meta := allMetadata[sourceFile]
		if meta == nil {
			continue
		}

		sourceValid, sourceIsAbsolute, sourceMessage := h.inspectStoredPath(sourceFile)
		issue := files.MetadataIssue{
			SourceFile:        sourceFile,
			SourceValid:       sourceValid,
			SourceIsAbsolute:  sourceIsAbsolute,
			SourceMessage:     sourceMessage,
			SuggestionsCount:  len(meta.Suggestions),
			CommentsCount:     len(meta.Comments),
			InvalidTestIssues: []files.MetadataTestIssue{},
		}

		for _, testRef := range meta.Tests {
			testValid, testIsAbsolute, testMessage := h.inspectStoredPath(testRef.TestFile)
			if testValid {
				continue
			}
			issue.InvalidTestIssues = append(issue.InvalidTestIssues, files.MetadataTestIssue{
				TestFile:   testRef.TestFile,
				TestName:   testRef.TestName,
				IsAbsolute: testIsAbsolute,
				Message:    testMessage,
			})
		}

		sort.Slice(issue.InvalidTestIssues, func(i, j int) bool {
			if issue.InvalidTestIssues[i].TestFile == issue.InvalidTestIssues[j].TestFile {
				return issue.InvalidTestIssues[i].TestName < issue.InvalidTestIssues[j].TestName
			}
			return issue.InvalidTestIssues[i].TestFile < issue.InvalidTestIssues[j].TestFile
		})

		if !issue.SourceValid || len(issue.InvalidTestIssues) > 0 {
			issues = append(issues, issue)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files.MetadataIssuesResponse{Issues: issues}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateSourcePath handles PUT /api/metadata/source-path
func (h *Handler) UpdateSourcePath(w http.ResponseWriter, r *http.Request) {
	var req files.UpdateSourcePathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.OldPath) == "" || strings.TrimSpace(req.NewPath) == "" {
		http.Error(w, "oldPath and newPath are required", http.StatusBadRequest)
		return
	}

	newPath, err := h.canonicalizeExistingPath(req.NewPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.metaStore.RenameSourcePath(req.OldPath, newPath); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTestPath handles PUT /api/metadata/test-path
func (h *Handler) UpdateTestPath(w http.ResponseWriter, r *http.Request) {
	var req files.UpdateTestPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SourceFile) == "" || strings.TrimSpace(req.TestFile) == "" || strings.TrimSpace(req.TestName) == "" || strings.TrimSpace(req.NewTestFile) == "" {
		http.Error(w, "sourceFile, testFile, testName, and newTestFile are required", http.StatusBadRequest)
		return
	}

	newTestFile, err := h.canonicalizeExistingPath(req.NewTestFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.metaStore.UpdateTestPath(req.SourceFile, req.TestFile, req.TestName, newTestFile); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteSourcePath handles DELETE /api/metadata/source-path
func (h *Handler) DeleteSourcePath(w http.ResponseWriter, r *http.Request) {
	var req files.DeleteSourcePathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Path) == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	if err := h.metaStore.DeleteSourcePath(req.Path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTestPath handles DELETE /api/metadata/test-path
func (h *Handler) DeleteTestPath(w http.ResponseWriter, r *http.Request) {
	var req files.DeleteTestPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.SourceFile) == "" || strings.TrimSpace(req.TestFile) == "" || strings.TrimSpace(req.TestName) == "" {
		http.Error(w, "sourceFile, testFile, and testName are required", http.StatusBadRequest)
		return
	}

	if err := h.metaStore.DeleteTestPath(req.SourceFile, req.TestFile, req.TestName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// buildFormattedExport creates a human-readable formatted export for AI agents
func buildFormattedExport(data files.ExportContextResponse) string {
	var b strings.Builder

	b.WriteString("# Code Review Export\n\n")
	b.WriteString(fmt.Sprintf("**File:** `%s`\n\n", data.SourceFile))

	if len(data.CodeContext) > 0 {
		b.WriteString("## Comments and Code Context\n\n")
		for i, block := range data.CodeContext {
			b.WriteString(fmt.Sprintf("### Issue %d (Line %d)\n\n", i+1, block.Comments[0].Line))
			b.WriteString(fmt.Sprintf("**Comment:** %s\n\n", block.Comments[0].Content))
			b.WriteString("**Code Context:**\n```\n")
			b.WriteString(block.Code)
			b.WriteString("```\n\n")
		}
	}

	if len(data.Tests) > 0 {
		b.WriteString("## Related Tests\n\n")
		for _, test := range data.Tests {
			b.WriteString(fmt.Sprintf("### %s\n", test.TestName))
			b.WriteString(fmt.Sprintf("- **Function:** %s\n", test.FunctionName))
			b.WriteString(fmt.Sprintf("- **Test File:** %s\n", test.TestFile))
			b.WriteString(fmt.Sprintf("- **Lines:** %d-%d\n", test.LineRange.Start, test.LineRange.End))
			if test.Comment != "" {
				b.WriteString(fmt.Sprintf("- **Description:** %s\n", test.Comment))
			}
			b.WriteString("\n")
		}
	}

	if len(data.Suggestions) > 0 {
		b.WriteString("## Test Suggestions\n\n")
		for _, sugg := range data.Suggestions {
			b.WriteString(fmt.Sprintf("### %s (Priority: %s)\n", sugg.SuggestedName, sugg.Priority))
			b.WriteString(fmt.Sprintf("- **Reason:** %s\n", sugg.Reason))
			b.WriteString(fmt.Sprintf("- **Target Lines:** %d-%d\n", sugg.TargetLines.Start, sugg.TargetLines.End))
			if sugg.TestSkeleton != "" {
				b.WriteString("\n**Suggested Test Skeleton:**\n```\n")
				b.WriteString(sugg.TestSkeleton)
				b.WriteString("\n```\n")
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// extractLines extracts lines from start to end (1-indexed, inclusive)
func extractLines(lines []string, start, end int) string {
	if start < 1 || end < 1 || start > len(lines) || end > len(lines) || start > end {
		return ""
	}

	// Convert to 0-indexed
	start--
	// end is inclusive, so no need to subtract 1

	selected := lines[start:end]
	return strings.Join(selected, "\n")
}

// GetOverview handles GET /api/overview
// Returns a global summary of all tests and functions across the codebase
func (h *Handler) GetOverview(w http.ResponseWriter, r *http.Request) {
	// Get all metadata from the store
	allMetadata := h.metaStore.GetAllMetadata()

	// Build the overview response
	response := files.OverviewResponse{
		TestsBySourceFile: make(map[string][]files.TestDetail),
	}

	// Track unique test files and source files
	testFilesSet := make(map[string]bool)
	sourceFilesSet := make(map[string]bool)

	// Group tests by function name to count unique functions
	functionTestsMap := make(map[string][]files.TestDetail)
	functionSourceMap := make(map[string]string) // functionName -> sourceFile

	for sourceFile, metadata := range allMetadata {
		if metadata == nil || len(metadata.Tests) == 0 {
			continue
		}

		sourceFilesSet[sourceFile] = true

		for _, testRef := range metadata.Tests {
			// Track test file
			testFilesSet[testRef.TestFile] = true

			// Build test detail
			detail := files.TestDetail{
				FunctionName: testRef.FunctionName,
				TestFile:     testRef.TestFile,
				TestName:     testRef.TestName,
				Comment:      testRef.Comment,
				LineRange:    testRef.LineRange,
				CoveredLines: testRef.CoveredLines,
				InputLines:   testRef.InputLines,
				OutputLines:  testRef.OutputLines,
			}

			// Read test file content
			testContent, err := h.fileService.ReadFile(testRef.TestFile)
			if err == nil {
				detail.Content = testContent.Content
				lines := strings.Split(testContent.Content, "\n")
				if testRef.InputLines.Start > 0 && testRef.InputLines.End > 0 {
					detail.InputData = extractLines(lines, testRef.InputLines.Start, testRef.InputLines.End)
				}
				if testRef.OutputLines.Start > 0 && testRef.OutputLines.End > 0 {
					detail.ExpectedOutput = extractLines(lines, testRef.OutputLines.Start, testRef.OutputLines.End)
				}
			}

			// Add to tests by source file
			response.TestsBySourceFile[sourceFile] = append(response.TestsBySourceFile[sourceFile], detail)

			// Group by function name
			funcKey := testRef.FunctionName + "::" + sourceFile
			functionTestsMap[funcKey] = append(functionTestsMap[funcKey], detail)
			functionSourceMap[funcKey] = sourceFile

			response.TotalTests++
		}
	}

	// Build function summaries
	for funcKey, tests := range functionTestsMap {
		// Extract function name from key (format: "functionName::sourceFile")
		functionName := funcKey
		for i := len(funcKey) - 1; i >= 0; i-- {
			if funcKey[i] == ':' && i > 0 && funcKey[i-1] == ':' {
				functionName = funcKey[:i-1]
				break
			}
		}

		sourceFile := functionSourceMap[funcKey]
		response.Functions = append(response.Functions, files.FunctionSummary{
			FunctionName: functionName,
			SourceFile:   sourceFile,
			TestCount:    len(tests),
			Tests:        tests,
		})
		response.TotalFunctions++
	}

	response.TotalSourceFiles = len(sourceFilesSet)
	response.TotalTestFiles = len(testFilesSet)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) inspectStoredPath(path string) (bool, bool, string) {
	isAbsolute := filepath.IsAbs(filepath.Clean(strings.TrimSpace(path)))

	resolvedPath, err := h.fileService.ResolvePath(path)
	if err != nil {
		return false, isAbsolute, err.Error()
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		if os.IsNotExist(err) {
			return false, isAbsolute, fmt.Sprintf("path %q does not exist under configured codebase root %q", path, h.fileService.BaseDir())
		}
		return false, isAbsolute, fmt.Sprintf("failed to inspect path %q: %v", path, err)
	}

	return true, isAbsolute, ""
}

func (h *Handler) canonicalizeExistingPath(path string) (string, error) {
	canonicalPath, err := h.fileService.CanonicalizePath(path)
	if err != nil {
		return "", err
	}

	resolvedPath, err := h.fileService.ResolvePath(canonicalPath)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(resolvedPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path %q does not exist under configured codebase root %q", canonicalPath, h.fileService.BaseDir())
		}
		return "", fmt.Errorf("failed to inspect path %q: %w", canonicalPath, err)
	}

	return canonicalPath, nil
}
