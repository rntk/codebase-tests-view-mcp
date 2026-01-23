import React, { useMemo, useState, useRef, useCallback } from 'react';
import type { MindMapNode, LayoutType, MindMapTransform } from '../../types';
import { MindMapControls } from './MindMapControls';
import { MindMapMinimap } from './MindMapMinimap';

interface MindMapProps {
  data: MindMapNode;
  onNodeClick?: (nodeId: string) => void;
}

interface PositionedNode {
  node: MindMapNode;
  x: number;
  y: number;
  isRoot: boolean;
}

export const MindMap: React.FC<MindMapProps> = ({ data, onNodeClick }) => {
  const children = useMemo(() => data.children || [], [data.children]);
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);

  // State
  const [layout, setLayout] = useState<LayoutType>('horizontal');
  const [transform, setTransform] = useState<MindMapTransform>({ x: 0, y: 0, scale: 1 });
  const [isPanning, setIsPanning] = useState(false);
  const [panStart, setPanStart] = useState({ x: 0, y: 0 });
  const [searchTerm, setSearchTerm] = useState('');
  const [showMinimap, setShowMinimap] = useState(true);

  // Layout configuration
  const nodeWidth = 200;
  const nodeHeight = 50;
  const hSpacing = 400;
  const vSpacing = 80;
  const xPadding = 50;
  const radialRadius = 250;

  // Calculate dimensions based on layout
  const { totalWidth, totalHeight, positions } = useMemo(() => {
    const positionedNodes: PositionedNode[] = [];
    let width: number;
    let height: number;

    if (layout === 'horizontal') {
      // Original horizontal layout
      const minHeight = 400;
      height = Math.max(minHeight, children.length * vSpacing + 100);
      width = xPadding * 2 + nodeWidth + hSpacing;

      const rootX = xPadding + nodeWidth / 2;
      const rootY = height / 2;
      const childX = rootX + hSpacing;

      positionedNodes.push({ node: data, x: rootX, y: rootY, isRoot: true });

      const totalChildrenHeight = (children.length - 1) * vSpacing;
      const startY = rootY - totalChildrenHeight / 2;

      children.forEach((child, index) => {
        const childY = children.length === 1 ? rootY : startY + index * vSpacing;
        positionedNodes.push({ node: child, x: childX, y: childY, isRoot: false });
      });
    } else if (layout === 'radial') {
      // Radial layout - nodes in a circle around root
      const size = radialRadius * 2 + nodeWidth + xPadding * 2;
      width = size;
      height = size;

      const centerX = size / 2;
      const centerY = size / 2;

      positionedNodes.push({ node: data, x: centerX, y: centerY, isRoot: true });

      if (children.length > 0) {
        const angleStep = (2 * Math.PI) / children.length;
        children.forEach((child, index) => {
          const angle = index * angleStep - Math.PI / 2; // Start from top
          const childX = centerX + radialRadius * Math.cos(angle);
          const childY = centerY + radialRadius * Math.sin(angle);
          positionedNodes.push({ node: child, x: childX, y: childY, isRoot: false });
        });
      }
    } else {
      // Clustered layout - group by first letter (simulated clustering)
      const clusters = new Map<string, MindMapNode[]>();
      children.forEach(child => {
        const key = child.label.charAt(0).toUpperCase();
        if (!clusters.has(key)) {
          clusters.set(key, []);
        }
        clusters.get(key)!.push(child);
      });

      const clusterCount = clusters.size;
      const nodesPerCluster = Math.ceil(children.length / Math.max(clusterCount, 1));
      height = Math.max(400, nodesPerCluster * vSpacing + 150);
      width = xPadding * 2 + nodeWidth + hSpacing + (clusterCount > 3 ? 200 : 0);

      const rootX = xPadding + nodeWidth / 2;
      const rootY = height / 2;

      positionedNodes.push({ node: data, x: rootX, y: rootY, isRoot: true });

      let clusterIndex = 0;
      clusters.forEach((clusterNodes) => {
        const clusterX = rootX + hSpacing + (clusterIndex % 2) * 100;
        const clusterStartY = rootY - ((clusterNodes.length - 1) * vSpacing) / 2;

        clusterNodes.forEach((child, nodeIndex) => {
          const childY = clusterStartY + nodeIndex * vSpacing + clusterIndex * 20;
          positionedNodes.push({ node: child, x: clusterX, y: childY, isRoot: false });
        });

        clusterIndex++;
      });
    }

    return { totalWidth: width, totalHeight: height, positions: positionedNodes };
  }, [data, children, layout]);

  // Search filtering
  const matchingNodeIds = useMemo(() => {
    if (!searchTerm.trim()) return new Set<string>();
    const term = searchTerm.toLowerCase();
    const matches = new Set<string>();

    if (data.label.toLowerCase().includes(term)) {
      matches.add(data.id);
    }
    children.forEach(child => {
      if (child.label.toLowerCase().includes(term) || child.edgeLabel?.toLowerCase().includes(term)) {
        matches.add(child.id);
      }
    });

    return matches;
  }, [searchTerm, data, children]);

  // Zoom handlers
  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault();
    const delta = e.deltaY > 0 ? 0.9 : 1.1;
    setTransform(prev => ({
      ...prev,
      scale: Math.min(Math.max(prev.scale * delta, 0.25), 3),
    }));
  }, []);

  // Pan handlers
  const handleMouseDown = useCallback((e: React.MouseEvent) => {
    if (e.button === 0) { // Left mouse button
      setIsPanning(true);
      setPanStart({ x: e.clientX - transform.x, y: e.clientY - transform.y });
    }
  }, [transform.x, transform.y]);

  const handleMouseMove = useCallback((e: React.MouseEvent) => {
    if (isPanning) {
      setTransform(prev => ({
        ...prev,
        x: e.clientX - panStart.x,
        y: e.clientY - panStart.y,
      }));
    }
  }, [isPanning, panStart]);

  const handleMouseUp = useCallback(() => {
    setIsPanning(false);
  }, []);

  // Reset view
  const handleResetView = useCallback(() => {
    setTransform({ x: 0, y: 0, scale: 1 });
  }, []);

  // Zoom controls
  const handleZoomIn = useCallback(() => {
    setTransform(prev => ({ ...prev, scale: Math.min(prev.scale * 1.2, 3) }));
  }, []);

  const handleZoomOut = useCallback(() => {
    setTransform(prev => ({ ...prev, scale: Math.max(prev.scale * 0.8, 0.25) }));
  }, []);

  // Get node style based on search
  const getNodeOpacity = (nodeId: string) => {
    if (!searchTerm.trim()) return 1;
    return matchingNodeIds.has(nodeId) ? 1 : 0.3;
  };

  // Render edge between two positions
  const renderEdge = (from: PositionedNode, to: PositionedNode) => {
    if (layout === 'radial') {
      // Straight line for radial
      return `M ${from.x} ${from.y} L ${to.x} ${to.y}`;
    }
    // Bezier curve for horizontal/clustered
    const midX = (from.x + nodeWidth / 2 + to.x - nodeWidth / 2) / 2;
    return `M ${from.x + nodeWidth / 2} ${from.y} C ${midX} ${from.y}, ${midX} ${to.y}, ${to.x - nodeWidth / 2} ${to.y}`;
  };

  const rootPosition = positions.find(p => p.isRoot);
  const childPositions = positions.filter(p => !p.isRoot);

  return (
    <div
      style={{
        width: '100%',
        background: 'var(--bg-secondary)',
        borderRadius: 'var(--radius-md)',
        border: '1px solid var(--border-color)',
        padding: 'var(--space-md)',
        position: 'relative',
      }}
    >
      <MindMapControls
        layout={layout}
        onLayoutChange={setLayout}
        onZoomIn={handleZoomIn}
        onZoomOut={handleZoomOut}
        onResetView={handleResetView}
        searchTerm={searchTerm}
        onSearchChange={setSearchTerm}
        showMinimap={showMinimap}
        onToggleMinimap={() => setShowMinimap(!showMinimap)}
        scale={transform.scale}
      />

      <div
        ref={containerRef}
        style={{
          overflow: 'hidden',
          cursor: isPanning ? 'grabbing' : 'grab',
          marginTop: 'var(--space-sm)',
        }}
        onWheel={handleWheel}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
      >
        <svg
          ref={svgRef}
          viewBox={`0 0 ${totalWidth} ${totalHeight}`}
          style={{
            width: '100%',
            maxWidth: '800px',
            height: 'auto',
            display: 'block',
            margin: '0 auto',
            transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`,
            transformOrigin: 'center center',
            transition: isPanning ? 'none' : 'transform 0.2s ease-out',
          }}
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
            <linearGradient id="highlightGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#f59e0b" />
              <stop offset="100%" stopColor="#d97706" />
            </linearGradient>
          </defs>

          {/* Edges */}
          {rootPosition && childPositions.map((childPos) => {
            const path = renderEdge(rootPosition, childPos);
            const labelX = layout === 'radial'
              ? (rootPosition.x + childPos.x) / 2
              : (rootPosition.x + nodeWidth / 2 + childPos.x - nodeWidth / 2) / 2;
            const labelY = (rootPosition.y + childPos.y) / 2;

            return (
              <g key={`edge-group-${childPos.node.id}`} style={{ opacity: getNodeOpacity(childPos.node.id) }}>
                <path
                  d={path}
                  fill="none"
                  stroke="var(--border-color)"
                  strokeWidth="2"
                  strokeOpacity="0.6"
                  style={{
                    transition: 'all 0.3s ease-out',
                  }}
                />

                {childPos.node.edgeLabel && (
                  <foreignObject
                    x={labelX - 100}
                    y={labelY - 12}
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
                      <span title={childPos.node.edgeLabel} style={{
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
                        {childPos.node.edgeLabel}
                      </span>
                    </div>
                  </foreignObject>
                )}
              </g>
            );
          })}

          {/* Root Node */}
          {rootPosition && (
            <g
              className="node root-node"
              onClick={() => onNodeClick?.(data.id)}
              style={{
                cursor: 'pointer',
                opacity: getNodeOpacity(data.id),
                transition: 'all 0.3s ease-out',
              }}
              filter="url(#shadow)"
            >
              <rect
                x={rootPosition.x - nodeWidth / 2}
                y={rootPosition.y - nodeHeight / 2}
                width={nodeWidth}
                height={nodeHeight}
                rx="8"
                fill={matchingNodeIds.has(data.id) && searchTerm ? 'url(#highlightGradient)' : 'url(#rootGradient)'}
                style={{ transition: 'all 0.3s ease-out' }}
              />
              <text
                x={rootPosition.x}
                y={rootPosition.y}
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
          )}

          {/* Child Nodes */}
          {childPositions.map((childPos) => (
            <g
              key={childPos.node.id}
              className="node child-node"
              onClick={() => onNodeClick?.(childPos.node.id)}
              style={{
                cursor: 'pointer',
                opacity: getNodeOpacity(childPos.node.id),
                transition: 'all 0.3s ease-out',
              }}
              filter="url(#shadow)"
            >
              <rect
                x={childPos.x - nodeWidth / 2}
                y={childPos.y - nodeHeight / 2}
                width={nodeWidth}
                height={nodeHeight}
                rx="8"
                fill={matchingNodeIds.has(childPos.node.id) && searchTerm ? 'url(#highlightGradient)' : 'url(#childGradient)'}
                style={{ transition: 'all 0.3s ease-out' }}
              />
              <text
                x={childPos.x}
                y={childPos.y}
                textAnchor="middle"
                dominantBaseline="middle"
                fill="white"
                fontSize="12"
                fontWeight="500"
                style={{ pointerEvents: 'none', userSelect: 'none' }}
              >
                {childPos.node.label.length > 25 ? childPos.node.label.substring(0, 22) + '...' : childPos.node.label}
              </text>
            </g>
          ))}
        </svg>
      </div>

      {showMinimap && (
        <MindMapMinimap
          positions={positions}
          totalWidth={totalWidth}
          totalHeight={totalHeight}
          transform={transform}
          nodeWidth={nodeWidth}
          nodeHeight={nodeHeight}
        />
      )}
    </div>
  );
};
