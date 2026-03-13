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
      <h2 className="section-title mt-0 mb-sm">
        Test Suggestions
      </h2>

      {!loading && !error && suggestions.length > 0 && (
        <div className="priority-stats">
          {highCount > 0 && (
            <span className="priority-stat--high">{highCount} high</span>
          )}
          {mediumCount > 0 && (
            <span className="priority-stat--medium">{mediumCount} medium</span>
          )}
          {lowCount > 0 && (
            <span className="priority-stat--low">{lowCount} low</span>
          )}
        </div>
      )}

      {loading && (
        <div className="loading-state">
          Loading suggestions...
        </div>
      )}

      {error && (
        <div className="error-state">
          Error: {error}
        </div>
      )}

      {!loading && !error && suggestions.length === 0 && (
        <div className="empty-state">
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
