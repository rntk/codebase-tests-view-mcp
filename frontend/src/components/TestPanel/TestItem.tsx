import React, { useEffect, useMemo, useRef, useState } from 'react';
import type { TestDetail } from '../../types';

const CONTEXT_LINE_COUNT = 10;

interface TestItemProps {
  test: TestDetail;
  isHighlighted?: boolean;
}

export const TestItem: React.FC<TestItemProps> = ({ test, isHighlighted }) => {
  const itemRef = useRef<HTMLDivElement>(null);
  const [showInputContext, setShowInputContext] = useState(false);
  const [showOutputContext, setShowOutputContext] = useState(false);
  const contentLines = useMemo(() => test.content.split(/\r?\n/), [test.content]);

  const renderContextLines = (startLine: number, endLine: number) => {
    const safeStart = Math.max(1, startLine - CONTEXT_LINE_COUNT);
    const safeEnd = Math.min(contentLines.length, endLine + CONTEXT_LINE_COUNT);
    const lineNumberWidth = String(safeEnd).length;
    const formatted = [];

    for (let lineNumber = safeStart; lineNumber <= safeEnd; lineNumber += 1) {
      const line = contentLines[lineNumber - 1] ?? '';
      formatted.push(`${String(lineNumber).padStart(lineNumberWidth, ' ')} | ${line}`);
    }

    return formatted.join('\n');
  };

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
        {test.inputLines && (
          <div><strong>Input Lines:</strong> {test.inputLines.start}-{test.inputLines.end}</div>
        )}
        {test.outputLines && (
          <div><strong>Output Lines:</strong> {test.outputLines.start}-{test.outputLines.end}</div>
        )}
      </div>

      {test.comment && (
        <div style={{ marginBottom: 'var(--space-md)' }}>
          <strong style={{ fontSize: '11px', color: 'var(--text-secondary)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Comment</strong>
          <div style={{
            marginTop: 'var(--space-xs)',
            fontSize: '12px',
            color: 'var(--text-secondary)',
            lineHeight: 1.4
          }}>
            {test.comment}
          </div>
        </div>
      )}

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
          {test.inputLines && (
            <div style={{ marginTop: 'var(--space-xs)' }}>
              <button
                type="button"
                onClick={() => setShowInputContext(prev => !prev)}
                style={{
                  border: 'none',
                  background: 'none',
                  color: 'var(--text-secondary)',
                  fontSize: '11px',
                  cursor: 'pointer',
                  padding: 0,
                  textDecoration: 'underline'
                }}
              >
                {showInputContext ? 'Hide' : 'Show'} input context (+/- {CONTEXT_LINE_COUNT} lines)
              </button>
            </div>
          )}
          {showInputContext && test.inputLines && (
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
              {renderContextLines(test.inputLines.start, test.inputLines.end)}
            </pre>
          )}
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
          {test.outputLines && (
            <div style={{ marginTop: 'var(--space-xs)' }}>
              <button
                type="button"
                onClick={() => setShowOutputContext(prev => !prev)}
                style={{
                  border: 'none',
                  background: 'none',
                  color: 'var(--text-secondary)',
                  fontSize: '11px',
                  cursor: 'pointer',
                  padding: 0,
                  textDecoration: 'underline'
                }}
              >
                {showOutputContext ? 'Hide' : 'Show'} output context (+/- {CONTEXT_LINE_COUNT} lines)
              </button>
            </div>
          )}
          {showOutputContext && test.outputLines && (
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
              {renderContextLines(test.outputLines.start, test.outputLines.end)}
            </pre>
          )}
        </div>
      )}
    </div>
  );
};
