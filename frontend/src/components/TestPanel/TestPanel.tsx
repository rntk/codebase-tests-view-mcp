import React from 'react';
import { TestItem } from './TestItem';
import type { TestDetail } from '../../types';

interface TestPanelProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
}

export const TestPanel: React.FC<TestPanelProps> = ({ tests, loading, error }) => {
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
          {tests.map((test, index) => (
            <TestItem key={`${test.testFile}-${test.testName}-${index}`} test={test} />
          ))}
        </div>
      )}
    </div>
  );
};
