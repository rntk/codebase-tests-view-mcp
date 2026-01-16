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
        <div style={{
            display: 'flex',
            alignItems: 'center',
            gap: '4px',
            marginBottom: '16px',
            fontSize: '13px',
            color: 'var(--text-secondary)',
            flexWrap: 'wrap'
        }}>
            <button
                onClick={handleRootClick}
                style={{
                    color: currentPath === '.' ? 'var(--text-primary)' : 'var(--accent-primary)',
                    fontWeight: currentPath === '.' ? '600' : '400',
                    padding: '2px 4px',
                    borderRadius: 'var(--radius-sm)'
                }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--accent-soft)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            >
                root
            </button>

            {segments.map((segment, index) => (
                <React.Fragment key={index}>
                    <span>/</span>
                    <button
                        onClick={() => handleSegmentClick(index)}
                        style={{
                            color: index === segments.length - 1 ? 'var(--text-primary)' : 'var(--accent-primary)',
                            fontWeight: index === segments.length - 1 ? '600' : '400',
                            padding: '2px 4px',
                            borderRadius: 'var(--radius-sm)'
                        }}
                        onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--accent-soft)'}
                        onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                    >
                        {segment}
                    </button>
                </React.Fragment>
            ))}
        </div>
    );
};
