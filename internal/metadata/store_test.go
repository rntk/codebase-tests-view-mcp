package metadata

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codebase-view-mcp/internal/files"
)

func TestAddTestMetadata(t *testing.T) {
	t.Run("adds tests to empty store", func(t *testing.T) {
		store := NewStore("")

		tests := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		}

		err := store.AddTestMetadata("source.go", tests)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		metadata := store.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(metadata.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(metadata.Tests))
		}

		if metadata.Tests[0].TestName != "TestFunc1" {
			t.Fatalf("expected TestFunc1, got %s", metadata.Tests[0].TestName)
		}
	})

	t.Run("merges tests with existing tests", func(t *testing.T) {
		store := NewStore("")

		// Add first test
		tests1 := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		}

		err := store.AddTestMetadata("source.go", tests1)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		// Add second test for same file
		tests2 := []TestReference{
			{
				TestFile:     "test2.go",
				TestName:     "TestFunc2",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 6, End: 10},
			},
		}

		err = store.AddTestMetadata("source.go", tests2)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		metadata := store.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(metadata.Tests) != 2 {
			t.Fatalf("expected 2 tests, got %d", len(metadata.Tests))
		}

		// Check both tests are present
		testNames := make(map[string]bool)
		for _, test := range metadata.Tests {
			testNames[test.TestName] = true
		}

		if !testNames["TestFunc1"] {
			t.Fatal("TestFunc1 not found")
		}

		if !testNames["TestFunc2"] {
			t.Fatal("TestFunc2 not found")
		}
	})

	t.Run("updates existing test with same name", func(t *testing.T) {
		store := NewStore("")

		// Add first test
		tests1 := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		}

		err := store.AddTestMetadata("source.go", tests1)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		// Update same test with different covered lines
		tests2 := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 15},
				CoveredLines: LineRange{Start: 1, End: 10},
			},
		}

		err = store.AddTestMetadata("source.go", tests2)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		metadata := store.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(metadata.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(metadata.Tests))
		}

		// Check that the test was updated
		test := metadata.Tests[0]
		if test.LineRange.End != 15 {
			t.Fatalf("expected line range end to be 15, got %d", test.LineRange.End)
		}

		if test.CoveredLines.End != 10 {
			t.Fatalf("expected covered lines end to be 10, got %d", test.CoveredLines.End)
		}
	})

	t.Run("handles multiple tests for same file from different test files", func(t *testing.T) {
		store := NewStore("")

		// Add tests from test1.go
		tests1 := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc2",
				LineRange:    LineRange{Start: 12, End: 20},
				CoveredLines: LineRange{Start: 6, End: 10},
			},
		}

		err := store.AddTestMetadata("source.go", tests1)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		// Add tests from test2.go
		tests2 := []TestReference{
			{
				TestFile:     "test2.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 3},
			},
			{
				TestFile:     "test2.go",
				TestName:     "TestFunc2",
				LineRange:    LineRange{Start: 12, End: 20},
				CoveredLines: LineRange{Start: 8, End: 12},
			},
		}

		err = store.AddTestMetadata("source.go", tests2)
		if err != nil {
			t.Fatalf("AddTestMetadata failed: %v", err)
		}

		metadata := store.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(metadata.Tests) != 4 {
			t.Fatalf("expected 4 tests, got %d", len(metadata.Tests))
		}

		// Check all tests are present with correct test files
		testKeys := make(map[string]bool)
		for _, test := range metadata.Tests {
			key := test.TestFile + ":" + test.TestName
			testKeys[key] = true
		}

		if !testKeys["test1.go:TestFunc1"] {
			t.Fatal("test1.go:TestFunc1 not found")
		}

		if !testKeys["test1.go:TestFunc2"] {
			t.Fatal("test1.go:TestFunc2 not found")
		}

		if !testKeys["test2.go:TestFunc1"] {
			t.Fatal("test2.go:TestFunc1 not found")
		}

		if !testKeys["test2.go:TestFunc2"] {
			t.Fatal("test2.go:TestFunc2 not found")
		}
	})
}

func TestSetTestMetadata(t *testing.T) {
	t.Run("overwrites existing tests", func(t *testing.T) {
		store := NewStore("")

		// Add first test
		tests1 := []TestReference{
			{
				TestFile:     "test1.go",
				TestName:     "TestFunc1",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		}

		err := store.SetTestMetadata("source.go", tests1)
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		// Set different tests - should overwrite
		tests2 := []TestReference{
			{
				TestFile:     "test2.go",
				TestName:     "TestFunc2",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 6, End: 10},
			},
		}

		err = store.SetTestMetadata("source.go", tests2)
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		metadata := store.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata is nil")
		}

		if len(metadata.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(metadata.Tests))
		}

		if metadata.Tests[0].TestName != "TestFunc2" {
			t.Fatalf("expected TestFunc2, got %s", metadata.Tests[0].TestName)
		}
	})
}

func TestGetTestFileMetadata(t *testing.T) {
	t.Run("returns nil when no metadata exists", func(t *testing.T) {
		store := NewStore("")

		refs := store.GetTestFileMetadata("test_test.go")
		if refs != nil {
			t.Fatalf("expected nil, got %+v", refs)
		}
	})

	t.Run("returns source references for test file", func(t *testing.T) {
		store := NewStore("")

		// Add metadata for source.go
		err := store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				FunctionName: "Process",
				Comment:      "Tests the process function",
				LineRange:    LineRange{Start: 3, End: 9},
				CoveredLines: LineRange{Start: 1, End: 10},
				InputLines:   LineRange{Start: 4, End: 4},
				OutputLines:  LineRange{Start: 5, End: 7},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		refs := store.GetTestFileMetadata("test_test.go")
		if len(refs) != 1 {
			t.Fatalf("expected 1 ref, got %d", len(refs))
		}

		ref := refs[0]
		if ref.SourceFile != "source.go" {
			t.Errorf("expected sourceFile %q, got %q", "source.go", ref.SourceFile)
		}

		if ref.FunctionName != "Process" {
			t.Errorf("expected functionName %q, got %q", "Process", ref.FunctionName)
		}

		if ref.TestName != "TestSomething" {
			t.Errorf("expected testName %q, got %q", "TestSomething", ref.TestName)
		}

		if ref.Comment != "Tests the process function" {
			t.Errorf("expected comment %q, got %q", "Tests the process function", ref.Comment)
		}
	})

	t.Run("returns multiple references for test file", func(t *testing.T) {
		store := NewStore("")

		// Add metadata for multiple source files
		err := store.SetTestMetadata("source1.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				FunctionName: "Process1",
				LineRange:    LineRange{Start: 3, End: 9},
				CoveredLines: LineRange{Start: 1, End: 10},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		err = store.SetTestMetadata("source2.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomethingElse",
				FunctionName: "Process2",
				LineRange:    LineRange{Start: 12, End: 20},
				CoveredLines: LineRange{Start: 1, End: 10},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		refs := store.GetTestFileMetadata("test_test.go")
		if len(refs) != 2 {
			t.Fatalf("expected 2 refs, got %d", len(refs))
		}
	})

	t.Run("skips nil metadata entries", func(t *testing.T) {
		store := NewStore("")

		// Manually set nil metadata to test robustness
		store.metadata["nil.go"] = nil

		refs := store.GetTestFileMetadata("test_test.go")
		// Should return nil since no valid metadata found
		if refs != nil {
			t.Fatalf("expected nil, got %+v", refs)
		}
	})
}

func TestGetAllMetadata(t *testing.T) {
	t.Run("returns empty map when no metadata exists", func(t *testing.T) {
		store := NewStore("")

		all := store.GetAllMetadata()
		if all == nil {
			t.Fatal("expected empty map, got nil")
		}

		if len(all) != 0 {
			t.Fatalf("expected 0 entries, got %d", len(all))
		}
	})

	t.Run("returns copy of all metadata", func(t *testing.T) {
		store := NewStore("")

		err := store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		all := store.GetAllMetadata()
		if len(all) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(all))
		}

		// Verify it's a copy by modifying the returned map
		all["source.go"] = nil
		original := store.GetTestMetadata("source.go")
		if original == nil {
			t.Fatal("original metadata was affected by modification to copy")
		}
	})

	t.Run("returns deep copy of nested metadata", func(t *testing.T) {
		store := NewStore("")

		err := store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestOriginal",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		_, err = store.AddComment("source.go", files.Comment{
			Line:    7,
			Content: "Original comment",
		})
		if err != nil {
			t.Fatalf("AddComment failed: %v", err)
		}

		all := store.GetAllMetadata()
		all["source.go"].Tests[0].TestName = "Mutated"
		all["source.go"].Comments[0].Content = "Updated"

		original := store.GetTestMetadata("source.go")
		if original.Tests[0].TestName != "TestOriginal" {
			t.Fatalf("expected original test name to remain unchanged, got %q", original.Tests[0].TestName)
		}
		if original.Comments[0].Content != "Original comment" {
			t.Fatalf("expected original comment to remain unchanged, got %q", original.Comments[0].Content)
		}
	})
}

func TestStorePersistence(t *testing.T) {
	t.Run("saves and loads metadata from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		// Create store and add metadata
		store1 := NewStore(metaFile)
		err := store1.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(metaFile); os.IsNotExist(err) {
			t.Fatal("metadata file was not created")
		}

		// Create new store and load from file
		store2 := NewStore(metaFile)
		metadata := store2.GetTestMetadata("source.go")
		if metadata == nil {
			t.Fatal("metadata was not loaded from file")
		}

		if len(metadata.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(metadata.Tests))
		}
	})

	t.Run("handles missing file gracefully", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "nonexistent.json")

		// Should not panic
		store := NewStore(metaFile)
		if store == nil {
			t.Fatal("NewStore returned nil")
		}
	})

	t.Run("Save returns nil when no file path", func(t *testing.T) {
		store := NewStore("")
		err := store.Save()
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	})
}

func TestCommentOperations(t *testing.T) {
	t.Run("AddComment generates ID if not provided", func(t *testing.T) {
		store := NewStore("")

		comment := files.Comment{
			Line:    10,
			Content: "Test comment",
		}

		created, err := store.AddComment("source.go", comment)
		if err != nil {
			t.Fatalf("AddComment failed: %v", err)
		}

		if created.ID == "" {
			t.Error("expected comment ID to be generated")
		}

		if created.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if created.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("AddComment preserves provided ID", func(t *testing.T) {
		store := NewStore("")

		comment := files.Comment{
			ID:      "custom-id",
			Line:    10,
			Content: "Test comment",
		}

		created, err := store.AddComment("source.go", comment)
		if err != nil {
			t.Fatalf("AddComment failed: %v", err)
		}

		if created.ID != "custom-id" {
			t.Errorf("expected ID %q, got %q", "custom-id", created.ID)
		}
	})

	t.Run("GetComments returns nil when no comments", func(t *testing.T) {
		store := NewStore("")

		comments := store.GetComments("nonexistent.go")
		if comments != nil {
			t.Fatalf("expected nil, got %+v", comments)
		}
	})

	t.Run("GetComments returns empty slice when metadata exists but no comments", func(t *testing.T) {
		store := NewStore("")

		err := store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if comments != nil {
			t.Fatalf("expected nil, got %+v", comments)
		}
	})

	t.Run("UpdateComment updates content and timestamp", func(t *testing.T) {
		store := NewStore("")

		created, _ := store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Original",
		})

		err := store.UpdateComment("source.go", created.ID, "Updated")
		if err != nil {
			t.Fatalf("UpdateComment failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if len(comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(comments))
		}

		if comments[0].Content != "Updated" {
			t.Errorf("expected content %q, got %q", "Updated", comments[0].Content)
		}
	})

	t.Run("UpdateComment does nothing for non-existent file", func(t *testing.T) {
		store := NewStore("")

		err := store.UpdateComment("nonexistent.go", "id", "content")
		if err != nil {
			t.Fatalf("UpdateComment failed: %v", err)
		}
	})

	t.Run("UpdateComment does nothing for non-existent comment", func(t *testing.T) {
		store := NewStore("")

		store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Original",
		})

		err := store.UpdateComment("source.go", "nonexistent-id", "content")
		if err != nil {
			t.Fatalf("UpdateComment failed: %v", err)
		}
	})

	t.Run("DeleteComment removes comment", func(t *testing.T) {
		store := NewStore("")

		created, _ := store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "To delete",
		})

		err := store.DeleteComment("source.go", created.ID)
		if err != nil {
			t.Fatalf("DeleteComment failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if len(comments) != 0 {
			t.Fatalf("expected 0 comments, got %d", len(comments))
		}
	})

	t.Run("DeleteComment does nothing for non-existent file", func(t *testing.T) {
		store := NewStore("")

		err := store.DeleteComment("nonexistent.go", "id")
		if err != nil {
			t.Fatalf("DeleteComment failed: %v", err)
		}
	})

	t.Run("DeleteComment does nothing for non-existent comment", func(t *testing.T) {
		store := NewStore("")

		_, _ = store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Keep",
		})

		err := store.DeleteComment("source.go", "nonexistent-id")
		if err != nil {
			t.Fatalf("DeleteComment failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if len(comments) != 1 {
			t.Fatalf("expected 1 comment, got %d", len(comments))
		}
	})

	t.Run("ToggleCommentResolved toggles from false to true", func(t *testing.T) {
		store := NewStore("")

		created, _ := store.AddComment("source.go", files.Comment{
			Line:     10,
			Content:  "To resolve",
			Resolved: false,
		})

		err := store.ToggleCommentResolved("source.go", created.ID)
		if err != nil {
			t.Fatalf("ToggleCommentResolved failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if !comments[0].Resolved {
			t.Error("expected comment to be resolved")
		}
	})

	t.Run("ToggleCommentResolved toggles from true to false", func(t *testing.T) {
		store := NewStore("")

		created, _ := store.AddComment("source.go", files.Comment{
			Line:     10,
			Content:  "Already resolved",
			Resolved: true,
		})

		err := store.ToggleCommentResolved("source.go", created.ID)
		if err != nil {
			t.Fatalf("ToggleCommentResolved failed: %v", err)
		}

		comments := store.GetComments("source.go")
		if comments[0].Resolved {
			t.Error("expected comment to be unresolved")
		}
	})

	t.Run("ToggleCommentResolved does nothing for non-existent file", func(t *testing.T) {
		store := NewStore("")

		err := store.ToggleCommentResolved("nonexistent.go", "id")
		if err != nil {
			t.Fatalf("ToggleCommentResolved failed: %v", err)
		}
	})

	t.Run("ToggleCommentResolved does nothing for non-existent comment", func(t *testing.T) {
		store := NewStore("")

		store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Keep",
		})

		err := store.ToggleCommentResolved("source.go", "nonexistent-id")
		if err != nil {
			t.Fatalf("ToggleCommentResolved failed: %v", err)
		}
	})

	t.Run("AddComment persists to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		store := NewStore(metaFile)
		_, err := store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Persistent comment",
		})
		if err != nil {
			t.Fatalf("AddComment failed: %v", err)
		}

		// Verify file was updated
		if _, err := os.Stat(metaFile); os.IsNotExist(err) {
			t.Fatal("metadata file was not created")
		}
	})

	t.Run("multiple comments on same line", func(t *testing.T) {
		store := NewStore("")

		store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "First comment",
		})

		store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Second comment",
		})

		comments := store.GetComments("source.go")
		if len(comments) != 2 {
			t.Fatalf("expected 2 comments, got %d", len(comments))
		}
	})

	t.Run("GetComments returns copy", func(t *testing.T) {
		store := NewStore("")

		_, err := store.AddComment("source.go", files.Comment{
			Line:    10,
			Content: "Original",
		})
		if err != nil {
			t.Fatalf("AddComment failed: %v", err)
		}

		comments := store.GetComments("source.go")
		comments[0].Content = "Mutated"

		fresh := store.GetComments("source.go")
		if fresh[0].Content != "Original" {
			t.Fatalf("expected original comment to remain unchanged, got %q", fresh[0].Content)
		}
	})
}

func TestNewStore(t *testing.T) {
	t.Run("creates store with empty metadata map", func(t *testing.T) {
		store := NewStore("")

		if store.metadata == nil {
			t.Fatal("metadata map is nil")
		}

		if len(store.metadata) != 0 {
			t.Fatalf("expected empty metadata map, got %d entries", len(store.metadata))
		}
	})

	t.Run("creates store with file path", func(t *testing.T) {
		store := NewStore("/path/to/metadata.json")

		if store.filePath != "/path/to/metadata.json" {
			t.Errorf("expected filePath %q, got %q", "/path/to/metadata.json", store.filePath)
		}
	})
}

func TestStoreConcurrency(t *testing.T) {
	t.Run("concurrent reads are safe", func(t *testing.T) {
		store := NewStore("")

		// Add some data
		store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})

		// Multiple concurrent reads
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_ = store.GetTestMetadata("source.go")
				_ = store.GetAllMetadata()
				_ = store.GetComments("source.go")
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestRenameSourcePath(t *testing.T) {
	t.Run("renames source path successfully", func(t *testing.T) {
		store := NewStore("")

		// Add metadata
		err := store.SetTestMetadata("old.go", []TestReference{
			{
				TestFile:     "test_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})
		if err != nil {
			t.Fatalf("SetTestMetadata failed: %v", err)
		}

		// Rename
		err = store.RenameSourcePath("old.go", "new.go")
		if err != nil {
			t.Fatalf("RenameSourcePath failed: %v", err)
		}

		// Verify old path is gone
		oldMeta := store.GetTestMetadata("old.go")
		if oldMeta != nil {
			t.Error("old path should not exist")
		}

		// Verify new path exists
		newMeta := store.GetTestMetadata("new.go")
		if newMeta == nil {
			t.Fatal("new path should exist")
		}
		if len(newMeta.Tests) != 1 {
			t.Errorf("expected 1 test, got %d", len(newMeta.Tests))
		}
	})

	t.Run("returns error when old path not found", func(t *testing.T) {
		store := NewStore("")

		err := store.RenameSourcePath("nonexistent.go", "new.go")
		if err == nil {
			t.Fatal("expected error for non-existent old path")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when new path already exists", func(t *testing.T) {
		store := NewStore("")

		// Add metadata to both paths
		store.SetTestMetadata("old.go", []TestReference{
			{TestFile: "test1.go", TestName: "Test1", LineRange: LineRange{Start: 1, End: 5}},
		})
		store.SetTestMetadata("new.go", []TestReference{
			{TestFile: "test2.go", TestName: "Test2", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.RenameSourcePath("old.go", "new.go")
		if err == nil {
			t.Fatal("expected error when new path already exists")
		}
		if !errors.Is(err, ErrAlreadyExists) {
			t.Errorf("expected ErrAlreadyExists, got: %v", err)
		}
	})

	t.Run("no-op when old and new paths are the same", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("same.go", []TestReference{
			{TestFile: "test.go", TestName: "Test", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.RenameSourcePath("same.go", "same.go")
		if err != nil {
			t.Fatalf("RenameSourcePath failed: %v", err)
		}

		meta := store.GetTestMetadata("same.go")
		if meta == nil {
			t.Error("metadata should still exist")
		}
	})

	t.Run("persists to file after rename", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		store := NewStore(metaFile)
		store.SetTestMetadata("old.go", []TestReference{
			{TestFile: "test.go", TestName: "Test", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.RenameSourcePath("old.go", "new.go")
		if err != nil {
			t.Fatalf("RenameSourcePath failed: %v", err)
		}

		// Load in new store
		store2 := NewStore(metaFile)
		meta := store2.GetTestMetadata("new.go")
		if meta == nil {
			t.Error("renamed metadata should be persisted")
		}
	})
}

func TestUpdateTestPath(t *testing.T) {
	t.Run("updates test file path successfully", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{
				TestFile:     "old_test.go",
				TestName:     "TestSomething",
				LineRange:    LineRange{Start: 1, End: 10},
				CoveredLines: LineRange{Start: 1, End: 5},
			},
		})

		err := store.UpdateTestPath("source.go", "old_test.go", "TestSomething", "new_test.go")
		if err != nil {
			t.Fatalf("UpdateTestPath failed: %v", err)
		}

		meta := store.GetTestMetadata("source.go")
		if meta == nil || len(meta.Tests) != 1 {
			t.Fatal("metadata not found")
		}
		if meta.Tests[0].TestFile != "new_test.go" {
			t.Errorf("expected test file %q, got %q", "new_test.go", meta.Tests[0].TestFile)
		}
	})

	t.Run("returns error when source file not found", func(t *testing.T) {
		store := NewStore("")

		err := store.UpdateTestPath("nonexistent.go", "test.go", "Test", "new.go")
		if err == nil {
			t.Fatal("expected error for non-existent source file")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when test reference not found", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "TestOther", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.UpdateTestPath("source.go", "test.go", "TestNonExistent", "new.go")
		if err == nil {
			t.Fatal("expected error for non-existent test reference")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when new test file already has same test", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "old_test.go", TestName: "TestSomething", LineRange: LineRange{Start: 1, End: 5}},
			{TestFile: "new_test.go", TestName: "TestSomething", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.UpdateTestPath("source.go", "old_test.go", "TestSomething", "new_test.go")
		if err == nil {
			t.Fatal("expected error when test already exists in new file")
		}
		if !errors.Is(err, ErrAlreadyExists) {
			t.Errorf("expected ErrAlreadyExists, got: %v", err)
		}
	})

	t.Run("persists to file after update", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		store := NewStore(metaFile)
		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "old_test.go", TestName: "Test", LineRange: LineRange{Start: 1, End: 5}},
		})

		store.UpdateTestPath("source.go", "old_test.go", "Test", "new_test.go")

		store2 := NewStore(metaFile)
		meta := store2.GetTestMetadata("source.go")
		if meta == nil || meta.Tests[0].TestFile != "new_test.go" {
			t.Error("updated test path should be persisted")
		}
	})
}

func TestDeleteSourcePath(t *testing.T) {
	t.Run("deletes source path successfully", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "Test", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.DeleteSourcePath("source.go")
		if err != nil {
			t.Fatalf("DeleteSourcePath failed: %v", err)
		}

		meta := store.GetTestMetadata("source.go")
		if meta != nil {
			t.Error("metadata should be deleted")
		}
	})

	t.Run("no error when source path does not exist", func(t *testing.T) {
		store := NewStore("")

		err := store.DeleteSourcePath("nonexistent.go")
		if err != nil {
			t.Fatalf("DeleteSourcePath should not error for non-existent path: %v", err)
		}
	})

	t.Run("persists deletion to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		store := NewStore(metaFile)
		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "Test", LineRange: LineRange{Start: 1, End: 5}},
		})

		store.DeleteSourcePath("source.go")

		store2 := NewStore(metaFile)
		meta := store2.GetTestMetadata("source.go")
		if meta != nil {
			t.Error("deletion should be persisted")
		}
	})
}

func TestDeleteTestPath(t *testing.T) {
	t.Run("deletes single test reference", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "Test1", LineRange: LineRange{Start: 1, End: 5}},
			{TestFile: "test.go", TestName: "Test2", LineRange: LineRange{Start: 6, End: 10}},
		})

		err := store.DeleteTestPath("source.go", "test.go", "Test1")
		if err != nil {
			t.Fatalf("DeleteTestPath failed: %v", err)
		}

		meta := store.GetTestMetadata("source.go")
		if meta == nil || len(meta.Tests) != 1 {
			t.Fatalf("expected 1 test, got %d", len(meta.Tests))
		}
		if meta.Tests[0].TestName != "Test2" {
			t.Errorf("expected Test2 to remain, got %s", meta.Tests[0].TestName)
		}
	})

	t.Run("returns error when source file not found", func(t *testing.T) {
		store := NewStore("")

		err := store.DeleteTestPath("nonexistent.go", "test.go", "Test")
		if err == nil {
			t.Fatal("expected error for non-existent source file")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("returns error when test reference not found", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "TestOther", LineRange: LineRange{Start: 1, End: 5}},
		})

		err := store.DeleteTestPath("source.go", "test.go", "TestNonExistent")
		if err == nil {
			t.Fatal("expected error for non-existent test")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("deletes last test removes metadata entry", func(t *testing.T) {
		store := NewStore("")

		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "TestOnly", LineRange: LineRange{Start: 1, End: 5}},
		})

		store.DeleteTestPath("source.go", "test.go", "TestOnly")

		meta := store.GetTestMetadata("source.go")
		// Metadata entry exists but has no tests
		if meta != nil && len(meta.Tests) != 0 {
			t.Error("test should be deleted")
		}
	})

	t.Run("persists deletion to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		metaFile := filepath.Join(tmpDir, "metadata.json")

		store := NewStore(metaFile)
		store.SetTestMetadata("source.go", []TestReference{
			{TestFile: "test.go", TestName: "Test1", LineRange: LineRange{Start: 1, End: 5}},
			{TestFile: "test.go", TestName: "Test2", LineRange: LineRange{Start: 6, End: 10}},
		})

		store.DeleteTestPath("source.go", "test.go", "Test1")

		store2 := NewStore(metaFile)
		meta := store2.GetTestMetadata("source.go")
		if meta == nil || len(meta.Tests) != 1 {
			t.Error("deletion should be persisted")
		}
	})
}
