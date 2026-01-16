import React from 'react';
import './CodeViewer.css';
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
    <div className="code-viewer-container">
      <div className="code-viewer-header">
        {filename}
      </div>
      <div className="code-viewer-content">
        <div className="line-numbers">
          {lines.map((_, index) => (
            <span key={index + 1} className="line-number">
              {index + 1}
            </span>
          ))}
        </div>
        <div className="code-lines">
          {lines.map((line, index) => {
            const lineNum = index + 1;
            const tests = lineToTests.get(lineNum);
            const isHighlighted = tests && tests.length > 0;

            return (
              <div
                key={lineNum}
                onClick={() => isHighlighted && handleLineClick(lineNum)}
                className={`code-line ${isHighlighted ? 'highlighted' : ''}`}
                title={isHighlighted ? `Covered by ${tests!.length} test(s). Click to view.` : undefined}
              >
                {line === '' ? '\u00a0' : line}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
