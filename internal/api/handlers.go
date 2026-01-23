package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
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

// ReadLines reads specific lines from a file
func ReadLines(filePath string, start, end int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	lineNum := 1

	for scanner.Scan() {
		if lineNum >= start && lineNum <= end {
			lines = append(lines, scanner.Text())
		}
		if lineNum > end {
			break
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}
