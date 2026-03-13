import React from 'react';
import type { TestDetail } from '../../types';

interface TestsListProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
  onTestClick: (test: TestDetail) => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

export const TestsList: React.FC<TestsListProps> = ({
  tests,
  loading,
  error,
  onTestClick,
  onSourceFileClick,
}) => {
  if (loading) {
    return (
      <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
        Loading tests...
      </div>
    );
  }

  if (error) {
    return (
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
    );
  }

  if (tests.length === 0) {
    return (
      <div
        style={{
          padding: 'var(--space-md)',
          color: 'var(--text-tertiary)',
          textAlign: 'center',
        }}
      >
        No tests found in this directory
      </div>
    );
  }

  return (
    <div className="tests-list">
      {tests.map((test, index) => (
        <div
          key={`${test.testFile}-${test.testName}-${index}`}
          style={{
            padding: 'var(--space-sm)',
            marginBottom: 'var(--space-sm)',
            backgroundColor: 'var(--bg-secondary)',
            borderRadius: 'var(--radius-md)',
            border: '1px solid var(--border-color)',
            cursor: 'pointer',
            transition: 'all 0.15s ease',
          }}
          onClick={() => onTestClick(test)}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = 'var(--bg-tertiary)';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = 'var(--bg-secondary)';
          }}
        >
          <div
            style={{
              fontSize: '13px',
              fontWeight: '600',
              color: 'var(--text-primary)',
              marginBottom: '4px',
            }}
          >
            {test.testName}
          </div>
          <div
            style={{
              fontSize: '12px',
              color: 'var(--text-secondary)',
              marginBottom: '4px',
            }}
          >
            <span
              style={{
                fontFamily: 'monospace',
                backgroundColor: 'var(--bg-tertiary)',
                padding: '2px 6px',
                borderRadius: '4px',
              }}
            >
              {test.functionName}
            </span>
          </div>
          <div
            style={{
              fontSize: '11px',
              color: 'var(--text-tertiary)',
              display: 'flex',
              alignItems: 'center',
              gap: '4px',
            }}
          >
            <span>in</span>
            <span
              style={{
                fontFamily: 'monospace',
                color: 'var(--accent-primary)',
                cursor: 'pointer',
                textDecoration: 'underline',
              }}
              onClick={(e) => {
                e.stopPropagation();
                onSourceFileClick(test.testFile, test.lineRange?.start ?? 1);
              }}
            >
              {test.testFile}
            </span>
            <span>
              (lines {test.lineRange?.start ?? '?'}-{test.lineRange?.end ?? '?'})
            </span>
          </div>
          {test.comment && (
            <div
              style={{
                fontSize: '11px',
                color: 'var(--text-secondary)',
                marginTop: 'var(--space-xs)',
                fontStyle: 'italic',
                borderTop: '1px solid var(--border-color)',
                paddingTop: 'var(--space-xs)',
              }}
            >
              {test.comment}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};
