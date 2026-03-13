import React from 'react';

interface BreadcrumbsProps {
    currentPath: string;
    onPathChange: (path: string) => void;
}

export const Breadcrumbs: React.FC<BreadcrumbsProps> = ({ currentPath, onPathChange }) => {
    const segments = currentPath === '.' ? [] : currentPath.split('/');

    const handleSegmentClick = (index: number) => {
        const newPath = segments.slice(0, index + 1).join('/');
        onPathChange(newPath || '.');
    };

    const handleRootClick = () => {
        onPathChange('.');
    };

    return (
        <div className="breadcrumbs">
            <button
                onClick={handleRootClick}
                className={`breadcrumb-btn ${currentPath === '.' ? 'breadcrumb-btn--root' : ''}`}
            >
                root
            </button>

            {segments.map((segment, index) => (
                <React.Fragment key={index}>
                    <span>/</span>
                    <button
                        onClick={() => handleSegmentClick(index)}
                        className={`breadcrumb-btn ${index === segments.length - 1 ? 'breadcrumb-btn--current' : ''}`}
                    >
                        {segment}
                    </button>
                </React.Fragment>
            ))}
        </div>
    );
};
