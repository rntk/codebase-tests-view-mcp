import React, { useState } from 'react';
import type { TestSuggestion } from '../../types';

interface SuggestionItemProps {
  suggestion: TestSuggestion;
}

export const SuggestionItem: React.FC<SuggestionItemProps> = ({ suggestion }) => {
  const [expanded, setExpanded] = useState(false);

  const getPriorityClasses = () => {
    switch (suggestion.priority) {
      case 'high':
        return 'suggestion-item--high';
      case 'medium':
        return 'suggestion-item--medium';
      case 'low':
        return 'suggestion-item--low';
      default:
        return 'suggestion-item--medium';
    }
  };

  const getPriorityBadgeClass = () => {
    switch (suggestion.priority) {
      case 'high':
        return 'priority-badge priority-high';
      case 'medium':
        return 'priority-badge priority-medium';
      case 'low':
        return 'priority-badge priority-low';
      default:
        return 'priority-badge priority-medium';
    }
  };

  return (
    <div
      className={`suggestion-item ${getPriorityClasses()}`}
    >
      <div className="suggestion-header">
        <span className={getPriorityBadgeClass()}>
          {suggestion.priority}
        </span>
        <h4 className="suggestion-title">
          {suggestion.suggestedName}
        </h4>
      </div>

      <div className="suggestion-meta">
        <strong>Target Lines:</strong> {suggestion.targetLines.start}-{suggestion.targetLines.end}
        {suggestion.functionName && (
          <span className="ml-md">
            <strong>Function:</strong> {suggestion.functionName}
          </span>
        )}
      </div>

      <div className="suggestion-reason">
        {suggestion.reason}
      </div>

      <button
        type="button"
        onClick={() => setExpanded(!expanded)}
        className="suggestion-expand-btn"
      >
        {expanded ? 'Hide test skeleton' : 'Show test skeleton'}
      </button>

      {expanded && (
        <pre className="suggestion-skeleton">
          {suggestion.testSkeleton}
        </pre>
      )}
    </div>
  );
};
