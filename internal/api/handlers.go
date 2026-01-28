package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	// Read file content
	fileContent, err := h.fileService.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Attach metadata if available
	metadata := h.metaStore.GetTestMetadata(path)
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

	// Get metadata for the file
	fileMeta := h.metaStore.GetTestMetadata(path)
	if fileMeta == nil || len(fileMeta.Tests) == 0 {
		// No tests found
		response := files.TestsResponse{
			SourceFile: path,
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
		SourceFile: path,
		Tests:      testDetails,
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

	// Get suggestions for the file
	suggestions := h.metaStore.GetSuggestions(path)
	if suggestions == nil {
		suggestions = []files.TestSuggestion{}
	}

	response := files.SuggestionsResponse{
		SourceFile:  path,
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

	comments := h.metaStore.GetComments(path)
	if comments == nil {
		comments = []files.Comment{}
	}

	response := files.CommentsResponse{
		SourceFile: path,
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

	created, err := h.metaStore.AddComment(path, comment)
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

	if err := h.metaStore.UpdateComment(path, commentID, strings.TrimSpace(req.Content)); err != nil {
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

	if err := h.metaStore.DeleteComment(path, commentID); err != nil {
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

	if err := h.metaStore.ToggleCommentResolved(path, commentID); err != nil {
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
	fileContent, err := h.fileService.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get comments
	comments := h.metaStore.GetComments(path)
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
		SourceFile:  path,
		CodeContext: codeBlocks,
	}

	// Include tests if requested
	if req.IncludeTests {
		testsMeta := h.metaStore.GetTestMetadata(path)
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
		response.Suggestions = h.metaStore.GetSuggestions(path)
	}

	// Build formatted string for easy copying
	response.Formatted = buildFormattedExport(response)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
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
