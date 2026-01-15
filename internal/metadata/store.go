package metadata

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Store manages test metadata storage
type Store struct {
	mu       sync.RWMutex
	metadata map[string]*FileMetadata // key: file path
	filePath string                    // path to JSON persistence file
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

// SetTestMetadata stores test metadata for a file
func (s *Store) SetTestMetadata(filePath string, tests []TestReference) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metadata[filePath] = &FileMetadata{Tests: tests}

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
