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
  depth: number;
}

export const MindMap: React.FC<MindMapProps> = ({ data, onNodeClick }) => {
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
  const { totalWidth, totalHeight, positions, edges } = useMemo(() => {
    const positionedNodes: PositionedNode[] = [];
    const edgePairs: Array<{ fromId: string; toId: string }> = [];

    const functions = data.children || [];
    const hasGrandchildren = functions.some((fn) => (fn.children?.length ?? 0) > 0);

    const addEdge = (fromId: string, toId: string) => {
      edgePairs.push({ fromId, toId });
    };

    const buildHorizontalLayout = (groups: MindMapNode[][]) => {
      const minHeight = 400;
      const clusterGapSlots = Math.max(groups.length - 1, 0);
      const totalSlots = groups.reduce((count, group) => {
        return count + group.reduce((sum, fn) => sum + Math.max(1, fn.children?.length ?? 0), 0);
      }, 0) + clusterGapSlots;

      const height = Math.max(minHeight, totalSlots * vSpacing + 100);
      const rootX = xPadding + nodeWidth / 2;
      const functionX = rootX + hSpacing;
      const testX = rootX + hSpacing * (hasGrandchildren ? 2 : 1);

      let cursorY = (height - (totalSlots - 1) * vSpacing) / 2;
      const functionPositions: PositionedNode[] = [];

      groups.forEach((group, groupIndex) => {
        group.forEach((fn) => {
          const tests = fn.children || [];
          if (tests.length > 0) {
            const firstTestY = cursorY;
            tests.forEach((test) => {
              positionedNodes.push({
                node: test,
                x: testX,
                y: cursorY,
                isRoot: false,
                depth: 2,
              });
              addEdge(fn.id, test.id);
              cursorY += vSpacing;
            });
            const lastTestY = cursorY - vSpacing;
            functionPositions.push({
              node: fn,
              x: functionX,
              y: (firstTestY + lastTestY) / 2,
              isRoot: false,
              depth: 1,
            });
          } else {
            functionPositions.push({
              node: fn,
              x: functionX,
              y: cursorY,
              isRoot: false,
              depth: 1,
            });
            cursorY += vSpacing;
          }
        });

        if (groupIndex < groups.length - 1) {
          cursorY += vSpacing;
        }
      });

      const rootY = functionPositions.length > 0
        ? functionPositions.reduce((sum, fn) => sum + fn.y, 0) / functionPositions.length
        : height / 2;

      positionedNodes.push({
        node: data,
        x: rootX,
        y: rootY,
        isRoot: true,
        depth: 0,
      });

      functionPositions.forEach((fn) => {
        positionedNodes.push(fn);
        addEdge(data.id, fn.node.id);
      });

      const width = xPadding * 2 + nodeWidth + hSpacing * (hasGrandchildren ? 2 : 1);
      return { width, height };
    };

    let width = 0;
    let height = 0;

    if (layout === 'radial') {
      const functionCount = functions.length;
      const functionRadius = radialRadius;
      const testRadius = radialRadius + 220;
      const outerRadius = hasGrandchildren ? testRadius : functionRadius;
      const size = outerRadius * 2 + nodeWidth + xPadding * 2;

      width = size;
      height = size;

      const centerX = size / 2;
      const centerY = size / 2;

      positionedNodes.push({ node: data, x: centerX, y: centerY, isRoot: true, depth: 0 });

      if (functionCount > 0) {
        const angleStep = (2 * Math.PI) / functionCount;
        functions.forEach((fn, index) => {
          const angle = index * angleStep - Math.PI / 2;
          const fnX = centerX + functionRadius * Math.cos(angle);
          const fnY = centerY + functionRadius * Math.sin(angle);
          positionedNodes.push({ node: fn, x: fnX, y: fnY, isRoot: false, depth: 1 });
          addEdge(data.id, fn.id);

          const tests = fn.children || [];
          if (tests.length > 0) {
            const sector = Math.min(angleStep * 0.8, Math.PI / 2);
            const startAngle = angle - sector / 2;
            const testStep = tests.length > 1 ? sector / (tests.length - 1) : 0;
            tests.forEach((test, testIndex) => {
              const testAngle = tests.length > 1 ? startAngle + testStep * testIndex : angle;
              const testX = centerX + testRadius * Math.cos(testAngle);
              const testY = centerY + testRadius * Math.sin(testAngle);
              positionedNodes.push({ node: test, x: testX, y: testY, isRoot: false, depth: 2 });
              addEdge(fn.id, test.id);
            });
          }
        });
      }
    } else {
      if (layout === 'clustered' && functions.length > 0) {
        const clusters = new Map<string, MindMapNode[]>();
        functions.forEach((fn) => {
          const key = fn.label.charAt(0).toUpperCase() || '#';
          if (!clusters.has(key)) {
            clusters.set(key, []);
          }
          clusters.get(key)!.push(fn);
        });
        const orderedGroups = Array.from(clusters.keys()).sort().map((key) => clusters.get(key)!);
        const layoutResult = buildHorizontalLayout(orderedGroups);
        width = layoutResult.width;
        height = layoutResult.height;
      } else {
        const layoutResult = buildHorizontalLayout([functions]);
        width = layoutResult.width;
        height = layoutResult.height;
      }
    }

    return { totalWidth: width, totalHeight: height, positions: positionedNodes, edges: edgePairs };
  }, [data, layout]);

  // Search filtering
  const matchingNodeIds = useMemo(() => {
    if (!searchTerm.trim()) return new Set<string>();
    const term = searchTerm.toLowerCase();
    const matches = new Set<string>();

    const visit = (node: MindMapNode) => {
      if (node.label.toLowerCase().includes(term) || node.edgeLabel?.toLowerCase().includes(term)) {
        matches.add(node.id);
      }
      node.children?.forEach(visit);
    };

    visit(data);
    return matches;
  }, [searchTerm, data]);

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

  const positionsById = useMemo(() => {
    const map = new Map<string, PositionedNode>();
    positions.forEach((pos) => {
      map.set(pos.node.id, pos);
    });
    return map;
  }, [positions]);
  const orderedPositions = useMemo(() => {
    return [...positions].sort((a, b) => a.depth - b.depth);
  }, [positions]);

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
            <linearGradient id="leafGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#0ea5e9" />
              <stop offset="100%" stopColor="#0284c7" />
            </linearGradient>
            <linearGradient id="highlightGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="#f59e0b" />
              <stop offset="100%" stopColor="#d97706" />
            </linearGradient>
          </defs>

          {/* Edges */}
          {edges.map((edge) => {
            const from = positionsById.get(edge.fromId);
            const to = positionsById.get(edge.toId);
            if (!from || !to) {
              return null;
            }
            const path = renderEdge(from, to);
            const labelX = layout === 'radial'
              ? (from.x + to.x) / 2
              : (from.x + nodeWidth / 2 + to.x - nodeWidth / 2) / 2;
            const labelY = (from.y + to.y) / 2;
            const edgeLabel = to.node.edgeLabel;

            return (
              <g key={`edge-group-${edge.fromId}-${edge.toId}`} style={{ opacity: getNodeOpacity(to.node.id) }}>
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

                {edgeLabel && (
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
                      <span title={edgeLabel} style={{
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
                        {edgeLabel}
                      </span>
                    </div>
                  </foreignObject>
                )}
              </g>
            );
          })}

          {/* Nodes */}
          {orderedPositions.map((pos) => {
            const isLeaf = !pos.node.children || pos.node.children.length === 0;
            const displayLabel = pos.isRoot
              ? (pos.node.label.length > 25 ? '...' + pos.node.label.slice(-22) : pos.node.label)
              : (pos.node.label.length > 25 ? pos.node.label.substring(0, 22) + '...' : pos.node.label);
            const fill = matchingNodeIds.has(pos.node.id) && searchTerm
              ? 'url(#highlightGradient)'
              : pos.isRoot
                ? 'url(#rootGradient)'
                : pos.depth === 1
                  ? 'url(#childGradient)'
                  : 'url(#leafGradient)';
            const fontSize = pos.isRoot ? '14' : pos.depth === 1 ? '12' : '11';
            const fontWeight = pos.isRoot ? '600' : '500';

            return (
              <g
                key={pos.node.id}
                className={`node ${pos.isRoot ? 'root-node' : pos.depth === 1 ? 'child-node' : 'leaf-node'}`}
                onClick={() => {
                  if (isLeaf) {
                    onNodeClick?.(pos.node.id);
                  }
                }}
                style={{
                  cursor: isLeaf ? 'pointer' : 'default',
                  opacity: getNodeOpacity(pos.node.id),
                  transition: 'all 0.3s ease-out',
                }}
                filter="url(#shadow)"
              >
                <rect
                  x={pos.x - nodeWidth / 2}
                  y={pos.y - nodeHeight / 2}
                  width={nodeWidth}
                  height={nodeHeight}
                  rx="8"
                  fill={fill}
                  style={{ transition: 'all 0.3s ease-out' }}
                />
                <text
                  x={pos.x}
                  y={pos.y}
                  textAnchor="middle"
                  dominantBaseline="middle"
                  fill="white"
                  fontSize={fontSize}
                  fontWeight={fontWeight}
                  style={{ pointerEvents: 'none', userSelect: 'none' }}
                >
                  {displayLabel}
                </text>
              </g>
            );
          })}
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
