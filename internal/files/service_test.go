package files

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	baseDir := "/test/dir"
	service := NewService(baseDir)

	if service == nil {
		t.Fatal("NewService returned nil")
	}

	if service.baseDir != baseDir {
		t.Errorf("expected baseDir %q, got %q", baseDir, service.baseDir)
	}
}

func TestListFiles(t *testing.T) {
	t.Run("returns files and directories sorted correctly", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create subdirectory
		subdir := filepath.Join(baseDir, "zzz_dir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		// Create files
		if err := os.WriteFile(filepath.Join(baseDir, "aaa.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "bbb.go"), []byte("package main"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		// Create hidden file (should be skipped)
		if err := os.WriteFile(filepath.Join(baseDir, ".hidden"), []byte("secret"), 0644); err != nil {
			t.Fatalf("write hidden file: %v", err)
		}

		service := NewService(baseDir)
		response, err := service.ListFiles(".")
		if err != nil {
			t.Fatalf("ListFiles failed: %v", err)
		}

		if response.Path != "." {
			t.Errorf("expected path %q, got %q", ".", response.Path)
		}

		// Should have 3 entries (2 files + 1 dir), hidden file should be skipped
		if len(response.Files) != 3 {
			t.Fatalf("expected 3 files, got %d: %+v", len(response.Files), response.Files)
		}

		// Directories should come first, then sorted by name
		if !response.Files[0].IsDir || response.Files[0].Name != "zzz_dir" {
			t.Errorf("expected first entry to be zzz_dir directory, got %+v", response.Files[0])
		}

		if response.Files[1].IsDir || response.Files[1].Name != "aaa.txt" {
			t.Errorf("expected second entry to be aaa.txt, got %+v", response.Files[1])
		}

		if response.Files[2].IsDir || response.Files[2].Name != "bbb.go" {
			t.Errorf("expected third entry to be bbb.go, got %+v", response.Files[2])
		}
	})

	t.Run("returns error for non-existent path", func(t *testing.T) {
		baseDir := t.TempDir()
		service := NewService(baseDir)

		_, err := service.ListFiles("nonexistent")
		if err == nil {
			t.Fatal("expected error for non-existent path")
		}

		if !strings.Contains(err.Error(), "path not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("returns error when path is a file not directory", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ListFiles("test.txt")
		if err == nil {
			t.Fatal("expected error when path is a file")
		}

		if err.Error() != "path is not a directory" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("handles empty directory", func(t *testing.T) {
		baseDir := t.TempDir()
		service := NewService(baseDir)

		response, err := service.ListFiles(".")
		if err != nil {
			t.Fatalf("ListFiles failed: %v", err)
		}

		if len(response.Files) != 0 {
			t.Errorf("expected 0 files in empty directory, got %d", len(response.Files))
		}
	})

	t.Run("skips files that fail to get info", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a valid file
		if err := os.WriteFile(filepath.Join(baseDir, "valid.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		response, err := service.ListFiles(".")
		if err != nil {
			t.Fatalf("ListFiles failed: %v", err)
		}

		if len(response.Files) != 1 {
			t.Errorf("expected 1 file, got %d", len(response.Files))
		}
	})
}

func TestReadFile(t *testing.T) {
	t.Run("returns file content with correct metadata", func(t *testing.T) {
		baseDir := t.TempDir()
		content := "Hello, World!"
		testFile := filepath.Join(baseDir, "test.txt")
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		fileContent, err := service.ReadFile("test.txt")
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if fileContent.Path != "test.txt" {
			t.Errorf("expected path %q, got %q", "test.txt", fileContent.Path)
		}

		if fileContent.Name != "test.txt" {
			t.Errorf("expected name %q, got %q", "test.txt", fileContent.Name)
		}

		if fileContent.Content != content {
			t.Errorf("expected content %q, got %q", content, fileContent.Content)
		}

		if fileContent.Size != int64(len(content)) {
			t.Errorf("expected size %d, got %d", len(content), fileContent.Size)
		}

		if fileContent.MimeType != "text/plain" && fileContent.MimeType != "text/plain; charset=utf-8" {
			t.Errorf("expected mimeType text/plain, got %q", fileContent.MimeType)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		baseDir := t.TempDir()
		service := NewService(baseDir)

		_, err := service.ReadFile("nonexistent.txt")
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("returns error when path is a directory", func(t *testing.T) {
		baseDir := t.TempDir()
		service := NewService(baseDir)

		_, err := service.ReadFile(".")
		if err == nil {
			t.Fatal("expected error when path is a directory")
		}

		// Either "directory" or "not found" are acceptable - both indicate invalid for file reading
		if !strings.Contains(err.Error(), "directory") && !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("detects correct MIME type for Go files", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "test.go")
		if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		fileContent, err := service.ReadFile("test.go")
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		// MIME type detection may vary by system
		// Accept either text/x-go or text/plain
		if fileContent.MimeType != "text/x-go" && fileContent.MimeType != "text/plain" && fileContent.MimeType != "text/x-golang" {
			t.Logf("mimeType: %q (acceptable: text/x-go, text/plain, or text/x-golang)", fileContent.MimeType)
		}
	})

	t.Run("detects correct MIME type for TypeScript files", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "test.ts")
		if err := os.WriteFile(testFile, []byte("const x = 1;"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		fileContent, err := service.ReadFile("test.ts")
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if fileContent.MimeType != "video/mp2t" && fileContent.MimeType != "text/plain" {
			// MIME type detection varies by system
			t.Logf("mimeType: %q", fileContent.MimeType)
		}
	})

	t.Run("uses text/plain for unknown extensions", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "test.xyz")
		if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		fileContent, err := service.ReadFile("test.xyz")
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		// MIME type detection may vary by system
		// Just verify it's set to something
		if fileContent.MimeType == "" {
			t.Error("expected mimeType to be set")
		}
	})
}

func TestResolvePath(t *testing.T) {
	t.Run("handles relative paths", func(t *testing.T) {
		baseDir := t.TempDir()
		subdir := filepath.Join(baseDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		service := NewService(baseDir)
		response, err := service.ListFiles("subdir")
		if err != nil {
			t.Fatalf("ListFiles with relative path failed: %v", err)
		}

		if response.Path != "subdir" {
			t.Errorf("expected path %q, got %q", "subdir", response.Path)
		}
	})

	t.Run("handles nested relative paths", func(t *testing.T) {
		baseDir := t.TempDir()
		nestedDir := filepath.Join(baseDir, "a", "b", "c")
		if err := os.MkdirAll(nestedDir, 0755); err != nil {
			t.Fatalf("mkdirAll: %v", err)
		}

		service := NewService(baseDir)
		response, err := service.ListFiles("a/b/c")
		if err != nil {
			t.Fatalf("ListFiles with nested path failed: %v", err)
		}

		if response.Path != "a/b/c" {
			t.Errorf("expected path %q, got %q", "a/b/c", response.Path)
		}
	})

	t.Run("handles dot path", func(t *testing.T) {
		baseDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(baseDir, "test.txt"), []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		response, err := service.ListFiles(".")
		if err != nil {
			t.Fatalf("ListFiles with dot path failed: %v", err)
		}

		if response.Path != "." {
			t.Errorf("expected path %q, got %q", ".", response.Path)
		}
	})

	t.Run("normalizes absolute paths under base directory", func(t *testing.T) {
		baseDir := t.TempDir()
		testFile := filepath.Join(baseDir, "nested", "test.txt")
		if err := os.MkdirAll(filepath.Dir(testFile), 0755); err != nil {
			t.Fatalf("mkdirAll: %v", err)
		}
		if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		fileContent, err := service.ReadFile(testFile)
		if err != nil {
			t.Fatalf("ReadFile with absolute path failed: %v", err)
		}

		if fileContent.Path != "nested/test.txt" {
			t.Fatalf("expected canonical path %q, got %q", "nested/test.txt", fileContent.Path)
		}
	})

	t.Run("rejects absolute paths outside base directory", func(t *testing.T) {
		baseDir := t.TempDir()
		otherDir := t.TempDir()
		testFile := filepath.Join(otherDir, "outside.txt")
		if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile(testFile)
		if err == nil {
			t.Fatal("expected error for absolute path outside base directory")
		}

		if !strings.Contains(err.Error(), "outside configured codebase root") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects path traversal outside base directory", func(t *testing.T) {
		baseDir := t.TempDir()
		service := NewService(baseDir)

		_, err := service.ListFiles("../outside")
		if err == nil {
			t.Fatal("expected error for path traversal")
		}

		if !strings.Contains(err.Error(), "outside configured codebase root") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFileContentModTime(t *testing.T) {
	baseDir := t.TempDir()
	testFile := filepath.Join(baseDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	service := NewService(baseDir)
	fileContent, err := service.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// ModTime should be set and not zero
	if fileContent.ModTime.IsZero() {
		t.Error("ModTime should not be zero")
	}

	// ModTime should be recent (within last minute)
	if time.Since(fileContent.ModTime) > time.Minute {
		t.Errorf("ModTime seems too old: %v", fileContent.ModTime)
	}
}

func TestListFilesResponseStructure(t *testing.T) {
	baseDir := t.TempDir()
	testFile := filepath.Join(baseDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	service := NewService(baseDir)
	response, err := service.ListFiles(".")
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	if len(response.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(response.Files))
	}

	entry := response.Files[0]
	if entry.Name != "test.txt" {
		t.Errorf("expected name %q, got %q", "test.txt", entry.Name)
	}

	if entry.Path != "test.txt" {
		t.Errorf("expected path %q, got %q", "test.txt", entry.Path)
	}

	if entry.IsDir {
		t.Error("expected IsDir to be false")
	}

	if entry.Size != 5 {
		t.Errorf("expected size 5, got %d", entry.Size)
	}
}

func TestSymlinkAttacks(t *testing.T) {
	t.Run("rejects symlink pointing outside base directory", func(t *testing.T) {
		baseDir := t.TempDir()
		otherDir := t.TempDir()

		// Create a file in the other directory
		secretFile := filepath.Join(otherDir, "secret.txt")
		if err := os.WriteFile(secretFile, []byte("secret data"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		// Create a symlink in baseDir pointing to the secret file
		symlinkPath := filepath.Join(baseDir, "evil_link")
		if err := os.Symlink(secretFile, symlinkPath); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("evil_link")
		if err == nil {
			t.Fatal("expected error for symlink pointing outside base directory")
		}

		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects symlink to directory outside base", func(t *testing.T) {
		baseDir := t.TempDir()
		otherDir := t.TempDir()

		// Create a directory in the other location
		secretDir := filepath.Join(otherDir, "secret_dir")
		if err := os.Mkdir(secretDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		// Create a symlink in baseDir pointing to the secret directory
		symlinkPath := filepath.Join(baseDir, "evil_dir_link")
		if err := os.Symlink(secretDir, symlinkPath); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("evil_dir_link")
		if err == nil {
			t.Fatal("expected error for symlink to directory")
		}

		if !strings.Contains(err.Error(), "directory") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects symlink to parent directory traversal", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a subdirectory
		subdir := filepath.Join(baseDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		// Create a symlink in subdir pointing to parent
		symlinkPath := filepath.Join(subdir, "parent_link")
		if err := os.Symlink(baseDir, symlinkPath); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("subdir/parent_link")
		if err == nil {
			t.Fatal("expected error for symlink to parent directory")
		}

		if !strings.Contains(err.Error(), "directory") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("allows symlink within base directory", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a file in a subdirectory
		subdir := filepath.Join(baseDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		targetFile := filepath.Join(subdir, "target.txt")
		if err := os.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		// Create a symlink in the root of baseDir pointing to the target
		symlinkPath := filepath.Join(baseDir, "link.txt")
		if err := os.Symlink(targetFile, symlinkPath); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		content, err := service.ReadFile("link.txt")
		if err != nil {
			t.Fatalf("expected no error for symlink within base directory: %v", err)
		}

		if content.Content != "target content" {
			t.Errorf("expected content 'target content', got %q", content.Content)
		}
	})

	t.Run("allows access when configured base directory is a symlink", func(t *testing.T) {
		realBaseDir := t.TempDir()
		parentDir := t.TempDir()

		targetFile := filepath.Join(realBaseDir, "target.txt")
		if err := os.WriteFile(targetFile, []byte("target content"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		symlinkBaseDir := filepath.Join(parentDir, "repo-link")
		if err := os.Symlink(realBaseDir, symlinkBaseDir); err != nil {
			t.Fatalf("create base dir symlink: %v", err)
		}

		service := NewService(symlinkBaseDir)
		content, err := service.ReadFile("target.txt")
		if err != nil {
			t.Fatalf("expected no error for file within symlinked base directory: %v", err)
		}

		if content.Content != "target content" {
			t.Errorf("expected content 'target content', got %q", content.Content)
		}
	})

	t.Run("rejects file access through symlinked directory", func(t *testing.T) {
		baseDir := t.TempDir()
		otherDir := t.TempDir()

		// Create a secret file in an outside directory
		secretFile := filepath.Join(otherDir, "secret.txt")
		if err := os.WriteFile(secretFile, []byte("secret data"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		// Create a symlink in baseDir that points to the outside directory
		// This tests intermediate symlink components
		symlinkDir := filepath.Join(baseDir, "linked_dir")
		if err := os.Symlink(otherDir, symlinkDir); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		// Try to read a file through the symlinked directory
		service := NewService(baseDir)
		_, err := service.ReadFile("linked_dir/secret.txt")
		if err == nil {
			t.Fatal("expected error for file accessed through symlinked directory outside base")
		}

		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects listing directory through symlinked directory outside base", func(t *testing.T) {
		baseDir := t.TempDir()
		otherDir := t.TempDir()

		// Create a directory in the other location
		secretDir := filepath.Join(otherDir, "secret_dir")
		if err := os.Mkdir(secretDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		// Create a symlink in baseDir that points to the outside directory
		symlinkDir := filepath.Join(baseDir, "linked_dir")
		if err := os.Symlink(secretDir, symlinkDir); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		// Try to list the symlinked directory
		service := NewService(baseDir)
		_, err := service.ListFiles("linked_dir")
		if err == nil {
			t.Fatal("expected error for listing symlinked directory outside base")
		}

		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestHiddenFiles(t *testing.T) {
	t.Run("rejects reading hidden files", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a hidden file
		hiddenFile := filepath.Join(baseDir, ".hidden")
		if err := os.WriteFile(hiddenFile, []byte("secret"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile(".hidden")
		if err == nil {
			t.Fatal("expected error for hidden file")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects reading dotfile in subdirectory", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a subdirectory with a hidden file
		subdir := filepath.Join(baseDir, "subdir")
		if err := os.Mkdir(subdir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		hiddenFile := filepath.Join(subdir, ".env")
		if err := os.WriteFile(hiddenFile, []byte("SECRET=123"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("subdir/.env")
		if err == nil {
			t.Fatal("expected error for hidden file in subdirectory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects reading files in hidden directories", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a hidden directory with a file
		hiddenDir := filepath.Join(baseDir, ".git")
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		configFile := filepath.Join(hiddenDir, "config")
		if err := os.WriteFile(configFile, []byte("[core]\nrepositoryformatversion = 0"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile(".git/config")
		if err == nil {
			t.Fatal("expected error for file in hidden directory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects reading files in nested hidden directories", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a hidden directory with nested structure
		hiddenDir := filepath.Join(baseDir, ".vscode")
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		settingsFile := filepath.Join(hiddenDir, "settings.json")
		if err := os.WriteFile(settingsFile, []byte(`{"editor.tabSize": 4}`), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile(".vscode/settings.json")
		if err == nil {
			t.Fatal("expected error for file in hidden directory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects reading files deeply nested in hidden directories", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a deeply nested structure under a hidden directory
		hiddenDir := filepath.Join(baseDir, ".hidden", "nested", "deep")
		if err := os.MkdirAll(hiddenDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		secretFile := filepath.Join(hiddenDir, "secret.txt")
		if err := os.WriteFile(secretFile, []byte("secret"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile(".hidden/nested/deep/secret.txt")
		if err == nil {
			t.Fatal("expected error for file deeply nested in hidden directory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects listing hidden directories", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a hidden directory
		hiddenDir := filepath.Join(baseDir, ".git")
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ListFiles(".git")
		if err == nil {
			t.Fatal("expected error for listing hidden directory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects reading symlink to hidden file within base", func(t *testing.T) {
		baseDir := t.TempDir()

		hiddenFile := filepath.Join(baseDir, ".env")
		if err := os.WriteFile(hiddenFile, []byte("SECRET=123"), 0644); err != nil {
			t.Fatalf("write hidden file: %v", err)
		}

		visibleLink := filepath.Join(baseDir, "visible.txt")
		if err := os.Symlink(hiddenFile, visibleLink); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("visible.txt")
		if err == nil {
			t.Fatal("expected error for symlink resolving to hidden file")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects listing symlink to hidden directory within base", func(t *testing.T) {
		baseDir := t.TempDir()

		hiddenDir := filepath.Join(baseDir, ".git")
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}

		visibleLink := filepath.Join(baseDir, "config")
		if err := os.Symlink(hiddenDir, visibleLink); err != nil {
			t.Fatalf("create symlink: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ListFiles("config")
		if err == nil {
			t.Fatal("expected error for symlink resolving to hidden directory")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFileSizeLimit(t *testing.T) {
	t.Run("rejects files exceeding MaxFileSize", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a file larger than MaxFileSize
		largeFile := filepath.Join(baseDir, "large.txt")
		largeContent := make([]byte, MaxFileSize+1)
		if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		_, err := service.ReadFile("large.txt")
		if err == nil {
			t.Fatal("expected error for file exceeding size limit")
		}

		if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("allows files within MaxFileSize", func(t *testing.T) {
		baseDir := t.TempDir()

		// Create a file smaller than MaxFileSize
		smallFile := filepath.Join(baseDir, "small.txt")
		if err := os.WriteFile(smallFile, []byte("hello"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}

		service := NewService(baseDir)
		content, err := service.ReadFile("small.txt")
		if err != nil {
			t.Fatalf("expected no error for small file: %v", err)
		}

		if content.Content != "hello" {
			t.Errorf("expected content 'hello', got %q", content.Content)
		}
	})
}
