import React from 'react';

interface PathInputProps {
  path: string;
  onChange: (path: string) => void;
}

export const PathInput: React.FC<PathInputProps> = ({ path, onChange }) => {
  return (
    <div className="path-input">
      <input
        type="text"
        value={path}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Enter path..."
        style={{
          width: '100%',
          padding: '8px',
          marginBottom: '16px',
          border: '1px solid #ddd',
          borderRadius: '4px',
          fontFamily: 'monospace',
        }}
      />
    </div>
  );
};
