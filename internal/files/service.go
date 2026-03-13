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

// BaseDir returns the configured filesystem root for this service.
func (s *Service) BaseDir() string {
	return s.baseDir
}

// ListFiles lists files and directories in the specified path
func (s *Service) ListFiles(path string) (*ListFilesResponse, error) {
	relativePath, err := s.CanonicalizePath(path)
	if err != nil {
		return nil, err
	}

	// Resolve the full path
	fullPath := s.resolveFSPath(relativePath)

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
			Path:    toSlash(filepath.Join(relativePath, entry.Name())),
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
		Path:  relativePath,
		Files: files,
	}, nil
}

// ReadFile reads the content of a file
func (s *Service) ReadFile(path string) (*FileContent, error) {
	relativePath, err := s.CanonicalizePath(path)
	if err != nil {
		return nil, err
	}

	// Resolve the full path
	fullPath := s.resolveFSPath(relativePath)

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
	mimeType := mime.TypeByExtension(filepath.Ext(relativePath))
	if mimeType == "" {
		mimeType = "text/plain"
	}

	return &FileContent{
		Path:     relativePath,
		Name:     filepath.Base(relativePath),
		Content:  string(content),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		MimeType: mimeType,
	}, nil
}

// CanonicalizePath resolves a path into the canonical repo-relative form.
func (s *Service) CanonicalizePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ".", nil
	}

	cleanPath := filepath.Clean(trimmed)
	if filepath.IsAbs(cleanPath) {
		relPath, err := filepath.Rel(s.baseDir, cleanPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve path %q relative to configured codebase root %q: %w", path, s.baseDir, err)
		}
		cleanPath = filepath.Clean(relPath)
	}

	if escapesBaseDir(cleanPath) {
		return "", fmt.Errorf("path %q is outside configured codebase root %q; submit repo-relative paths like %q", path, s.baseDir, "internal/foo.go")
	}

	return toSlash(cleanPath), nil
}

// ResolvePath returns the on-disk absolute path for a canonical or user-supplied path.
func (s *Service) ResolvePath(path string) (string, error) {
	relativePath, err := s.CanonicalizePath(path)
	if err != nil {
		return "", err
	}
	return s.resolveFSPath(relativePath), nil
}

func (s *Service) resolveFSPath(path string) string {
	if path == "." {
		return s.baseDir
	}
	return filepath.Join(s.baseDir, filepath.FromSlash(path))
}

func escapesBaseDir(path string) bool {
	if path == ".." {
		return true
	}
	prefix := ".." + string(filepath.Separator)
	return strings.HasPrefix(path, prefix)
}

func toSlash(path string) string {
	if path == "." {
		return "."
	}
	return filepath.ToSlash(path)
}
