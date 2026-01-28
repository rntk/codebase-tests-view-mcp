package files

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Service handles file system operations
type Service struct {
	baseDir string
}

// NewService creates a new file service
func NewService(baseDir string) *Service {
	return &Service{
		baseDir: baseDir,
	}
}

// ListFiles lists files and directories in the specified path
func (s *Service) ListFiles(path string) (*ListFilesResponse, error) {
	// Resolve the full path
	fullPath := s.resolvePath(path)

	// Check if path exists and is a directory
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory")
	}

	// Read directory contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Convert to FileEntry
	var files []FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip hidden files (starting with .)
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		files = append(files, FileEntry{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	// Sort: directories first, then by name
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	return &ListFilesResponse{
		Path:  path,
		Files: files,
	}, nil
}

// ReadFile reads the content of a file
func (s *Service) ReadFile(path string) (*FileContent, error) {
	// Resolve the full path
	fullPath := s.resolvePath(path)

	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file")
	}

	// Read file content
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Determine MIME type
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = "text/plain"
	}

	return &FileContent{
		Path:     path,
		Name:     filepath.Base(path),
		Content:  string(content),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		MimeType: mimeType,
	}, nil
}

// resolvePath resolves a relative path to an absolute path within baseDir
func (s *Service) resolvePath(path string) string {
	if path == "" || path == "." {
		return s.baseDir
	}

	// Clean the path to prevent directory traversal
	cleanPath := filepath.Clean(path)

	// If it's already absolute, use it directly (for development)
	if filepath.IsAbs(cleanPath) {
		return cleanPath
	}

	// Otherwise, join with baseDir
	return filepath.Join(s.baseDir, cleanPath)
}
