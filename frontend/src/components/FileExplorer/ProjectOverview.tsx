import React, { useState, useMemo } from 'react';
import type { OverviewResponse, TestDetail } from '../../types';

interface ProjectOverviewProps {
  overview: OverviewResponse | null;
  loading: boolean;
  error: string | null;
  onTestClick: (test: TestDetail) => void;
  onSourceFileClick: (sourceFile: string, line: number) => void;
}

export const ProjectOverview: React.FC<ProjectOverviewProps> = ({
  overview,
  loading,
  error,
  onTestClick,
  onSourceFileClick,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedFunctions, setExpandedFunctions] = useState<Set<string>>(new Set());
  const [activeView, setActiveView] = useState<'functions' | 'tests'>('functions');

  // Compute all tests from overview - must be before any early returns
  const allTests = useMemo(() => {
    if (!overview) return [];
    const tests: TestDetail[] = [];
    Object.values(overview.testsBySourceFile).forEach((fileTests) => {
      tests.push(...fileTests);
    });
    return tests;
  }, [overview?.testsBySourceFile]);

  // Filter functions based on search query
  const filteredFunctions = useMemo(() => {
    if (!overview) return [];
    if (!searchQuery.trim()) return overview.functions;
    const query = searchQuery.toLowerCase();
    return overview.functions.filter(
      (fn) =>
        fn.functionName.toLowerCase().includes(query) ||
        fn.sourceFile.toLowerCase().includes(query)
    );
  }, [overview?.functions, searchQuery]);

  // Filter tests based on search query
  const filteredTests = useMemo(() => {
    if (!searchQuery.trim()) return allTests;
    const query = searchQuery.toLowerCase();
    return allTests.filter(
      (test) =>
        test.testName.toLowerCase().includes(query) ||
        test.functionName.toLowerCase().includes(query) ||
        test.testFile.toLowerCase().includes(query)
    );
  }, [allTests, searchQuery]);

  const toggleFunction = (functionKey: string) => {
    setExpandedFunctions((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(functionKey)) {
        newSet.delete(functionKey);
      } else {
        newSet.add(functionKey);
      }
      return newSet;
    });
  };

  // Early returns after all hooks
  if (loading) {
    return (
      <div className="loading-state">
        Loading project overview...
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

  if (!overview || overview.totalTests === 0) {
    return (
      <div className="empty-state">
        No tests found in the project.
        <br />
        <span className="text-sm text-secondary">
          Tests will appear here once metadata is submitted.
        </span>
      </div>
    );
  }

  return (
    <div className="project-overview">
      {/* Statistics Cards */}
      <div className="overview-stats">
        <div className="stat-card">
          <div className="stat-value">{overview.totalTests}</div>
          <div className="stat-label">Total Tests</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{overview.totalFunctions}</div>
          <div className="stat-label">Functions with Tests</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{overview.totalSourceFiles}</div>
          <div className="stat-label">Source Files</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{overview.totalTestFiles}</div>
          <div className="stat-label">Test Files</div>
        </div>
      </div>

      {/* Search and View Toggle */}
      <div className="overview-controls">
        <div className="search-input-wrapper" style={{ maxWidth: '100%', flex: 1 }}>
          <span className="search-icon">🔍</span>
          <input
            type="text"
            className="search-input"
            placeholder="Search functions, tests, or files..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button
              className="search-clear-btn"
              onClick={() => setSearchQuery('')}
            >
              ×
            </button>
          )}
        </div>
        <div className="view-toggle">
          <button
            className={`btn-control ${activeView === 'functions' ? 'btn-control--active' : ''}`}
            onClick={() => setActiveView('functions')}
          >
            Functions ({filteredFunctions.length})
          </button>
          <button
            className={`btn-control ${activeView === 'tests' ? 'btn-control--active' : ''}`}
            onClick={() => setActiveView('tests')}
          >
            Tests ({filteredTests.length})
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="overview-content">
        {activeView === 'functions' ? (
          <div className="functions-list">
            {filteredFunctions.length === 0 ? (
              <div className="empty-state">
                No functions match your search
              </div>
            ) : (
              filteredFunctions.map((fn) => {
                const functionKey = `${fn.sourceFile}::${fn.functionName}`;
                const isExpanded = expandedFunctions.has(functionKey);

                return (
                  <div key={functionKey} className="function-accordion">
                    <div
                      className="function-header"
                      onClick={() => toggleFunction(functionKey)}
                    >
                      <span
                        className={`function-toggle ${isExpanded ? 'function-toggle--expanded' : ''}`}
                      >
                        ▶
                      </span>
                      <span className="function-name">{fn.functionName}</span>
                      <span className="function-meta">
                        <span className="inline-code text-sm">{fn.sourceFile}</span>
                        <span className="badge badge-secondary">
                          {fn.testCount} test{fn.testCount > 1 ? 's' : ''}
                        </span>
                      </span>
                    </div>

                    {isExpanded && (
                      <div className="function-tests">
                        {fn.tests.map((test, index) => (
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
              })
            )}
          </div>
        ) : (
          <div className="tests-list">
            {filteredTests.length === 0 ? (
              <div className="empty-state">
                No tests match your search
              </div>
            ) : (
              filteredTests.map((test, index) => (
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
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
};
