import React from 'react';
import type { TestReference } from '../../types';

interface CodeViewerProps {
  content: string;
  filename: string;
  testReferences?: TestReference[];
  onLineClick?: (testId: string) => void;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({
  content,
  filename,
  testReferences = [],
  onLineClick
}) => {
  const lines = content.split(/\r?\n/);

  // Build a map of line numbers to test references
  const lineToTests = new Map<number, TestReference[]>();
  testReferences.forEach(test => {
    for (let i = test.coveredLines.start; i <= test.coveredLines.end; i++) {
      if (!lineToTests.has(i)) {
        lineToTests.set(i, []);
      }
      lineToTests.get(i)!.push(test);
    }
  });

  const handleLineClick = (lineNum: number) => {
    const tests = lineToTests.get(lineNum);
    if (tests && tests.length > 0 && onLineClick) {
      // Highlight all tests covering this line
      tests.forEach(test => {
        const testId = `${test.testFile}:${test.testName}`;
        onLineClick(testId);
      });
    }
  };

  return (
    <div>
      <h3 style={{ marginTop: 0, marginBottom: '12px', fontSize: '16px', fontFamily: 'monospace' }}>
        {filename}
      </h3>
      <div
        style={{
          backgroundColor: '#f5f5f5',
          padding: '16px',
          borderRadius: '4px',
          overflow: 'auto',
          fontSize: '13px',
          lineHeight: '1.5',
          fontFamily: 'monospace',
          whiteSpace: 'pre',
        }}
      >
        {lines.map((line, index) => {
          const lineNum = index + 1;
          const tests = lineToTests.get(lineNum);
          const isHighlighted = tests && tests.length > 0;

          return (
            <div
              key={lineNum}
              onClick={() => isHighlighted && handleLineClick(lineNum)}
              style={{
                backgroundColor: isHighlighted ? '#fef3c7' : 'transparent',
                borderLeft: isHighlighted ? '3px solid #f59e0b' : 'none',
                paddingLeft: isHighlighted ? '8px' : '0',
                marginLeft: isHighlighted ? '-11px' : '0',
                cursor: isHighlighted ? 'pointer' : 'default',
                transition: 'background-color 0.2s',
              }}
              onMouseEnter={(e) => {
                if (isHighlighted) {
                  e.currentTarget.style.backgroundColor = '#fde68a';
                }
              }}
              onMouseLeave={(e) => {
                if (isHighlighted) {
                  e.currentTarget.style.backgroundColor = '#fef3c7';
                }
              }}
              title={isHighlighted ? `Covered by ${tests!.length} test(s). Click to view.` : undefined}
            >
              {line === '' ? '\u00a0' : line}
            </div>
          );
        })}
      </div>
    </div>
  );
};
