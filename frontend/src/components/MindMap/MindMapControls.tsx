import React from 'react';
import type { LayoutType } from '../../types';

interface MindMapControlsProps {
  layout: LayoutType;
  onLayoutChange: (layout: LayoutType) => void;
  onZoomIn: () => void;
  onZoomOut: () => void;
  onResetView: () => void;
  searchTerm: string;
  onSearchChange: (term: string) => void;
  showMinimap: boolean;
  onToggleMinimap: () => void;
  scale: number;
}

const buttonStyle: React.CSSProperties = {
  padding: '4px 8px',
  border: '1px solid var(--border-color)',
  borderRadius: 'var(--radius-sm)',
  backgroundColor: 'var(--bg-primary)',
  color: 'var(--text-secondary)',
  cursor: 'pointer',
  fontSize: '12px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  minWidth: '28px',
  height: '28px',
  transition: 'all 0.15s ease',
};

const activeButtonStyle: React.CSSProperties = {
  ...buttonStyle,
  backgroundColor: 'var(--accent-primary)',
  color: 'white',
  borderColor: 'var(--accent-primary)',
};

export const MindMapControls: React.FC<MindMapControlsProps> = ({
  layout,
  onLayoutChange,
  onZoomIn,
  onZoomOut,
  onResetView,
  searchTerm,
  onSearchChange,
  showMinimap,
  onToggleMinimap,
  scale,
}) => {
  return (
    <div
      style={{
        display: 'flex',
        flexWrap: 'wrap',
        gap: 'var(--space-sm)',
        alignItems: 'center',
        paddingBottom: 'var(--space-sm)',
        borderBottom: '1px solid var(--border-color)',
      }}
    >
      {/* Layout toggle */}
      <div style={{ display: 'flex', gap: '2px' }}>
        <button
          type="button"
          style={layout === 'horizontal' ? activeButtonStyle : buttonStyle}
          onClick={() => onLayoutChange('horizontal')}
          title="Horizontal Layout"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <line x1="5" y1="12" x2="19" y2="12" />
            <polyline points="12 5 19 12 12 19" />
          </svg>
        </button>
        <button
          type="button"
          style={layout === 'radial' ? activeButtonStyle : buttonStyle}
          onClick={() => onLayoutChange('radial')}
          title="Radial Layout"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="12" cy="12" r="10" />
            <circle cx="12" cy="12" r="3" />
          </svg>
        </button>
        <button
          type="button"
          style={layout === 'clustered' ? activeButtonStyle : buttonStyle}
          onClick={() => onLayoutChange('clustered')}
          title="Clustered Layout"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="3" y="3" width="7" height="7" />
            <rect x="14" y="3" width="7" height="7" />
            <rect x="3" y="14" width="7" height="7" />
            <rect x="14" y="14" width="7" height="7" />
          </svg>
        </button>
      </div>

      <div style={{ width: '1px', height: '20px', backgroundColor: 'var(--border-color)' }} />

      {/* Zoom controls */}
      <div style={{ display: 'flex', gap: '2px', alignItems: 'center' }}>
        <button
          type="button"
          style={buttonStyle}
          onClick={onZoomOut}
          title="Zoom Out"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
            <line x1="8" y1="11" x2="14" y2="11" />
          </svg>
        </button>
        <span
          style={{
            fontSize: '11px',
            color: 'var(--text-tertiary)',
            minWidth: '40px',
            textAlign: 'center',
          }}
        >
          {Math.round(scale * 100)}%
        </span>
        <button
          type="button"
          style={buttonStyle}
          onClick={onZoomIn}
          title="Zoom In"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
            <line x1="11" y1="8" x2="11" y2="14" />
            <line x1="8" y1="11" x2="14" y2="11" />
          </svg>
        </button>
        <button
          type="button"
          style={buttonStyle}
          onClick={onResetView}
          title="Reset View"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" />
            <path d="M3 3v5h5" />
          </svg>
        </button>
      </div>

      <div style={{ width: '1px', height: '20px', backgroundColor: 'var(--border-color)' }} />

      {/* Search */}
      <div style={{ position: 'relative', flex: '1', minWidth: '120px', maxWidth: '200px' }}>
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          style={{
            position: 'absolute',
            left: '8px',
            top: '50%',
            transform: 'translateY(-50%)',
            color: 'var(--text-tertiary)',
          }}
        >
          <circle cx="11" cy="11" r="8" />
          <line x1="21" y1="21" x2="16.65" y2="16.65" />
        </svg>
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search tests..."
          style={{
            width: '100%',
            padding: '4px 8px 4px 28px',
            border: '1px solid var(--border-color)',
            borderRadius: 'var(--radius-sm)',
            fontSize: '12px',
            backgroundColor: 'var(--bg-primary)',
            color: 'var(--text-primary)',
            height: '28px',
            outline: 'none',
          }}
        />
        {searchTerm && (
          <button
            type="button"
            onClick={() => onSearchChange('')}
            style={{
              position: 'absolute',
              right: '4px',
              top: '50%',
              transform: 'translateY(-50%)',
              border: 'none',
              background: 'none',
              cursor: 'pointer',
              padding: '2px',
              color: 'var(--text-tertiary)',
            }}
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        )}
      </div>

      {/* Minimap toggle */}
      <button
        type="button"
        style={showMinimap ? activeButtonStyle : buttonStyle}
        onClick={onToggleMinimap}
        title={showMinimap ? 'Hide Minimap' : 'Show Minimap'}
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <rect x="12" y="12" width="8" height="8" rx="1" fill="currentColor" fillOpacity="0.3" />
        </svg>
      </button>
    </div>
  );
};
