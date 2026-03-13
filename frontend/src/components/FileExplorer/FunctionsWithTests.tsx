import React, { useState, useMemo } from 'react';
import type { TestDetail } from '../../types';
import { groupTestsByFunction } from '../../utils/testUtils';

interface FunctionWithTests {
  functionName: string;
  tests: TestDetail[];
}

interface FunctionsWithTestsProps {
  tests: TestDetail[];
  loading: boolean;
  error: string | null;
  onTestClick: (test: TestDetail) => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

export const FunctionsWithTests: React.FC<FunctionsWithTestsProps> = ({
  tests,
  loading,
  error,
  onTestClick,
  onSourceFileClick,
}) => {
  const [expandedFunctions, setExpandedFunctions] = useState<Set<string>>(new Set());

  if (loading) {
    return (
      <div className="loading-state">
        Loading functions...
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

  // Group tests by function name using shared utility
  const functionsMap = useMemo(() => groupTestsByFunction(tests), [tests]);

  const functions: FunctionWithTests[] = useMemo(
    () =>
      Array.from(functionsMap.entries()).map(([functionName, functionTests]) => ({
        functionName,
        tests: functionTests,
      })),
    [functionsMap]
  );

  if (functions.length === 0) {
    return (
      <div className="empty-state">
        No functions with tests found
      </div>
    );
  }

  const toggleFunction = (functionName: string) => {
    setExpandedFunctions((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(functionName)) {
        newSet.delete(functionName);
      } else {
        newSet.add(functionName);
      }
      return newSet;
    });
  };

  return (
    <div className="functions-with-tests">
      {functions.map(({ functionName, tests: functionTests }) => {
        const isExpanded = expandedFunctions.has(functionName);

        return (
          <div
            key={functionName}
            className="function-accordion"
          >
            {/* Function header */}
            <div
              className="function-header"
              onClick={() => toggleFunction(functionName)}
            >
              <span
                className={`function-toggle ${isExpanded ? 'function-toggle--expanded' : ''}`}
              >
                ▶
              </span>
              <span className="function-name">
                {functionName}
              </span>
              <span className="function-test-count">
                {functionTests.length} test{functionTests.length > 1 ? 's' : ''}
              </span>
            </div>

            {/* Tests list */}
            {isExpanded && (
              <div className="function-tests">
                {functionTests.map((test, index) => (
                  <div
                    key={`${test.testFile}-${test.testName}-${index}`}
                    className="test-item-card"
                    onClick={() => onTestClick(test)}
                  >
                    <div className="test-item-title">
                      {test.testName}
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
            )}
          </div>
        );
      })}
    </div>
  );
};
