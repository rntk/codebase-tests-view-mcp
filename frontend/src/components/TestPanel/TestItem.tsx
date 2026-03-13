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

  const renderContextBlock = (startLine: number, endLine: number) => {
    const safeStart = Math.max(1, startLine - CONTEXT_LINE_COUNT);
    const safeEnd = Math.min(contentLines.length, endLine + CONTEXT_LINE_COUNT);
    const lineNumberWidth = String(safeEnd).length;
    const lines = [];

    for (let lineNumber = safeStart; lineNumber <= safeEnd; lineNumber += 1) {
      const line = contentLines[lineNumber - 1] ?? '';
      const isTarget = lineNumber >= startLine && lineNumber <= endLine;
      lines.push(
        <div key={lineNumber} className={`code-line-context ${isTarget ? 'code-line-context--target' : ''}`}>
          <span className="code-line-number" style={{ width: `${lineNumberWidth}ch` }}>
            {lineNumber}
          </span>
          <span className="whitespace-pre-wrap break-all">{line}</span>
        </div>
      );
    }

    return lines;
  };

  useEffect(() => {
    if (isHighlighted && itemRef.current) {
      itemRef.current.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }, [isHighlighted]);

  return (
    <div
      ref={itemRef}
      className={`test-detail-item ${isHighlighted ? 'test-detail-item--highlighted' : ''}`}
    >
      <h4 className="test-detail-title">
        {test.testName}
      </h4>
      <div className="test-detail-meta">
        <div className="break-all"><strong>File:</strong> {test.testFile}</div>
        <div><strong>Test Lines:</strong> {test.lineRange.start}-{test.lineRange.end}</div>
        <div><strong>Target Lines:</strong> {test.coveredLines.start}-{test.coveredLines.end}</div>
        {test.inputData && test.inputLines && (
          <div><strong>Input Lines:</strong> {test.inputLines.start}-{test.inputLines.end}</div>
        )}
        {test.expectedOutput && test.outputLines && (
          <div><strong>Output Lines:</strong> {test.outputLines.start}-{test.outputLines.end}</div>
        )}
      </div>

      {test.comment && (
        <div className="test-detail-section">
          <strong className="test-detail-section-title">Comment</strong>
          <div className="test-detail-section-content">
            {test.comment}
          </div>
        </div>
      )}

      {test.inputData && (
        <div className="test-detail-section">
          <div className="test-detail-actions">
            <strong className="test-detail-section-title">Input Data</strong>
            {test.inputLines && (
              <button
                type="button"
                onClick={() => setShowInputContext(prev => !prev)}
                className="test-detail-context-btn"
                title={showInputContext ? 'Show only the input data' : 'Show input data with surrounding lines from the test file'}
              >
                {showInputContext ? 'Show Raw Data' : 'Show Context'}
              </button>
            )}
          </div>
          <div className="code-block">
            {showInputContext && test.inputLines
              ? renderContextBlock(test.inputLines.start, test.inputLines.end)
              : test.inputData
            }
          </div>
        </div>
      )}

      {test.expectedOutput && (
        <div className="test-detail-section">
          <div className="test-detail-actions">
            <strong className="test-detail-section-title">Expected Result</strong>
            {test.outputLines && (
              <button
                type="button"
                onClick={() => setShowOutputContext(prev => !prev)}
                className="test-detail-context-btn"
                title={showOutputContext ? 'Show only the expected result' : 'Show expected result with surrounding lines from the test file'}
              >
                {showOutputContext ? 'Show Raw Data' : 'Show Context'}
              </button>
            )}
          </div>
          <div className="code-block">
            {showOutputContext && test.outputLines
              ? renderContextBlock(test.outputLines.start, test.outputLines.end)
              : test.expectedOutput
            }
          </div>
        </div>
      )}
    </div>
  );
};
