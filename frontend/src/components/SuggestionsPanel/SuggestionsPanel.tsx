import React from 'react';
import { SuggestionItem } from './SuggestionItem';
import type { TestSuggestion } from '../../types';

interface SuggestionsPanelProps {
  suggestions: TestSuggestion[];
  loading: boolean;
  error: string | null;
}

export const SuggestionsPanel: React.FC<SuggestionsPanelProps> = ({
  suggestions,
  loading,
  error,
}) => {
  // Sort suggestions by priority (high > medium > low)
  const sortedSuggestions = [...suggestions].sort((a, b) => {
    const priorityOrder = { high: 0, medium: 1, low: 2 };
    return (priorityOrder[a.priority] ?? 2) - (priorityOrder[b.priority] ?? 2);
  });

  const highCount = suggestions.filter(s => s.priority === 'high').length;
  const mediumCount = suggestions.filter(s => s.priority === 'medium').length;
  const lowCount = suggestions.filter(s => s.priority === 'low').length;

  return (
    <div className="suggestions-panel-container">
      <h2
        style={{
          marginTop: 0,
          marginBottom: 'var(--space-sm)',
          fontSize: '16px',
          fontWeight: '600',
          color: 'var(--text-primary)',
        }}
      >
        Test Suggestions
      </h2>

      {!loading && !error && suggestions.length > 0 && (
        <div
          style={{
            display: 'flex',
            gap: 'var(--space-sm)',
            marginBottom: 'var(--space-md)',
            fontSize: '11px',
          }}
        >
          {highCount > 0 && (
            <span style={{ color: '#dc2626' }}>{highCount} high</span>
          )}
          {mediumCount > 0 && (
            <span style={{ color: '#d97706' }}>{mediumCount} medium</span>
          )}
          {lowCount > 0 && (
            <span style={{ color: '#16a34a' }}>{lowCount} low</span>
          )}
        </div>
      )}

      {loading && (
        <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
          Loading suggestions...
        </div>
      )}

      {error && (
        <div
          style={{
            padding: 'var(--space-md)',
            color: 'var(--error)',
            backgroundColor: '#fef2f2',
            borderRadius: 'var(--radius-md)',
            fontSize: '13px',
          }}
        >
          Error: {error}
        </div>
      )}

      {!loading && !error && suggestions.length === 0 && (
        <div
          style={{
            padding: 'var(--space-lg) var(--space-md)',
            color: 'var(--text-tertiary)',
            textAlign: 'center',
            backgroundColor: 'var(--bg-primary)',
            borderRadius: 'var(--radius-md)',
            border: '1px dashed var(--border-color)',
          }}
        >
          No test suggestions for this file
        </div>
      )}

      {!loading && !error && sortedSuggestions.length > 0 && (
        <div className="suggestion-items">
          {sortedSuggestions.map((suggestion, index) => (
            <SuggestionItem
              key={`${suggestion.suggestedName}-${index}`}
              suggestion={suggestion}
            />
          ))}
        </div>
      )}
    </div>
  );
};
