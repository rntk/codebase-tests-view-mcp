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

// MaxFileSize limits the maximum file size that can be read (10MB)
const MaxFileSize = 10 * 1024 * 1024

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

	// Check for hidden directories in the path
	// e.g., ".git", ".vscode", "src/.hidden"
	if containsHiddenComponent(relativePath) {
		return nil, fmt.Errorf("path not found: %w", os.ErrNotExist)
	}

	// Resolve the full path
	fullPath := s.resolveFSPath(relativePath)

	// Validate path (check for symlink escapes) - expect a directory
	validatedPath, err := s.resolveAndValidatePath(fullPath, true)
	if err != nil {
		return nil, err
	}

	// Check if path exists and is a directory
	info, err := os.Stat(validatedPath)
	if err != nil {
		return nil, fmt.Errorf("path not found: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory")
	}

	// Read directory contents
	entries, err := os.ReadDir(validatedPath)
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
		if isHiddenFile(entry.Name()) {
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

	// Check for hidden files or files in hidden directories
	// e.g., ".env", ".git/config", "src/.hidden/file.txt"
	if containsHiddenComponent(relativePath) {
		return nil, fmt.Errorf("file not found: %w", os.ErrNotExist)
	}

	// Resolve the full path
	fullPath := s.resolveFSPath(relativePath)

	// Validate path (check for symlink escapes) - expect a file
	validatedPath, err := s.resolveAndValidatePath(fullPath, false)
	if err != nil {
		return nil, err
	}

	// Get file info (already validated as file in resolveAndValidatePath)
	info, err := os.Stat(validatedPath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Check file size limit
	if info.Size() > MaxFileSize {
		return nil, fmt.Errorf("file exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	// Read file content
	file, err := os.Open(validatedPath)
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

// isHiddenFile checks if a file or directory name is hidden (starts with .)
func isHiddenFile(name string) bool {
	return strings.HasPrefix(name, ".")
}

// containsHiddenComponent checks if any component in the path is hidden
// e.g., ".git/config", "src/.env", ".vscode/settings.json" all return true
func containsHiddenComponent(path string) bool {
	// Clean the path and split into components
	cleanPath := filepath.Clean(path)

	// Handle special cases
	if cleanPath == "." {
		return false
	}

	// Split path into components and check each one
	parts := strings.Split(cleanPath, string(filepath.Separator))
	for _, part := range parts {
		if part != "" && isHiddenFile(part) {
			return true
		}
	}
	return false
}

// resolveAndValidatePath resolves ALL symlinks (including intermediate components)
// and validates the final resolved path stays within baseDir.
// Returns the resolved real path and an error if the path is invalid.
// For files: validates that path is not a directory and symlinks don't escape
// For directories: validates path exists and symlinks don't escape to outside base
func (s *Service) resolveAndValidatePath(fullPath string, expectDir bool) (string, error) {
	// First, resolve ALL symlinks in the path (including intermediate components)
	// This is critical for security: a path like "subdir/file.txt" where subdir is
	// a symlink to /tmp/outside would otherwise bypass our checks.
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path not found: %w", err)
		}
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// Get info about the resolved path
	info, err := os.Stat(realPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path not found: %w", err)
		}
		return "", fmt.Errorf("failed to access path: %w", err)
	}

	// Check if resolved path type matches expectation
	if expectDir && !info.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}
	if !expectDir && info.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file")
	}

	realBaseDir, err := s.resolveRealBaseDir()
	if err != nil {
		return "", err
	}

	// Verify the fully resolved path is within the canonical baseDir
	relPath, err := filepath.Rel(realBaseDir, realPath)
	if err != nil {
		return "", fmt.Errorf("failed to validate path: %w", err)
	}

	if escapesBaseDir(relPath) {
		if expectDir {
			return "", fmt.Errorf("path not found: %w", os.ErrNotExist)
		}
		return "", fmt.Errorf("file not found: %w", os.ErrNotExist)
	}
	if containsHiddenComponent(relPath) {
		if expectDir {
			return "", fmt.Errorf("path not found: %w", os.ErrNotExist)
		}
		return "", fmt.Errorf("file not found: %w", os.ErrNotExist)
	}

	return realPath, nil
}

func (s *Service) resolveRealBaseDir() (string, error) {
	realBaseDir, err := filepath.EvalSymlinks(s.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("configured codebase root not found: %w", err)
		}
		return "", fmt.Errorf("failed to resolve configured codebase root: %w", err)
	}

	return realBaseDir, nil
}

func toSlash(path string) string {
	if path == "." {
		return "."
	}
	return filepath.ToSlash(path)
}
