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
      className="modal-backdrop"
    >
      <div className="modal-content">
        {/* Header */}
        <div className="panel-header">
          <h2 className="panel-title">
            Export for AI Agent
          </h2>
          <button
            onClick={onClose}
            className="btn-icon btn-icon--close"
          >
            ×
          </button>
        </div>

        {/* Content */}
        <div className="panel-content">
          {loading && (
            <div className="empty-state">
              Generating export...
            </div>
          )}

          {!loading && !exportData && (
            <div className="empty-state">
              No data to export
            </div>
          )}

          {!loading && exportData && (
            <>
              <div className="tip-box">
                <strong>Tip:</strong> Copy this formatted context and paste it directly into your AI
                coding assistant. It includes code context, comments, and test information.
              </div>

              <div className="code-block-dark">
                <div className="code-block-header">
                  <span className="code-block-filename">
                    {exportData.sourceFile}
                  </span>
                  <button
                    onClick={handleCopy}
                    className={`btn ${copied ? 'btn-success' : 'btn-primary'}`}
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
                <pre className="code-block-content">
                  {exportData.formatted}
                </pre>
              </div>

              {/* Stats */}
              <div className="export-stats">
                <span>
                  <strong>{exportData.codeContext.length}</strong> comment blocks
                </span>
                {exportData.tests && (
                  <span>
                    <strong>{exportData.tests.length}</strong> tests
                  </span>
                )}
              </div>
            </>
          )}
        </div>

        {/* Footer */}
        <div className="panel-footer">
          <button
            onClick={onClose}
            className="btn btn-secondary"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
