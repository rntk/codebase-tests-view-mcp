import React, { useEffect, useRef, useState } from 'react';
import './ThreePanel.css';

interface ThreePanelProps {
  left: React.ReactNode;
  center: React.ReactNode;
  right: React.ReactNode;
}

export const ThreePanel: React.FC<ThreePanelProps> = ({ left, center, right }) => {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const dragRef = useRef<{
    type: 'left' | 'right';
    startX: number;
    startLeft: number;
    startRight: number;
    containerWidth: number;
  } | null>(null);

  const [leftWidth, setLeftWidth] = useState(260);
  const [rightWidth, setRightWidth] = useState(420);

  useEffect(() => {
    const handlePointerMove = (event: PointerEvent) => {
      if (!dragRef.current) return;
      const { type, startX, startLeft, startRight, containerWidth } = dragRef.current;
      const deltaX = event.clientX - startX;
      const minLeft = 200;
      const minRight = 260;
      const minCenter = 320;

      if (type === 'left') {
        const maxLeft = containerWidth - minCenter - startRight;
        const nextLeft = Math.max(minLeft, Math.min(maxLeft, startLeft + deltaX));
        setLeftWidth(nextLeft);
      } else {
        const maxRight = containerWidth - minCenter - startLeft;
        const nextRight = Math.max(minRight, Math.min(maxRight, startRight - deltaX));
        setRightWidth(nextRight);
      }
    };

    const handlePointerUp = () => {
      dragRef.current = null;
    };

    window.addEventListener('pointermove', handlePointerMove);
    window.addEventListener('pointerup', handlePointerUp);

    return () => {
      window.removeEventListener('pointermove', handlePointerMove);
      window.removeEventListener('pointerup', handlePointerUp);
    };
  }, []);

  const handlePointerDown = (type: 'left' | 'right') => (event: React.PointerEvent) => {
    event.preventDefault();
    const containerWidth = containerRef.current?.getBoundingClientRect().width ?? 0;
    dragRef.current = {
      type,
      startX: event.clientX,
      startLeft: leftWidth,
      startRight: rightWidth,
      containerWidth,
    };
  };

  return (
    <div
      className="three-panel"
      ref={containerRef}
      style={{
        gridTemplateColumns: `${leftWidth}px 6px 1fr 6px ${rightWidth}px`,
      }}
    >
      <div className="panel panel-left">{left}</div>
      <div
        className="resizer resizer-left"
        onPointerDown={handlePointerDown('left')}
        role="separator"
        aria-orientation="vertical"
      />
      <div className="panel panel-center">{center}</div>
      <div
        className="resizer resizer-right"
        onPointerDown={handlePointerDown('right')}
        role="separator"
        aria-orientation="vertical"
      />
      <div className="panel panel-right">{right}</div>
    </div>
  );
};
