import React from 'react';
import type { TestDetail } from '../../types';

interface TestsListProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
  onTestClick: (test: TestDetail) => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

export const TestsList: React.FC<TestsListProps> = ({
  tests,
  loading,
  error,
  onTestClick,
  onSourceFileClick,
}) => {
  if (loading) {
    return (
      <div className="loading-state">
        Loading tests...
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-state">
        Error: {error}
      </div>
    );
  }

  if (tests.length === 0) {
    return (
      <div className="empty-state">
        No tests found in this directory
      </div>
    );
  }

  return (
    <div className="tests-list">
      {tests.map((test, index) => (
        <div
          key={`${test.testFile}-${test.testName}-${index}`}
          className="test-item-card"
          onClick={() => onTestClick(test)}
        >
          <div className="test-item-title">
            {test.testName}
          </div>
          <div className="test-item-meta">
            <span className="inline-code">
              {test.functionName}
            </span>
          </div>
          <div className="test-item-details">
            <span>in</span>
            <span
              className="inline-code-link"
              onClick={(e) => {
                e.stopPropagation();
                onSourceFileClick(test.testFile, test.lineRange?.start ?? 1);
              }}
            >
              {test.testFile}
            </span>
            <span>
              (lines {test.lineRange?.start ?? '?'}-{test.lineRange?.end ?? '?'})
            </span>
          </div>
          {test.comment && (
            <div className="test-item-comment">
              {test.comment}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};
