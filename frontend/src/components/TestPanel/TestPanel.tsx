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
      <h2 className="section-title mt-0 mb-md">
        Related Tests
      </h2>

      {loading && (
        <div className="loading-state">
          Loading tests...
        </div>
      )}

      {error && (
        <div className="error-state">
          Error: {error}
        </div>
      )}

      {!loading && !error && tests.length === 0 && (
        <div className="empty-state">
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
