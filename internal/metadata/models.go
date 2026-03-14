package metadata

import "codebase-view-mcp/internal/files"

// Re-export types from files package for clarity
type (
	FileMetadata  = files.FileMetadata
	TestReference = files.TestReference
	LineRange     = files.LineRange
)
