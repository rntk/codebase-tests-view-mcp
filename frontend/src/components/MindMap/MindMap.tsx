import React from 'react';
import type { MindMapNode } from '../../types';

interface MindMapProps {
  data: MindMapNode;
  onNodeClick?: (nodeId: string) => void;
}

export const MindMap: React.FC<MindMapProps> = ({ data, onNodeClick }) => {
  const width = 400;
  const height = 300;
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = 100;

  const children = data.children || [];
  const angleStep = children.length > 0 ? (2 * Math.PI) / children.length : 0;

  return (
    <svg width={width} height={height} style={{ display: 'block', margin: '0 auto' }}>
      {/* Central node (source file) */}
      <g
        className="node central-node"
        onClick={() => onNodeClick?.(data.id)}
        style={{ cursor: 'pointer' }}
      >
        <circle cx={centerX} cy={centerY} r={40} fill="#2196F3" />
        <text
          x={centerX}
          y={centerY}
          textAnchor="middle"
          dominantBaseline="middle"
          fill="white"
          fontSize="12"
          fontWeight="bold"
        >
          {data.label.length > 10 ? data.label.substring(0, 10) + '...' : data.label}
        </text>
      </g>

      {/* Child nodes (tests) with connecting lines */}
      {children.map((child, index) => {
        const angle = angleStep * index - Math.PI / 2;
        const x = centerX + radius * Math.cos(angle);
        const y = centerY + radius * Math.sin(angle);

        return (
          <g key={child.id}>
            <line
              x1={centerX}
              y1={centerY}
              x2={x}
              y2={y}
              stroke="#999"
              strokeWidth="2"
            />
            <g
              className="node child-node"
              onClick={() => onNodeClick?.(child.id)}
              style={{ cursor: 'pointer' }}
            >
              <circle cx={x} cy={y} r={30} fill="#4CAF50" />
              <text
                x={x}
                y={y}
                textAnchor="middle"
                dominantBaseline="middle"
                fill="white"
                fontSize="10"
              >
                {child.label.length > 8 ? child.label.substring(0, 8) + '...' : child.label}
              </text>
            </g>
          </g>
        );
      })}
    </svg>
  );
};
