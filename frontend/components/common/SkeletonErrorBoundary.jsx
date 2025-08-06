import React from 'react';

/**
 * Error boundary component for skeleton loading
 * Falls back to text loading indicators if skeleton rendering fails
 */
class SkeletonErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    // Log the error for debugging purposes
    console.warn('Skeleton loading failed, falling back to text indicator:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      // Fallback to the provided fallback UI or default text loading
      return this.props.fallback || (
        <div className="flex justify-center items-center py-8">
          <div style={{ color: 'var(--secondary-text)' }}>
            {this.props.loadingText || 'Loading...'}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default SkeletonErrorBoundary;