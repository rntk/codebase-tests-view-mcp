import React, { useMemo } from 'react';
import './CodeViewer.css';
import type { TestReference, CoverageDepth, Comment } from '../../types';

interface CodeViewerProps {
  content: string;
  filename: string;
  testReferences?: TestReference[];
  coverageDepth?: CoverageDepth;
  comments?: Comment[];
  onLineClick?: (testId: string) => void;
  selectedLine?: number | null;
  onLineSelect?: (line: number) => void;
  onLineDoubleClick?: (line: number) => void;
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
  comments = [],
  onLineClick,
  selectedLine,
  onLineSelect,
  onLineDoubleClick,
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

  // Build a map of line numbers to comments
  const lineToComments = useMemo(() => {
    const map = new Map<number, Comment[]>();
    comments.forEach(comment => {
      if (!comment.resolved) {
        if (!map.has(comment.line)) {
          map.set(comment.line, []);
        }
        map.get(comment.line)!.push(comment);
      }
    });
    return map;
  }, [comments]);

  // Calculate max depth for normalization
  const maxDepth = useMemo(() => {
    let max = 0;
    Object.values(coverageDepth).forEach(tests => {
      if (tests.length > max) max = tests.length;
    });
    return max;
  }, [coverageDepth]);

  const handleLineClick = (lineNum: number) => {
    // Notify about line selection for filtering
    if (onLineSelect) {
      onLineSelect(lineNum);
    }

    const tests = lineToTests.get(lineNum);
    if (tests && tests.length > 0 && onLineClick) {
      // Highlight all tests covering this line
      tests.forEach(test => {
        const testId = `${test.testFile}:${test.testName}`;
        onLineClick(testId);
      });
    }
  };

  const handleLineDoubleClick = (lineNum: number) => {
    if (onLineDoubleClick) {
      onLineDoubleClick(lineNum);
    }
  };

  const hasCoverageData = maxDepth > 0;
  const hasComments = comments.length > 0;

  // Scroll to selected line on mount or when selectedLine changes
  const selectedLineRef = React.useRef<HTMLDivElement>(null);
  React.useEffect(() => {
    if (selectedLine && selectedLineRef.current) {
      selectedLineRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
      });
    }
  }, [selectedLine, content]);

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
        
        {/* Coverage gutter */}
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

        {/* Comment gutter */}
        {hasComments && (
          <div className="comment-gutter" style={{
            display: 'flex',
            flexDirection: 'column',
            minWidth: '24px',
            padding: '4px 0',
          }}>
            {lines.map((_, index) => {
              const lineNum = index + 1;
              const lineComments = lineToComments.get(lineNum);
              const hasComment = lineComments && lineComments.length > 0;
              
              return (
                <div
                  key={lineNum}
                  style={{
                    height: '20px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                  title={hasComment ? `Comment: ${lineComments![0].content.substring(0, 50)}...` : undefined}
                >
                  {hasComment && (
                    <span style={{
                      width: '8px',
                      height: '8px',
                      borderRadius: '50%',
                      backgroundColor: 'var(--warning)',
                      cursor: 'pointer',
                    }}>
                      {lineComments!.length > 1 && (
                        <span style={{
                          position: 'absolute',
                          fontSize: '9px',
                          color: 'var(--warning)',
                          marginLeft: '10px',
                          marginTop: '-2px',
                        }}>
                          {lineComments!.length}
                        </span>
                      )}
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        )}

        <div className="code-lines">
          {lines.map((line, index) => {
            const lineNum = index + 1;
            const tests = lineToTests.get(lineNum);
            const isHighlighted = tests && tests.length > 0;
            const isSelected = selectedLine === lineNum;
            const coverageTests = coverageDepth[lineNum];
            const hasCoverage = coverageTests && coverageTests.length > 0;
            const lineComments = lineToComments.get(lineNum);
            const hasComment = lineComments && lineComments.length > 0;

            return (
              <div
                key={lineNum}
                ref={isSelected ? selectedLineRef : null}
                onClick={() => handleLineClick(lineNum)}
                onDoubleClick={() => handleLineDoubleClick(lineNum)}
                className={`code-line ${isHighlighted ? 'highlighted' : ''} ${hasCoverage && !isHighlighted ? 'covered' : ''} ${isSelected ? 'selected' : ''} ${hasComment ? 'has-comment' : ''}`}
                title={isHighlighted ? `Covered by ${tests!.length} test(s). Click to filter/view.` : hasCoverage ? `Covered by ${coverageTests.length} test(s)` : hasComment ? `Double-click to add comment. Has ${lineComments!.length} comment(s).` : 'Double-click to add comment'}
                style={{
                  backgroundColor: hasComment ? 'rgba(245, 158, 11, 0.1)' : undefined,
                }}
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
