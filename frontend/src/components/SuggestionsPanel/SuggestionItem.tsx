import React, { useState } from 'react';
import type { TestSuggestion } from '../../types';

interface SuggestionItemProps {
  suggestion: TestSuggestion;
}

const priorityColors: Record<string, { bg: string; text: string; border: string }> = {
  high: { bg: '#fef2f2', text: '#dc2626', border: '#fca5a5' },
  medium: { bg: '#fffbeb', text: '#d97706', border: '#fcd34d' },
  low: { bg: '#f0fdf4', text: '#16a34a', border: '#86efac' },
};

export const SuggestionItem: React.FC<SuggestionItemProps> = ({ suggestion }) => {
  const [expanded, setExpanded] = useState(false);
  const colors = priorityColors[suggestion.priority] || priorityColors.medium;

  return (
    <div
      style={{
        marginBottom: 'var(--space-md)',
        padding: 'var(--space-md)',
        border: `1px solid ${colors.border}`,
        borderRadius: 'var(--radius-md)',
        backgroundColor: 'var(--bg-primary)',
        transition: 'all 0.2s ease',
      }}
    >
      <div style={{ display: 'flex', alignItems: 'flex-start', gap: 'var(--space-sm)' }}>
        <span
          style={{
            padding: '2px 8px',
            borderRadius: '12px',
            fontSize: '10px',
            fontWeight: '600',
            textTransform: 'uppercase',
            backgroundColor: colors.bg,
            color: colors.text,
            border: `1px solid ${colors.border}`,
            flexShrink: 0,
          }}
        >
          {suggestion.priority}
        </span>
        <h4
          style={{
            margin: 0,
            fontSize: '14px',
            fontWeight: '600',
            color: 'var(--text-primary)',
            flex: 1,
          }}
        >
          {suggestion.suggestedName}
        </h4>
      </div>

      <div
        style={{
          fontSize: '12px',
          color: 'var(--text-tertiary)',
          marginTop: 'var(--space-sm)',
          marginBottom: 'var(--space-sm)',
        }}
      >
        <strong>Target Lines:</strong> {suggestion.targetLines.start}-{suggestion.targetLines.end}
        {suggestion.functionName && (
          <span style={{ marginLeft: 'var(--space-md)' }}>
            <strong>Function:</strong> {suggestion.functionName}
          </span>
        )}
      </div>

      <div
        style={{
          fontSize: '13px',
          color: 'var(--text-secondary)',
          lineHeight: 1.5,
          marginBottom: 'var(--space-sm)',
        }}
      >
        {suggestion.reason}
      </div>

      <button
        type="button"
        onClick={() => setExpanded(!expanded)}
        style={{
          border: 'none',
          background: 'none',
          color: 'var(--accent-primary)',
          fontSize: '12px',
          cursor: 'pointer',
          padding: 0,
          textDecoration: 'underline',
        }}
      >
        {expanded ? 'Hide test skeleton' : 'Show test skeleton'}
      </button>

      {expanded && (
        <pre
          style={{
            backgroundColor: 'var(--bg-secondary)',
            padding: 'var(--space-sm)',
            borderRadius: 'var(--radius-sm)',
            fontSize: '11px',
            marginTop: 'var(--space-sm)',
            border: '1px solid var(--border-color)',
            overflow: 'auto',
            fontFamily: 'var(--font-mono)',
            whiteSpace: 'pre-wrap',
            maxHeight: '300px',
          }}
        >
          {suggestion.testSkeleton}
        </pre>
      )}
    </div>
  );
};
