import React, { useMemo } from 'react';
import './CodeViewer.css';
import type { TestReference, CoverageDepth, Comment, SourceReference } from '../../types';

interface CodeViewerProps {
  content: string;
  filename: string;
  testReferences?: TestReference[];
  coverageDepth?: CoverageDepth;
  comments?: Comment[];
  sourceReferences?: SourceReference[];
  onLineClick?: (testId: string) => void;
  onSourceRefClick?: (sourceFile: string, line: number) => void;
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
  sourceReferences = [],
  onLineClick,
  onSourceRefClick,
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

  // Build maps for source reference line types (test body / input / output)
  type SourceRefLineType = 'test-body' | 'test-input' | 'test-output';
  const lineToSourceRef = useMemo(() => {
    const map = new Map<number, { type: SourceRefLineType; ref: SourceReference }[]>();
    const addRange = (start: number, end: number, type: SourceRefLineType, ref: SourceReference) => {
      if (!start || !end) return;
      for (let i = start; i <= end; i++) {
        if (!map.has(i)) map.set(i, []);
        map.get(i)!.push({ type, ref });
      }
    };
    sourceReferences.forEach(ref => {
      // inputLines and outputLines take priority over lineRange for coloring
      addRange(ref.lineRange.start, ref.lineRange.end, 'test-body', ref);
      if (ref.inputLines?.start && ref.inputLines?.end) {
        addRange(ref.inputLines.start, ref.inputLines.end, 'test-input', ref);
      }
      if (ref.outputLines?.start && ref.outputLines?.end) {
        addRange(ref.outputLines.start, ref.outputLines.end, 'test-output', ref);
      }
    });
    return map;
  }, [sourceReferences]);

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

    // Handle source reference navigation (test file → source)
    const srcRefs = lineToSourceRef.get(lineNum);
    if (srcRefs && srcRefs.length > 0 && onSourceRefClick) {
      const first = srcRefs[0];
      onSourceRefClick(first.ref.sourceFile, first.ref.coveredLines.start);
    }
  };

  const handleLineDoubleClick = (lineNum: number) => {
    if (onLineDoubleClick) {
      onLineDoubleClick(lineNum);
    }
  };

  const hasCoverageData = maxDepth > 0;
  const hasComments = comments.length > 0;
  const hasSourceRefs = sourceReferences.length > 0;

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
        {hasSourceRefs && (
          <span className="source-ref-indicator">
            Source links enabled
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
          <div className="comment-gutter">
            {lines.map((_, index) => {
              const lineNum = index + 1;
              const lineComments = lineToComments.get(lineNum);
              const hasComment = lineComments && lineComments.length > 0;

              return (
                <div
                  key={lineNum}
                  className="comment-gutter-item"
                  title={hasComment ? `Comment: ${lineComments![0].content.substring(0, 50)}...` : undefined}
                >
                  {hasComment && (
                    <span className="comment-indicator">
                      {lineComments!.length > 1 && (
                        <span className="comment-count">
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
            const srcRefs = lineToSourceRef.get(lineNum);

            // Determine source ref class: output > input > body (most specific wins)
            let srcRefClass = '';
            let srcRefTitle = '';
            if (srcRefs && srcRefs.length > 0) {
              const hasOutput = srcRefs.some(r => r.type === 'test-output');
              const hasInput = srcRefs.some(r => r.type === 'test-input');
              if (hasOutput) {
                srcRefClass = 'test-output';
                srcRefTitle = `Expected output — click to go to source: ${srcRefs[0].ref.sourceFile}`;
              } else if (hasInput) {
                srcRefClass = 'test-input';
                srcRefTitle = `Test input — click to go to source: ${srcRefs[0].ref.sourceFile}`;
              } else {
                srcRefClass = 'test-body';
                srcRefTitle = `Tests ${srcRefs[0].ref.functionName} in ${srcRefs[0].ref.sourceFile} — click to navigate`;
              }
            }

            return (
              <div
                key={lineNum}
                ref={isSelected ? selectedLineRef : null}
                onClick={() => handleLineClick(lineNum)}
                onDoubleClick={() => handleLineDoubleClick(lineNum)}
                className={`code-line ${isHighlighted ? 'highlighted' : ''} ${hasCoverage && !isHighlighted ? 'covered' : ''} ${isSelected ? 'selected' : ''} ${hasComment ? 'has-comment' : ''} ${!isHighlighted && !isSelected && srcRefClass ? srcRefClass : ''}`}
                title={isHighlighted ? `Covered by ${tests!.length} test(s). Click to filter/view.` : srcRefClass ? srcRefTitle : hasCoverage ? `Covered by ${coverageTests.length} test(s)` : hasComment ? `Double-click to add comment. Has ${lineComments!.length} comment(s).` : 'Double-click to add comment'}
                style={{
                  backgroundColor: hasComment && !srcRefClass ? 'rgba(245, 158, 11, 0.1)' : undefined,
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
