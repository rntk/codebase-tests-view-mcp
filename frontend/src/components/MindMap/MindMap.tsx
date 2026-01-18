import React, { useMemo } from 'react';
import type { MindMapNode } from '../../types';

interface MindMapProps {
  data: MindMapNode;
  onNodeClick?: (nodeId: string) => void;
}

export const MindMap: React.FC<MindMapProps> = ({ data, onNodeClick }) => {
  const children = useMemo(() => data.children || [], [data.children]);

  // Layout configuration
  const nodeWidth = 200;
  const nodeHeight = 50;
  const hSpacing = 400; // Horizontal space between root center and child center
  const vSpacing = 80;  // Vertical space between children centers
  const xPadding = 50;  // Padding from the sides

  // Calculate dimensions
  // Height depends on number of children
  const minHeight = 400;
  const calculatedHeight = Math.max(minHeight, children.length * vSpacing + 100);

  // Width: padding + node/2 + spacing + node/2 + padding
  const totalWidth = xPadding * 2 + nodeWidth + hSpacing;

  // Root Position (Left)
  const rootX = xPadding + nodeWidth / 2;
  const rootY = calculatedHeight / 2;

  // Children Column Position (Right)
  const childX = rootX + hSpacing;

  return (
    <div style={{
      width: '100%',
      overflowX: 'auto',
      background: 'var(--bg-secondary)',
      borderRadius: 'var(--radius-md)',
      border: '1px solid var(--border-color)',
      padding: 'var(--space-md)'
    }}>
      <svg
        viewBox={`0 0 ${totalWidth} ${calculatedHeight}`}
        style={{ width: '100%', maxWidth: '800px', height: 'auto', display: 'block', margin: '0 auto' }}
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
          <linearGradient id="rootGradient" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#3b82f6" />
            <stop offset="100%" stopColor="#2563eb" />
          </linearGradient>
          <linearGradient id="childGradient" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#10b981" />
            <stop offset="100%" stopColor="#059669" />
          </linearGradient>
        </defs>

        {/* Edges */}
        {children.map((child, index) => {
          // Distribute children vertically centered around rootY
          const totalChildrenHeight = (children.length - 1) * vSpacing;
          const startY = rootY - totalChildrenHeight / 2;
          const childY = children.length === 1 ? rootY : startY + index * vSpacing;

          // Control points for cubic bezier
          const startPoint = { x: rootX + nodeWidth / 2, y: rootY };
          const endPoint = { x: childX - nodeWidth / 2, y: childY };

          // Midpoint for control points (curve handling)
          const midX = (startPoint.x + endPoint.x) / 2;

          const path = `M ${startPoint.x} ${startPoint.y} 
                        C ${midX} ${startPoint.y}, ${midX} ${childY}, ${endPoint.x} ${childY}`;

          // Label Position (roughly center of the path)
          const labelX = midX;
          const labelY = (startPoint.y + childY) / 2;

          return (
            <g key={`edge-group-${child.id}`}>
              <path
                d={path}
                fill="none"
                stroke="var(--border-color)"
                strokeWidth="2"
                strokeOpacity="0.6"
              />

              {child.edgeLabel && (
                <foreignObject
                  x={labelX - 100} // Centered horizontally (200px width)
                  y={labelY - 12}  // Centered vertically (roughly)
                  width="200"
                  height="30"
                  style={{ overflow: 'visible' }}
                >
                  <div style={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    width: '100%',
                    height: '100%'
                  }}>
                    <span title={child.edgeLabel} style={{
                      background: 'var(--bg-primary)',
                      border: '1px solid var(--border-color)',
                      padding: '2px 8px',
                      borderRadius: '12px',
                      fontSize: '11px',
                      color: 'var(--text-secondary)',
                      boxShadow: '0 1px 2px rgba(0,0,0,0.05)',
                      whiteSpace: 'nowrap',
                      maxWidth: '180px',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis'
                    }}>
                      {child.edgeLabel}
                    </span>
                  </div>
                </foreignObject>
              )}
            </g>
          );
        })}

        {/* Root Node */}
        <g
          className="node root-node"
          onClick={() => onNodeClick?.(data.id)}
          style={{ cursor: 'pointer' }}
          filter="url(#shadow)"
        >
          <rect
            x={rootX - nodeWidth / 2}
            y={rootY - nodeHeight / 2}
            width={nodeWidth}
            height={nodeHeight}
            rx="8"
            fill="url(#rootGradient)"
          />
          <text
            x={rootX}
            y={rootY}
            textAnchor="middle"
            dominantBaseline="middle"
            fill="white"
            fontSize="14"
            fontWeight="600"
            style={{ pointerEvents: 'none', userSelect: 'none' }}
          >
            {data.label.length > 25 ? '...' + data.label.slice(-22) : data.label}
          </text>
        </g>

        {/* Child Nodes */}
        {children.map((child, index) => {
          const totalChildrenHeight = (children.length - 1) * vSpacing;
          const startY = rootY - totalChildrenHeight / 2;
          const childY = children.length === 1 ? rootY : startY + index * vSpacing;

          return (
            <g
              key={child.id}
              className="node child-node"
              onClick={() => onNodeClick?.(child.id)}
              style={{ cursor: 'pointer' }}
              filter="url(#shadow)"
            >
              <rect
                x={childX - nodeWidth / 2}
                y={childY - nodeHeight / 2}
                width={nodeWidth}
                height={nodeHeight}
                rx="8"
                fill="url(#childGradient)"
              />
              <text
                x={childX}
                y={childY}
                textAnchor="middle"
                dominantBaseline="middle"
                fill="white"
                fontSize="12"
                fontWeight="500"
                style={{ pointerEvents: 'none', userSelect: 'none' }}
              >
                {child.label.length > 25 ? child.label.substring(0, 22) + '...' : child.label}
              </text>
            </g>
          );
        })}
      </svg>
    </div>
  );
};
