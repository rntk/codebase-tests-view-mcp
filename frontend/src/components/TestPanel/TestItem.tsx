import React, { useEffect, useRef } from 'react';
import type { TestDetail } from '../../types';

interface TestItemProps {
  test: TestDetail;
  isHighlighted?: boolean;
}

export const TestItem: React.FC<TestItemProps> = ({ test, isHighlighted }) => {
  const itemRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isHighlighted && itemRef.current) {
      itemRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }, [isHighlighted]);

  return (
    <div
      ref={itemRef}
      style={{
        marginBottom: '24px',
        padding: '12px',
        border: isHighlighted ? '2px solid #f59e0b' : '1px solid #ddd',
        borderRadius: '4px',
        backgroundColor: isHighlighted ? '#fef3c7' : 'white',
        transition: 'all 0.3s',
      }}
    >
      <h4 style={{ margin: '0 0 8px 0', fontSize: '14px', fontWeight: 'bold' }}>
        {test.testName}
      </h4>
      <div style={{ fontSize: '12px', color: '#666', marginBottom: '8px' }}>
        <div>File: {test.testFile}</div>
        <div>Lines: {test.lineRange.start}-{test.lineRange.end}</div>
        <div>Covers: {test.coveredLines.start}-{test.coveredLines.end}</div>
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
