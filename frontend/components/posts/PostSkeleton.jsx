import React from 'react';

const PostSkeleton = ({ style = {} }) => {
  return (
    <div 
      className="rounded-xl p-4 animate-pulse post-skeleton" 
      style={{ 
        backgroundColor: 'var(--primary-background)',
        ...style 
      }}
      role="status"
      aria-label="Loading post content"
      aria-live="polite"
    >
      {/* Screen reader only text */}
      <span className="sr-only">Loading post content, please wait...</span>
      {/* Header */}
      <div className="flex items-center gap-2 mb-3">
        <div 
          className="w-10 h-10 rounded-full loading-skeleton"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
        <div className="flex-1">
          <div 
            className="h-4 w-32 rounded loading-skeleton mb-1"
            style={{ backgroundColor: 'var(--secondary-background)' }}
            aria-hidden="true"
          />
          <div 
            className="h-3 w-20 rounded loading-skeleton"
            style={{ backgroundColor: 'var(--secondary-background)' }}
            aria-hidden="true"
          />
        </div>
      </div>
      
      {/* Content */}
      <div className="mb-4">
        <div 
          className="h-4 w-full rounded loading-skeleton mb-2"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
        <div 
          className="h-4 w-3/4 rounded loading-skeleton mb-2"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
        <div 
          className="h-4 w-1/2 rounded loading-skeleton"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
      </div>
      
      {/* Actions */}
      <div className="flex items-center gap-4 pt-3 border-t" style={{ borderColor: 'var(--border-color)' }}>
        <div 
          className="h-8 w-16 rounded loading-skeleton"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
        <div 
          className="h-8 w-20 rounded loading-skeleton"
          style={{ backgroundColor: 'var(--secondary-background)' }}
          aria-hidden="true"
        />
      </div>
    </div>
  );
};

export default PostSkeleton;