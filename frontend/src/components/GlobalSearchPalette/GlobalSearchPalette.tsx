import { useState, useEffect, useCallback, useRef } from 'react';
import { search } from '../../api/client';
import type { SearchResult } from '../../types';

interface GlobalSearchPaletteProps {
  isOpen: boolean;
  onClose: () => void;
  onResultSelect: (result: SearchResult) => void;
}

export function GlobalSearchPalette({ isOpen, onClose, onResultSelect }: GlobalSearchPaletteProps) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const requestIdRef = useRef(0); // Track the latest request ID to ignore stale responses

  // Focus input when opened
  useEffect(() => {
    if (isOpen && inputRef.current) {
      inputRef.current.focus();
    }
  }, [isOpen]);

  // Reset state when closed
  useEffect(() => {
    if (!isOpen) {
      setQuery('');
      setResults([]);
      setSelectedIndex(0);
      setIsLoading(false);
    }
  }, [isOpen]);

  // Search with debounce
  useEffect(() => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    const currentQuery = query.trim();
    if (!currentQuery) {
      ++requestIdRef.current; // invalidate any in-flight request
      setResults([]);
      return;
    }

    // Increment request ID to track this specific request
    const requestId = ++requestIdRef.current;

    debounceRef.current = setTimeout(async () => {
      setIsLoading(true);
      try {
        const response = await search(currentQuery);
        // Only update if this is still the latest request
        // This prevents stale responses from overwriting newer results
        if (requestId === requestIdRef.current && isOpen) {
          setResults(response.results.slice(0, 50));
          setSelectedIndex(0);
        }
      } catch (error) {
        console.error('Search failed:', error);
        // Only set empty results if this is still the latest request
        if (requestId === requestIdRef.current) {
          setResults([]);
        }
      } finally {
        // Only clear loading if this is still the latest request
        if (requestId === requestIdRef.current) {
          setIsLoading(false);
        }
      }
    }, 150);

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [query, isOpen]);

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown' && results.length > 0) {
      e.preventDefault();
      setSelectedIndex(prev => Math.min(prev + 1, results.length - 1));
    } else if (e.key === 'ArrowUp' && results.length > 0) {
      e.preventDefault();
      setSelectedIndex(prev => Math.max(prev - 1, 0));
    } else if (e.key === 'Enter' && results[selectedIndex]) {
      e.preventDefault();
      onResultSelect(results[selectedIndex]);
      onClose();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
    }
  }, [results, selectedIndex, onResultSelect, onClose]);

  // Handle result click
  const handleResultClick = useCallback((result: SearchResult) => {
    onResultSelect(result);
    onClose();
  }, [onResultSelect, onClose]);

  // Get icon based on result type
  const getResultIcon = (type: SearchResult['type']) => {
    switch (type) {
      case 'file':
        return '📄';
      case 'function':
        return '⚡';
      case 'test':
        return '🧪';
      default:
        return '📄';
    }
  };

  if (!isOpen) {
    return null;
  }

  return (
    <div className="search-palette-overlay" onClick={onClose}>
      <div className="search-palette" onClick={e => e.stopPropagation()}>
        <div className="search-palette__header">
          <span className="search-palette__icon">🔍</span>
          <input
            ref={inputRef}
            type="text"
            className="search-palette__input"
            placeholder="Search files, functions, tests... (Cmd+P)"
            value={query}
            onChange={e => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
          />
          {isLoading && (
            <div className="search-palette__loading">Loading...</div>
          )}
        </div>

        <div className="search-palette__results">
          {results.length === 0 && query.trim() && !isLoading && (
            <div className="search-palette__empty">No results found</div>
          )}

          {results.map((result, index) => (
            <div
              key={`${result.type}-${result.path}-${result.line}-${result.title}-${index}`}
              className={`search-palette__result ${index === selectedIndex ? 'search-palette__result--selected' : ''}`}
              onClick={() => handleResultClick(result)}
              onMouseEnter={() => setSelectedIndex(index)}
            >
              <div className="search-palette__result-icon">
                {getResultIcon(result.type)}
              </div>
              <div className="search-palette__result-content">
                <div className="search-palette__result-title">
                  {result.title}
                  <span className="search-palette__result-type">
                    {result.type}
                  </span>
                </div>
                <div className="search-palette__result-subtitle">
                  {result.subtitle}
                  {result.line > 0 && `:${result.line}`}
                </div>
              </div>
            </div>
          ))}
        </div>

        {results.length > 0 && (
          <div className="search-palette__footer">
            <span className="search-palette__shortcut">↑↓</span>
            <span className="search-palette__shortcut-label">navigate</span>
            <span className="search-palette__shortcut">↵</span>
            <span className="search-palette__shortcut-label">select</span>
            <span className="search-palette__shortcut">esc</span>
            <span className="search-palette__shortcut-label">close</span>
          </div>
        )}
      </div>
    </div>
  );
}
