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
        className="input input-mono"
      />
    </div>
  );
};
