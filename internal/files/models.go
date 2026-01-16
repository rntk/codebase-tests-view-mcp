package files

import "time"

// FileEntry represents a file or directory in a listing
type FileEntry struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size,omitempty"`
	ModTime time.Time `json:"modTime"`
}

// FileContent represents file content with metadata
type FileContent struct {
	Path     string        `json:"path"`
	Name     string        `json:"name"`
	Content  string        `json:"content"`
	Size     int64         `json:"size"`
	ModTime  time.Time     `json:"modTime"`
	MimeType string        `json:"mimeType"`
	Metadata *FileMetadata `json:"metadata,omitempty"`
}

// FileMetadata contains test-related metadata for a file
type FileMetadata struct {
	Tests []TestReference `json:"tests,omitempty"`
}

// TestReference links a source file to its tests
type TestReference struct {
	TestFile     string    `json:"testFile"`
	TestName     string    `json:"testName"`
	LineRange    LineRange `json:"lineRange"`
	CoveredLines LineRange `json:"coveredLines"`
	InputLines   LineRange `json:"inputLines,omitempty"`
	OutputLines  LineRange `json:"outputLines,omitempty"`
}

// LineRange specifies a range of lines
type LineRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// ListFilesResponse for GET /api/files
type ListFilesResponse struct {
	Path  string      `json:"path"`
	Files []FileEntry `json:"files"`
}

// FileResponse for GET /api/files/{path}
type FileResponse struct {
	File FileContent `json:"file"`
}

// TestDetail contains full test information
type TestDetail struct {
	TestFile       string    `json:"testFile"`
	TestName       string    `json:"testName"`
	Content        string    `json:"content"`
	LineRange      LineRange `json:"lineRange"`
	CoveredLines   LineRange `json:"coveredLines"`
	InputData      string    `json:"inputData,omitempty"`
	InputLines     LineRange `json:"inputLines,omitempty"`
	ExpectedOutput string    `json:"expectedOutput,omitempty"`
	OutputLines    LineRange `json:"outputLines,omitempty"`
}

// TestsResponse for GET /api/files/{path}/tests
type TestsResponse struct {
	SourceFile string       `json:"sourceFile"`
	Tests      []TestDetail `json:"tests"`
}
