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
        marginBottom: 'var(--space-md)',
        padding: 'var(--space-md)',
        border: isHighlighted ? '2px solid var(--warning)' : '1px solid var(--border-color)',
        borderRadius: 'var(--radius-md)',
        backgroundColor: isHighlighted ? '#fffbeb' : 'var(--bg-primary)',
        transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
        boxShadow: isHighlighted ? 'var(--shadow-md)' : 'var(--shadow-sm)',
        transform: isHighlighted ? 'scale(1.02)' : 'scale(1)',
      }}
    >
      <h4 style={{
        margin: '0 0 var(--space-sm) 0',
        fontSize: '14px',
        fontWeight: '600',
        color: 'var(--text-primary)'
      }}>
        {test.testName}
      </h4>
      <div style={{
        fontSize: '12px',
        color: 'var(--text-tertiary)',
        marginBottom: 'var(--space-md)',
        display: 'flex',
        flexDirection: 'column',
        gap: '2px'
      }}>
        <div style={{ wordBreak: 'break-all' }}><strong>File:</strong> {test.testFile}</div>
        <div><strong>Test Lines:</strong> {test.lineRange.start}-{test.lineRange.end}</div>
        <div><strong>Target Lines:</strong> {test.coveredLines.start}-{test.coveredLines.end}</div>
      </div>

      {test.inputData && (
        <div style={{ marginBottom: 'var(--space-md)' }}>
          <strong style={{ fontSize: '11px', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Input Data</strong>
          <pre style={{
            backgroundColor: 'var(--bg-secondary)',
            padding: 'var(--space-sm)',
            borderRadius: 'var(--radius-sm)',
            fontSize: '11px',
            marginTop: 'var(--space-xs)',
            border: '1px solid var(--border-color)',
            overflow: 'auto',
            fontFamily: 'var(--font-mono)'
          }}>
            {test.inputData}
          </pre>
        </div>
      )}

      {test.expectedOutput && (
        <div>
          <strong style={{ fontSize: '11px', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Expected Result</strong>
          <pre style={{
            backgroundColor: 'var(--bg-secondary)',
            padding: 'var(--space-sm)',
            borderRadius: 'var(--radius-sm)',
            fontSize: '11px',
            marginTop: 'var(--space-xs)',
            border: '1px solid var(--border-color)',
            overflow: 'auto',
            fontFamily: 'var(--font-mono)'
          }}>
            {test.expectedOutput}
          </pre>
        </div>
      )}
    </div>
  );
};
