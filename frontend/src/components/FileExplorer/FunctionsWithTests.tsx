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
      <div style={{ padding: 'var(--space-md)', color: 'var(--text-tertiary)' }}>
        Loading functions...
      </div>
    );
  }

  if (error) {
    return (
      <div
        style={{
          padding: 'var(--space-md)',
          color: 'var(--error)',
          backgroundColor: '#fef2f2',
          borderRadius: 'var(--radius-md)',
          fontSize: '13px',
        }}
      >
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
      <div
        style={{
          padding: 'var(--space-md)',
          color: 'var(--text-tertiary)',
          textAlign: 'center',
        }}
      >
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
            style={{
              marginBottom: 'var(--space-sm)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-md)',
              overflow: 'hidden',
            }}
          >
            {/* Function header */}
            <div
              style={{
                padding: 'var(--space-sm)',
                backgroundColor: 'var(--bg-secondary)',
                display: 'flex',
                alignItems: 'center',
                gap: 'var(--space-sm)',
                cursor: 'pointer',
                userSelect: 'none',
              }}
              onClick={() => toggleFunction(functionName)}
            >
              <span
                style={{
                  fontSize: '12px',
                  color: 'var(--text-tertiary)',
                  transform: isExpanded ? 'rotate(90deg)' : 'rotate(0deg)',
                  transition: 'transform 0.15s ease',
                }}
              >
                ▶
              </span>
              <span
                style={{
                  fontFamily: 'monospace',
                  fontSize: '13px',
                  fontWeight: '600',
                  color: 'var(--text-primary)',
                  backgroundColor: 'var(--bg-tertiary)',
                  padding: '2px 8px',
                  borderRadius: '4px',
                }}
              >
                {functionName}
              </span>
              <span
                style={{
                  fontSize: '11px',
                  color: 'var(--text-tertiary)',
                  marginLeft: 'auto',
                }}
              >
                {functionTests.length} test{functionTests.length > 1 ? 's' : ''}
              </span>
            </div>

            {/* Tests list */}
            {isExpanded && (
              <div
                style={{
                  padding: 'var(--space-sm)',
                  backgroundColor: 'var(--bg-primary)',
                  borderTop: '1px solid var(--border-color)',
                }}
              >
                {functionTests.map((test, index) => (
                  <div
                    key={`${test.testFile}-${test.testName}-${index}`}
                    style={{
                      padding: 'var(--space-xs)',
                      marginBottom: 'var(--space-xs)',
                      backgroundColor: 'var(--bg-secondary)',
                      borderRadius: 'var(--radius-sm)',
                      border: '1px solid var(--border-color)',
                      cursor: 'pointer',
                      transition: 'all 0.15s ease',
                    }}
                    onClick={() => onTestClick(test)}
                    onMouseEnter={(e) => {
                      e.currentTarget.style.backgroundColor = 'var(--bg-tertiary)';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.backgroundColor = 'var(--bg-secondary)';
                    }}
                  >
                    <div
                      style={{
                        fontSize: '12px',
                        fontWeight: '500',
                        color: 'var(--text-primary)',
                        marginBottom: '4px',
                      }}
                    >
                      {test.testName}
                    </div>
                    <div
                      style={{
                        fontSize: '11px',
                        color: 'var(--text-tertiary)',
                        display: 'flex',
                        alignItems: 'center',
                        gap: '4px',
                      }}
                    >
                      <span>in</span>
                      <span
                        style={{
                          fontFamily: 'monospace',
                          color: 'var(--accent-primary)',
                          cursor: 'pointer',
                          textDecoration: 'underline',
                        }}
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
                      <div
                        style={{
                          fontSize: '10px',
                          color: 'var(--text-secondary)',
                          marginTop: 'var(--space-xs)',
                          fontStyle: 'italic',
                          borderTop: '1px solid var(--border-color)',
                          paddingTop: 'var(--space-xs)',
                        }}
                      >
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
