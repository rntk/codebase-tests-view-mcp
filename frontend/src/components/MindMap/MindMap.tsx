import React, { useMemo } from 'react';
import type { MindMapNode } from '../../types';

interface MindMapProps {
  data: MindMapNode;
  onNodeClick?: (nodeId: string) => void;
}

export const MindMap: React.FC<MindMapProps> = ({ data, onNodeClick }) => {
  const width = 600;
  const height = 400;
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = 140;

  const children = useMemo(() => data.children || [], [data.children]);
  const angleStep = children.length > 0 ? (2 * Math.PI) / children.length : 0;

  return (
    <div style={{
      width: '100%',
      overflow: 'hidden',
      background: 'var(--bg-secondary)',
      borderRadius: 'var(--radius-md)',
      border: '1px solid var(--border-color)',
      padding: 'var(--space-md)'
    }}>
      <svg
        viewBox={`0 0 ${width} ${height}`}
        style={{ width: '100%', height: 'auto', display: 'block' }}
      >
        <defs>
          <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
            <feGaussianBlur in="SourceAlpha" stdDeviation="3" />
            <feOffset dx="0" dy="2" result="offsetblur" />
            <feComponentTransfer>
              <feFuncA type="linear" slope="0.1" />
            </feComponentTransfer>
            <feMerge>
              <feMergeNode />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
          <linearGradient id="mainNodeGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#60a5fa" />
            <stop offset="100%" stopColor="#2563eb" />
          </linearGradient>
          <linearGradient id="childNodeGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#34d399" />
            <stop offset="100%" stopColor="#059669" />
          </linearGradient>
        </defs>

        {/* Child nodes lines */}
        {children.map((child, index) => {
          const angle = angleStep * index - Math.PI / 2;
          const x = centerX + radius * Math.cos(angle);
          const y = centerY + radius * Math.sin(angle);

          return (
            <line
              key={`line-${child.id}`}
              x1={centerX}
              y1={centerY}
              x2={x}
              y2={y}
              stroke="var(--border-color)"
              strokeWidth="1.5"
              strokeDasharray="4 2"
            />
          );
        })}

        {/* Central node */}
        <g
          className="node central-node"
          onClick={() => onNodeClick?.(data.id)}
          style={{ cursor: 'pointer' }}
          filter="url(#shadow)"
        >
          <circle cx={centerX} cy={centerY} r={50} fill="url(#mainNodeGradient)" />
          <text
            x={centerX}
            y={centerY}
            textAnchor="middle"
            dominantBaseline="middle"
            fill="white"
            fontSize="14"
            fontWeight="600"
            style={{ pointerEvents: 'none' }}
          >
            {data.label.length > 15 ? data.label.substring(0, 12) + '...' : data.label}
          </text>
        </g>

        {/* Child nodes */}
        {children.map((child, index) => {
          const angle = angleStep * index - Math.PI / 2;
          const x = centerX + radius * Math.cos(angle);
          const y = centerY + radius * Math.sin(angle);

          return (
            <g
              key={child.id}
              className="node child-node"
              onClick={() => onNodeClick?.(child.id)}
              style={{ cursor: 'pointer' }}
              filter="url(#shadow)"
            >
              <circle cx={x} cy={y} r={35} fill="url(#childNodeGradient)" />
              <text
                x={x}
                y={y}
                textAnchor="middle"
                dominantBaseline="middle"
                fill="white"
                fontSize="11"
                fontWeight="500"
                style={{ pointerEvents: 'none' }}
              >
                {child.label.length > 10 ? child.label.substring(0, 8) + '...' : child.label}
              </text>
            </g>
          );
        })}
      </svg>
    </div>
  );
};
