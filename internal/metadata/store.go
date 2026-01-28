package metadata

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"codebase-view-mcp/internal/files"
	"github.com/google/uuid"
)

// Store manages test metadata storage
type Store struct {
	mu       sync.RWMutex
	metadata map[string]*FileMetadata // key: file path
	filePath string                   // path to JSON persistence file
}

// NewStore creates a new metadata store
func NewStore(persistPath string) *Store {
	store := &Store{
		metadata: make(map[string]*FileMetadata),
		filePath: persistPath,
	}

	if persistPath != "" {
		if err := store.load(); err != nil {
			log.Printf("Warning: failed to load metadata from %s: %v", persistPath, err)
		}
	}

	return store
}

// SetTestMetadata stores test metadata for a file (overwrites existing data)
func (s *Store) SetTestMetadata(filePath string, tests []TestReference) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metadata[filePath] = &FileMetadata{Tests: tests}

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}

// AddTestMetadata adds test metadata to a file, merging with existing tests
func (s *Store) AddTestMetadata(filePath string, tests []TestReference) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.metadata[filePath]
	if existing == nil {
		// No existing metadata, create new
		s.metadata[filePath] = &FileMetadata{Tests: tests}
	} else {
		// Merge with existing tests
		// Use a map to deduplicate based on testFile+testName
		testMap := make(map[string]TestReference)

		// Add existing tests first
		for _, test := range existing.Tests {
			key := test.TestFile + ":" + test.TestName
			testMap[key] = test
		}

		// Add/update with new tests
		for _, test := range tests {
			key := test.TestFile + ":" + test.TestName
			testMap[key] = test
		}

		// Convert map back to slice
		mergedTests := make([]TestReference, 0, len(testMap))
		for _, test := range testMap {
			mergedTests = append(mergedTests, test)
		}

		s.metadata[filePath].Tests = mergedTests
	}

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}

// GetTestMetadata retrieves test metadata for a file
func (s *Store) GetTestMetadata(filePath string) *FileMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.metadata[filePath]
}

// GetAllMetadata returns all stored metadata
func (s *Store) GetAllMetadata() map[string]*FileMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	copy := make(map[string]*FileMetadata, len(s.metadata))
	for k, v := range s.metadata {
		copy[k] = v
	}

	return copy
}

// load loads metadata from the JSON file
func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's OK
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return json.Unmarshal(data, &s.metadata)
}

// saveUnsafe saves metadata to the JSON file (must be called with lock held)
func (s *Store) saveUnsafe() error {
	data, err := json.MarshalIndent(s.metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// Save explicitly saves metadata to disk
func (s *Store) Save() error {
	if s.filePath == "" {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.saveUnsafe()
}

// AddSuggestions adds test suggestions to a file, merging with existing suggestions
func (s *Store) AddSuggestions(filePath string, suggestions []TestSuggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.metadata[filePath]
	if existing == nil {
		// No existing metadata, create new
		s.metadata[filePath] = &FileMetadata{Suggestions: suggestions}
	} else {
		// Merge with existing suggestions
		// Use a map to deduplicate based on suggestedName
		suggestionMap := make(map[string]TestSuggestion)

		// Add existing suggestions first
		for _, sugg := range existing.Suggestions {
			suggestionMap[sugg.SuggestedName] = sugg
		}

		// Add/update with new suggestions
		for _, sugg := range suggestions {
			suggestionMap[sugg.SuggestedName] = sugg
		}

		// Convert map back to slice
		mergedSuggestions := make([]TestSuggestion, 0, len(suggestionMap))
		for _, sugg := range suggestionMap {
			mergedSuggestions = append(mergedSuggestions, sugg)
		}

		s.metadata[filePath].Suggestions = mergedSuggestions
	}

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}

// GetSuggestions retrieves test suggestions for a file
func (s *Store) GetSuggestions(filePath string) []TestSuggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	meta := s.metadata[filePath]
	if meta == nil {
		return nil
	}

	return meta.Suggestions
}

// ==================== COMMENT METHODS ====================

// AddComment adds a new comment to a file
func (s *Store) AddComment(filePath string, comment files.Comment) (files.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate ID if not provided
	if comment.ID == "" {
		comment.ID = uuid.New().String()
	}

	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	existing := s.metadata[filePath]
	if existing == nil {
		// No existing metadata, create new
		s.metadata[filePath] = &FileMetadata{
			Comments: []files.Comment{comment},
		}
	} else {
		// Append to existing comments
		s.metadata[filePath].Comments = append(existing.Comments, comment)
	}

	if s.filePath != "" {
		if err := s.saveUnsafe(); err != nil {
			return comment, err
		}
	}

	return comment, nil
}

// UpdateComment updates an existing comment's content
func (s *Store) UpdateComment(filePath string, commentID string, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.metadata[filePath]
	if existing == nil {
		return nil // No metadata for this file
	}

	for i, comment := range existing.Comments {
		if comment.ID == commentID {
			s.metadata[filePath].Comments[i].Content = content
			s.metadata[filePath].Comments[i].UpdatedAt = time.Now()
			break
		}
	}

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}

// DeleteComment removes a comment from a file
func (s *Store) DeleteComment(filePath string, commentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.metadata[filePath]
	if existing == nil {
		return nil // No metadata for this file
	}

	// Filter out the comment to delete
	filtered := make([]files.Comment, 0, len(existing.Comments))
	for _, comment := range existing.Comments {
		if comment.ID != commentID {
			filtered = append(filtered, comment)
		}
	}

	s.metadata[filePath].Comments = filtered

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}

// GetComments retrieves all comments for a file
func (s *Store) GetComments(filePath string) []files.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	meta := s.metadata[filePath]
	if meta == nil {
		return nil
	}

	return meta.Comments
}

// ToggleCommentResolved toggles the resolved status of a comment
func (s *Store) ToggleCommentResolved(filePath string, commentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.metadata[filePath]
	if existing == nil {
		return nil // No metadata for this file
	}

	for i, comment := range existing.Comments {
		if comment.ID == commentID {
			s.metadata[filePath].Comments[i].Resolved = !comment.Resolved
			s.metadata[filePath].Comments[i].UpdatedAt = time.Now()
			break
		}
	}

	if s.filePath != "" {
		return s.saveUnsafe()
	}

	return nil
}
