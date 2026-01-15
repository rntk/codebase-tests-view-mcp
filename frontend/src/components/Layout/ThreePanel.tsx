import React from 'react';
import './ThreePanel.css';

interface ThreePanelProps {
  left: React.ReactNode;
  center: React.ReactNode;
  right: React.ReactNode;
}

export const ThreePanel: React.FC<ThreePanelProps> = ({ left, center, right }) => {
  return (
    <div className="three-panel">
      <div className="panel panel-left">{left}</div>
      <div className="panel panel-center">{center}</div>
      <div className="panel panel-right">{right}</div>
    </div>
  );
};
