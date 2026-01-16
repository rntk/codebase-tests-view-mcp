import React from 'react';
import { TestItem } from './TestItem';
import type { TestDetail } from '../../types';

interface TestPanelProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
  highlightedTestIds?: Set<string>;
}

export const TestPanel: React.FC<TestPanelProps> = ({ tests, loading, error, highlightedTestIds }) => {
  return (
    <div className="test-panel-container">
      <h2 style={{
        marginTop: 0,
        marginBottom: 'var(--space-md)',
        fontSize: '16px',
        fontWeight: '600',
        color: 'var(--text-primary)'
      }}>
        Related Tests
      </h2>

      {loading && (
        <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
          Loading tests...
        </div>
      )}

      {error && (
        <div style={{
          padding: 'var(--space-md)',
          color: 'var(--error)',
          backgroundColor: '#fef2f2',
          borderRadius: 'var(--radius-md)',
          fontSize: '13px'
        }}>
          Error: {error}
        </div>
      )}

      {!loading && !error && tests.length === 0 && (
        <div style={{
          padding: 'var(--space-lg) var(--space-md)',
          color: 'var(--text-tertiary)',
          textAlign: 'center',
          backgroundColor: 'var(--bg-primary)',
          borderRadius: 'var(--radius-md)',
          border: '1px dashed var(--border-color)'
        }}>
          No tests found for this file
        </div>
      )}

      {!loading && !error && tests.length > 0 && (
        <div className="test-items">
          {tests.map((test, index) => {
            const testId = `${test.testFile}:${test.testName}`;
            const isHighlighted = highlightedTestIds?.has(testId) ?? false;
            return (
              <TestItem
                key={`${test.testFile}-${test.testName}-${index}`}
                test={test}
                isHighlighted={isHighlighted}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};
