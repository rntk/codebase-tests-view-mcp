import React, { useState } from 'react';
import type { ExportContextResponse } from '../../types';

interface ExportModalProps {
  isOpen: boolean;
  onClose: () => void;
  exportData: ExportContextResponse | null;
  loading: boolean;
}

export const ExportModal: React.FC<ExportModalProps> = ({
  isOpen,
  onClose,
  exportData,
  loading,
}) => {
  const [copied, setCopied] = useState(false);

  if (!isOpen) return null;

  const handleCopy = async () => {
    if (exportData?.formatted) {
      await navigator.clipboard.writeText(exportData.formatted);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div
      onClick={handleBackdropClick}
      style={{
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.5)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
    >
      <div
        style={{
          backgroundColor: 'var(--bg-primary)',
          borderRadius: 'var(--radius-lg)',
          width: '90%',
          maxWidth: '800px',
          maxHeight: '90vh',
          display: 'flex',
          flexDirection: 'column',
          boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
        }}
      >
        {/* Header */}
        <div
          style={{
            padding: 'var(--space-md) var(--space-lg)',
            borderBottom: '1px solid var(--border-color)',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <h2
            style={{
              margin: 0,
              fontSize: '18px',
              fontWeight: '600',
              color: 'var(--text-primary)',
            }}
          >
            Export for AI Agent
          </h2>
          <button
            onClick={onClose}
            style={{
              backgroundColor: 'transparent',
              border: 'none',
              fontSize: '24px',
              color: 'var(--text-secondary)',
              cursor: 'pointer',
              padding: '4px 8px',
              lineHeight: 1,
            }}
          >
            ×
          </button>
        </div>

        {/* Content */}
        <div
          style={{
            padding: 'var(--space-lg)',
            overflow: 'auto',
            flex: 1,
          }}
        >
          {loading && (
            <div
              style={{
                textAlign: 'center',
                padding: 'var(--space-xl)',
                color: 'var(--text-tertiary)',
              }}
            >
              Generating export...
            </div>
          )}

          {!loading && !exportData && (
            <div
              style={{
                textAlign: 'center',
                padding: 'var(--space-xl)',
                color: 'var(--text-tertiary)',
              }}
            >
              No data to export
            </div>
          )}

          {!loading && exportData && (
            <>
              <div
                style={{
                  marginBottom: 'var(--space-md)',
                  padding: 'var(--space-md)',
                  backgroundColor: 'var(--bg-secondary)',
                  borderRadius: 'var(--radius-md)',
                  fontSize: '13px',
                  color: 'var(--text-secondary)',
                }}
              >
                <strong>Tip:</strong> Copy this formatted context and paste it directly into your AI
                coding assistant. It includes code context, comments, test information, and
                suggestions.
              </div>

              <div
                style={{
                  position: 'relative',
                  backgroundColor: '#1e1e1e',
                  borderRadius: 'var(--radius-md)',
                  overflow: 'hidden',
                }}
              >
                <div
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: 'var(--space-sm) var(--space-md)',
                    backgroundColor: '#2d2d2d',
                    borderBottom: '1px solid #3d3d3d',
                  }}
                >
                  <span
                    style={{
                      fontSize: '12px',
                      color: '#9e9e9e',
                      fontFamily: 'var(--font-mono)',
                    }}
                  >
                    {exportData.sourceFile}
                  </span>
                  <button
                    onClick={handleCopy}
                    style={{
                      padding: '6px 12px',
                      backgroundColor: copied ? 'var(--success)' : 'var(--accent-primary)',
                      color: 'white',
                      border: 'none',
                      borderRadius: 'var(--radius-sm)',
                      fontSize: '12px',
                      cursor: 'pointer',
                      display: 'flex',
                      alignItems: 'center',
                      gap: '4px',
                    }}
                  >
                    {copied ? (
                      <>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <polyline points="20,6 9,17 4,12" />
                        </svg>
                        Copied!
                      </>
                    ) : (
                      <>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                        </svg>
                        Copy to Clipboard
                      </>
                    )}
                  </button>
                </div>
                <pre
                  style={{
                    margin: 0,
                    padding: 'var(--space-md)',
                    fontSize: '13px',
                    lineHeight: 1.6,
                    color: '#d4d4d4',
                    fontFamily: 'var(--font-mono)',
                    overflow: 'auto',
                    maxHeight: '50vh',
                    whiteSpace: 'pre-wrap',
                    wordWrap: 'break-word',
                  }}
                >
                  {exportData.formatted}
                </pre>
              </div>

              {/* Stats */}
              <div
                style={{
                  display: 'flex',
                  gap: 'var(--space-md)',
                  marginTop: 'var(--space-md)',
                  fontSize: '12px',
                  color: 'var(--text-secondary)',
                }}
              >
                <span>
                  <strong>{exportData.codeContext.length}</strong> comment blocks
                </span>
                {exportData.tests && (
                  <span>
                    <strong>{exportData.tests.length}</strong> tests
                  </span>
                )}
                {exportData.suggestions && (
                  <span>
                    <strong>{exportData.suggestions.length}</strong> suggestions
                  </span>
                )}
              </div>
            </>
          )}
        </div>

        {/* Footer */}
        <div
          style={{
            padding: 'var(--space-md) var(--space-lg)',
            borderTop: '1px solid var(--border-color)',
            display: 'flex',
            justifyContent: 'flex-end',
          }}
        >
          <button
            onClick={onClose}
            style={{
              padding: '8px 16px',
              backgroundColor: 'var(--bg-secondary)',
              color: 'var(--text-primary)',
              border: '1px solid var(--border-color)',
              borderRadius: 'var(--radius-sm)',
              fontSize: '14px',
              cursor: 'pointer',
            }}
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
