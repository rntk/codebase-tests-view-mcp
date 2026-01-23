import React, { useMemo } from 'react';
import './CodeViewer.css';
import type { TestReference, CoverageDepth } from '../../types';

interface CodeViewerProps {
  content: string;
  filename: string;
  testReferences?: TestReference[];
  coverageDepth?: CoverageDepth;
  onLineClick?: (testId: string) => void;
}

// Get heatmap color based on coverage intensity
const getHeatmapColor = (intensity: number): string => {
  // Green gradient from light to dark
  if (intensity <= 0) return 'transparent';
  if (intensity <= 0.25) return '#dcfce7'; // green-100
  if (intensity <= 0.5) return '#86efac';  // green-300
  if (intensity <= 0.75) return '#22c55e'; // green-500
  return '#16a34a'; // green-600
};

interface CoverageBarProps {
  depth: number;
  maxDepth: number;
  tests: string[];
}

const CoverageBar: React.FC<CoverageBarProps> = ({ depth, maxDepth, tests }) => {
  const intensity = maxDepth > 0 ? depth / maxDepth : 0;
  const color = getHeatmapColor(intensity);
  const widthPercent = maxDepth > 0 ? (depth / maxDepth) * 100 : 0;

  return (
    <div
      className="coverage-bar-container"
      title={depth > 0 ? `Covered by ${depth} test${depth > 1 ? 's' : ''}: ${tests.join(', ')}` : 'Not covered'}
    >
      <div
        className="coverage-bar"
        style={{
          backgroundColor: color,
          width: `${widthPercent}%`,
        }}
      />
    </div>
  );
};

export const CodeViewer: React.FC<CodeViewerProps> = ({
  content,
  filename,
  testReferences = [],
  coverageDepth = {},
  onLineClick
}) => {
  const lines = content.split(/\r?\n/);

  // Build a map of line numbers to test references
  const lineToTests = useMemo(() => {
    const map = new Map<number, TestReference[]>();
    testReferences.forEach(test => {
      for (let i = test.coveredLines.start; i <= test.coveredLines.end; i++) {
        if (!map.has(i)) {
          map.set(i, []);
        }
        map.get(i)!.push(test);
      }
    });
    return map;
  }, [testReferences]);

  // Calculate max depth for normalization
  const maxDepth = useMemo(() => {
    let max = 0;
    Object.values(coverageDepth).forEach(tests => {
      if (tests.length > max) max = tests.length;
    });
    return max;
  }, [coverageDepth]);

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

  const hasCoverageData = maxDepth > 0;

  return (
    <div className="code-viewer-container">
      <div className="code-viewer-header">
        {filename}
        {hasCoverageData && (
          <span className="coverage-indicator">
            Coverage depth enabled
          </span>
        )}
      </div>
      <div className="code-viewer-content">
        <div className="line-numbers">
          {lines.map((_, index) => (
            <span key={index + 1} className="line-number">
              {index + 1}
            </span>
          ))}
        </div>
        {hasCoverageData && (
          <div className="coverage-gutter">
            {lines.map((_, index) => {
              const lineNum = index + 1;
              const testsForLine = coverageDepth[lineNum] || [];
              return (
                <CoverageBar
                  key={lineNum}
                  depth={testsForLine.length}
                  maxDepth={maxDepth}
                  tests={testsForLine}
                />
              );
            })}
          </div>
        )}
        <div className="code-lines">
          {lines.map((line, index) => {
            const lineNum = index + 1;
            const tests = lineToTests.get(lineNum);
            const isHighlighted = tests && tests.length > 0;
            const coverageTests = coverageDepth[lineNum];
            const hasCoverage = coverageTests && coverageTests.length > 0;

            return (
              <div
                key={lineNum}
                onClick={() => isHighlighted && handleLineClick(lineNum)}
                className={`code-line ${isHighlighted ? 'highlighted' : ''} ${hasCoverage && !isHighlighted ? 'covered' : ''}`}
                title={isHighlighted ? `Covered by ${tests!.length} test(s). Click to view.` : hasCoverage ? `Covered by ${coverageTests.length} test(s)` : undefined}
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
