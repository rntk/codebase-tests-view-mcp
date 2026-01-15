import React from 'react';
import type { TestDetail } from '../../types';

interface TestItemProps {
  test: TestDetail;
}

export const TestItem: React.FC<TestItemProps> = ({ test }) => {
  return (
    <div style={{ marginBottom: '24px', padding: '12px', border: '1px solid #ddd', borderRadius: '4px' }}>
      <h4 style={{ margin: '0 0 8px 0', fontSize: '14px', fontWeight: 'bold' }}>
        {test.testName}
      </h4>
      <div style={{ fontSize: '12px', color: '#666', marginBottom: '8px' }}>
        <div>File: {test.testFile}</div>
        <div>Lines: {test.lineRange.start}-{test.lineRange.end}</div>
      </div>

      {test.inputData && (
        <div style={{ marginBottom: '8px' }}>
          <strong style={{ fontSize: '12px' }}>Input Data:</strong>
          <pre style={{
            backgroundColor: '#f5f5f5',
            padding: '8px',
            borderRadius: '4px',
            fontSize: '11px',
            marginTop: '4px',
          }}>
            {test.inputData}
          </pre>
        </div>
      )}

      {test.expectedOutput && (
        <div>
          <strong style={{ fontSize: '12px' }}>Expected Output:</strong>
          <pre style={{
            backgroundColor: '#f5f5f5',
            padding: '8px',
            borderRadius: '4px',
            fontSize: '11px',
            marginTop: '4px',
          }}>
            {test.expectedOutput}
          </pre>
        </div>
      )}
    </div>
  );
};
