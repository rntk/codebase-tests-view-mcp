import React from 'react';
import type { MindMapNode, MindMapTransform } from '../../types';

interface PositionedNode {
  node: MindMapNode;
  x: number;
  y: number;
  isRoot: boolean;
}

interface MindMapMinimapProps {
  positions: PositionedNode[];
  totalWidth: number;
  totalHeight: number;
  transform: MindMapTransform;
  nodeWidth: number;
  nodeHeight: number;
}

export const MindMapMinimap: React.FC<MindMapMinimapProps> = ({
  positions,
  totalWidth,
  totalHeight,
  transform,
  nodeWidth,
  nodeHeight,
}) => {
  const minimapWidth = 120;
  const minimapHeight = 80;

  // Scale factor to fit the graph in the minimap
  const scaleX = minimapWidth / totalWidth;
  const scaleY = minimapHeight / totalHeight;
  const minimapScale = Math.min(scaleX, scaleY) * 0.9;

  // Calculate viewport rectangle
  const viewportWidth = (800 / transform.scale) * minimapScale;
  const viewportHeight = (400 / transform.scale) * minimapScale;
  const viewportX = (totalWidth * minimapScale - viewportWidth) / 2 - (transform.x / transform.scale) * minimapScale;
  const viewportY = (totalHeight * minimapScale - viewportHeight) / 2 - (transform.y / transform.scale) * minimapScale;

  return (
    <div className="mind-map-minimap">
      <svg
        viewBox={`0 0 ${totalWidth * minimapScale} ${totalHeight * minimapScale}`}
        className="mind-map-minimap-svg"
      >
        {/* Background */}
        <rect
          x="0"
          y="0"
          width={totalWidth * minimapScale}
          height={totalHeight * minimapScale}
          fill="var(--bg-secondary)"
        />

        {/* Nodes */}
        {positions.map((pos) => (
          <rect
            key={pos.node.id}
            x={(pos.x - nodeWidth / 2) * minimapScale}
            y={(pos.y - nodeHeight / 2) * minimapScale}
            width={nodeWidth * minimapScale}
            height={nodeHeight * minimapScale}
            rx={2}
            fill={pos.isRoot ? '#3b82f6' : '#10b981'}
            opacity={0.8}
          />
        ))}

        {/* Viewport indicator */}
        <rect
          x={Math.max(0, viewportX)}
          y={Math.max(0, viewportY)}
          width={viewportWidth}
          height={viewportHeight}
          fill="none"
          stroke="var(--accent-primary)"
          strokeWidth="1.5"
          strokeDasharray="3 2"
          opacity={0.8}
        />
      </svg>
    </div>
  );
};
