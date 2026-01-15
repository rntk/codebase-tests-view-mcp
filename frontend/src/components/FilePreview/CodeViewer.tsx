import React from 'react';

interface CodeViewerProps {
  content: string;
  filename: string;
}

export const CodeViewer: React.FC<CodeViewerProps> = ({ content, filename }) => {
  return (
    <div>
      <h3 style={{ marginTop: 0, marginBottom: '12px', fontSize: '16px', fontFamily: 'monospace' }}>
        {filename}
      </h3>
      <pre
        style={{
          backgroundColor: '#f5f5f5',
          padding: '16px',
          borderRadius: '4px',
          overflow: 'auto',
          fontSize: '13px',
          lineHeight: '1.5',
          fontFamily: 'monospace',
        }}
      >
        {content}
      </pre>
    </div>
  );
};
