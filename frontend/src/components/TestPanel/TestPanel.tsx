import React from 'react';
import { TestItem } from './TestItem';
import type { TestDetail } from '../../types';

interface TestPanelProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
  highlightedTestId?: string | null;
}

export const TestPanel: React.FC<TestPanelProps> = ({ tests, loading, error, highlightedTestId }) => {
  return (
    <div>
      <h2 style={{ marginTop: 0, marginBottom: '16px', fontSize: '18px' }}>
        Related Tests
      </h2>

      {loading && <div>Loading tests...</div>}
      {error && <div style={{ color: 'red' }}>Error: {error}</div>}

      {!loading && !error && tests.length === 0 && (
        <div style={{ color: '#666' }}>
          No tests found for this file
        </div>
      )}

      {!loading && !error && tests.length > 0 && (
        <div>
          {tests.map((test, index) => {
            const testId = `${test.testFile}:${test.testName}`;
            const isHighlighted = testId === highlightedTestId;
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
