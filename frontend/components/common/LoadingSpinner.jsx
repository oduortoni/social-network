import React from 'react';

const LoadingSpinner = ({ size = 'md', text = 'Loading...' }) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8'
  };

  return (
    <div className="flex items-center justify-center gap-2 py-4">
      <div 
        className={`${sizeClasses[size]} animate-spin rounded-full border-2 border-t-transparent`}
        style={{ borderColor: 'var(--secondary-text)', borderTopColor: 'transparent' }}
      />
      <span style={{ color: 'var(--secondary-text)' }} className="text-sm">
        {text}
      </span>
    </div>
  );
};

export default LoadingSpinner;